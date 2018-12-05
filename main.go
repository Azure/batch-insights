package main

import "github.com/Azure/batch-insights/pkg"

func main() {
	batchinsights.PrintSystemInfo()
	batchinsights.ListenForStats()
}
