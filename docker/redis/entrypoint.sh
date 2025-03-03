#!/bin/bash

PORT=${REDIS_CLUSTER_START_PORT:-7001}
NODES=${REDIS_CLUSTER_NODES:-3}

# Computed vars
ENDPORT=$((PORT+NODES))
HOSTS=""

echo "Starting Redis Cluster with $NODES nodes"

# Loop to create individual redis instances
while [ $((PORT < ENDPORT)) != "0" ]; do
    echo "Starting $PORT"
    redis-server --port $PORT \
      --cluster-enabled yes \
      --cluster-config-file nodes-${PORT}.conf \
      --appendonly yes \
      --dbfilename dump-${PORT}.rdb \
      --logfile ${PORT}.log \
      --protected-mode no \
      --cluster-announce-ip 0.0.0.0 \
      --cluster-announce-port ${PORT} \
      --daemonize yes
    HOSTS="$HOSTS 127.0.0.1:$PORT"
    PORT=$((PORT+1))
done

sleep 2

# Uses redis cli to link the redis created above into a cluster
echo "yes" | redis-cli --cluster create $HOSTS --cluster-replicas 0

# Prints some info about the cluster
redis-cli --cluster check 127.0.0.1:${REDIS_CLUSTER_START_PORT:-7001}

# Tails logs for all the redis instances
tail -f *.log
