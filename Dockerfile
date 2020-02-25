FROM golang:alpine

WORKDIR $GOPATH/src/modbus-server
COPY . $GOPATH/src/modbus-server
RUN go build main.go

EXPOSE 5020
ENTRYPOINT ["./main"]