package tradingdata

import (
	"testing"
	"time"

	tradingdatapb "github.com/zhs007/tradingdataserv/proto"
)

func TestMakeTradeDataChunkKey(t *testing.T) {
	str := makeTradeDataChunkKey("market", "usdbtc", 1558354301123)
	if str != "market:usdbtc:20190520" {
		t.Fatalf("TestMakeTradeDataChunkKey %v", str)

		return
	}

	t.Logf("TestMakeTradeDataChunkKey is OK")
}

func TestInsertTradeData(t *testing.T) {
	var lst []*tradingdatapb.TradeInfo
	t1 := &tradingdatapb.TradeInfo{
		Id:      "1",
		Curtime: 1,
	}

	lst1, err := insertTradeData(lst, -1, t1)
	if err != nil || len(lst1) != 1 {
		t.Fatalf("TestInsertTradeData 1 %v", lst1)

		return
	}

	t2 := &tradingdatapb.TradeInfo{
		Id:      "2",
		Curtime: 2,
	}

	lst21, err := insertTradeData(lst1, 0, t2)
	if err != nil || len(lst21) != 2 || lst21[0].Curtime != 2 ||
		lst21[1].Curtime != 1 {

		t.Fatalf("TestInsertTradeData 21 %v", lst21)

		return
	}

	t3 := &tradingdatapb.TradeInfo{
		Id:      "3",
		Curtime: 3,
	}

	lst213, err := insertTradeData(lst21, len(lst21), t3)
	if err != nil || len(lst213) != 3 || lst213[0].Curtime != 2 ||
		lst21[1].Curtime != 1 || lst213[2].Curtime != 3 {

		t.Fatalf("TestInsertTradeData 213 %v", lst213)

		return
	}

	t4 := &tradingdatapb.TradeInfo{
		Id:      "4",
		Curtime: 4,
	}

	lst2143, err := insertTradeData(lst213, len(lst213)-1, t4)
	if err != nil || len(lst2143) != 4 || lst2143[0].Curtime != 2 ||
		lst2143[1].Curtime != 1 || lst2143[2].Curtime != 4 || lst2143[3].Curtime != 3 {

		t.Fatalf("TestInsertTradeData 2143 %v", lst2143)

		return
	}

	t5 := &tradingdatapb.TradeInfo{
		Id:      "5",
		Curtime: 5,
	}

	lst25143, err := insertTradeData(lst2143, 1, t5)
	if err != nil || len(lst25143) != 5 || lst25143[0].Curtime != 2 ||
		lst25143[1].Curtime != 5 || lst25143[2].Curtime != 1 ||
		lst25143[3].Curtime != 4 || lst25143[4].Curtime != 3 {

		t.Fatalf("TestInsertTradeData 25143 %v", lst25143)

		return
	}

	t6 := &tradingdatapb.TradeInfo{
		Id:      "6",
		Curtime: 6,
	}

	lst251436, err := insertTradeData(lst25143, 100, t6)
	if err != nil || len(lst251436) != 6 || lst251436[0].Curtime != 2 ||
		lst251436[1].Curtime != 5 || lst251436[2].Curtime != 1 ||
		lst251436[3].Curtime != 4 || lst251436[4].Curtime != 3 || lst251436[5].Curtime != 6 {

		t.Fatalf("TestInsertTradeData 251436 %v", lst251436)

		return
	}

	t.Logf("TestInsertTradeData is OK")
}

