package tradingdata

import (
	"context"
	"net"

	"go.uber.org/zap"

	"github.com/zhs007/jarviscore"
	jarvisbase "github.com/zhs007/jarviscore/base"
	tradingdatapb "github.com/zhs007/tradingdataserv/proto"
	"google.golang.org/grpc"
)

// tradingDataServ
type tradingDataServ struct {
	node     jarviscore.JarvisNode
	lis      net.Listener
	grpcServ *grpc.Server
}

// newTradingDataServ -
func newTradingDataServ(node jarviscore.JarvisNode, cfg *Config) (*tradingDataServ, error) {
	lis, err := net.Listen("tcp", cfg.BindAddr)
	if err != nil {
		jarvisbase.Error("newTradingDataServ", zap.Error(err))

		return nil, err
	}

	jarvisbase.Info("newTradingDataServ:Listen",
		zap.String("addr", cfg.BindAddr))

	grpcServ := grpc.NewServer()

	s := &tradingDataServ{
		node:     node,
		lis:      lis,
		grpcServ: grpcServ,
	}

	tradingdatapb.RegisterTradingDataServiceServer(grpcServ, s)

	return s, nil
}

// Start -
func (s *tradingDataServ) Start(ctx context.Context) error {
	return s.grpcServ.Serve(s.lis)
}

// Stop -
func (s *tradingDataServ) Stop() {
	s.lis.Close()

	return
}

func (s *tradingDataServ) SendTradeData(context.Context,
	*tradingdatapb.TradeDataChunk) (*tradingdatapb.ReplySendTradeData, error) {

	return &tradingdatapb.ReplySendTradeData{}, nil
}
