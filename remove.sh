#!/usr/bin/env bash

output=$(go run main.go 2>&1)
line=0
re="app/algorithms/stop_words.go:[0-9]*:3:"
lines_to_delete="NR!="
separator=" && NR!="
for i in $output; do
   [[ $i =~ $re ]] || continue
   line=${i#*:}
   line=${line%:*}
   line=${line%:*}
   lines_to_delete="$lines_to_delete$line$separator"
   echo ${lines_to_delete}
done

lines_to_delete=${lines_to_delete%\&\&*}
lines_to_delete=${lines_to_delete%\&\&*}
echo $lines_to_delete
awk "$lines_to_delete { print }" app/algorithms/stop_words.go > app/algorithms/stop_words1.go
rm app/algorithms/stop_words.go
mv app/algorithms/stop_words1.go app/algorithms/stop_words.go

