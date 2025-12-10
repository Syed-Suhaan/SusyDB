package main

import "fmt"

func getCommand(workloadType string, key string, j int) string {
	val := "x"
	switch workloadType {
	case "set":
		return fmt.Sprintf("SET %s %s\r\n", key, val)
	case "setex":
		return fmt.Sprintf("SETEX %s 60 %s\r\n", key, val)
	case "get":
		return fmt.Sprintf("GET %s\r\n", key)
	case "incr":
		return fmt.Sprintf("INCR %s\r\n", key)
	case "hash":
		if j%2 == 0 {
			return fmt.Sprintf("HSET %s field1 %s\r\n", key, val)
		} else {
			return fmt.Sprintf("HGET %s field1\r\n", key)
		}
	case "mixed":
		if j%2 == 0 {
			return fmt.Sprintf("SET %s %s\r\n", key, val)
		} else {
			return fmt.Sprintf("GET %s\r\n", key)
		}
	case "publish":
		// Publish to a static channel 'bench_chan'
		return fmt.Sprintf("PUBLISH bench_chan %s\r\n", val)
	default:
		return ""
	}
}
