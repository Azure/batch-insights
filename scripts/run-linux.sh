#!/bin/bash
set -e;
wget -o ./batch-insights "$BATCH_INSIGHTS_DOWNLOAD_URL";
chmod +x ./batch-insights;
./batch-insights > node-stats.log &