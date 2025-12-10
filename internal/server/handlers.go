package server

import (
	"net"

	"github.com/Syed-Suhaan/SusyDB/pkg/core"
)

// CommandHandler is the function signature for all command handlers.
// It returns a response byte slice. If nil, it means the handler handled writing itself (e.g. Subscribe).
type CommandHandler func(conn net.Conn, store *core.KVStore, parts []string) []byte

// Handlers map maps command strings to their handler functions.
var Handlers = map[string]CommandHandler{
	"SET":       handleSet,
	"SETEX":     handleSetEx,
	"GET":       handleGet,
	"INCR":      handleIncr,
	"INCRBY":    handleIncrBy,
	"HSET":      handleHSet,
	"HGET":      handleHGet,
	"HGETALL":   handleHGetAll,
	"HDEL":      handleHDel,
	"DEL":       handleDel,
	"INFO":      handleInfo,
	"PING":      handlePing,
	"PUBLISH":   handlePublish,
	"SUBSCRIBE": handleSubscribe,
}
