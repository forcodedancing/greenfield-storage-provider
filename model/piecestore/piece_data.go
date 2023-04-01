package piecestore

import (
	"math"

	"github.com/bnb-chain/greenfield-storage-provider/model"
	merrors "github.com/bnb-chain/greenfield-storage-provider/model/errors"
	storagetypes "github.com/bnb-chain/greenfield/x/storage/types"
)

// GetPieceDataSize return the size of piece data
func GetPieceDataSize(object *storagetypes.ObjectInfo, maxSegmentSize uint64, dataChunkNum uint32,
	dataType model.DateType, idx uint32) (uint64, error) {
	if object.GetPayloadSize() == 0 {
		return 0, merrors.ErrZeroDataSize
	}
	switch dataType {
	case model.SegmentType:
		segmentCnt := ComputeSegmentCount(object.GetPayloadSize(), maxSegmentSize)
		if idx > segmentCnt-1 {
			return 0, merrors.ErrOverflowDataSize
		}
		if idx == segmentCnt-1 {
			return object.GetPayloadSize() - maxSegmentSize*(uint64(idx)), nil
		}
		return maxSegmentSize, nil
	case model.PieceType:
		segmentCnt := ComputeSegmentCount(object.GetPayloadSize(), maxSegmentSize)
		if idx > segmentCnt-1 {
			return 0, merrors.ErrOverflowDataSize
		}
		segmentSize := maxSegmentSize
		if idx == segmentCnt-1 {
			// ignore the size of a small amount of data filled by ec encoding
			segmentSize = object.GetPayloadSize() - maxSegmentSize*(uint64(idx))
		}
		return uint64(math.Ceil(float64(segmentSize) / float64(dataChunkNum))), nil
	default:
		return 0, merrors.ErrUnsupportedDataType
	}
}
