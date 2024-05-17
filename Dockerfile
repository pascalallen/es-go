FROM --platform=linux/arm64 golang:1.22

LABEL org.opencontainers.image.source=https://github.com/pascalallen/es-go
LABEL org.opencontainers.image.description="Container image for es-go"
LABEL org.opencontainers.image.licenses=MIT

WORKDIR /app

ADD . /app

COPY scripts/wait-for-it.sh /usr/bin/wait-for-it.sh

RUN chmod +x /usr/bin/wait-for-it.sh

ENV GOCACHE=/root/.cache/go-build

RUN --mount=type=cache,target="/root/.cache/go-build" CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -C cmd/es-go -o /es-go

CMD /usr/bin/wait-for-it.sh eventstore:$EVENTSTORE_HTTP_PORT \
    && /usr/bin/wait-for-it.sh $RABBITMQ_HOST:$RABBITMQ_PORT \
    && /es-go
