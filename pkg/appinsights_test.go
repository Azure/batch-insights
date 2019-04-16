package batchinsights_test

import (
	"github.com/Azure/batch-insights/pkg"
	"github.com/Microsoft/ApplicationInsights-Go/appinsights"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetMetricID(t *testing.T) {
	metric := appinsights.NewMetricTelemetry("Disk usage", 134)
	metric.Properties["Some #"] = "4"
	metric.Properties["Other #"] = "5"

	assert.Ok(t == "Disk usage/Other #=5,Some #=4" || t == "Disk usage/Some #=4,Other #=5", batchinsights.GetMetricID(metric))

	metric = appinsights.NewMetricTelemetry("Disk IO", 543)
	assert.Equal(t, "Disk IO/", batchinsights.GetMetricID(metric))
}
