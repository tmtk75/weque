#!/usr/bin/env bash
c=$(cat .counter)
c=$((c + 1))

./weque parse-event user-event 2> stderr | tee a-$c.txt
#tee a-$c.txt
#set | grep -i index | tee b-$c.txt
echo $c > .counter
