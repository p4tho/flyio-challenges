#!/bin/sh

set -e

go build
../maelstrom/maelstrom test -w broadcast --bin 3d --node-count 25 --time-limit 20 --rate 100 --latency 100
rm -r store
