package main

import (
	"fmt"

	"./ipTracker"
)

func main() {
	ipTracker.Clear()
	ipTracker.RequestHandled("192.168.0.1")
	ipTracker.RequestHandled("10.0.0.1")
	ipTracker.RequestHandled("8.8.8.8")
	ipTracker.RequestHandled("8.8.8.8")
	ipTracker.RequestHandled("8.8.8.8")
	ipTracker.RequestHandled("192.168.0.1")
	topEntries := ipTracker.Top100()

	for _, entry := range topEntries {
		fmt.Println(entry)
	}
}
