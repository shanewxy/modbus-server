FROM golang:alpine as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .

FROM golang:alpine
RUN apk update && apk add socat
COPY --from=builder /build/main /app/
WORKDIR /app
EXPOSE 5020
CMD ["./main"]
