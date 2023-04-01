FROM golang:1.19-alpine AS builder
RUN mkdir /build
ADD . /build
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main ./cmd

FROM scratch
COPY --from=builder /build/main /app/
COPY ./config/config.yaml /app/config/config.yaml
WORKDIR /app
ENTRYPOINT ["./main"]
EXPOSE 8080