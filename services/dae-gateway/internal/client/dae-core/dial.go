package daecore

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type DialConfig struct {
	Addr             string
	Insecure         bool
	ServerName       string
	KeepAliveTime    time.Duration
	KeepAliveTimeout time.Duration
}

func Dial(ctx context.Context, cfg DialConfig, extra ...grpc.DialOption) (*grpc.ClientConn, error) {
	if cfg.KeepAliveTime <= 0 {
		cfg.KeepAliveTime = 30 * time.Second
	}
	if cfg.KeepAliveTimeout <= 0 {
		cfg.KeepAliveTimeout = 10 * time.Second
	}

	var creds credentials.TransportCredentials
	if !cfg.Insecure {
		creds = credentials.NewTLS(nil)
	} else {
		creds = insecure.NewCredentials()
	}

	opt := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                cfg.KeepAliveTime,
			Timeout:             cfg.KeepAliveTimeout,
			PermitWithoutStream: true,
		}),
	}

	opt = append(opt, extra...)

	return grpc.NewClient(cfg.Addr, opt...)
}
