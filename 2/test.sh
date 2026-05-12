#!/bin/sh

set -e

go build
../maelstrom/maelstrom test -w unique-ids --bin 2 --time-limit 30 --rate 1000 --node-count 3 --availability total --nemesis partition
rm -r store
