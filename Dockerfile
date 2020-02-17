FROM golang:buster AS build-env
WORKDIR /go/src
COPY . animated-robot

RUN go get -v animated-robot/...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app animated-robot/cmd

FROM scratch
WORKDIR /app
COPY --from=build-env /go/src/app .

CMD ["/app/app"]
