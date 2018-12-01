#!/usr/bin/env bash
set -e

echo $$
sleep ${1-0}

cat<<EOF | tee -a ./test-handler.log
-- `date`
ACTION:           ${ACTION}
REPOSITORY_NAME:  ${REPOSITORY_NAME}
OWNER_NAME:       ${OWNER_NAME}
EVENT:            ${EVENT}
DELIVERY:         ${DELIVERY}
REF:              ${REF}
AFTER:            ${AFTER}
BEFORE:           ${BEFORE}
CREATED:          ${CREATED}
DELETED:          ${DELETED}
PUSHER_NAME:      ${PUSHER_NAME}
RELEASE_TAG_NAME: ${RELEASE_TAG_NAME}
EOF
