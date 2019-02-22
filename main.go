package main

import (
	"flag"
	"fmt"
	"github.com/Azure/batch-insights/pkg"
	"os"
	"strings"
)

func parseListArgs(value string) []string {
	return strings.Split(value, ",")
}

func getenv(key string) *string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return nil
	}
	return &value
}

func main() {
	disableArg := flag.String("disable", "", "List of metrics to disable")
	processArg := flag.String("process", "", "List of process name to watch")

	envConfig := batchinsights.UserConfig{
		InstrumentationKey: getenv("APP_INSIGHTS_INSTRUMENTATION_KEY"),
		PoolID:             getenv("AZ_BATCH_POOL_ID"),
		NodeID:             getenv("AZ_BATCH_NODE_ID"),
	}
	processEnv := getenv("AZ_BATCH_MONITOR_PROCESSES")
	if processEnv != nil {
		envConfig.Process = parseListArgs(*processEnv)
	}
	argsConfig := batchinsights.UserConfig{
		PoolID:             flag.String("poolID", "", "Batch pool ID"),
		NodeID:             flag.String("nodeID", "", "Batch node ID"),
		Aggregation:        flag.Int("aggregation", 1, "Aggregation in minutes"),
		InstrumentationKey: flag.String("instkey", "", "Application Insights instrumentation KEY"),
	}

	flag.Parse()
	if processArg != nil {
		argsConfig.Process = parseListArgs(*processArg)
	}
	if disableArg != nil {
		argsConfig.Disable = parseListArgs(*disableArg)
	}

	config := envConfig.Merge(argsConfig)
	config.Print()

	fmt.Printf("%v\n", flag.Args())

	// if len(os.Args) > 3 {
	// 	appInsightsKey = os.Args[3]
	// }

	// if len(os.Args) > 4 {
	// 	processNamesStr = os.Args[4]
	// }

	// processNames := strings.Split(processNamesStr, ",")
	// for i := range processNames {
	// 	processNames[i] = strings.TrimSpace(processNames[i])
	// }

	// batchinsights.PrintSystemInfo()
	// fmt.Printf("   Pool ID: %s\n", poolId)
	// fmt.Printf("   Node ID: %s\n", nodeId)

	// hiddenKey := "-"
	// if appInsightsKey != "" {
	// 	hiddenKey = "xxxxx"
	// }

	// fmt.Printf("   Instrumentation Key: %s\n", hiddenKey)

	// fmt.Printf("   Monitoring processes: %s\n", strings.Join(processNames, ", "))

	// batchinsights.ListenForStats(poolId, nodeId, appInsightsKey, processNames)
}
