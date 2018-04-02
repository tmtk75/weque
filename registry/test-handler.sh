#!/usr/bin/env bash
set -e

echo $$

cat<<EOF | tee -a ./test-handler.log
-- `date`
EVENT_ID:   ${EVENT_ID}
REPOSITORY: ${REPOSITORY}
DIGEST:     ${DIGEST}
URL:        ${URL}
TAG:        ${TAG}
REQUEST_ID: ${REQUEST_ID}
ADDR:       ${ADDR}
HOST:       ${HOST}
USER_AGENT: ${USER_AGENT}
EOF
