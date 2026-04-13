#!/bin/bash
#
#    SPDX-License-Identifier: Apache-2.0
#
# Redirecting console.log to log file.
# Please visit ./logs/app to view the application logs and visit the ./logs/db to view the Database logs and visit the ./logs/console for the console.log

echo "************************************************************************************"
echo "**************************** Hyperledger Explorer **********************************"
echo "************************************************************************************"

export LOG_LEVEL_APP=${LOG_LEVEL_APP:-debug}
export LOG_LEVEL_DB=${LOG_LEVEL_DB:-debug}
export LOG_LEVEL_CONSOLE=${LOG_LEVEL_CONSOLE:-info}
export LOG_CONSOLE_STDOUT=${LOG_CONSOLE_STDOUT:-false}

export DISCOVERY_AS_LOCALHOST=${DISCOVERY_AS_LOCALHOST:-false}
export EXPLORER_APP_ROOT=${EXPLORER_APP_ROOT:-dist}
export PORT=${PORT:-8081}

#clear cache wallet
rm -rf wallet/*
rm -rf logs/*
mkdir -p logs/console

echo "Server starting in background..."
nohup node ${EXPLORER_APP_ROOT}/main.js name - hyperledger-explorer > logs/console/run.log 2>&1 &
EXPLORER_PID=$!

echo "Server running in background. PID: $EXPLORER_PID"
echo "Check logs/console/run.log for output."
