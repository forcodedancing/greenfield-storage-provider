package tasknode

import (
	"context"
	"encoding/hex"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/bnb-chain/greenfield-common/go/hash"
	"github.com/bnb-chain/greenfield-common/go/redundancy"
	storagetypes "github.com/bnb-chain/greenfield/x/storage/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	merrors "github.com/bnb-chain/greenfield-storage-provider/model/errors"
	"github.com/bnb-chain/greenfield-storage-provider/model/piecestore"
	"github.com/bnb-chain/greenfield-storage-provider/pkg/log"
	"github.com/bnb-chain/greenfield-storage-provider/service/tasknode/types"
	servicetypes "github.com/bnb-chain/greenfield-storage-provider/service/types"
	"github.com/bnb-chain/greenfield-storage-provider/store/sqldb"
)

const (
	// ReplicateFactor defines the redundancy of replication
	// TODO:: will update to （1, 2] on main net
	ReplicateFactor = 1
	// GetApprovalTimeout defines the timeout of getting secondary sp approval
	GetApprovalTimeout = 10
	// LoadSegmentBackOff defines the backoff time of loading segment
	LoadSegmentBackOff = 1
	// SealObjectTimeoutHeight defines the timeout height of listening sealobject
	SealObjectTimeoutHeight = 10
)

var _ types.TaskNodeServiceServer = &TaskNode{}

// ReplicateObject call AsyncReplicateObject non-blocking upstream services
func (taskNode *TaskNode) ReplicateObject(ctx context.Context, req *types.ReplicateObjectRequest) (
	resp *types.ReplicateObjectResponse, err error) {
	taskNode.spDB.UpdateJobState(req.GetObjectInfo().Id.Uint64(),
		servicetypes.JobState_JOB_STATE_ALLOC_SECONDARY_DOING)
	var (
		object     *storagetypes.ObjectInfo
		rCtx       *ReplicateContext
		replicates int
	)
	defer func() {
		if err != nil {
			taskNode.spDB.UpdateJobState(req.GetObjectInfo().Id.Uint64(),
				servicetypes.JobState_JOB_STATE_ALLOC_SECONDARY_ERROR)
			return
		}
		taskNode.spDB.UpdateJobState(req.GetObjectInfo().Id.Uint64(),
			servicetypes.JobState_JOB_STATE_ALLOC_SECONDARY_DONE)
	}()

	object = req.GetObjectInfo()
	if object == nil {
		log.Errorw("the object pointer dangling")
		err = merrors.ErrDanglingPointer
		return
	}
	logger := log.WithValue(ctx, "object_id", object.Id.String())

	scope, err := taskNode.rcScope.BeginSpan()
	if err != nil {
		log.CtxErrorw(logger, "failed to begin resource span", "error", err)
		return
	}
	params, err := taskNode.spDB.GetStorageParams()
	if err != nil {
		log.CtxErrorw(logger, "failed to query sp params", "error", err)
		return
	}
	replicates = int(params.GetRedundantDataChunkNum() + params.GetRedundantParityChunkNum())
	spList, err := taskNode.getApproval(object, replicates, replicates*ReplicateFactor, GetApprovalTimeout)
	if err != nil {
		log.CtxErrorw(logger, "failed to get sps to replicate", "error", err)
		return
	}
	rCtx = NewReplicateContext(logger, object, scope, params, spList)
	go taskNode.asyncReplicateObject(rCtx)
	log.CtxDebugw(rCtx.logger, "start to async replicate object to sps")
	return
}

// asyncReplicateObject replicate an object payload to other storage providers.
func (taskNode *TaskNode) asyncReplicateObject(ctx *ReplicateContext) {
	taskNode.spDB.UpdateJobState(ctx.GetObject().Id.Uint64(),
		servicetypes.JobState_JOB_STATE_REPLICATE_OBJECT_DOING)
	defer func() {
		ctx.Cancel()
		if ctx.Err() != nil {
			taskNode.spDB.UpdateJobState(ctx.GetObject().Id.Uint64(),
				servicetypes.JobState_JOB_STATE_REPLICATE_OBJECT_ERROR)
			return
		}
		taskNode.spDB.UpdateJobState(ctx.GetObject().Id.Uint64(),
			servicetypes.JobState_JOB_STATE_REPLICATE_OBJECT_DONE)
		log.CtxInfow(ctx.logger, "succeed to load and encode segments")
		taskNode.signSealObjectMsg(ctx)
	}()

	for {
		// reselect SPs and try again
		if err := ctx.Retry(); err != nil {
			log.CtxErrorw(ctx.logger, "failed to reselect sp to replicate", "error", err)
		}
		if ctx.Done() || ctx.Err() != nil {
			return
		}
		// async to load and encode segment
		go taskNode.loadAndEncodeSegment(ctx)
		for msg := range ctx.WaitReplicateMsg() {
			for i, replicateIdx := range msg.replicates {
				// skip the failed replication, wait to try again
				if ctx.HasRetryReplicate(replicateIdx) {
					continue
				}
				sp, err := ctx.GetApprovalSP(replicateIdx)
				if err != nil {
					// no need to specifically release resources, the occurrence of an error will release all resources
					log.CtxErrorw(ctx.logger, "failed to get approval sp", "replicate_idx", replicateIdx, "error", err)
					return
				}
				checksum := hash.GenerateChecksum(msg.data[i])
				err = sp.Gateway().ReplicatePieceData(ctx.GetObject(), replicateIdx,
					msg.segment, sp.Approval(), checksum, msg.data[i])
				ctx.ReleasePieceResource(msg.segment, 1)
				if err != nil {
					log.CtxErrorw(ctx.logger, "failed to replicate object to sp",
						"domain", sp.SP().GetEndpoint(), "error", err)
					// add replicate idx to retry list
					ctx.AddRetryReplicate(replicateIdx)
				}
				ctx.ReleasePieceResource(msg.segment, 1)
				log.CtxDebugw(ctx.logger, "succeed to replicate one segment", "segment_idx", replicateIdx)
			}
		}
	}
}

