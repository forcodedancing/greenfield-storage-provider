package types

import (
	"errors"

	"github.com/bnb-chain/greenfield-storage-provider/pkg/rcmgr"
	tqueue "github.com/bnb-chain/greenfield-storage-provider/pkg/taskqueue"
)

var _ tqueue.Task = &Task{}

func (m *Task) Key() tqueue.TKey {
	return ""
}

func (m *Task) Type() tqueue.TType {
	return tqueue.TypeTaskUnknown
}

func (m *Task) GetPriority() tqueue.TPriority {
	if m == nil {
		return tqueue.TPriority(tqueue.UnKnownTaskPriority)
	}
	return tqueue.TPriority(m.GetTaskPriority())
}

func (m *Task) SetPriority(prio tqueue.TPriority) {
	if m == nil {
		return
	}
	m.TaskPriority = int32(prio)
}

func (m *Task) LimitEstimate() rcmgr.Limit {
	return rcmgr.InfinitesimalLimit()
}

func (m *Task) SetCreateTime(time int64) {
	if m == nil {
		return
	}
	m.CreateTime = time
}

func (m *Task) SetUpdateTime(time int64) {
	if m == nil {
		return
	}
	m.UpdateTime = time
}

func (m *Task) SetTimeout(time int64) {
	if m == nil {
		return
	}
	m.Timeout = time
}

func (m *Task) IncRetry() bool {
	if m == nil {
		return false
	}
	m.Retry++
	return m.GetRetry() <= m.GetRetryLimit()
}

func (m *Task) DeIncRetry() {
	if m == nil {
		return
	}
	m.Retry--
	if m.Retry < 0 {
		m.Retry = 0
	}
}

func (m *Task) SetRetry(retry int64) {
	if m == nil {
		return
	}
	m.Retry = retry
}

func (m *Task) SetRetryLimit(limit int64) {
	if m == nil {
		return
	}
	m.RetryLimit = limit
}

func (m *Task) Error() error {
	if m == nil {
		return nil
	}
	if len(m.ErrMsg) == 0 {
		return nil
	}
	return errors.New(m.ErrMsg)
}

func (m *Task) SetError(err error) {
	if m == nil {
		return
	}
	m.ErrMsg = err.Error()
}
