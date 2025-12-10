package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"
)

var (
	host        = flag.String("h", "localhost:7379", "SusyDB host address")
	concurrency = flag.Int("c", 50, "Number of concurrent connections")
	requests    = flag.Int("n", 100000, "Total number of requests")
	workload    = flag.String("test", "mixed", "Workload type: set, get, mixed, setex, incr, hash")
)

func main() {
	flag.Parse()

	fmt.Printf("Benchmarking %s | Test: %s | Clients: %d | Reqs: %d\n", *host, *workload, *concurrency, *requests)

	// Pre-warm / Pre-populate for GET tests
	if *workload == "get" {
		fmt.Println("Pre-populating keys for GET test...")
		populateKeys(*requests / 10)
	}

	start := time.Now()
	var wg sync.WaitGroup
	latencyChan := make(chan time.Duration, *requests*2)

	requestsPerWorker := *requests / *concurrency

	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			conn, err := net.Dial("tcp", *host)
			if err != nil {
				fmt.Printf("Worker %d connect error: %v\n", workerID, err)
				return
			}
			defer conn.Close()
			reader := bufio.NewReader(conn)

			for j := 0; j < requestsPerWorker; j++ {
				key := fmt.Sprintf("key:%d", rand.Intn(10000))
				cmd := getCommand(*workload, key, j)

				t0 := time.Now()
				fmt.Fprint(conn, cmd)
				_, err = reader.ReadString('\n')
				if err != nil {
					return
				}
				latencyChan <- time.Since(t0)
			}
		}(i)
	}

	wg.Wait()
	close(latencyChan)
	totalDuration := time.Since(start)

	var latencies []time.Duration
	for l := range latencyChan {
		latencies = append(latencies, l)
	}

	printReport(*workload, totalDuration, latencies)
}

func populateKeys(count int) {
	conn, err := net.Dial("tcp", *host)
	if err != nil {
		return
	}
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for i := 0; i < count; i++ {
		key := fmt.Sprintf("key:%d", i)
		fmt.Fprintf(conn, "SET %s value_payload\r\n", key)
		reader.ReadString('\n')
	}
}
