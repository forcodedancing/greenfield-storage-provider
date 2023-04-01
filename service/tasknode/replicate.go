package tasknode

import (
	"context"
	"errors"
	"sync"

	sptypes "github.com/bnb-chain/greenfield/x/sp/types"
	storagetypes "github.com/bnb-chain/greenfield/x/storage/types"

	"github.com/bnb-chain/greenfield-storage-provider/model"
	merrors "github.com/bnb-chain/greenfield-storage-provider/model/errors"
	"github.com/bnb-chain/greenfield-storage-provider/model/piecestore"
	"github.com/bnb-chain/greenfield-storage-provider/pkg/log"
	p2ptypes "github.com/bnb-chain/greenfield-storage-provider/pkg/p2p/types"
	"github.com/bnb-chain/greenfield-storage-provider/pkg/rcmgr"
	gatewayclient "github.com/bnb-chain/greenfield-storage-provider/service/gateway/client"
)

const (
	// LoadSegmentRetry defines the retry number of loading segment
	LoadSegmentRetry = 3
)

// ReplicateMessage defines the replication info between load/encode and replicate goroutines
type ReplicateMessage struct {
	segment    uint32
	replicates []uint32
	data       [][]byte
}

// ApprovalSP defines the approved SP, including SP info, accepted approval and SP gateway client
type ApprovalSP struct {
	sp       *sptypes.StorageProvider
	approval *p2ptypes.GetApprovalResponse
	gateway  *gatewayclient.GatewayClient
}

// NewApprovalSP return an instance of ApprovalSP
func NewApprovalSP(sp *sptypes.StorageProvider, approval *p2ptypes.GetApprovalResponse) (*ApprovalSP, error) {
	if sp == nil || approval == nil {
		return nil, errors.New("sp or approval is nil")
	}
	gateway, err := gatewayclient.NewGatewayClient(sp.GetEndpoint())
	if err != nil {
		return nil, err
	}
	approvalSP := &ApprovalSP{
		sp:       sp,
		approval: approval,
		gateway:  gateway,
	}
	return approvalSP, nil
}

// SP returns the approved SP info
func (sp *ApprovalSP) SP() *sptypes.StorageProvider {
	return sp.sp
}

// Approval returns the approval of the approved SP
func (sp *ApprovalSP) Approval() *p2ptypes.GetApprovalResponse {
	return sp.approval
}

// Gateway return the gateway client of the approved SP
func (sp *ApprovalSP) Gateway() *gatewayclient.GatewayClient {
	return sp.gateway
}

// ReplicateContext defines the context of replicating object
type ReplicateContext struct {
	object        *storagetypes.ObjectInfo
	scope         rcmgr.ResourceScopeSpan
	storageParams *storagetypes.Params

	sp       map[uint32]*ApprovalSP
	backupSP []*ApprovalSP

	replicating    map[uint32]struct{} // the replicate idx being replicated
	replicateRetry map[uint32]struct{} // the replicate idx has been failed, and waiting to try again
	replicateCh    chan *ReplicateMessage
	mux            sync.RWMutex

	loadSegmentRetry int
	logger           context.Context
	innerErr         error
}

// NewReplicateContext returns an instance of ReplicateContext
func NewReplicateContext(logger context.Context, object *storagetypes.ObjectInfo, scope rcmgr.ResourceScopeSpan,
	params *storagetypes.Params, sp []*ApprovalSP) *ReplicateContext {
	ctx := &ReplicateContext{
		object:           object,
		scope:            scope,
		storageParams:    params,
		backupSP:         sp,
		replicating:      make(map[uint32]struct{}),
		replicateRetry:   make(map[uint32]struct{}),
		loadSegmentRetry: LoadSegmentRetry,
		logger:           logger,
	}
	for i := uint32(0); i < ctx.storageParams.GetRedundantDataChunkNum()+
		ctx.storageParams.GetRedundantParityChunkNum(); i++ {
		ctx.replicateRetry[i] = struct{}{}
	}
	return ctx
}

// Retry change the replicateRetry set to the replicating set, clears replicateRetry,
// and picks up approved SP to replicate
func (ctx *ReplicateContext) Retry() error {
	ctx.mux.Lock()
	defer ctx.mux.Unlock()
	// has been finished
	if len(ctx.replicateRetry) == 0 {
		return ctx.innerErr
	}
	ctx.replicating = ctx.replicateRetry
	ctx.replicateRetry = make(map[uint32]struct{})
	if len(ctx.backupSP) < len(ctx.replicating) {
		ctx.innerErr = merrors.ErrDepletedSP
		return ctx.innerErr
	}
	// delete the information of successful replication
	for replicateIdx, _ := range ctx.sp {
		if _, ok := ctx.replicating[replicateIdx]; !ok {
			delete(ctx.sp, replicateIdx)
			continue
		}
	}
	// pick up SP from backup SP set
	for replicateIdx, _ := range ctx.replicating {
		ctx.sp[replicateIdx] = ctx.backupSP[0]
		ctx.backupSP = ctx.backupSP[1:]
	}
	ctx.replicateCh = make(chan *ReplicateMessage)
	return ctx.innerErr
}

// Done returns an indication of whether the replication is finished
func (ctx *ReplicateContext) Done() bool {
	ctx.mux.RLock()
	defer ctx.mux.RUnlock()
	return len(ctx.replicateRetry) == 0
}

