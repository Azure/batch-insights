package main

import "github.com/Azure/batch-insights/lib"

func main() {
	batchinsights.PrintSystemInfo()
	batchinsights.ListenForStats()
}
