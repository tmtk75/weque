#!/usr/bin/env bash
set -e

echo $$
sleep ${1-3}

log_path=./handler.log
echo "-- `date`" | tee -a $log_path
echo "REPOSITORY_NAME: $REPOSITORY_NAME" | tee -a $log_path
echo "OWNER_NAME:      $OWNER_NAME"      | tee -a $log_path
echo "EVENT:           $EVENT"           | tee -a $log_path
echo "DELIVERY:        $DELIVERY"        | tee -a $log_path
echo "REF:             $REF"             | tee -a $log_path
echo "AFTER:           $AFTER"           | tee -a $log_path
echo "BEFORE:          $BEFORE"          | tee -a $log_path
echo "CREATED:         $CREATED"         | tee -a $log_path
echo "DELETED:         $DELETED"         | tee -a $log_path
echo "PUSHER_NAME:     $PUSHER_NAME"     | tee -a $log_path
echo "" | tee -a $log_path