// Err return the inner error of replicate context
func (ctx *ReplicateContext) Err() error {
	ctx.mux.RLock()
	defer ctx.mux.RUnlock()
	return ctx.innerErr
}

// TerminateReplicate terminates the replication and assigns err to replicate context inner err
func (ctx *ReplicateContext) TerminateReplicate(err error) {
	ctx.mux.Lock()
	defer ctx.mux.Unlock()
	ctx.innerErr = err
	close(ctx.replicateCh)
}

// FinishReplicate closes the channel of transferring replicate msg, WaitReplicateMsg receives the close notification
func (ctx *ReplicateContext) FinishReplicate() {
	close(ctx.replicateCh)
}

// NotifyReplicate transfers the msg by the channel and WaitReplicateMsg receives the msg
func (ctx *ReplicateContext) NotifyReplicate(msg *ReplicateMessage) {
	ctx.replicateCh <- msg
}

// WaitReplicateMsg returns the channel of transferring replicate msg
func (ctx *ReplicateContext) WaitReplicateMsg() chan *ReplicateMessage {
	return ctx.replicateCh
}

// Cancel release the resource of replication object
func (ctx *ReplicateContext) Cancel() {
	ctx.mux.Lock()
	defer ctx.mux.Unlock()
	ctx.scope.Done()
	close(ctx.replicateCh)
}

// GetObject return the object info of replication
func (ctx *ReplicateContext) GetObject() *storagetypes.ObjectInfo {
	return ctx.object
}

// GetStorageParams returns the storage params of the replication
func (ctx *ReplicateContext) GetStorageParams() *storagetypes.Params {
	return ctx.storageParams
}

// GetLoadSegmentRetry returns the retry number of loading segment from piece store
func (ctx *ReplicateContext) GetLoadSegmentRetry() int {
	return ctx.loadSegmentRetry
}

// GetReplicateDataSize return the piece data size of segment piece or ec piece
func (ctx *ReplicateContext) GetReplicateDataSize(idx uint32) (size uint64, err error) {
	if ctx.object.GetRedundancyType() == storagetypes.REDUNDANCY_EC_TYPE {
		size, err = piecestore.GetPieceDataSize(ctx.object, ctx.storageParams.GetMaxSegmentSize(),
			ctx.storageParams.GetRedundantDataChunkNum(), model.PieceType, idx)
	} else {
		size, err = piecestore.GetPieceDataSize(ctx.object, ctx.storageParams.GetMaxSegmentSize(),
			ctx.storageParams.GetRedundantDataChunkNum(), model.SegmentType, idx)
	}
	return
}

// GetApprovalSP returns the approved SP by replicate idx
func (ctx *ReplicateContext) GetApprovalSP(idx uint32) (*ApprovalSP, error) {
	ctx.mux.RLock()
	defer ctx.mux.RUnlock()
	sp, ok := ctx.sp[idx]
	if !ok {
		ctx.innerErr = merrors.ErrDanglingSP
		return nil, ctx.innerErr
	}
	return sp, nil
}

// GetReplicate returns the set of replicating idx
func (ctx *ReplicateContext) GetReplicate() map[uint32]struct{} {
	ctx.mux.Lock()
	defer ctx.mux.Unlock()
	shadow := make(map[uint32]struct{})
	for r, s := range ctx.replicating {
		shadow[r] = s
	}
	return shadow
}

// AddRetryReplicate adds the replicate idx to the retry set, waits try again
func (ctx *ReplicateContext) AddRetryReplicate(idx uint32) {
	ctx.mux.Lock()
	defer ctx.mux.Unlock()
	ctx.replicateRetry[idx] = struct{}{}
}

// HasRetryReplicate returns an indication of whether the replicate idx among replicate retry set
func (ctx *ReplicateContext) HasRetryReplicate(idx uint32) bool {
	ctx.mux.RLock()
	defer ctx.mux.RUnlock()
	_, ok := ctx.replicateRetry[idx]
	return ok
}

// ReservePieceResource reserves the memory from resource manager according to object info and the cnt
func (ctx *ReplicateContext) ReservePieceResource(segIdx uint32, cnt int) error {
	size, _ := ctx.GetReplicateDataSize(segIdx)
	err := ctx.scope.ReserveMemory(int(size)*cnt, rcmgr.ReservationPriorityAlways)
	if err != nil {
		return err
	}
	var state string
	rcmgr.ResrcManager().ViewSystem(func(scope rcmgr.ResourceScope) error {
		state = scope.Stat().String()
		return nil
	})
	log.CtxDebugw(ctx.logger, "reserve memory from resource manager ", "reserve_size", size,
		"system state", state)
	return nil
}

// ReleasePieceResource releases the memory from resource manager according to object info and the cnt
func (ctx *ReplicateContext) ReleasePieceResource(segIdx uint32, cnt int) {
	size, _ := ctx.GetReplicateDataSize(segIdx)
	ctx.scope.ReleaseMemory(int(size) * cnt)
	var state string
	rcmgr.ResrcManager().ViewSystem(func(scope rcmgr.ResourceScope) error {
		state = scope.Stat().String()
		return nil
	})
	log.CtxDebugw(ctx.logger, "release memory to resource manager ", "release_size", size,
		"system state", state)
}
