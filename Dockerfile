FROM --platform=$BUILDPLATFORM golang:1.26.1 AS build

ARG VERSION
ARG REVISION
ARG BUILDTIME
ARG BUILDER

ARG TARGETOS
ARG TARGETARCH

ENV CGO_ENABLED=0

WORKDIR /src
COPY . .

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -ldflags "-X github.com/vint-go/vint/cli.version=${VERSION} -X github.com/vint-go/vint/cli.commit=${REVISION} -X github.com/vint-go/vint/cli.date=${BUILDTIME} -X github.com/vint-go/vint/cli.builtBy=${BUILDER}"

FROM scratch

COPY --from=build /src/vint /vint

ENTRYPOINT ["/vint"]
