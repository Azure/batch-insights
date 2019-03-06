# BatchInsights Configuration

Batch Insights provides various configuration option(Version `1.2.0` and above).


#### `--poolID <value>` 
Pool ID. Override pool ID provided by the `AZ_BATCH_POOL_ID` environment variable
#### `--nodeID <value>` 
Node ID. Override node ID provided by the `AZ_BATCH_NODE_ID` environment variable
#### `--instKey <value>` 
Instrumentation key. Application Insights instrumentation key to emit the metrics
#### `--disable <value>` 
Comma separated list of metrics to disable. e.g. `--disable networkIO,diskUsage`

Available metrics names:
    - diskIO
    - diskUsage
    - networkIO
    - memory
    - CPU
    - GPU

* `--aggregation <value>` Number in minutes to aggregate the data locally. Defaults to 1 minute 

Example: `--agregation 5` to aggregate for 5 minutes

#### `--processes <value>` 
Comma separated list of processes to monitor.

Example: `--processes notepad.exe,explorer.exe`
