FROM alpine:latest

WORKDIR /

COPY socialite socialite

COPY config.yaml config.yaml

EXPOSE 8080

CMD ["./socialite"]

# docker build -t socialite:1 -f build/Dockerfile .