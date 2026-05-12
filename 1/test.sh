#!/bin/sh

set -e

go build
../maelstrom/maelstrom test -w echo --bin 1 --node-count 1 --time-limit 10
rm -r store
