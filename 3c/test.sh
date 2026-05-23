#!/bin/sh

set -e

go build
../maelstrom/maelstrom test -w broadcast --bin 3c --node-count 5 --time-limit 20 --rate 10 --nemesis partition
rm -r store
