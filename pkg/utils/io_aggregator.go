package utils

import (
	"time"
)

type IOAggregator struct {
	lastTimestamp *time.Time
	lastRead      uint64
	lastWrite     uint64
}

type IOStats struct {
	ReadBps  uint64
	WriteBps uint64
}

func (aggregator *IOAggregator) UpdateAggregates(currentRead uint64, currentWrite uint64) IOStats {
	var now = time.Now()
	var readBps uint64
	var writeBps uint64

	if aggregator.lastTimestamp != nil {

		var delta = now.Sub(*aggregator.lastTimestamp).Seconds()
		readBps = uint64(float64(currentRead-aggregator.lastRead) / delta)
		writeBps = uint64(float64(currentWrite-aggregator.lastWrite) / delta)
	}

	aggregator.lastTimestamp = &now
	aggregator.lastRead = currentRead
	aggregator.lastWrite = currentWrite

	return IOStats{
		ReadBps:  readBps,
		WriteBps: writeBps,
	}
}
