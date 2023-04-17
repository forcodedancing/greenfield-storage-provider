package manager

import (
	"context"
	"time"

	merrors "github.com/bnb-chain/greenfield-storage-provider/model/errors"
	"github.com/bnb-chain/greenfield-storage-provider/pkg/log"
	tqueuetypes "github.com/bnb-chain/greenfield-storage-provider/pkg/taskqueue/types"
	"github.com/bnb-chain/greenfield-storage-provider/service/manager/types"
)

var _ types.ManagerServiceServer = &Manager{}

// AskUploadObject asks to create object to SP manager.
func (m *Manager) AskUploadObject(ctx context.Context, req *types.AskUploadObjectRequest) (
	*types.AskUploadObjectResponse, error) {
	uploading := m.pqueue.GetUploadingTasksCount()
	allow := uploading < m.config.UploadNumber
	resp := &types.AskUploadObjectResponse{
		Allow: allow,
	}
	if !allow {
		log.CtxWarnw(ctx, "refuse upload object", "uploading", uploading, "limit", m.config.UploadNumber)
	}
	return resp, nil
}

// CreateUploadObjectTask asks to upload object to SP manager.
func (m *Manager) CreateUploadObjectTask(ctx context.Context, req *types.CreateUploadObjectTaskRequest) (
	*types.CreateUploadObjectTaskResponse, error) {
	resp := &types.CreateUploadObjectTaskResponse{}
	task := req.GetUploadObjectTask()
	if task == nil {
		return resp, merrors.ErrDanglingTaskPointer
	}
	ctx = log.WithValue(ctx, "object_id", string(task.Key()))
	err := m.pqueue.PushTask(task)
	if err != nil {
		log.CtxErrorw(ctx, "failed to create upload object task", "error", err)
	}
	return resp, err
}

// DoneUploadObjectTask notifies the manager the upload object task has been done.
func (m *Manager) DoneUploadObjectTask(ctx context.Context, req *types.DoneUploadObjectTaskRequest) (
	*types.DoneUploadObjectTaskResponse, error) {
	resp := &types.DoneUploadObjectTaskResponse{}
	task := req.GetUploadObjectTask()
	if task == nil {
		return resp, merrors.ErrDanglingTaskPointer
	}
	ctx = log.WithValue(ctx, "object_id", string(task.Key()))
	m.pqueue.PopTask(task.Key())
	replicateTask := task.TransferReplicatePieceTask()
	err := m.pqueue.PushTask(replicateTask)
	if err != nil {
		log.CtxErrorw(ctx, "failed to push replicate task to queue", "error", err)
	} else {
		log.CtxDebugw(ctx, "succeed to push replicate task to queue")
	}
	return resp, err
}

// DoneReplicatePieceTask notifies the manager the replicate piece task has been done.
func (m *Manager) DoneReplicatePieceTask(ctx context.Context, req *types.DoneReplicatePieceTaskRequest) (
	*types.DoneReplicatePieceTaskResponse, error) {
	resp := &types.DoneReplicatePieceTaskResponse{}
	task := req.GetReplicatePieceTask()
	if task == nil {
		return resp, merrors.ErrDanglingTaskPointer
	}
	ctx = log.WithValue(ctx, "object_id", string(task.Key()))
	m.pqueue.PopTask(task.Key())
	sealTask := task.TransferSealObjectTask()
	err := m.pqueue.PushTask(sealTask)
	if err != nil {
		log.CtxErrorw(ctx, "failed to push seal object task to queue", "error", err)
	} else {
		log.CtxDebugw(ctx, "succeed to push seal object task to queue")
	}
	return resp, err
}

// DoneSealObjectTask notifies the manager the seal object task has been done.
func (m *Manager) DoneSealObjectTask(ctx context.Context, req *types.DoneSealObjectTaskRequest) (
	*types.DoneSealObjectTaskResponse, error) {
	resp := &types.DoneSealObjectTaskResponse{}
	task := req.GetSealObjectTask()
	if task == nil {
		return resp, merrors.ErrDanglingTaskPointer
	}
	ctx = log.WithValue(ctx, "object_id", string(task.Key()))
	m.pqueue.PopTask(task.Key())
	log.CtxDebugw(ctx, "succeed to seal object on chain")
	return resp, nil
}

// AskTask asks the task to execute
func (m *Manager) AskTask(ctx context.Context, req *types.AskTaskRequest) (*types.AskTaskResponse, error) {
	resp := &types.AskTaskResponse{}
	limit := req.GetLimit().TransferRcmgrLimits()
	task := m.pqueue.PopTaskWithLimit(limit)
	if task == nil {
		return resp, nil
	}
	defer func() {
		task.SetUpdateTime(time.Now().Unix())
		m.pqueue.PushTask(task)
	}()
	switch t := (task).(type) {
	case *tqueuetypes.ReplicatePieceTask:
		resp.Task = &types.AskTaskResponse_ReplicatePieceTask{
			ReplicatePieceTask: t,
		}
	case *tqueuetypes.SealObjectTask:
		resp.Task = &types.AskTaskResponse_SealObjectTask{
			SealObjectTask: t,
		}
	case *tqueuetypes.GCObjectTask:
		resp.Task = &types.AskTaskResponse_GcObjectTask{
			GcObjectTask: t,
		}
	default:
		log.CtxErrorw(ctx, "task node does not support task type", "task_type", task.Type())
		return resp, merrors.ErrUnsupportedDispatchTaskType
	}
	resp.HasTask = true
	return resp, nil
}
