package manager

import (
	"sync"
	"sync/atomic"

	merrors "github.com/bnb-chain/greenfield-storage-provider/model/errors"
	"github.com/bnb-chain/greenfield-storage-provider/pkg/log"
	"github.com/bnb-chain/greenfield-storage-provider/pkg/rcmgr"
	tqueue "github.com/bnb-chain/greenfield-storage-provider/pkg/taskqueue"
	"github.com/bnb-chain/greenfield-storage-provider/pkg/taskqueue/queue"
	tqueuetypes "github.com/bnb-chain/greenfield-storage-provider/pkg/taskqueue/types"
)

const (
	PriorityQueueManage    = "priority-queue-manager"
	SubQueueUpload         = "sub-queue-upload"
	SubQueueReplicatePiece = "sub-queue-replicate-piece"
	SubQueueSealObject     = "sub-queue-seal-object"
	SubQueueGCObject       = "sub-queue-gc-object"

	SubQueueUploadCap         = 1024
	SubQueueReplicatePieceCap = 1024
	SubQueueSealPieceCap      = 1024
	SubQueueGCObjectCap       = 4
)

type MPQueue struct {
	pqueue      tqueue.TPriorityQueueWithLimit
	supportType map[tqueue.TType]struct{}
	gcRunning   atomic.Int32
	mux         sync.RWMutex
}

func NewMPQueue() *MPQueue {
	uploadStrategy := queue.NewStrategy()
	uploadStrategy.SetCollectionFunc(queue.DefaultGCTasksByTimeout)
	uploadQueue := queue.NewTaskQueue(SubQueueUpload, SubQueueUploadCap, uploadStrategy, false)

	replicateStrategy := queue.NewStrategy()
	replicateStrategy.SetCollectionFunc(queue.DefaultGCTasksByRetry)
	replicateQueue := queue.NewTaskQueue(SubQueueReplicatePiece, SubQueueReplicatePieceCap, replicateStrategy, true)

	sealStrategy := queue.NewStrategy()
	sealStrategy.SetCollectionFunc(queue.DefaultGCTasksByRetry)
	sealQueue := queue.NewTaskQueue(SubQueueSealObject, SubQueueSealPieceCap, sealStrategy, true)

	gcObjectStrategy := queue.NewStrategy()
	gcObjectStrategy.SetCollectionFunc(queue.DefaultGCTasksByRetry)
	gcObjectQueue := queue.NewTaskQueue(SubQueueGCObject, SubQueueGCObjectCap, gcObjectStrategy, true)

	subQueues := map[tqueue.TPriority]tqueue.TQueueWithLimit{
		tqueuetypes.GetTaskPriorityMap().GetPriority(tqueue.TypeTaskUpload):         uploadQueue,
		tqueuetypes.GetTaskPriorityMap().GetPriority(tqueue.TypeTaskReplicatePiece): replicateQueue,
		tqueuetypes.GetTaskPriorityMap().GetPriority(tqueue.TypeTaskSealObject):     sealQueue,
		tqueuetypes.GetTaskPriorityMap().GetPriority(tqueue.TypeTaskGCObject):       gcObjectQueue,
	}

	mstrategy := queue.NewStrategy()
	mstrategy.SetPickUpFunc(queue.DefaultPickUpTaskByPriority)
	pqueue := queue.NewTaskPriorityQueue(PriorityQueueManage, subQueues, mstrategy, true)

	supportTask := map[tqueue.TType]struct{}{
		tqueue.TypeTaskUpload:         struct{}{},
		tqueue.TypeTaskReplicatePiece: struct{}{},
		tqueue.TypeTaskSealObject:     struct{}{},
		tqueue.TypeTaskGCObject:       struct{}{},
	}
	return &MPQueue{
		pqueue:      pqueue,
		supportType: supportTask,
	}
}

func (q *MPQueue) PopTaskWithLimit(limit rcmgr.Limit) tqueue.Task {
	return q.pqueue.PopWithLimit(limit)
}

func (q *MPQueue) PushTask(task tqueue.Task) error {
	q.mux.RLock()
	_, ok := q.supportType[task.Type()]
	q.mux.RUnlock()
	if !ok {
		return merrors.ErrUnsupportedTaskType
	}
	return q.pqueue.Push(task)
}

func (q *MPQueue) PopTask(key tqueue.TKey) tqueue.Task {
	return q.pqueue.PopWithKey(key)
}

func (q *MPQueue) GCTask() {
	if q.gcRunning.Swap(1) == 1 {
		log.Debugw("manager priority queue gc is running")
		return
	}
	defer q.gcRunning.Swap(0)
	q.pqueue.RunStrategy()
}

func (q *MPQueue) GetUploadingTasksCount() int {
	return q.getUploadPrimaryTasksCount() + q.getReplicatePieceTasksCount() + q.getSealObjectTasksCount()
}

func (q *MPQueue) getUploadPrimaryTasksCount() int {
	prio := tqueuetypes.GetTaskPriorityMap().GetPriority(tqueue.TypeTaskUpload)
	return q.pqueue.SubQueueLen(prio)
}

func (q *MPQueue) getReplicatePieceTasksCount() int {
	prio := tqueuetypes.GetTaskPriorityMap().GetPriority(tqueue.TypeTaskReplicatePiece)
	return q.pqueue.SubQueueLen(prio)
}

func (q *MPQueue) getSealObjectTasksCount() int {
	prio := tqueuetypes.GetTaskPriorityMap().GetPriority(tqueue.TypeTaskSealObject)
	return q.pqueue.SubQueueLen(prio)
}
