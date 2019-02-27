#!/bin/bash
set -e;
echo "RUn"
wget -O ./batch-insights "$BATCH_INSIGHTS_DOWNLOAD_URL";
chmod +x ./batch-insights;
./batch-insights $AZ_BATCH_INSIGHTS_ARGS  > batch-insights.log &