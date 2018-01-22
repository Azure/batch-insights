set -e

apt-get update
apt-get -y install python-dev python-pip
pip install psutil python-dateutil applicationinsights
wget --no-cache https://raw.githubusercontent.com/Azure/batch-insights/master/nodestats.py
python --version
python nodestats.py > node-stats.log 2>&1 &