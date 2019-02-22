package batchinsights_test

import (
	"github.com/Azure/batch-insights/pkg"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildConfig(t *testing.T) {
	pool1 := "pool-1"
	node1 := "node-1"

	input := batchinsights.UserConfig{
		PoolID:  &pool1,
		NodeID:  &node1,
		Process: []string{"foo.exe", "bar"},
	}
	result := batchinsights.BuildConfig(input)

	assert.Equal(t, "pool-1", *result.PoolID)
	assert.Equal(t, "node-1", *result.NodeID)
	assert.Equal(t, []string{"foo.exe", "bar"}, result.Process)
	assert.Equal(t, false, result.Disable.DiskIO)
	assert.Equal(t, false, result.Disable.NetworkIO)
	assert.Equal(t, false, result.Disable.DiskUsage)
	assert.Equal(t, false, result.Disable.CPU)
	assert.Equal(t, false, result.Disable.Memory)
	assert.Equal(t, false, result.Disable.GPU)

	result = batchinsights.BuildConfig(batchinsights.UserConfig{
		PoolID:  &pool1,
		NodeID:  &node1,
		Disable: []string{"diskIO", "cpu"},
	})

	// assert.Equal(t, []string{}, result.Process)
	assert.Equal(t, true, result.Disable.DiskIO)
	assert.Equal(t, false, result.Disable.NetworkIO)
	assert.Equal(t, false, result.Disable.DiskUsage)
	assert.Equal(t, true, result.Disable.CPU)
	assert.Equal(t, false, result.Disable.Memory)
	assert.Equal(t, false, result.Disable.GPU)
}
