#!/usr/bin/env bash

set -eou pipefail
#set -x  # useful for debugging

docker_cleanup() {
    echo "cleaning up existing network and containers..."
    CONTAINERS='key'
    docker ps | grep -E ${CONTAINERS} | awk '{print $1}' | xargs -I {} docker stop {} || true
    docker ps -a | grep -E ${CONTAINERS} | awk '{print $1}' | xargs -I {} docker rm {} || true
    docker network list | grep ${CONTAINERS} | awk '{print $2}' | xargs -I {} docker network rm {} || true
}

# optional settings (generally defaults should be fine, but sometimes useful for debugging)
KEY_LOG_LEVEL="${KEY_LOG_LEVEL:-INFO}"  # or DEBUG
KEY_TIMEOUT="${KEY_TIMEOUT:-5}"  # 10, or 20 for really sketchy network

# local and filesystem constants
LOCAL_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# container command constants
KEY_IMAGE="gcr.io/elxir-core-infra/key:snapshot" # develop

echo
echo "cleaning up from previous runs..."
docker_cleanup

echo
echo "creating key docker network..."
docker network create key

echo
echo "starting key..."
port=10100
name="key-0"
docker run --name "${name}" --net=key -d -p ${port}:${port} ${KEY_IMAGE} \
    start \
    --logLevel "${KEY_LOG_LEVEL}" \
    --serverPort ${port}
key_addrs="${name}:${port}"
key_containers="${name}"

echo
echo "testing key health..."
docker run --rm --net=key ${KEY_IMAGE} test health \
    --keys "${key_addrs}" \
    --logLevel "${KEY_LOG_LEVEL}"

echo
echo "testing key ..."
docker run --rm --net=key ${KEY_IMAGE} test io \
    --keys "${key_addrs}" \
    --logLevel "${KEY_LOG_LEVEL}"

echo
echo "cleaning up..."
docker_cleanup

echo
echo "All tests passed."
