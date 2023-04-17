package queue

import (
	"math/rand"
	"sort"
	"sync"
	"time"

	tqueue "github.com/bnb-chain/greenfield-storage-provider/pkg/taskqueue"
)

var _ tqueue.TQueueStrategy = &Strategy{}

type Strategy struct {
	pick func([]tqueue.Task) tqueue.Task
	gc   func(tqueue.TQueueOnStrategy, []tqueue.TKey)
	mux  sync.RWMutex
}

func NewStrategy() tqueue.TQueueStrategy {
	return &Strategy{
		pick: nil,
		gc:   nil,
	}
}

func (s *Strategy) SetPickUpFunc(pick func([]tqueue.Task) tqueue.Task) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.pick = pick
}

func (s *Strategy) SetCollectionFunc(gc func(queue tqueue.TQueueOnStrategy, keys []tqueue.TKey)) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.gc = gc
}

func (s *Strategy) RunPickUp(tasks []tqueue.Task) tqueue.Task {
	s.mux.RLock()
	defer s.mux.RUnlock()
	if s.pick == nil {
		// TODO:: add warn log
		return nil
	}
	return s.pick(tasks)
}

func (s *Strategy) RunCollection(queue tqueue.TQueueOnStrategy, keys []tqueue.TKey) {
	s.mux.RLock()
	defer s.mux.RUnlock()
	if s.gc == nil {
		// TODO:: add warn log
		return
	}
	s.gc(queue, keys)
}

func DefaultGCTasksByRetry(queue tqueue.TQueueOnStrategy, keys []tqueue.TKey) {
	for _, key := range keys {
		if queue.RetryFailed(key) && queue.Timeout(key) {
			queue.PopWithKey(key)
		}
	}
}

func DefaultGCTasksByTimeout(queue tqueue.TQueueOnStrategy, keys []tqueue.TKey) {
	for _, key := range keys {
		if queue.Timeout(key) {
			queue.PopWithKey(key)
		}
	}
}

func DefaultPickUpTaskByPriority(tasks []tqueue.Task) tqueue.Task {
	if len(tasks) == 0 {
		return nil
	}
	if len(tasks) == 1 {
		return tasks[0]
	}
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].GetPriority() < tasks[j].GetPriority()
	})
	var totalPrio int
	for _, task := range tasks {
		totalPrio += int(task.GetPriority())
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	var task tqueue.Task
	totalPrio = 0
	randPrio := r.Intn(totalPrio + 1)
	for _, t := range tasks {
		totalPrio += int(t.GetPriority())
		if totalPrio >= randPrio {
			task = t
			break
		}
	}
	return task
}
