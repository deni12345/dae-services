package daecore

import (
	"context"
	"time"

	pb "github.com/deni12345/dae-services/proto/gen"
	"google.golang.org/grpc"
)

type Client struct {
	Health pb.HealthServiceClient
	User   pb.UsersServiceClient
	Sheet  pb.SheetsServiceClient
	Order  pb.OrdersServiceClient

	defaultTimeOut time.Duration
	conn           *grpc.ClientConn
}

type Config struct {
	Addr           string
	Insecure       bool
	ConnectTimeout time.Duration
	DefaultTimeOut time.Duration
}

func NewClient(ctx context.Context, cfg Config) (*Client, error) {
	conn, err := Dial(ctx, DialConfig{
		Addr:     cfg.Addr,
		Insecure: cfg.Insecure,
	})
	if err != nil {
		return nil, err
	}

	return New(conn, cfg.DefaultTimeOut)
}

func New(conn *grpc.ClientConn, defaultTimeout time.Duration) (*Client, error) {
	if defaultTimeout <= 0 {
		defaultTimeout = 10 * time.Second
	}

	return &Client{
		Health: pb.NewHealthServiceClient(conn),
		User:   pb.NewUsersServiceClient(conn),
		Sheet:  pb.NewSheetsServiceClient(conn),
		Order:  pb.NewOrdersServiceClient(conn),

		defaultTimeOut: defaultTimeout,
		conn:           conn,
	}, nil
}

func (c *Client) Close() error {
	if c != nil {
		return c.conn.Close()
	}
	return nil
}

func withTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if _, has := ctx.Deadline(); has {
		return ctx, func() {}
	}
	return context.WithTimeout(ctx, timeout)
}
