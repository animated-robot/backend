FROM golang:1.14-alpine AS build-env
WORKDIR /go/src/animated-robot
COPY . .
WORKDIR /go/src/animated-robot/cmd
RUN go get -v
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -installsuffix cgo -o app .

FROM scratch
WORKDIR /app
COPY --from=build-env /go/src/animated-robot/cmd/app .

CMD ["/app/app"]
