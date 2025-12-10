package server

import (
	"fmt"
	"net"
	"strings"

	"github.com/Syed-Suhaan/SusyDB/pkg/core"
)

func handleHSet(conn net.Conn, store *core.KVStore, parts []string) []byte {
	if len(parts) < 4 {
		return []byte("-ERR wrong number of arguments for 'hset' command\r\n")
	}
	key := parts[1]
	field := parts[2]
	value := strings.Join(parts[3:], " ")
	err := store.HSet(key, field, value)
	if err != nil {
		return []byte("-WARN " + err.Error() + "\r\n")
	}
	return []byte("+OK\r\n")
}

func handleHGet(conn net.Conn, store *core.KVStore, parts []string) []byte {
	if len(parts) < 3 {
		return []byte("-ERR wrong number of arguments for 'hget' command\r\n")
	}
	key := parts[1]
	field := parts[2]
	val, ok, err := store.HGet(key, field)
	if err != nil {
		return []byte("-WARN " + err.Error() + "\r\n")
	} else if !ok {
		return []byte("$-1\r\n")
	}
	return []byte(fmt.Sprintf("%s\r\n", val))
}

func handleHGetAll(conn net.Conn, store *core.KVStore, parts []string) []byte {
	if len(parts) < 2 {
		return []byte("-ERR wrong number of arguments for 'hgetall' command\r\n")
	}
	key := parts[1]
	hash, ok, err := store.HGetAll(key)
	if err != nil {
		return []byte("-WARN " + err.Error() + "\r\n")
	} else if !ok {
		return []byte("*0\r\n")
	}
	var sb strings.Builder
	for k, v := range hash {
		sb.WriteString(fmt.Sprintf("%s\n%s\n", k, v))
	}
	return []byte(sb.String())
}

func handleHDel(conn net.Conn, store *core.KVStore, parts []string) []byte {
	if len(parts) < 3 {
		return []byte("-ERR wrong number of arguments for 'hdel' command\r\n")
	}
	key := parts[1]
	field := parts[2]
	deleted, err := store.HDel(key, field)
	if err != nil {
		return []byte("-WARN " + err.Error() + "\r\n")
	} else if deleted {
		return []byte(":1\r\n")
	}
	return []byte(":0\r\n")
}
