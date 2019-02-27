package main

import (
	"flag"
	"github.com/Azure/batch-insights/pkg"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

func parseListArgs(value string) []string {
	names := strings.Split(value, ",")
	for i := range names {
		names[i] = strings.TrimSpace(names[i])
	}
	return names
}

func getenv(key string) *string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return nil
	}
	return &value
}

func initLogger() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:    false,
		DisableTimestamp: true,
		ForceColors:      true,
	})
}

func main() {
	initLogger()
	disableArg := flag.String("disable", "", "List of metrics to disable")
	processArg := flag.String("processes", "", "List of process name to watch")

	envConfig := batchinsights.UserConfig{
		InstrumentationKey: getenv("APP_INSIGHTS_INSTRUMENTATION_KEY"),
		PoolID:             getenv("AZ_BATCH_POOL_ID"),
		NodeID:             getenv("AZ_BATCH_NODE_ID"),
	}
	processEnv := getenv("AZ_BATCH_MONITOR_PROCESSES")
	if processEnv != nil {
		envConfig.Processes = parseListArgs(*processEnv)
	}
	argsConfig := batchinsights.UserConfig{
		PoolID:             flag.String("poolID", "", "Batch pool ID"),
		NodeID:             flag.String("nodeID", "", "Batch node ID"),
		Aggregation:        flag.Int("aggregation", 1, "Aggregation in minutes"),
		InstrumentationKey: flag.String("instKey", "", "Application Insights instrumentation KEY"),
	}

	flag.Parse()
	if processArg != nil {
		argsConfig.Processes = parseListArgs(*processArg)
	}
	if disableArg != nil {
		argsConfig.Disable = parseListArgs(*disableArg)
	}

	config := envConfig.Merge(argsConfig)

	positionalArgs := flag.Args()
	if len(positionalArgs) > 0 {
		log.Warn("Using postional arguments for Node ID, PoolID, KEY and  Process names is deprecated. Use --poolID, --nodeID, --instKey, --process")
		log.Warn("It will be removed in 2.0.0")
		config.PoolID = &positionalArgs[0]
	}

	if len(positionalArgs) > 1 {
		config.NodeID = &positionalArgs[1]
	}

	if len(positionalArgs) > 2 {
		config.InstrumentationKey = &positionalArgs[2]
	}

	if len(positionalArgs) > 3 {
		config.Processes = parseListArgs(positionalArgs[3])
	}

	config.Print()

	computedConfig, err := batchinsights.ValidateAndBuildConfig(config)

	if err != nil {
		log.Error("Invalid config", err)
		os.Exit(2)
	}

	computedConfig.Print()
	batchinsights.PrintSystemInfo()
	batchinsights.ListenForStats(computedConfig)
}
