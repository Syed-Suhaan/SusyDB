package server

import (
	"fmt"
	"net"

	"github.com/Syed-Suhaan/susydb/pkg/core"
)

func handleInfo(conn net.Conn, store *core.KVStore, parts []string) []byte {
	info := store.Info()
	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(info), info))
}

func handlePing(conn net.Conn, store *core.KVStore, parts []string) []byte {
	return []byte("+PONG\r\n")
}
