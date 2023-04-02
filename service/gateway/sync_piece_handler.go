package gateway

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"

	storagetypes "github.com/bnb-chain/greenfield/x/storage/types"

	"github.com/bnb-chain/greenfield-storage-provider/model"
	"github.com/bnb-chain/greenfield-storage-provider/pkg/log"
	p2ptypes "github.com/bnb-chain/greenfield-storage-provider/pkg/p2p/types"
	receivertypes "github.com/bnb-chain/greenfield-storage-provider/service/receiver/types"
	"github.com/bnb-chain/greenfield-storage-provider/util"
)

// syncPieceHandler handle sync piece data request
func (gateway *Gateway) replicatePieceHandler(w http.ResponseWriter, r *http.Request) {
	var (
		err               error
		errDescription    *errorDescription
		reqContext        *requestContext
		objectInfo        = storagetypes.ObjectInfo{}
		segmentIdx        uint32
		replicateIdx      uint32
		checksum          []byte
		replicateApproval = &p2ptypes.GetApprovalResponse{}
		size              int
		readN             int
		buf               = make([]byte, model.DefaultStreamBufSize)
	)

	reqContext = newRequestContext(r)
	defer func() {
		if errDescription != nil {
			_ = errDescription.errorResponse(w, reqContext)
		}
		if errDescription != nil && errDescription.statusCode == http.StatusOK {
			log.Errorf("action(%v) statusCode(%v) %v", replicatePieceRouterName, errDescription.statusCode, reqContext.generateRequestDetail())
		} else {
			log.Infof("action(%v) statusCode(200) %v", replicatePieceRouterName, reqContext.generateRequestDetail())
		}
	}()

	if gateway.receiver == nil {
		log.Error("failed to sync data due to not config receiver")
		errDescription = NotExistComponentError
		return
	}

	objectInfoMsg, err := hex.DecodeString(r.Header.Get(model.GnfdObjectInfoHeader))
	if err != nil {
		log.Errorw("failed to parse object info header", "object_info", r.Header.Get(model.GnfdObjectInfoHeader))
		errDescription = InvalidHeader
		return
	}
	if storagetypes.ModuleCdc.UnmarshalJSON(objectInfoMsg, &objectInfo) != nil {
		log.Errorw("failed to unmarshal object info header", "object_info", r.Header.Get(model.GnfdObjectInfoHeader))
		errDescription = InvalidHeader
		return
	}
	if segmentIdx, err = util.StringToUint32(r.Header.Get(model.GnfdSegmentIdxHeader)); err != nil {
		log.Errorw("failed to parse segment_size header", "segment_size", r.Header.Get(model.GnfdSegmentIdxHeader))
		errDescription = InvalidHeader
		return
	}
	if replicateIdx, err = util.StringToUint32(r.Header.Get(model.GnfdReplicaIdxHeader)); err != nil {
		log.Errorw("failed to parse replica_idx header", "replica_idx", r.Header.Get(model.GnfdReplicaIdxHeader))
		errDescription = InvalidHeader
		return
	}
	checksum = []byte(r.Header.Get(model.GnfdReplicateDataChecksum))
	if err = json.Unmarshal([]byte(r.Header.Get(model.GnfdReplicateApproval)), replicateApproval); err != nil {
		log.Errorw("failed to parse replicate_approval header", "replicate_approval", r.Header.Get(model.GnfdReplicateApproval))
		errDescription = InvalidHeader
		return
	}
	if err = gateway.verifyReplicateApproval(replicateApproval); err != nil {
		log.Errorw("failed to verify replicate_approval header", "replicate_approval", r.Header.Get(model.GnfdReplicateApproval))
		errDescription = InvalidHeader
		return
	}
	// TODO:: add S3 signature, the field be signed include: objectInfo, segmentIdx, replicateIdx, checksum

	stream, err := gateway.receiver.ReplicatePiece(context.Background())
	if err != nil {
		log.Errorw("failed to replicate piece", "error", err)
		errDescription = InternalError
		return
	}
	for {
		readN, err = r.Body.Read(buf)
		if err != nil && err != io.EOF {
			log.Errorw("failed to sync piece due to reader error", "error", err)
			errDescription = InternalError
			return
		}
		if readN > 0 {
			if err = stream.Send(&receivertypes.ReceivePieceRequest{
				ObjectInfo:   &objectInfo,
				SegmentIdx:   segmentIdx,
				ReplicateIdx: replicateIdx,
				Checksum:     checksum,
				PieceData:    buf[:readN],
			}); err != nil {
				log.Errorw("failed to send stream", "error", err)
				errDescription = InternalError
				return
			}
			size += readN
		}
		if err == io.EOF {
			if size == 0 {
				log.Errorw("failed to replicate piece due to payload is empty")
				errDescription = InvalidPayload
				return
			}
			_, err = stream.CloseAndRecv()
			if err != nil {
				log.Errorw("failed to replicate piece due to stream close", "error", err)
				errDescription = InternalError
				return
			}
			// succeed to replicate piece
			break
		}
	}
}

func (gateway *Gateway) getIntegrityHashHandler(w http.ResponseWriter, r *http.Request) {
	var (
		errDescription *errorDescription
		reqContext     *requestContext
		objectInfo     = storagetypes.ObjectInfo{}
	)

	reqContext = newRequestContext(r)
	defer func() {
		if errDescription != nil {
			_ = errDescription.errorResponse(w, reqContext)
		}
		if errDescription != nil && errDescription.statusCode == http.StatusOK {
			log.Errorf("action(%v) statusCode(%v) %v", replicatePieceRouterName, errDescription.statusCode, reqContext.generateRequestDetail())
		} else {
			log.Infof("action(%v) statusCode(200) %v", replicatePieceRouterName, reqContext.generateRequestDetail())
		}
	}()

	if gateway.receiver == nil {
		log.Error("failed to get integrity hash due to not config receiver")
		errDescription = NotExistComponentError
		return
	}

	objectInfoMsg, err := hex.DecodeString(r.Header.Get(model.GnfdObjectInfoHeader))
	if err != nil {
		log.Errorw("failed to parse object info header", "object_info", r.Header.Get(model.GnfdObjectInfoHeader))
		errDescription = InvalidHeader
		return
	}
	if storagetypes.ModuleCdc.UnmarshalJSON(objectInfoMsg, &objectInfo) != nil {
		log.Errorw("failed to unmarshal object info header", "object_info", r.Header.Get(model.GnfdObjectInfoHeader))
		errDescription = InvalidHeader
		return
	}
	//TODO:: verify objectInfo by S3 signature

	integrityHash, signature, err := gateway.receiver.DoneReplicate(context.Background(), &objectInfo)
	if err != nil {
		errDescription = makeErrorDescription(err)
		return
	}
	w.Header().Set(model.GnfdIntegrityHashHeader, hex.EncodeToString(integrityHash))
	w.Header().Set(model.GnfdIntegrityHashSignatureHeader, hex.EncodeToString(signature))
}
