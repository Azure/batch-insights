set -e

apt-get update
apt-get -y install python-dev python-pip
pip install psutil python-dateutil applicationinsights==0.11.3
wget --no-cache https://raw.githubusercontent.com/Azure/batch-insights/master/nodestats.py
python --version
python nodestats.py > batch-insights.log 2>&1 &
