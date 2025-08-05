package client

import (
	"context"
	"fmt"

	pb "apigw/client/proto"
	"apigw/internal/app/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// TicketServiceClient represents a client for the ticket service
type OrderServiceClient struct {
	client pb.OrderServiceClient
	conn   *grpc.ClientConn
}

// NewOrderServiceClient creates a new order service client
func NewOrderServiceClient(cfg *config.OrderServiceConfig) (*OrderServiceClient, error) {
	address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	conn, err := grpc.NewClient(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                cfg.GRPC.KeepaliveTime,
			Timeout:             cfg.GRPC.KeepaliveTimeout,
			PermitWithoutStream: cfg.GRPC.KeepalivePermitWithoutStream,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ticket service: %w", err)
	}

	client := pb.NewOrderServiceClient(conn)

	return &OrderServiceClient{
		client: client,
		conn:   conn,
	}, nil
}

// Close closes the gRPC connection
func (c *OrderServiceClient) Close() error {
	return c.conn.Close()
}

// PurchaseTicket purchases a ticket for the specified event and user
func (c *OrderServiceClient) PurchaseTicket(ctx context.Context, req *pb.PurchaseRequest) (*pb.PurchaseResponse, error) {
	return c.client.PurchaseTicket(ctx, req)
}
