# Stage 1 - build executable in go container
FROM golang:latest as builder

WORKDIR $GOPATH/src/wbL0/
COPY . .

RUN export CGO_ENABLED=0 && make build

# Stage 2 - build final image
FROM alpine:latest

# Copy our static executable
COPY --from=builder /go/src/wbL0/bin/wbl0 go/bin/wbl0

# Run the binary.
ENTRYPOINT ["go/bin/wbl0"]
