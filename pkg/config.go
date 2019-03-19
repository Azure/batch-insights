package batchinsights

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// DefaultAggregationTime default time range where metrics are preaggregated locally
const DefaultAggregationTime = time.Duration(1) * time.Minute

// DefaultSamplingRate default time between metrics sampling
const DefaultSamplingRate = time.Duration(5) * time.Second

// UserConfig config provided by the user either via command line, file or environemnt variable.
type UserConfig struct {
	PoolID             *string
	NodeID             *string
	InstrumentationKey *string  // Application insights instrumentation key
	Processes          []string // List of process names to watch
	Aggregation        *int     // Local aggregation of data in minutes (default: 1)
	Disable            []string // List of metrics to disable
}

// Print print the config to console
func (config UserConfig) Print() {
	fmt.Printf("User configuration:\n")
	fmt.Printf("   Pool ID: %s\n", *config.PoolID)
	fmt.Printf("   Node ID: %s\n", *config.NodeID)
	if config.InstrumentationKey != nil {
		fmt.Printf("   Instrumentation Key: %s\n", hideSecret(*config.InstrumentationKey))
	}
	fmt.Printf("   Aggregation: %d\n", *config.Aggregation)
	fmt.Printf("   Disable: %v\n", config.Disable)
	fmt.Printf("   Monitoring processes: %v\n", config.Processes)
}

// Merge with another config
func (config UserConfig) Merge(other UserConfig) UserConfig {
	if other.PoolID != nil && *other.PoolID != "" {
		config.PoolID = other.PoolID
	}
	if other.NodeID != nil && *other.NodeID != "" {
		config.NodeID = other.NodeID
	}
	if other.InstrumentationKey != nil && *other.InstrumentationKey != "" {
		config.InstrumentationKey = other.InstrumentationKey
	}
	if other.Aggregation != nil {
		config.Aggregation = other.Aggregation
	}
	if len(other.Processes) > 0 {
		config.Processes = other.Processes
	}
	if len(other.Disable) > 0 {
		config.Disable = other.Disable
	}
	return config
}

// DisableConfig config showing which feature are disabled
type DisableConfig struct {
	DiskIO    bool `json:"diskIO"`
	DiskUsage bool `json:"diskUsage"`
	NetworkIO bool `json:"networkIO"`
	GPU       bool `json:"gpu"`
	CPU       bool `json:"cpu"`
	Memory    bool `json:"memory"`
}

func (d DisableConfig) String() string {
	s, _ := json.Marshal(d)
	return string(s)
}

// Config General config batch insights takes as input
type Config struct {
	PoolID             string
	NodeID             string
	InstrumentationKey string
	Processes          []string
	Aggregation        time.Duration
	SamplingRate       time.Duration
	Disable            DisableConfig
}

// Print print the config to console
func (config Config) Print() {
	fmt.Printf("BatchInsights configuration:\n")
	fmt.Printf("   Pool ID: %s\n", config.PoolID)
	fmt.Printf("   Node ID: %s\n", config.NodeID)
	fmt.Printf("   Instrumentation Key: %s\n", hideSecret(config.InstrumentationKey))
	fmt.Printf("   Aggregation: %v\n", config.Aggregation)
	fmt.Printf("   Sampling rate: %d\n", config.SamplingRate)
	fmt.Printf("   Disable: %+v\n", config.Disable)
	fmt.Printf("   Monitoring processes: %v\n", config.Processes)
}

// ValidateAndBuildConfig Convert Batch insights user config into config taken by the library
func ValidateAndBuildConfig(userConfig UserConfig) (Config, error) {
	aggregation := parseAggregation(userConfig.Aggregation)

	if userConfig.PoolID == nil {
		return Config{}, errors.New("Pool ID must be specified")
	}
	if userConfig.PoolID == nil {
		return Config{}, errors.New("Node ID must be specified")
	}
	key := ""
	if userConfig.InstrumentationKey != nil {
		key = *userConfig.InstrumentationKey
	}
	return Config{
		PoolID:             *userConfig.PoolID,
		NodeID:             *userConfig.NodeID,
		InstrumentationKey: key,
		Processes:          userConfig.Processes,
		Aggregation:        aggregation,
		Disable:            parseDisableConfig(userConfig.Disable),
		SamplingRate:       DefaultSamplingRate,
	}, nil
}

func parseAggregation(value *int) time.Duration {
	if value == nil {
		return DefaultAggregationTime
	}
	return time.Duration(*value) * time.Minute
}

func parseDisableConfig(values []string) DisableConfig {
	disableMap := make(map[string]bool)
	for _, key := range values {
		disableMap[strings.ToLower(key)] = true
	}
	return DisableConfig{
		DiskIO:    disableMap["diskio"],
		DiskUsage: disableMap["diskusage"],
		NetworkIO: disableMap["networkio"],
		GPU:       disableMap["gpu"],
		CPU:       disableMap["cpu"],
		Memory:    disableMap["memory"],
	}
}

// Hide a secret
func hideSecret(secret string) string {
	if secret == "" {
		return "-"
	}
	return "xxxxx"
}
