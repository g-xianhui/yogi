#!/bin/sh

# server
protoc --go_out=../server/src/msg/ *.proto

# client
protoc --go_out=../client/msg/ *.proto
