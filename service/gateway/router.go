package gateway

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/bnb-chain/greenfield-storage-provider/model"
	"github.com/bnb-chain/greenfield-storage-provider/pkg/log"
)

const (
	approvalRouterName             = "GetApproval"
	putObjectRouterName            = "PutObject"
	getObjectRouterName            = "GetObject"
	challengeRouterName            = "Challenge"
	syncPieceRouterName            = "SyncPiece"
	getUserBucketsRouterName       = "GetUserBuckets"
	listObjectsByBucketRouterName  = "ListObjectsByBucketName"
	getBucketReadQuotaRouterName   = "GetBucketReadQuota"
	listBucketReadRecordRouterName = "ListBucketReadRecord"
)

const (
	createBucketApprovalAction = "CreateBucket"
	createObjectApprovalAction = "CreateObject"
)

// notFoundHandler log not found request info.
func (g *Gateway) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	log.Errorw("not found handler", "header", r.Header, "host", r.Host, "url", r.URL)
	if _, err := io.ReadAll(r.Body); err != nil {
		log.Errorw("failed to read the unknown request", "error", err)
	}
	if err := NoRouter.errorResponse(w, &requestContext{}); err != nil {
		log.Errorw("failed to response the unknown request", "error", err)
	}
}

// registerHandler is used to register mux handlers.
func (g *Gateway) registerHandler(r *mux.Router) {
	// bucket router, virtual-hosted style
	bucketRouter := r.Host("{bucket:.+}." + g.config.Domain).Subrouter()
	bucketRouter.NewRoute().
		Name(putObjectRouterName).
		Methods(http.MethodPut).
		Path("/{object:.+}").
		HandlerFunc(g.putObjectHandler)
	bucketRouter.NewRoute().
		Name(getObjectRouterName).
		Methods(http.MethodGet).
		Path("/{object:.+}").
		HandlerFunc(g.getObjectHandler)
	bucketRouter.NewRoute().
		Name(getBucketReadQuotaRouterName).
		Methods(http.MethodGet).
		Queries(model.GetBucketReadQuotaQuery, "",
			model.GetBucketReadQuotaMonthQuery, "{year_month}").
		HandlerFunc(g.getBucketReadQuotaHandler)
	bucketRouter.NewRoute().
		Name(listBucketReadRecordRouterName).
		Methods(http.MethodGet).
		Queries(model.ListBucketReadRecordQuery, "",
			model.ListBucketReadRecordMaxRecordsQuery, "{max_records}",
			model.StartTimestampUs, "{start_ts}",
			model.EndTimestampUs, "{end_ts}").
		HandlerFunc(g.listBucketReadRecordHandler)
	bucketRouter.NewRoute().
		Name(listObjectsByBucketRouterName).
		Methods(http.MethodGet).
		Path("/").
		HandlerFunc(g.listObjectsByBucketNameHandler)
	bucketRouter.NotFoundHandler = http.HandlerFunc(g.notFoundHandler)

	// bucket list router, virtual-hosted style
	bucketListRouter := r.Host(g.config.Domain).Subrouter()
	bucketListRouter.NewRoute().
		Name(getUserBucketsRouterName).
		Methods(http.MethodGet).
		Path("/").
		HandlerFunc(g.getUserBucketsHandler)
	// admin router, path style, new router will prefer use virtual-hosted style
	r.Path(model.GetApprovalPath).
		Name(approvalRouterName).
		Methods(http.MethodGet).
		Queries(model.ActionQuery, "{action}").
		HandlerFunc(g.getApprovalHandler)
	r.Path(model.ChallengePath).
		Name(challengeRouterName).
		Methods(http.MethodGet).
		HandlerFunc(g.challengeHandler)
	// sync piece to receiver
	r.Path(model.SyncPath).
		Name(syncPieceRouterName).
		Methods(http.MethodPut).
		HandlerFunc(g.syncPieceHandler)
	// TODO Barry Delete this path
	r.Name("getPaymentByBucketName").
		Methods(http.MethodGet).
		Path("/payment/buckets_name/{bucket:.+}/{is_full_list:.+}").
		HandlerFunc(g.getPaymentByBucketNameHandler)
	r.Name("getPaymentByBucketID").
		Methods(http.MethodGet).
		Path("/payment/buckets_id/{bucket:.+}/{is_full_list:.+}").
		HandlerFunc(g.getPaymentByBucketIDHandler)
	r.Name("getUserBucketsCountRouterName").
		Methods(http.MethodGet).
		Path("/accounts/{account_id:.+}/buckets/count").
		HandlerFunc(g.getUserBucketsCountHandler)
	r.Name("getBucketByBucketNameRouterName").
		Methods(http.MethodGet).
		Path("/buckets/name/{bucket_name:.+}/{is_full_list:.+}").
		HandlerFunc(g.getBucketByBucketNameHandler)
	//TODO barry remove all sp internal handlers
	r.Name("getBucketByBucketIDRouterName").
		Methods(http.MethodGet).
		Path("/buckets/id/{bucket_id:.+}/{is_full_list:.+}").
		HandlerFunc(g.getBucketByBucketIDHandler)
	r.Name("listDeletedObjectsByBlockNumberRangeRouterName").
		Methods(http.MethodGet).
		Path("/delete/{start:.+}/{end:.+}/{is_full_list:.+}").
		HandlerFunc(g.listDeletedObjectsByBlockNumberRangeHandler)
	r.NotFoundHandler = http.HandlerFunc(g.notFoundHandler)
}
