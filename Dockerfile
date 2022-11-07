FROM --platform=${BUILDPLATFORM} golang:1.14-alpine AS builder

RUN \
  echo 'nobody:x:65534:65534:nobody:/:' > /tmp/passwd && \
  echo 'nobody:x:65534:' > /tmp/group

RUN apk add --no-cache git

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY internal/ ./internal/
COPY stratus.go ./

ARG TARGETARCH
ARG TARGETOS

RUN CGO_ENABLED=0 GOARCH=${TARGETARCH} GOOS=${TARGETOS} \
  go build -installsuffix 'static' -o /app .

FROM gcr.io/distroless/base:debug AS final-base

RUN ["/busybox/sh", "-c", "ln -s /busybox/sh /bin/sh"]

COPY --from=builder /tmp/group /tmp/passwd /etc/

COPY --from=builder /app /bin/stratus

USER nobody:nobody

ENTRYPOINT ["/bin/sh"]

FROM gcr.io/distroless/static:latest AS final-static

COPY --from=builder /tmp/group /tmp/passwd /etc/

COPY --from=builder /app /bin/stratus

USER nobody:nobody

ENTRYPOINT ["/bin/stratus"]
CMD ["--help"]
