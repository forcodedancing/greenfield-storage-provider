package task

import (
	"github.com/bnb-chain/greenfield-storage-provider/pkg/rcmgr"
	storagetypes "github.com/bnb-chain/greenfield/x/storage/types"
)

type TKey string

type Task interface {
	Key() TKey
	Type() TType
	GetCreateTime() int64
	SetCreateTime(int64)
	GetUpdateTime() int64
	SetUpdateTime(int64)
	GetTimeout() int64
	SetTimeout(int64)
	GetPriority() TPriority
	SetPriority(TPriority)
	IncRetry() bool
	DeIncRetry()
	GetRetry() int64
	SetRetry(int64)
	GetRetryLimit() int64
	SetRetryLimit(int64)
	LimitEstimate() rcmgr.Limit
	Error() error
	SetError(err error)
}

type TType int32

const (
	TypeTaskUnknown TType = iota
	TypeTaskUpload
	TypeTaskReplicatePiece
	TypeTaskSealObject
	TypeTaskReceivePiece
	TypeTaskDownloadObject
	TypeTaskGCObject
	TypeTaskGCZombiePiece
	TypeTaskGCStore
)

type ObjectTask interface {
	Task
	GetObject() *storagetypes.ObjectInfo
}

type UploadObjectTask interface {
	ObjectTask
	TransferReplicatePieceTask() ReplicatePieceTask
}

type ReplicatePieceTask interface {
	ObjectTask
	TransferSealObjectTask() SealObjectTask
}

type SealObjectTask interface {
	ObjectTask
}

type ReceivePieceTask interface {
	ObjectTask
	GetReplicateIdx() uint32
	SetReplicateIdx(uint32)
}

type DownloadObjectTask interface {
	ObjectTask
	GetNeedIntegrity() bool
	GetSize() uint64
	SetLow(uint64)
	GetLow() uint64
	SetHigh(uint64)
	GetHigh() uint64
}

type GCTask interface {
	Task
	IsRunning() bool
	StopRunning()
}

type GCObjectTask interface {
	GCTask
	SetStartBlockNumber(uint64)
	GetStartBlockNumber() uint64
	SetEndBlockNumber(uint64)
	GetEndBlockNumber() uint64
	GetGCObjectProcess() (uint64, uint64)
	SetGCObjectProcess(uint64, uint64)
}

type GCZombiePieceTask interface {
	GCTask
	GetGCZombiePieceStatus() uint64
	SetGCZombiePieceStatus(uint64)
}

type GCStoreTask interface {
	GCTask
	GetGCStoreStatus() (uint64, uint64)
	SetGCStoreStatus(uint64, uint64)
}

type TQueue interface {
	Top() Task
	Pop() Task
	Push(Task) error
	Has(TKey) bool
	Len() int
	Cap() int
}

type TQueueOnStrategy interface {
	TQueue
	Timeout(TKey) bool
	RetryFailed(TKey) bool
	PopWithKey(TKey) Task
	RunStrategy()
	SetStrategy(TQueueStrategy)
}

type TQueueWithLimit interface {
	TQueueOnStrategy
	PopWithLimit(rcmgr.Limit) Task
	GetSupportPickUpByLimit() bool
	SetSupportPickUpByLimit(bool)
}

type TQueueStrategy interface {
	SetPickUpFunc(func([]Task) Task)
	RunPickUp([]Task) Task
	SetCollectionFunc(func(TQueueOnStrategy, []TKey))
	RunCollection(TQueueOnStrategy, []TKey)
}

type TPriority uint8

const (
	UnSchedulingPriority      = 0
	UnKnownTaskPriority       = 0
	MaxTaskPriority           = 255
	DefaultLargerTaskPriority = 170
	DefaultSmallerPriority    = 85
)

type TPriorityLevel int32

const (
	TLowPriorityLevel TPriorityLevel = iota
	THighPriorityLevel
)

type TTaskPriority interface {
	GetPriority(TType) TPriority
	SetPriority(TType, TPriority)
	GetAllPriorities() map[TType]TPriority
	SetAllPriorities(map[TType]TPriority)
	GetLowLevelPriority() TPriority
	SetLowLevelPriority(TPriority) error
	GetHighLevelPriority() TPriority
	SetHighLevelPriority(TPriority) error
	HighLevelPriority(TPriority) bool
	LowLevelPriority(TPriority) bool
	SupportTask(TType) bool
}

type TPriorityQueueWithLimit interface {
	TQueueWithLimit
	SubQueueLen(TPriority) int
	GetPriority() []TPriority
	SetPriorityQueue(TPriority, TQueueWithLimit) error
}
