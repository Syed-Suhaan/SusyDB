package main

import (
	"fmt"
	"os"
	"sort"
	"time"
)

func printReport(workloadType string, totalDuration time.Duration, latencies []time.Duration) {
	sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })

	totalOps := len(latencies)
	if totalOps == 0 {
		fmt.Println("No successful operations.")
		os.Exit(1)
	}

	rps := float64(totalOps) / totalDuration.Seconds()
	p50 := latencies[totalOps*50/100]
	p99 := latencies[totalOps*99/100]

	fmt.Println("\n------------------------------------------------")
	fmt.Printf("Summary (%s):\n", workloadType)
	fmt.Printf("  Total Ops:   %d\n", totalOps)
	fmt.Printf("  Duration:    %v\n", totalDuration)
	fmt.Printf("  Throughput:  %.2f requests/sec\n", rps)
	fmt.Printf("  P50 Latency: %v\n", p50)
	fmt.Printf("  P99 Latency: %v\n", p99)
	fmt.Println("------------------------------------------------")
}
