package service

import (
	"context"
	"crypto/tls"
	"fmt"
	"gitlab.top.slotssprite.com/my/api-layout/party/grpcx"
	"log"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	kgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"gitlab.top.slotssprite.com/my/api-layout/internal/conf"
	"gitlab.top.slotssprite.com/my/api-layout/party/util/xcolor"
)

var (
	EndpointNilError = errors.New("endpoint is nil")

	// _rpcLayoutConn ...
	_rpcLayoutConn      *grpc.ClientConn
	_rpcLayoutConnMutex = new(sync.Mutex)
)

// GetRpcLayoutConn ...
func GetRpcLayoutConn() *grpc.ClientConn {
	_rpcLayoutConnMutex.Lock()
	defer _rpcLayoutConnMutex.Unlock()
	if _rpcLayoutConn != nil {
		return _rpcLayoutConn
	}

	endpoint := viper.GetString(conf.PathRpcLayoutEndpoint)
	callTimeout := viper.GetInt(conf.PathRpcLayoutTimeout)
	dialWithCredentials := viper.GetBool(conf.PathRpcLayoutDialWithCredentials)
	_rpcLayoutConn = dial(endpoint, callTimeout, dialWithCredentials)
	return _rpcLayoutConn
}

// dial ...
func dial(endpoint string, callTimeout int, dialWithCredentials bool) (conn *grpc.ClientConn) {
	if endpoint == "" {
		panic(EndpointNilError)
	}

	opts := []kgrpc.ClientOption{
		kgrpc.WithEndpoint(endpoint),
		kgrpc.WithMiddleware(recovery.Recovery(), tracing.Client(), grpcx.ClientBreaker(), metadata.Client()),
		// kgrpc.WithOptions(grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32), grpc.MaxCallSendMsgSize(math.MaxInt32))),
	}
	if callTimeout >= 0 {
		opts = append(opts, kgrpc.WithTimeout(time.Duration(callTimeout)*time.Second))
	}

	var err error
	if dialWithCredentials {
		opts = append(opts, kgrpc.WithTLSConfig(&tls.Config{}))
		conn, err = kgrpc.Dial(context.Background(), opts...)
	} else {
		conn, err = kgrpc.DialInsecure(context.Background(), opts...)
	}
	if err != nil {
		panic(err)
	}
	log.Printf(xcolor.YELLOW, fmt.Sprintf("Connecting at %s", endpoint))
	return
}
