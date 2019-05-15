package tradingdata

import (
	"github.com/zhs007/ankadb"
	jarvisbase "github.com/zhs007/jarviscore/base"
	"go.uber.org/zap"
)

// tradingDataDB - TradingData database
type tradingDataDB struct {
	ankaDB ankadb.AnkaDB
}

// newDTDataDB - new dtdata db
func newDTDataDB(dbpath string, httpAddr string, engine string) (*tradingDataDB, error) {
	cfg := ankadb.NewConfig()

	cfg.AddrHTTP = httpAddr
	cfg.PathDBRoot = dbpath
	cfg.ListDB = append(cfg.ListDB, ankadb.DBConfig{
		Name:   TradingDataDBName,
		Engine: engine,
		PathDB: TradingDataDBName,
	})

	ankaDB, err := ankadb.NewAnkaDB(cfg, nil)
	if ankaDB == nil {
		jarvisbase.Error("newDTDataDB", zap.Error(err))

		return nil, err
	}

	jarvisbase.Info("newDTDataDB", zap.String("dbpath", dbpath),
		zap.String("httpAddr", httpAddr), zap.String("engine", engine))

	db := &tradingDataDB{
		ankaDB: ankaDB,
	}

	return db, nil
}
