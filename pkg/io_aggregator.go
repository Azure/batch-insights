package batchinsights

import (
	"time"
)

type IOAggregator struct {
	lastTimestamp *time.Time
	lastRead      uint64
	lastWrite     uint64
}

type IOStats struct {
	readBps  uint64
	writeBps uint64
}

func (aggregator *IOAggregator) UpdateAggregates(currentRead uint64, currentWrite uint64) IOStats {
	// fmt.Printf("Update agg %v\n", currentRead)
	var now = time.Now()
	var readBps uint64
	var writeBps uint64

	// fmt.Printf("NOT NIL %v %v\n", aggregator.lastTimestamp != nil, aggregator.lastRead)
	if aggregator.lastTimestamp != nil {

		var delta = now.Sub(*aggregator.lastTimestamp).Seconds()
		readBps = uint64(float64(currentRead-aggregator.lastRead) / delta)
		writeBps = uint64(float64(currentWrite-aggregator.lastWrite) / delta)
	}

	aggregator.lastTimestamp = &now
	aggregator.lastRead = currentRead
	aggregator.lastWrite = currentWrite

	// fmt.Printf("NOT NIL2 %v\n", aggregator.lastTimestamp != nil)

	return IOStats{
		readBps:  readBps,
		writeBps: writeBps,
	}
}
