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
	db       *tradingDataDB
}

// newTradingDataServ -
func newTradingDataServ(node jarviscore.JarvisNode, cfg *Config) (*tradingDataServ, error) {
	lis, err := net.Listen("tcp", cfg.BindAddr)
	if err != nil {
		jarvisbase.Error("newTradingDataServ:Listen",
			zap.Error(err))

		return nil, err
	}

	jarvisbase.Info("newTradingDataServ:Listen",
		zap.String("addr", cfg.BindAddr))

	grpcServ := grpc.NewServer()

	db, err := newTradingDataDB(cfg.AnkaDB.DBPath, cfg.AnkaDB.HTTPAddr, cfg.AnkaDB.Engine)
	if err != nil {
		jarvisbase.Error("newTradingDataServ:newTradingDataDB",
			zap.Error(err))

		return nil, err
	}

	s := &tradingDataServ{
		node:     node,
		lis:      lis,
		grpcServ: grpcServ,
		db:       db,
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

func (s *tradingDataServ) SendTradeData(ctx context.Context,
	chunk *tradingdatapb.TradeDataChunk) (*tradingdatapb.ReplySendTradeData, error) {

	lastlst := chunk.Trades
	for {
		retlst, lastlst, ts, err := getTradeDataWithDay(&lastlst)
		if err != nil {
			jarvisbase.Error("tradingDataServ.SendTradeData:getTradeDataWithDay",
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
