package main

import (
	"fmt"
	"github.com/Azure/batch-insights/pkg"
	"log"
	"os"
	"text/template"

	"github.com/mindprince/gonvml"
)

func main() {
	gpuTest()
	var appInsightsKey = os.Getenv("APP_INSIGHTS_INSTRUMENTATION_KEY")
	var poolId = os.Getenv("AZ_BATCH_POOL_ID")
	var nodeId = os.Getenv("AZ_BATCH_NODE_ID")

	if len(os.Args) > 2 {
		poolId = os.Args[1]
		nodeId = os.Args[2]
	}

	if len(os.Args) > 3 {
		appInsightsKey = os.Args[3]
	}

	batchinsights.PrintSystemInfo()
	fmt.Printf("   Pool ID: %s\n", poolId)
	fmt.Printf("   Node ID: %s\n", nodeId)

	hiddenKey := "-"
	if appInsightsKey != "" {
		hiddenKey = "xxxxx"
	}
	fmt.Printf("   Instrumentation Key: %s\n", hiddenKey)
	batchinsights.ListenForStats(poolId, nodeId, appInsightsKey)
}

const (
	DEVICEINFO = `UUID           : {{.UUID}}
Model          : {{or .Model "N/A"}}
Path           : {{.Path}}
Power          : {{if .Power}}{{.Power}} W{{else}}N/A{{end}}
Memory	       : {{if .Memory}}{{.Memory}} MiB{{else}}N/A{{end}}
CPU Affinity   : {{if .CPUAffinity}}NUMA node{{.CPUAffinity}}{{else}}N/A{{end}}
Bus ID         : {{.PCI.BusID}}
BAR1           : {{if .PCI.BAR1}}{{.PCI.BAR1}} MiB{{else}}N/A{{end}}
Bandwidth      : {{if .PCI.Bandwidth}}{{.PCI.Bandwidth}} MB/s{{else}}N/A{{end}}
Cores          : {{if .Clocks.Cores}}{{.Clocks.Cores}} MHz{{else}}N/A{{end}}
Memory         : {{if .Clocks.Memory}}{{.Clocks.Memory}} MHz{{else}}N/A{{end}}
P2P Available  : {{if not .Topology}}None{{else}}{{range .Topology}}
    	       	 {{.BusID}} - {{(.Link.String)}}{{end}}{{end}}
---------------------------------------------------------------------
`
)

func gpuTest() {
	err := gonvml.Initialize()
	if err != nil {
		fmt.Println("Error while loading nvml")
		fmt.Println(err)
		return
	}
	defer gonvml.Shutdown()
	fmt.Printf("Initialize() took %v\n", time.Since(start))

	driverVersion, err := gonvml.SystemDriverVersion()
	if err != nil {
		fmt.Printf("SystemDriverVersion() error: %v\n", err)
		return
	}
	fmt.Printf("SystemDriverVersion(): %v\n", driverVersion)

	numDevices, err := gonvml.DeviceCount()
	if err != nil {
		fmt.Printf("DeviceCount() error: %v\n", err)
		return
	}
	fmt.Printf("DeviceCount(): %v\n", numDevices)

	for i := 0; i < int(numDevices); i++ {
		dev, err := gonvml.DeviceHandleByIndex(uint(i))
		if err != nil {
			fmt.Printf("\tDeviceHandleByIndex() error: %v\n", err)
			return
		}

		minorNumber, err := dev.MinorNumber()
		if err != nil {
			fmt.Printf("\tdev.MinorNumber() error: %v\n", err)
			return
		}
		fmt.Printf("\tminorNumber: %v\n", minorNumber)

		uuid, err := dev.UUID()
		if err != nil {
			fmt.Printf("\tdev.UUID() error: %v\n", err)
			return
		}
		fmt.Printf("\tuuid: %v\n", uuid)

		name, err := dev.Name()
		if err != nil {
			fmt.Printf("\tdev.Name() error: %v\n", err)
			return
		}
		fmt.Printf("\tname: %v\n", name)

		totalMemory, usedMemory, err := dev.MemoryInfo()
		if err != nil {
			fmt.Printf("\tdev.MemoryInfo() error: %v\n", err)
			return
		}
		fmt.Printf("\tmemory.total: %v, memory.used: %v\n", totalMemory, usedMemory)

		gpuUtilization, memoryUtilization, err := dev.UtilizationRates()
		if err != nil {
			fmt.Printf("\tdev.UtilizationRates() error: %v\n", err)
			return
		}
		fmt.Printf("\tutilization.gpu: %v, utilization.memory: %v\n", gpuUtilization, memoryUtilization)

		powerDraw, err := dev.PowerUsage()
		if err != nil {
			fmt.Printf("\tdev.PowerUsage() error: %v\n", err)
			return
		}
		fmt.Printf("\tpower.draw: %v\n", powerDraw)

		averagePowerDraw, err := dev.AveragePowerUsage(10 * time.Second)
		if err != nil {
			fmt.Printf("\tdev.AveragePowerUsage() error: %v\n", err)
			return
		}
		fmt.Printf("\taverage power.draw for last 10s: %v\n", averagePowerDraw)

		averageGPUUtilization, err := dev.AverageGPUUtilization(10 * time.Second)
		if err != nil {
			fmt.Printf("\tdev.AverageGPUUtilization() error: %v\n", err)
			return
		}
		fmt.Printf("\taverage utilization.gpu for last 10s: %v\n", averageGPUUtilization)

		temperature, err := dev.Temperature()
		if err != nil {
			fmt.Printf("\tdev.Temperature() error: %v\n", err)
			return
		}
		fmt.Printf("\ttemperature.gpu: %v C\n", temperature)

		fanSpeed, err := dev.FanSpeed()
		if err != nil {
			fmt.Printf("\tdev.FanSpeed() error: %v\n", err)
			return
		}
		fmt.Printf("\tfan.speed: %v%%\n", fanSpeed)

		encoderUtilization, _, err := dev.EncoderUtilization()
		if err != nil {
			fmt.Printf("\tdev.EncoderUtilization() error: %v\n", err)
			return
		}
		fmt.Printf("\tutilization.encoder: %d\n", encoderUtilization)

		decoderUtilization, _, err := dev.DecoderUtilization()
		if err != nil {
			fmt.Printf("\tdev.DecoderUtilization() error: %v\n", err)
			return
		}
		fmt.Printf("\tutilization.decoder: %d\n", decoderUtilization)
		fmt.Println()
	}
}
