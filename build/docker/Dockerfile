#
# Build binary
#
FROM golang:1.14.9-alpine3.11 AS builder

ARG APP_NAME
ARG VERSION
ARG GIT_COMMIT
ARG BUILD_TIME
ARG GOOS_TYPE

COPY ./ /project

WORKDIR /project/cmd/userservice

ENV GOSUMDB=off
RUN CGO_ENABLED=0 GOOS=${GOOS_TYPE} go build -mod vendor \
    -ldflags "-X main.version=$VERSION -X main.commitHash=$GIT_COMMIT -X main.buildTime=$BUILD_TIME" \
    -a -v -installsuffix cgo -o /${APP_NAME}

#
# Build Backend Docker Image
#
FROM archlinux/base

ARG APP_NAME
ARG VERSION
ARG GIT_COMMIT
ARG BUILD_TIME

LABEL name="${APP_NAME}-backend" buidTime="${BUILD_TIME}" gitCommit="${GIT_COMMIT}" version="${VERSION}"

COPY --from=builder /${APP_NAME} /
RUN chmod +x ./${APP_NAME}

ENTRYPOINT ["./usersvc"]
