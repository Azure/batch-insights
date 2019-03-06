set -e

branch=$BATCH_INSIGHTS_BRANCH
echo "Running Batch insights dev script for linux from branch $branch"

apt-get update  
apt-get install -y git binutils bison build-essential

export GOROOT=/usr/local/go
if [ -d "$GOROOT" ]; then rm -rf $GOROOT; fi

wget https://dl.google.com/go/go1.11.linux-amd64.tar.gz
tar -xvf go1.11.linux-amd64.tar.gz
mv go /usr/local
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH

echo GO version $(go version)

git clone https://github.com/Azure/batch-insights -b $branch

cd batch-insights
go build

./batch-insights $AZ_BATCH_INSIGHTS_ARGS > $AZ_BATCH_TASK_WORKING_DIR/batch-insights.log 2>&1 &
