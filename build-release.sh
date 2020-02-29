#!/bin/bash

docker build -t animated-robot-build -f Dockerfile.release .

docker run --rm -it -v $(pwd):/build animated-robot-build cp /go/src/app /build


docker rmi animated-robot-build
