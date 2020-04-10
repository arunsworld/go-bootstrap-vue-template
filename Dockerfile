FROM golang:1.13.4-alpine AS buildenv

COPY . /tmp/app/
WORKDIR /tmp/app/cmd/

ARG BUILD_COMMAND='env GOOS=linux GOARCH=amd64 go build -mod vendor'

ARG GIT_COMMIT="unknown"

RUN mkdir bin && \
    cd app && \
    ${BUILD_COMMAND} -ldflags "-s -w -X main.GitCommit=$GIT_COMMIT" . \
    && mv app ../bin

# truepay
FROM alpine:latest

COPY docker-entrypoint.sh cmd/app/pwd.csv /bin/
RUN chmod 755 /bin/docker-entrypoint.sh

COPY --from=buildenv /tmp/app/cmd/bin/ /bin/
ADD web /bin/web/

WORKDIR /bin

ENTRYPOINT ["docker-entrypoint.sh"]