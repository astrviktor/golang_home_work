package internalgrpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
)

func (s *Server) WithLoggingInterceptor(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	now := "[" + time.Now().Format(s.logger.GetTimeFormat()) + "]"
	start := time.Now()

	h, err := handler(ctx, req)

	s.logger.Info(fmt.Sprint(now, info.FullMethod, "GRPC", time.Since(start)))

	return h, err
}
