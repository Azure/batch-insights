set -e

apt-get update  
apt-get install -y git binutils bison build-essential

export GOROOT=/usr/local/go
if [ -d "$GOROOT" ]; then rm -rf $GOROOT; fi

wget https://dl.google.com/go/go1.11.linux-amd64.tar.gz
tar -xvf go1.11.linux-amd64.tar.gz
mv go /usr/local
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH

echo GO version $(go version)

git clone https://github.com/Azure/batch-insights

cd batch-insights
go build

./batch-insights > $AZ_BATCH_TASK_WORKING_DIR/node-stats.log 2>&1 &
