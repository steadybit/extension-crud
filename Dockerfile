# syntax=docker/dockerfile:1

##
## Build
##
FROM --platform=$BUILDPLATFORM golang:1.20-alpine AS build

ARG TARGETOS TARGETARCH
ARG NAME
ARG VERSION
ARG REVISION

WORKDIR /app

RUN apk add build-base
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build \
    -ldflags="\
    -X 'github.com/steadybit/extension-kit/extbuild.ExtensionName=${NAME}' \
    -X 'github.com/steadybit/extension-kit/extbuild.Version=${VERSION}' \
    -X 'github.com/steadybit/extension-kit/extbuild.Revision=${REVISION}'" \
    -o ./extension
RUN make licenses-report

##
## Runtime
##
FROM alpine:3.18

LABEL "steadybit.com.discovery-disabled"="true"

ARG USERNAME=steadybit
ARG USER_UID=1000

RUN adduser -u $USER_UID -D $USERNAME

USER $USERNAME

WORKDIR /

COPY --from=build /app/extension /extension
COPY --from=build /app/licenses /licenses

EXPOSE 8080

ENTRYPOINT ["/extension"]
