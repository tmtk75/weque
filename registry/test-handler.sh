#!/usr/bin/env bash
set -e

echo $$

cat<<EOF | tee -a ./test-handler.log
-- `date`
EVENT_ID:   ${EVENT_ID}
TIMESTAMP   ${TIMESTAMP}
ACTION      ${ACTION}
REPOSITORY: ${REPOSITORY}
DIGEST:     ${DIGEST}
URL:        ${URL}
TAG:        ${TAG}
REQUEST_ID: ${REQUEST_ID}
ADDR:       ${ADDR}
HOST:       ${HOST}
METHOD:     ${METHOD}
USER_AGENT: ${USER_AGENT}
ACTOR_NAME: ${ACTOR_NAME}
EOF
