#!/bin/bash
#
#    SPDX-License-Identifier: Apache-2.0
#

EXPLORER_PROCESS_ID=$(ps aux | grep "[n]ame - hyperledger-explorer" | awk '{print $2}')

if [ -n "$EXPLORER_PROCESS_ID" ]; then
    echo "Stopping node process hyperledger-explorer, id: $EXPLORER_PROCESS_ID"
    kill -SIGTERM $EXPLORER_PROCESS_ID
    echo "Server stopped."
else
    echo "No process name hyperledger-explorer found."
fi
