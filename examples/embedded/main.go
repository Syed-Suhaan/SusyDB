package main

import (
	"fmt"
	"time"

	"github.com/Syed-Suhaan/susydb/pkg/core"
)

func main() {
	fmt.Println("SusyDB Embedded Mode Demo")

	// Initialize the store directly
	store := core.NewKVStore()

	// 1. Set key
	fmt.Println("1. Setting key 'greeting'...")
	store.Set("greeting", "Hello from Embedded SusyDB!", 0)

	// 2. Get key
	val, ok, _ := store.Get("greeting")
	if ok {
		fmt.Printf("2. Got value: %s\n", val)
	}

	// 3. Counter
	fmt.Println("3. Incrementing counter...")
	newVal, _ := store.IncrBy("my_counter", 1)
	fmt.Printf("   Counter is now: %d\n", newVal)

	// 4. TTL
	fmt.Println("4. Setting volatile key (1s TTL)...")
	store.Set("volatile", "I will disappear", 1)

	time.Sleep(1500 * time.Millisecond)

	_, ok, _ = store.Get("volatile")
	if !ok {
		fmt.Println("5. Volatile key expired as expected")
	} else {
		fmt.Println("Volatile key failed to expire!")
	}
}
