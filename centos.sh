set -e
yum -y install epel-release
yum -y install gcc python-pip python-devel

echo "Python version:"
python --version
echo "Pip version:"
pip --version
pip install psutil python-dateutil applicationinsights

wget --no-cache https://raw.githubusercontent.com/Azure/batch-insights/master/nodestats.py
python nodestats.py > node-stats.log 2>&1 &