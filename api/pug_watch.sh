#!/bin/bash

make build-pug

# Monitor for changes in the directories
inotifywait -m -r -e close_write,create,delete "../client" |
while read -r directory event file
do
  path="$directory$file"
  if [[ $path == *".pug" ]]; then
    make build-pug
  fi
done
