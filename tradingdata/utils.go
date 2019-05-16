package tradingdata

import (
	"fmt"
	"time"

	tradingdatapb "github.com/zhs007/tradingdataserv/proto"
)

func makeTradeDataChunkKey(market string, symbol string, ts int64) string {
	tm := time.Unix(ts, 0)

	return fmt.Sprintf("%v:%v:%v", market, symbol, tm.Format("20060102150405"))
}

func insertTradeData(lst []*tradingdatapb.TradeInfo, index int,
	trade *tradingdatapb.TradeInfo) ([]*tradingdatapb.TradeInfo, error) {

	if index <= 0 {
		return append([]*tradingdatapb.TradeInfo{trade}, lst[0:]...), nil
	}

	if index >= len(lst) {
		return append(lst, trade), nil
	}

	lstf := append(lst[0:index], trade)

	return append(lstf, lst[index:]...), nil
}

func insert2TradeDataChunk(chunk *tradingdatapb.TradeDataChunk, trade *tradingdatapb.TradeInfo) error {
	for i, v := range chunk.Trades {
		if v.Curtime == trade.Curtime {
			if v.Id == trade.Id {
				return ErrInvalidTradeDataID
			}
		}

		if v.Curtime > trade.Curtime {
			lst, err := insertTradeData(chunk.Trades, i, trade)
			if err != nil {
				return err
			}

			chunk.Trades = lst

			return nil
		}
	}

	lst, err := insertTradeData(chunk.Trades, len(chunk.Trades), trade)
	if err != nil {
		return err
	}

	chunk.Trades = lst

	return nil
}
