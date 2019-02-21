package batchinsights

import (
	"strings"
	"time"
)

// UserConfig config provided by the user either via command line, file or environemnt variable.
type UserConfig struct {
	PoolID             string
	NodeID             string
	InstrumentationKey *string  // Application insights instrumentation key
	Process            []string // List of process names to watch
	Aggregation        *int     // Local aggregation of data in minutes (default: 1)
	Disable            []string // List of metrics to disable
}

// DisableConfig config showing which feature are disabled
type DisableConfig struct {
	DiskIO    bool
	DiskUsage bool
	NetworkIO bool
	GPU       bool
	CPU       bool
	Memory    bool
}

// Config General config batch insights takes as input
type Config struct {
	PoolID             string
	NodeID             string
	InstrumentationKey *string
	Process            []string
	Aggregation        time.Duration
	Disable            *DisableConfig
}

// BuildConfig Convert Batch insights user config into config taken by the library
func BuildConfig(userConfig UserConfig) Config {
	return Config{
		PoolID:             userConfig.PoolID,
		NodeID:             userConfig.NodeID,
		InstrumentationKey: userConfig.InstrumentationKey,
		Process:            userConfig.Process,
		Aggregation:        parseAggregation(userConfig.Aggregation),
		Disable:            parseDisableConfig(userConfig.Disable),
	}
}

func parseAggregation(value *int) time.Duration {
	if value == nil {
		return time.Duration(1) * time.Minute
	}
	return time.Duration(*value) * time.Minute
}

func parseDisableConfig(values []string) *DisableConfig {
	disableMap := make(map[string]bool)
	for _, key := range values {
		disableMap[strings.ToLower(key)] = true
	}
	return &DisableConfig{
		DiskIO:    disableMap["diskio"],
		DiskUsage: disableMap["diskusage"],
		NetworkIO: disableMap["networkio"],
		GPU:       disableMap["gpu"],
		CPU:       disableMap["cpu"],
		Memory:    disableMap["memory"],
	}
}
