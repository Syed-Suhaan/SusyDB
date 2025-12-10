package server

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/Syed-Suhaan/SusyDB/pkg/core"
)

func handleSet(conn net.Conn, store *core.KVStore, parts []string) []byte {
	if len(parts) < 3 {
		return []byte("-ERR wrong number of arguments for 'set' command\r\n")
	}
	key := parts[1]
	value := strings.Join(parts[2:], " ")
	if err := store.Set(key, value, 0); err != nil {
		return []byte("-ERR " + err.Error() + "\r\n")
	}
	return []byte("+OK\r\n")
}

func handleSetEx(conn net.Conn, store *core.KVStore, parts []string) []byte {
	if len(parts) < 4 {
		return []byte("-ERR wrong number of arguments for 'setex' command\r\n")
	}
	key := parts[1]
	secondsStr := parts[2]
	value := strings.Join(parts[3:], " ")
	seconds, err := strconv.ParseInt(secondsStr, 10, 64)
	if err != nil {
		return []byte("-ERR invalid expire time in 'setex' command\r\n")
	}
	if err := store.Set(key, value, seconds); err != nil {
		return []byte("-ERR " + err.Error() + "\r\n")
	}
	return []byte("+OK\r\n")
}

func handleGet(conn net.Conn, store *core.KVStore, parts []string) []byte {
	if len(parts) < 2 {
		return []byte("-ERR wrong number of arguments for 'get' command\r\n")
	}
	key := parts[1]
	val, ok, err := store.Get(key)
	if err != nil {
		return []byte("-WARN " + err.Error() + "\r\n")
	} else if !ok {
		return []byte("$-1\r\n")
	}
	return []byte(fmt.Sprintf("%s\r\n", val))
}

func handleIncr(conn net.Conn, store *core.KVStore, parts []string) []byte {
	if len(parts) < 2 {
		return []byte("-ERR wrong number of arguments for 'incr' command\r\n")
	}
	key := parts[1]
	newVal, err := store.IncrBy(key, 1)
	if err != nil {
		return []byte("-ERR " + err.Error() + "\r\n")
	}
	return []byte(fmt.Sprintf(":%d\r\n", newVal))
}

func handleIncrBy(conn net.Conn, store *core.KVStore, parts []string) []byte {
	if len(parts) < 3 {
		return []byte("-ERR wrong number of arguments for 'incrby' command\r\n")
	}
	key := parts[1]
	deltaStr := parts[2]
	delta, err := strconv.ParseInt(deltaStr, 10, 64)
	if err != nil {
		return []byte("-ERR value is not an integer or out of range\r\n")
	}

	newVal, err := store.IncrBy(key, delta)
	if err != nil {
		return []byte("-ERR " + err.Error() + "\r\n")
	}
	return []byte(fmt.Sprintf(":%d\r\n", newVal))
}

func handleDel(conn net.Conn, store *core.KVStore, parts []string) []byte {
	if len(parts) < 2 {
		return []byte("-ERR wrong number of arguments for 'del' command\r\n")
	}
	key := parts[1]
	store.Delete(key)
	return []byte("+OK\r\n")
}