func TestInsert2TradeDataChunk(t *testing.T) {
	chunk := &tradingdatapb.TradeDataChunk{}
	t2 := &tradingdatapb.TradeInfo{
		Id:      "2",
		Curtime: 2,
	}

	err := insert2TradeDataChunk(chunk, t2)
	if err != nil || len(chunk.Trades) != 1 {
		t.Fatalf("TestInsert2TradeDataChunk 2 %v", chunk)

		return
	}

	t5 := &tradingdatapb.TradeInfo{
		Id:      "5",
		Curtime: 5,
	}

	err = insert2TradeDataChunk(chunk, t5)
	if err != nil || len(chunk.Trades) != 2 || chunk.Trades[0].Curtime != 2 ||
		chunk.Trades[1].Curtime != 5 {

		t.Fatalf("TestInsert2TradeDataChunk 25 %v", chunk)

		return
	}

	t3 := &tradingdatapb.TradeInfo{
		Id:      "3",
		Curtime: 3,
	}

	err = insert2TradeDataChunk(chunk, t3)
	if err != nil || len(chunk.Trades) != 3 || chunk.Trades[0].Curtime != 2 ||
		chunk.Trades[1].Curtime != 3 || chunk.Trades[2].Curtime != 5 {

		t.Fatalf("TestInsert2TradeDataChunk 235 %v", chunk)

		return
	}

	t1 := &tradingdatapb.TradeInfo{
		Id:      "1",
		Curtime: 1,
	}

	err = insert2TradeDataChunk(chunk, t1)
	if err != nil || len(chunk.Trades) != 4 || chunk.Trades[0].Curtime != 1 ||
		chunk.Trades[1].Curtime != 2 || chunk.Trades[2].Curtime != 3 ||
		chunk.Trades[3].Curtime != 5 {

		t.Fatalf("TestInsert2TradeDataChunk 1235 %v", chunk)

		return
	}

	t6 := &tradingdatapb.TradeInfo{
		Id:      "6",
		Curtime: 6,
	}

	err = insert2TradeDataChunk(chunk, t6)
	if err != nil || len(chunk.Trades) != 5 || chunk.Trades[0].Curtime != 1 ||
		chunk.Trades[1].Curtime != 2 || chunk.Trades[2].Curtime != 3 ||
		chunk.Trades[3].Curtime != 5 || chunk.Trades[4].Curtime != 6 {

		t.Fatalf("TestInsert2TradeDataChunk 12356 %v", chunk)

		return
	}

	t4 := &tradingdatapb.TradeInfo{
		Id:      "4",
		Curtime: 4,
	}

	err = insert2TradeDataChunk(chunk, t4)
	if err != nil || len(chunk.Trades) != 6 || chunk.Trades[0].Curtime != 1 ||
		chunk.Trades[1].Curtime != 2 || chunk.Trades[2].Curtime != 3 ||
		chunk.Trades[3].Curtime != 4 || chunk.Trades[4].Curtime != 5 ||
		chunk.Trades[5].Curtime != 6 {

		t.Fatalf("TestInsert2TradeDataChunk 123456 %v", chunk)

		return
	}

	err = insert2TradeDataChunk(chunk, t3)
	if err != ErrInvalidTradeDataID {
		t.Fatalf("TestInsert2TradeDataChunk 1233456 %v", chunk)

		return
	}

	t7 := &tradingdatapb.TradeInfo{
		Id:      "7",
		Curtime: 3,
	}

	err = insert2TradeDataChunk(chunk, t7)
	if err != nil || len(chunk.Trades) != 7 || chunk.Trades[0].Curtime != 1 ||
		chunk.Trades[1].Curtime != 2 || chunk.Trades[2].Curtime != 3 ||
		chunk.Trades[3].Curtime != 3 ||
		chunk.Trades[4].Curtime != 4 || chunk.Trades[5].Curtime != 5 ||
		chunk.Trades[6].Curtime != 6 {

		t.Fatalf("TestInsert2TradeDataChunk 1233456 %v", chunk)

		return
	}

	t.Logf("TestInsert2TradeDataChunk is OK")
}

func TestGetTradeDataWithDay(t *testing.T) {
	lst0 := []*tradingdatapb.TradeInfo{
		&tradingdatapb.TradeInfo{
			Curtime: 1558354301123,
			Id:      "1",
		},
		&tradingdatapb.TradeInfo{
			Curtime: 1558444102000,
			Id:      "2",
		},
		&tradingdatapb.TradeInfo{
			Curtime: 1558444333000,
			Id:      "3",
		},
	}

	lsttmstr := []string{
		"20190520",
		"20190521",
	}

	lstlen := []int{
		1,
		2,
	}

	i := 0
	for {
		retlst, lastlst, ts, err := getTradeDataWithDay(lst0)
		if err != nil {
			t.Fatalf("TestGetTradeDataWithDay getTradeDataWithDay %v %v", i, lst0)

			return
		}

		tm := time.Unix(ts/1000, 0)
		if tm.Format("20060102") != lsttmstr[i] {
			t.Fatalf("TestGetTradeDataWithDay tm %v %v %v", lsttmstr[i], tm.Format("20060102"), ts)

			return
		}

		if len(retlst) != lstlen[i] {
			t.Fatalf("TestGetTradeDataWithDay len %v %v", len(retlst), lstlen[i])

			return
		}

		if len(lastlst) == 0 {
			break
		}

		if retlst == nil {
			t.Fatalf("TestGetTradeDataWithDay retlst nil")

			return
		}

		i = i + 1
		lst0 = lastlst
	}

	t.Logf("TestGetTradeDataWithDay is OK")
}