// loadAndEncodeSegment loads segments from piece store and encodes the segments according to redundancy type
func (taskNode *TaskNode) loadAndEncodeSegment(ctx *ReplicateContext) {
	var (
		err        error
		msg        = &ReplicateMessage{}
		segmentCnt = piecestore.ComputeSegmentCount(ctx.GetObject().GetPayloadSize(),
			ctx.GetStorageParams().GetMaxSegmentSize())
		replicates = int(ctx.GetStorageParams().GetRedundantDataChunkNum() +
			ctx.GetStorageParams().GetRedundantParityChunkNum())
	)
	defer func() {
		if err != nil {
			ctx.TerminateReplicate(err)
		}
	}()

	for segIdx := uint32(0); segIdx < segmentCnt; segIdx++ {
		key := piecestore.EncodeSegmentPieceKey(ctx.GetObject().Id.Uint64(), segIdx)
		var (
			segmentData     []byte
			encodeData      [][]byte
			releasePieceCnt int
		)
		err = ctx.ReservePieceResource(segIdx, replicates)
		if err != nil {
			log.CtxErrorw(ctx.logger, "failed to reserve resource", "error", err)
			return
		}
		for retry := 0; retry < ctx.GetLoadSegmentRetry(); retry++ {
			err = nil
			segmentData, err = taskNode.pieceStore.GetSegment(context.Background(), key, 0, 0)
			if err != nil {
				// TODO:: add backoff strategy
				time.Sleep(time.Duration(LoadSegmentBackOff) * time.Second)
				continue
			}
		}
		if err != nil {
			log.CtxErrorw(ctx.logger, "failed to load segment from piece store", "error", err)
			return
		}
		msg.segment = segIdx
		for doing, _ := range ctx.GetReplicate() {
			msg.replicates = append(msg.replicates, doing)
		}
		switch ctx.GetObject().GetRedundancyType() {
		case storagetypes.REDUNDANCY_EC_TYPE:
			encodeData, err = redundancy.EncodeRawSegment(segmentData,
				int(ctx.GetStorageParams().GetRedundantDataChunkNum()),
				int(ctx.GetStorageParams().GetRedundantParityChunkNum()))
			if err != nil {
				log.CtxErrorw(ctx.logger, "failed to encode segment", "error", err)
				return
			}
			for ecIdx, data := range encodeData {
				if _, ok := ctx.GetReplicate()[uint32(ecIdx)]; !ok {
					releasePieceCnt++
					continue
				}
				msg.data = append(msg.data, data)
			}
			ctx.NotifyReplicate(msg)
		case storagetypes.REDUNDANCY_REPLICA_TYPE:
			for i := 0; i < len(msg.replicates); i++ {
				msg.data = append(msg.data, segmentData)
			}
			releasePieceCnt = replicates - len(msg.replicates)
			ctx.NotifyReplicate(msg)
		default:
			log.CtxErrorw(ctx.logger, "unsupported redundancy type",
				"redundancy_type", ctx.GetObject().GetRedundancyType())
			err = merrors.ErrUnsupportedRedundancyType
			return
		}
		if releasePieceCnt > 0 {
			ctx.ReleasePieceResource(segIdx, releasePieceCnt)
		}
	}
	ctx.FinishReplicate()
}

