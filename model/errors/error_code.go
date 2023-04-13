package errors

const (
	// RPCErrCode defines storage provider rpc error code
	RPCErrCode = 10000
	// ErrorCodeBadRequest defines bad request error code
	ErrorCodeBadRequest = 40001
	// ErrorCodeNotFound defines not found error code
	ErrorCodeNotFound = 40004
	// ErrorCodeInternalError defines internal error code
	ErrorCodeInternalError = 50001
)

const (
	// deinfe common error code, from 10000 to 14999
	MismatchIntegrityHashErrCode    = 11000
	CacheMissedErrCode              = 11001
	DanglingPointerErrCode          = 11002
	PayloadStreamErrCode            = 11003
	ResourceMgrBeginSpanErrCode     = 11004
	ComputePieceSizeErrCode         = 11005
	ResourceMgrReserveMemoryErrCode = 11006
	ObjectNotFoundErrCode           = 11007
	HexDecodeStringErrCode          = 11008
	StringToByteSliceErrCode        = 11009
	ParseStringToIntErrCode         = 11010
	// sp database error, from 15000 to 19999
	DBRecordNotFound                    = 15000
	DBUnknownAddressTypeErrCode         = 15001
	DBGetStorageParamsErrCode           = 15100
	DBGetBucketTrafficErrCode           = 15102
	DBGetBucketReadRecordErrCode        = 15103
	DBQuotaNotEnoughErrCode             = 15104
	DBQueryInBucketTrafficTableErrCode  = 15106
	DBInsertInBucketTrafficTableErrCode = 15107
	DBUpdateInBucketTrafficTableErrCode = 15108
	DBInsertInReadRecordTableErrCode    = 15109
	DBQueryInSPInfoTableErrCode         = 20000
	DBDeleteInSPInfoTableErrCode        = 20001
	DBInsertInSPInfoTableErrCode        = 20002
	DBQueryOwnSPInSPInfoTableErrCode    = 20003
	DBInsertOwnSPInSPInfoTableErrCode   = 20004
	DBUpdateOwnSPInSPInfoTableErrCode   = 20005
	DBQueryInStorageParamsTableErrCode  = 20006
	DBInsertInStorageParamsTableErrCode = 20007
	DBUpdateInStorageParamsTableErrCode = 20008
	DBQueryInJobTableErrCode            = 20009
	DBInsertInJobTableErrCode           = 20010
	DBUpdateInJobTableErrCode           = 20011
	DBQueryInObjectTableErrCode         = 20012
	DBInsertInObjectTableErrCode        = 20013
	DBUpdateInObjectTableErrCode        = 20014
	DBQueryInIntegrityMetaTableErrCode  = 20016
	DBInsertInIntegrityMetaTableErrCode = 20017

	// uploader service error, from 20000 to 20999
	UploaderMismatchChecksumNumErrCode = 20100

	// challenge service error, from 21000 to 21999

	// downloader service error, from 22000 to 22999
	DownloaderInvalidPieceInfoParamsErrCode = 22100

	// signer service error, from 23000 to 23999
	SignerSignIntegrityHashErrCode = 23100

	// p2p service, from 24000 to 24999

	// receiver service error code, from 25000 to 25999

	// task node service error, from 26000 to 26999

	// piece store error code, from 27000 to 27999
	PieceStorePutObjectError = 27100
	PieceStoreGetObjectError = 27101
)
