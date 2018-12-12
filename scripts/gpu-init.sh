set -e
# apt remove nvidia-cuda-toolkit
# apt remove nvidia-*
apt update
apt-key adv --fetch-keys  http://developer.download.nvidia.com/compute/cuda/repos/ubuntu1804/x86_64/7fa2af80.pub
bash -c 'echo "deb http://developer.download.nvidia.com/compute/cuda/repos/ubuntu1804/x86_64 /" > /etc/apt/sources.list.d/cuda.list'
apt update
apt install -y nvidia-driver-410 --no-install-recommends
apt install -y cuda-10-0 --no-install-recommends
apt-get install -y git binutils bison build-essential --no-install-recommends
