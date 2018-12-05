#!/bin/bash
set -e;
wget "$BATCH_INSIGHTS_DOWNLOAD_URL" -o ./batch-insights;
chmod +x ./batch-insights;
./batch-insights > node-stats.log &