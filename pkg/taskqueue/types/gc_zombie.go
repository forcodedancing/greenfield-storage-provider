package types

import (
	"strconv"
	"time"

	"github.com/bnb-chain/greenfield-storage-provider/pkg/rcmgr"
	tqueue "github.com/bnb-chain/greenfield-storage-provider/pkg/taskqueue"
	storagetypes "github.com/bnb-chain/greenfield/x/storage/types"
)

var _ tqueue.GCZombiePieceTask = &GCZombiePieceTask{}

func NewGCZombiePieceTask(object *storagetypes.ObjectInfo) (*GCZombiePieceTask, error) {
	if object == nil {
		return nil, ErrObjectDangling
	}
	task := &Task{
		CreateTime:   time.Now().Unix(),
		UpdateTime:   time.Now().Unix(),
		RetryLimit:   DefaultGCZombiePieceTaskRetryLimit,
		Timeout:      DefaultGCZombiePieceTimeout,
		TaskPriority: int32(GetTaskPriorityMap().GetPriority(tqueue.TypeTaskGCZombiePiece)),
	}
	return &GCZombiePieceTask{
		Task: task,
	}, nil
}

func (m *GCZombiePieceTask) Key() tqueue.TKey {
	if m == nil {
		return ""
	}
	return tqueue.TKey("GCZombie" + "-" + strconv.FormatInt(m.GetCreateTime(), 10))
}

func (m *GCZombiePieceTask) Type() tqueue.TType {
	return tqueue.TypeTaskGCZombiePiece
}

func (m *GCZombiePieceTask) LimitEstimate() rcmgr.Limit {
	limit := &rcmgr.BaseLimit{}
	prio := m.GetPriority()
	if GetTaskPriorityMap().HighLevelPriority(prio) {
		limit.TasksHighPriority = 1
	} else if GetTaskPriorityMap().LowLevelPriority(prio) {
		limit.TasksLowPriority = 1
	} else {
		limit.TasksMediumPriority = 1
	}
	return limit
}

func (m *GCZombiePieceTask) GetCreateTime() int64 {
	if m == nil {
		return 0
	}
	if m.GetTask() == nil {
		return 0
	}
	return m.GetTask().GetCreateTime()
}

func (m *GCZombiePieceTask) SetCreateTime(time int64) {
	if m == nil {
		return
	}
	if m.GetTask() == nil {
		return
	}
	m.GetTask().SetCreateTime(time)
}

func (m *GCZombiePieceTask) GetUpdateTime() int64 {
	if m == nil {
		return 0
	}
	if m.GetTask() == nil {
		return 0
	}
	return m.GetTask().GetUpdateTime()
}

func (m *GCZombiePieceTask) SetUpdateTime(time int64) {
	if m == nil {
		return
	}
	if m.GetTask() == nil {
		return
	}
	m.GetTask().SetUpdateTime(time)
}

func (m *GCZombiePieceTask) GetTimeout() int64 {
	if m == nil {
		return 0
	}
	if m.GetTask() == nil {
		return 0
	}
	return m.GetTask().GetTimeout()
}

func (m *GCZombiePieceTask) SetTimeout(time int64) {
	if m == nil {
		return
	}
	if m.GetTask() == nil {
		return
	}
	m.GetTask().SetTimeout(time)
}

func (m *GCZombiePieceTask) GetPriority() tqueue.TPriority {
	if m == nil {
		return tqueue.TPriority(0)
	}
	if m.GetTask() == nil {
		return tqueue.TPriority(0)
	}
	return m.GetTask().GetPriority()
}

func (m *GCZombiePieceTask) IncRetry() bool {
	if m == nil {
		return false
	}
	if m.GetTask() == nil {
		return false
	}
	return m.GetTask().GetRetry() <= m.GetTask().GetRetryLimit()
}

func (m *GCZombiePieceTask) DeIncRetry() {
	if m == nil {
		return
	}
	if m.GetTask() == nil {
		return
	}
	m.GetTask().DeIncRetry()
}

func (m *GCZombiePieceTask) SetPriority(prio tqueue.TPriority) {
	if m == nil {
		return
	}
	if m.GetTask() == nil {
		return
	}
	m.GetTask().SetPriority(prio)
}

func (m *GCZombiePieceTask) GetRetry() int64 {
	if m == nil {
		return 0
	}
	if m.GetTask() == nil {
		return 0
	}
	return m.GetTask().GetRetry()
}

func (m *GCZombiePieceTask) SetRetry(retry int64) {
	if m == nil {
		return
	}
	if m.GetTask() == nil {
		return
	}
	m.GetTask().SetRetry(retry)
}

func (m *GCZombiePieceTask) GetRetryLimit() int64 {
	if m == nil {
		return 0
	}
	if m.GetTask() == nil {
		return 0
	}
	return m.GetTask().GetRetryLimit()
}

func (m *GCZombiePieceTask) SetRetryLimit(limit int64) {
	if m == nil {
		return
	}
	if m.GetTask() == nil {
		return
	}
	m.GetTask().SetRetryLimit(limit)
}

func (m *GCZombiePieceTask) Error() error {
	if m == nil {
		return nil
	}
	if m.GetTask() == nil {
		return nil
	}
	return m.GetTask().Error()
}

func (m *GCZombiePieceTask) SetError(err error) {
	if m == nil {
		return
	}
	if m.GetTask() == nil {
		return
	}
	m.GetTask().SetError(err)
}

func (m *GCZombiePieceTask) GetGCZombiePieceStatus() uint64 {
	if m == nil {
		return 0
	}
	return m.GetDeleteCount()
}

func (m *GCZombiePieceTask) SetGCZombiePieceStatus(del uint64) {
	if m == nil {
		return
	}
	m.DeleteCount = del
}

func (m *GCZombiePieceTask) IsRunning() bool {
	if m == nil {
		return false
	}
	return m.GetRunning()
}

func (m *GCZombiePieceTask) StopRunning() {
	if m == nil {
		return
	}
	m.Running = false
}
