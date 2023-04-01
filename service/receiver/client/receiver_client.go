package client

import (
	"context"

	storagetypes "github.com/bnb-chain/greenfield/x/storage/types"
	"google.golang.org/grpc"

	"github.com/bnb-chain/greenfield-storage-provider/pkg/log"
	"github.com/bnb-chain/greenfield-storage-provider/pkg/metrics"
	mwgrpc "github.com/bnb-chain/greenfield-storage-provider/pkg/middleware/grpc"
	"github.com/bnb-chain/greenfield-storage-provider/service/receiver/types"
	utilgrpc "github.com/bnb-chain/greenfield-storage-provider/util/grpc"
)

// ReceiverClient is a receiver gRPC service client wrapper
type ReceiverClient struct {
	address  string
	conn     *grpc.ClientConn
	receiver types.ReceiverServiceClient
}

// NewReceiverClient return a ReceiverClient instance
func NewReceiverClient(address string) (*ReceiverClient, error) {
	options := utilgrpc.GetDefaultClientOptions()
	if metrics.GetMetrics().Enabled() {
		options = append(options, mwgrpc.GetDefaultClientInterceptor()...)
	}
	conn, err := grpc.DialContext(context.Background(), address, options...)
	if err != nil {
		log.Errorw("failed to dial receiver", "error", err)
		return nil, err
	}
	client := &ReceiverClient{
		address:  address,
		conn:     conn,
		receiver: types.NewReceiverServiceClient(conn),
	}
	return client, nil
}

// Close the receiver gPRC connection
func (client *ReceiverClient) Close() error {
	return client.conn.Close()
}

// ReplicatePiece replicates a piece data with object info
func (client *ReceiverClient) ReplicatePiece(
	ctx context.Context,
	opts ...grpc.CallOption) (types.ReceiverService_ReceivePieceClient, error) {
	return client.receiver.ReceivePiece(ctx, opts...)
}

// DoneReplicate completes a replication of an object payload data
func (client *ReceiverClient) DoneReplicate(ctx context.Context, object *storagetypes.ObjectInfo, opts ...grpc.CallOption) (
	[]byte, []byte, error) {
	rea := &types.DoneReplicateRequest{ObjectInfo: object}
	resp, err := client.receiver.DoneReplicate(ctx, rea, opts...)
	if err != nil {
		return []byte{}, []byte{}, err
	}
	return resp.GetIntegrityHash(), resp.GetSignature(), nil
}
