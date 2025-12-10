package main

import (
	"flag"
	"fmt"

	"github.com/Syed-Suhaan/susydb/internal/server"
	"github.com/Syed-Suhaan/susydb/pkg/core"
)

func main() {
	addr := flag.String("addr", ":7379", "Server address")
	flag.Parse()

	// 1. Initialize the Store
	store := core.NewKVStore()

	// 2. Start the Garbage Collector
	fmt.Println("🧹 Starting Background Garbage Collector...")
	store.StartGC()

	// 3. Start the TCP Server
	server.Start(store, *addr)
}