// signSealObjectMsg collects the replicate SPs' integrityHash and signature, and verify the signatures,
// use to organize into MsgSealObject
func (taskNode *TaskNode) signSealObjectMsg(ctx *ReplicateContext) {
	var (
		sealMsg    = &storagetypes.MsgSealObject{}
		replicates = ctx.GetStorageParams().GetRedundantDataChunkNum() +
			ctx.GetStorageParams().GetRedundantParityChunkNum()
		err error
	)
	sealMsg.Operator = taskNode.config.SpOperatorAddress
	sealMsg.BucketName = ctx.GetObject().GetBucketName()
	sealMsg.ObjectName = ctx.GetObject().GetObjectName()
	sealMsg.SecondarySpAddresses = make([]string, replicates)
	sealMsg.SecondarySpSignatures = make([][]byte, replicates)
	taskNode.spDB.UpdateJobState(ctx.GetObject().Id.Uint64(),
		servicetypes.JobState_JOB_STATE_SIGN_OBJECT_DOING)

	defer func() {
		if err != nil {
			taskNode.spDB.UpdateJobState(ctx.GetObject().Id.Uint64(),
				servicetypes.JobState_JOB_STATE_SIGN_OBJECT_ERROR)
			return
		}
		taskNode.spDB.UpdateJobState(ctx.GetObject().Id.Uint64(),
			servicetypes.JobState_JOB_STATE_SIGN_OBJECT_DONE)
		taskNode.sealObject(ctx, sealMsg)
	}()

	for rIdx := uint32(0); rIdx < replicates; rIdx++ {
		approval, err := ctx.GetApprovalSP(rIdx)
		if err != nil {
			ctx.TerminateReplicate(err)
			return
		}
		integrityHash, signature, err := approval.Gateway().GetReplicateIntegrityHash(ctx.GetObject())
		if err != nil {
			return
		}
		log.CtxDebugw(ctx.logger, "receive the sp response", "replica_idx", rIdx, "domain", approval.SP().GetEndpoint(),
			"integrity_hash", hex.EncodeToString(integrityHash), "signature", hex.EncodeToString(signature))

		// verify the signature
		msg := storagetypes.NewSecondarySpSignDoc(approval.SP().GetOperator(),
			sdkmath.NewUint(ctx.GetObject().Id.Uint64()), integrityHash).GetSignBytes()
		approvalAddr, err := sdk.AccAddressFromHexUnsafe(approval.SP().GetApprovalAddress())
		if err != nil {
			log.CtxErrorw(ctx.logger, "failed to parser sp operator address", "endpoint", approval.SP().GetEndpoint(),
				"sp", approval.SP().GetApprovalAddress(), "error", err)
			return
		}
		err = storagetypes.VerifySignature(approvalAddr, sdk.Keccak256(msg), signature)
		if err != nil {
			log.CtxErrorw(ctx.logger, "failed to verify sp signature", "endpoint", approval.SP().GetEndpoint(),
				"sp", approval.SP().GetApprovalAddress(), "error", err)
			return
		}
		sealMsg.GetSecondarySpAddresses()[rIdx] = approval.SP().GetOperator().String()
		sealMsg.GetSecondarySpSignatures()[rIdx] = signature
	}
	log.CtxErrorw(ctx.logger, "succeed to sign the seal object msg by signer")
}

// sealObject sends MsgSealObject to signer and broadcasts to the chain
func (taskNode *TaskNode) sealObject(ctx *ReplicateContext, msg *storagetypes.MsgSealObject) {
	taskNode.spDB.UpdateJobState(ctx.GetObject().Id.Uint64(),
		servicetypes.JobState_JOB_STATE_SEAL_OBJECT_DOING)
	var err error
	defer func() {
		if err != nil {
			taskNode.spDB.UpdateJobState(ctx.GetObject().Id.Uint64(),
				servicetypes.JobState_JOB_STATE_SEAL_OBJECT_ERROR)
			return
		}
		taskNode.spDB.UpdateJobState(ctx.GetObject().Id.Uint64(),
			servicetypes.JobState_JOB_STATE_SEAL_OBJECT_DONE)
	}()

	// seal object
	_, err = taskNode.signer.SealObjectOnChain(ctx.logger, msg)
	if err != nil {
		log.CtxErrorw(ctx.logger, "failed to sign object by signer", "error", err)
		return
	}
	// wait the result of sealing object
	err = taskNode.chain.ListenObjectSeal(ctx.logger, ctx.GetObject().GetBucketName(),
		ctx.GetObject().GetObjectName(), SealObjectTimeoutHeight)
	if err != nil {
		log.CtxErrorw(ctx.logger, "failed to seal object on chain", "error", err)
		return
	}
	log.CtxErrorw(ctx.logger, "succeed to seal object on chain")
}

// getApproval ask approval whether willing to store the payload data by P2P protocol
func (taskNode *TaskNode) getApproval(objectInfo *storagetypes.ObjectInfo, low, high int, timeout int64) ([]*ApprovalSP, error) {
	var spList []*ApprovalSP
	approvals, _, err := taskNode.p2p.GetApproval(context.Background(), objectInfo, int64(high), timeout)
	if err != nil {
		return spList, err
	}
	if len(approvals) < low {
		return spList, merrors.ErrInsufficientApprovalSP
	}
	for spOpAddr, approval := range approvals {
		sp, err := taskNode.spDB.GetSpByAddress(spOpAddr, sqldb.OperatorAddressType)
		if err != nil {
			log.Errorw("failed to get sp info from db", "object_id", objectInfo.Id.Uint64(), "error", err)
			continue
		}
		approvalSP, err := NewApprovalSP(sp, approval)
		if err != nil {
			log.Errorw("failed to new approval sp", "object_id", objectInfo.Id.Uint64(), "error", err)
			continue
		}
		spList = append(spList, approvalSP)
	}
	if len(spList) < low {
		return spList, merrors.ErrInsufficientSP
	}
	return spList, nil
}
