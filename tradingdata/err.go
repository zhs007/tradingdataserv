package tradingdata

import "errors"

var (
	// ErrNoBindAddress - no bind address
	ErrNoBindAddress = errors.New("no bind address")
	// ErrNoAnkaDBConfig - no ankadb config
	ErrNoAnkaDBConfig = errors.New("no ankadb config")
	// ErrNoHTTPServerAddr - no http server address
	ErrNoHTTPServerAddr = errors.New("no http server address")
	// ErrInvalidInsertIndex - invalid insert index
	ErrInvalidInsertIndex = errors.New("invalid insert index")
	// ErrInvalidTradeDataID - invalid tradedata id
	ErrInvalidTradeDataID = errors.New("invalid tradedata id")
)
