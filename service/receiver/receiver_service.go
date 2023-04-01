package receiver

import (
	"bytes"
	"context"
	"io"

	"github.com/bnb-chain/greenfield-common/go/hash"
	"github.com/bnb-chain/greenfield-storage-provider/model/piecestore"
	"github.com/bnb-chain/greenfield-storage-provider/pkg/rcmgr"
	storagetypes "github.com/bnb-chain/greenfield/x/storage/types"

	merrors "github.com/bnb-chain/greenfield-storage-provider/model/errors"
	"github.com/bnb-chain/greenfield-storage-provider/pkg/log"
	"github.com/bnb-chain/greenfield-storage-provider/service/receiver/types"
)

var _ types.ReceiverServiceServer = &Receiver{}

// ReceivePiece receives a piece data and store it into piece store.
func (receiver *Receiver) ReceivePiece(stream types.ReceiverService_ReceivePieceServer) error {
	var (
		init   = true
		req    *types.ReceivePieceRequest
		object *storagetypes.ObjectInfo
		ctx    context.Context
		data   []byte
	)
	scope, err := receiver.rcScope.BeginSpan()
	if err != nil {
		return err
	}
	defer func() {
		scope.Done()
	}()

	for {
		req, err = stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Debugw("receive payload exception", "error", err)
			return err
		}
		if init {
			init = false
			object = req.GetObjectInfo()
			if object == nil {
				log.Errorw("object pointer dangling")
				return merrors.ErrDanglingPointer
			}
			ctx = log.WithValue(ctx, "object_id", object.Id.String())
		}
		err = scope.ReserveMemory(len(req.GetPieceData()), rcmgr.ReservationPriorityAlways)
		if err != nil {
			log.CtxErrorw(ctx, "failed to reserve memory", "error", err)
			return err
		}
		data = append(data, req.GetPieceData()...)
	}
	checksum := hash.GenerateChecksum(data)
	// TODO:: write SPDB, key is objectID, segmentIdx, replicateIdx
	if !bytes.Equal(checksum, req.GetChecksum()) {
		return merrors.ErrChecksumMismatch
	}
	log.CtxDebugw(ctx, "succeed to receive piece data", "segment_idx", req.GetSegmentIdx(),
		"replicate_idx", req.GetReplicateIdx())
	return nil
}

// DoneReplicate returns the integrity hash of all pieces and the signature of signing the integrity hash
func (receiver *Receiver) DoneReplicate(ctx context.Context, req *types.DoneReplicateRequest) (
	*types.DoneReplicateResponse, error) {
	var (
		// TODO:: get all pieces checksums from SPDB by order
		checksums  [][]byte
		replicates int
		object     = req.GetObjectInfo()
		resp       = &types.DoneReplicateResponse{}
	)
	if object == nil {
		log.Errorw("object pointer dangling")
		return resp, merrors.ErrDanglingPointer
	}
	ctx = log.WithValue(ctx, "object_id", object.Id.String())
	params, err := receiver.spDB.GetStorageParams()
	if err != nil {
		log.CtxErrorw(ctx, "failed to get storage params from db", "error", err)
		return resp, err
	}
	replicates = int(piecestore.ComputeSegmentCount(object.GetPayloadSize(), params.GetMaxSegmentSize()))
	if replicates != len(checksums) {
		return resp, merrors.ErrPieceCntMismatch
	}
	integrityHash, signature, err := receiver.signer.SignIntegrityHash(ctx, object.Id.Uint64(), checksums)
	if err != nil {
		log.CtxErrorw(ctx, "failed to sign integrity hash", "error", err)
		return resp, err
	}
	resp.IntegrityHash = integrityHash
	resp.Signature = signature
	return resp, nil
}
