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

// GRPCServer -
type GRPCServer struct {
	node     jarviscore.JarvisNode
	lis      net.Listener
	grpcServ *grpc.Server
	db       *tradingDataDB
}

// newGRPCServer -
func newGRPCServer(node jarviscore.JarvisNode, cfg *Config, db *tradingDataDB) (*GRPCServer, error) {
	lis, err := net.Listen("tcp", cfg.BindAddr)
	if err != nil {
		jarvisbase.Error("newGRPCServer:Listen",
			zap.Error(err))

		return nil, err
	}

	jarvisbase.Info("newGRPCServer:Listen",
		zap.String("addr", cfg.BindAddr))

	grpcServ := grpc.NewServer()

	s := &GRPCServer{
		node:     node,
		lis:      lis,
		grpcServ: grpcServ,
		db:       db,
	}

	tradingdatapb.RegisterTradingDataServiceServer(grpcServ, s)

	return s, nil
}

// Start -
func (s *GRPCServer) Start(ctx context.Context) error {
	return s.grpcServ.Serve(s.lis)
}

// Stop -
func (s *GRPCServer) Stop() {
	s.lis.Close()

	return
}

// SendTradeData -
func (s *GRPCServer) SendTradeData(ctx context.Context,
	chunk *tradingdatapb.TradeDataChunk) (*tradingdatapb.ReplySendTradeData, error) {

	lastlst := chunk.Trades
	for {
		retlst, lastlst, ts, err := getTradeDataWithDay(&lastlst)
		if err != nil {
			jarvisbase.Error("GRPCServer.SendTradeData:getTradeDataWithDay",
				zap.Error(err))

			return &tradingdatapb.ReplySendTradeData{
				Nums:    0,
				ErrInfo: err.Error(),
			}, nil
		}

		if retlst != nil {
			s.db.addTradeData(ctx, ts, &tradingdatapb.TradeDataChunk{
				Market: chunk.Market,
				Symbol: chunk.Symbol,
				Trades: retlst,
			})
		}

		if len(lastlst) == 0 {
			break
		}
	}

	return &tradingdatapb.ReplySendTradeData{
		Nums: int32(len(chunk.Trades)),
	}, nil
}
