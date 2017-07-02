FROM golang:1.8-alpine as base

WORKDIR /app/prebid-server/

COPY go.mod go.sum ./
RUN go mod download

COPY rest ./rest
COPY storage ./storage
COPY config.go ./config.go
COPY main.go ./main.go

RUN go mod vendor

FROM base as minimal-dropbox-builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app/prebid-server
COPY --from=base /app/prebid-server /app/prebid-server/
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH GOPROXY="https://proxy.golang.org" go build -o minimal-dropbox-build -mod=vendor ./cmd/minimal-dropbox

FROM scratch as minimal-dropbox

WORKDIR /usr/local/bin
COPY --from=minimal-dropbox-builder /app/prebid-server/minimal-dropbox-build ./minimal-dropbox

ENTRYPOINT ["./minimal-dropbox"]
