package tradingdata

import (
	"fmt"
	"time"

	tradingdatapb "github.com/zhs007/tradingdataserv/proto"
)

func makeTradeDataChunkKey(market string, symbol string, ts int64) string {
	tm := time.Unix(ts/1000, 0)

	return fmt.Sprintf("%v:%v:%v", market, symbol, tm.Format("20060102"))
}

func insertTradeData(lst []*tradingdatapb.TradeInfo, index int,
	trade *tradingdatapb.TradeInfo) ([]*tradingdatapb.TradeInfo, error) {

	if index <= 0 {
		return append([]*tradingdatapb.TradeInfo{trade}, lst[0:]...), nil
	}

	if index >= len(lst) {
		return append(lst, trade), nil
	}

	tmp := append([]*tradingdatapb.TradeInfo{}, lst[index:]...)
	lstf := append(lst[0:index], trade)

	// return lstf, nil
	return append(lstf, tmp[0:]...), nil
}

func insert2TradeDataChunk(chunk *tradingdatapb.TradeDataChunk,
	trade *tradingdatapb.TradeInfo) error {

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

func getTradeDataWithDay(lst []*tradingdatapb.TradeInfo) ([]*tradingdatapb.TradeInfo,
	[]*tradingdatapb.TradeInfo, int64, error) {

	if len(lst) <= 0 {
		return nil, nil, 0, nil
	}

	var lastlst []*tradingdatapb.TradeInfo
	var retlst []*tradingdatapb.TradeInfo

	ts := lst[0].Curtime
	tm := time.Unix(ts/1000, 0)
	td := tm.Format("20060102")

	for _, v := range lst {
		cts := v.Curtime
		ctm := time.Unix(cts/1000, 0)
		ctd := ctm.Format("20060102")

		if ctd == td {
			retlst = append(retlst, v)
		} else {
			lastlst = append(lastlst, v)
		}
	}

	return retlst, lastlst, ts, nil
}
