#!/bin/bash

docker build -t animated-robot-build -f Dockerfile.release .

docker run --name animated-robot-build-release --rm -it -v $(pwd):/build animated-robot-build cp /go/src/animated-robot/cmd/app /build

docker rmi animated-robot-build

