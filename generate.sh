#!/bin/bash

protoc --go_out=plugins=grpc:. pbs/*.proto
# protoc pbs/jarvis/*.proto --go_out=plugins=grpc:.