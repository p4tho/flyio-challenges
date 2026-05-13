#!/bin/sh

set -e

go build
../maelstrom/maelstrom test -w broadcast --bin 3a --node-count 1 --time-limit 20 --rate 10
rm -r store
