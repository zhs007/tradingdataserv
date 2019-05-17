package tradingdata

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/zhs007/ankadb"
	jarvisbase "github.com/zhs007/jarviscore/base"
	tradingdatapb "github.com/zhs007/tradingdataserv/proto"
	"go.uber.org/zap"
)

// tradingDataDB - TradingData database
type tradingDataDB struct {
	ankaDB ankadb.AnkaDB
}

// newTradingDataDB - new tradingdata db
func newTradingDataDB(dbpath string, httpAddr string, engine string) (*tradingDataDB, error) {
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
		jarvisbase.Error("newTradingDataDB", zap.Error(err))

		return nil, err
	}

	jarvisbase.Info("newTradingDataDB", zap.String("dbpath", dbpath),
		zap.String("httpAddr", httpAddr), zap.String("engine", engine))

	db := &tradingDataDB{
		ankaDB: ankaDB,
	}

	return db, nil
}

func (db *tradingDataDB) addTradeData(ctx context.Context, ts int64,
	chunk *tradingdatapb.TradeDataChunk) (int, error) {

	if len(chunk.Trades) <= 0 {
		return len(chunk.Trades), nil
	}

	key := makeTradeDataChunkKey(chunk.Market, chunk.Symbol, ts)

	srcbuf, err := db.ankaDB.Get(ctx, TradingDataDBName, key)
	if err != nil {
		jarvisbase.Warn("tradingDataDB.addTradeData:Get",
			zap.Error(err))

		return 0, err
	}

	srctrunk := &tradingdatapb.TradeDataChunk{}

	err = proto.Unmarshal(srcbuf, srctrunk)
	if err != nil {
		jarvisbase.Warn("tradingDataDB.addTradeData:Unmarshal",
			zap.Error(err))

		return 0, err
	}

	nums := 0
	for _, v := range chunk.Trades {
		err := insert2TradeDataChunk(srctrunk, v)
		if err == nil {
			nums = nums + 1
			// jarvisbase.Warn("tradingDataDB.addTradeData:insert2TradeDataChunk",
			// 	zap.Error(err))

			// return 0, err
		}
	}

	buf, err := proto.Marshal(srctrunk)
	if err != nil {
		jarvisbase.Warn("tradingDataDB.addTradeData:Marshal",
			zap.Error(err))

		return 0, err
	}

	db.ankaDB.Set(ctx, TradingDataDBName, key, buf)

	return nums, nil
}

// getTradeData -
func (db *tradingDataDB) getTradeData(ctx context.Context, market string, symbol string,
	ts int64) (*tradingdatapb.TradeDataChunk, error) {

	key := makeTradeDataChunkKey(market, symbol, ts)

	buf, err := db.ankaDB.Get(ctx, TradingDataDBName, key)
	if err != nil {
		if err == ankadb.ErrNotFoundKey {
			return nil, nil
		}

		jarvisbase.Warn("tradingDataDB.getTradeData:Get", zap.Error(err))

		return nil, err
	}

	data := &tradingdatapb.TradeDataChunk{}

	err = proto.Unmarshal(buf, data)
	if err != nil {
		jarvisbase.Warn("tradingDataDB.getTradeData:Unmarshal", zap.Error(err))

		return nil, err
	}

	return data, err
}

// clearTradeDataID -
func (db *tradingDataDB) clearTradeDataID(ctx context.Context, market string, symbol string,
	ts int64) error {

	key := makeTradeDataChunkKey(market, symbol, ts)

	buf, err := db.ankaDB.Get(ctx, TradingDataDBName, key)
	if err != nil {
		if err == ankadb.ErrNotFoundKey {
			return nil
		}

		jarvisbase.Warn("tradingDataDB.clearTradeDataID:Get", zap.Error(err))

		return err
	}

	data := &tradingdatapb.TradeDataChunk{}

	err = proto.Unmarshal(buf, data)
	if err != nil {
		jarvisbase.Warn("tradingDataDB.clearTradeDataID:Unmarshal", zap.Error(err))

		return err
	}

	for _, v := range data.Trades {
		v.Id = ""
	}

	nbuf, err := proto.Marshal(data)
	if err != nil {
		jarvisbase.Warn("tradingDataDB.clearTradeDataID:Marshal",
			zap.Error(err))

		return err
	}

	db.ankaDB.Set(ctx, TradingDataDBName, key, nbuf)

	return err
}
