package tradingdata

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang/protobuf/proto"
	jarvisbase "github.com/zhs007/jarviscore/base"
	"go.uber.org/zap"
)

func replyMsg(w http.ResponseWriter, msg proto.Message) {
	jsonBytes, err := json.Marshal(msg)
	if err != nil {
		jarvisbase.Warn("replyMsg:Marshal",
			zap.Error(err))

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	w.Write(jsonBytes)
}

// HTTPServer -
type HTTPServer struct {
	addr string
	serv *http.Server
	db   *tradingDataDB
}

func (s *HTTPServer) onGetTradeData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	market := r.URL.Query().Get("market")
	symbol := r.URL.Query().Get("symbol")
	daytime := r.URL.Query().Get("daytime")

	ts, err := time.Parse("20060102", daytime)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	chunk, err := s.db.getTradeData(r.Context(), market, symbol, ts.UnixNano()/1000000)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	replyMsg(w, chunk)
}

// HTTPServer -
func newHTTPServer(addr string, db *tradingDataDB) (*HTTPServer, error) {
	s := &HTTPServer{
		addr: addr,
		serv: nil,
		db:   db,
	}

	return s, nil
}

func (s *HTTPServer) start(ctx context.Context) error {

	mux := http.NewServeMux()
	mux.HandleFunc("/gettradedata", func(w http.ResponseWriter, r *http.Request) {
		s.onGetTradeData(w, r)
	})

	// fsh := http.FileServer(http.Dir("./www/static"))
	// mux.Handle("/", http.StripPrefix("/", fsh))

	server := &http.Server{
		Addr:         s.addr,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		Handler:      mux,
	}

	s.serv = server

	err := server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (s *HTTPServer) stop() {
	if s.serv != nil {
		s.serv.Close()
	}

	return
}
