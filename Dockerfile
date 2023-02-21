FROM registry.entrade.com.vn/dockerhub/library/alpine:3

RUN apk add --no-cache tzdata

COPY .go/bin/runner /bin/runner
COPY config.yaml /app/config.yaml
COPY migrations /app/migrations

WORKDIR /app
EXPOSE 8080

ENTRYPOINT ["/bin/runner"]
