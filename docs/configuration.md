# BatchInsights Configuration

Batch Insights provides various configuration option(Version `1.2.0` and above).


#### `--poolID <value>` 
Pool ID. Override pool ID provided by the `AZ_BATCH_POOL_ID` environment variable
#### `--nodeID <value>` 
Node ID. Override node ID provided by the `AZ_BATCH_NODE_ID` environment variable
#### `--instKey <value>` 
Instrumentation key. Application insights instrumentation key to emit the metrics
#### `--disable <value>` 
List of metrics comma separated to disable. e.g. `--disable networkIO,diskUsage`

Available options
    - diskIO
    - diskUsage
    - networkIO
    - memory
    - CPU
    - GPU

* `--aggregation <value>` Number in minutes to aggregate the data locally. Default to 1 minute 

Example: `--agregation 5` to aggreate for 5 minutes

#### `--processes <value>` 
List of process names to monitor comma separated.

Example: `--processes notepad.exe,explorer.exe`
