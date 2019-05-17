package tradingdata

import (
	"context"

	"go.uber.org/zap"

	"github.com/zhs007/jarviscore"
	jarvisbase "github.com/zhs007/jarviscore/base"
)

// Serv - trading data server
type Serv struct {
	grpcServ *GRPCServer
	httpServ *HTTPServer
	db       *tradingDataDB
}

// NewServ -
func NewServ(node jarviscore.JarvisNode, cfg *Config) (*Serv, error) {

	db, err := newTradingDataDB(cfg.AnkaDB.DBPath, cfg.AnkaDB.HTTPAddr, cfg.AnkaDB.Engine)
	if err != nil {
		jarvisbase.Error("NewServ:newTradingDataDB",
			zap.Error(err))

		return nil, err
	}

	grpcserv, err := newGRPCServer(node, cfg, db)
	if err != nil {
		jarvisbase.Error("NewServ:newGRPCServer",
			zap.Error(err))

		return nil, err
	}

	httpserv, err := newHTTPServer(cfg.HTTPAddr, db)
	if err != nil {
		jarvisbase.Error("NewServ:newGRPCServer",
			zap.Error(err))

		return nil, err
	}

	s := &Serv{
		grpcServ: grpcserv,
		httpServ: httpserv,
		db:       db,
	}

	return s, nil
}

// Start -
func (s *Serv) Start(ctx context.Context) {
	go s.grpcServ.Start(ctx)
	go s.httpServ.start(ctx)
}

// Stop -
func (s *Serv) Stop() {
	s.grpcServ.Stop()
	s.httpServ.stop()
}
