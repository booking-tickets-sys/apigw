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

// UserServiceClient represents a client for the user service
type UserServiceClient struct {
	client pb.UserServiceClient
	conn   *grpc.ClientConn
}

// NewUserServiceClient creates a new user service client
func NewUserServiceClient(cfg *config.UserServiceConfig) (*UserServiceClient, error) {
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
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}

	client := pb.NewUserServiceClient(conn)

	return &UserServiceClient{
		client: client,
		conn:   conn,
	}, nil
}

// Close closes the gRPC connection
func (c *UserServiceClient) Close() error {
	return c.conn.Close()
}

// Register registers a new user
func (c *UserServiceClient) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	return c.client.Register(ctx, req)
}

// Login authenticates a user
func (c *UserServiceClient) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	return c.client.Login(ctx, req)
}

// RefreshToken refreshes an access token
func (c *UserServiceClient) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	return c.client.RefreshToken(ctx, req)
}
