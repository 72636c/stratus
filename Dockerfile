FROM golang:1.12-alpine AS builder

RUN \
  echo 'nobody:x:65534:65534:nobody:/:' > /tmp/passwd && \
  echo 'nobody:x:65534:' > /tmp/group

RUN apk add --no-cache git

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY internal/ ./internal/
COPY stratus.go ./

RUN CGO_ENABLED=0 go build -installsuffix 'static' -o /app .

FROM gcr.io/distroless/static:latest AS final-static

COPY --from=builder /tmp/group /tmp/passwd /etc/

COPY --from=builder /app /bin/stratus

USER nobody:nobody

ENTRYPOINT ["/bin/stratus"]
CMD ["--help"]

FROM gcr.io/distroless/base:debug AS final-busybox

COPY --from=builder /tmp/group /tmp/passwd /etc/

COPY --from=builder /app /bin/stratus

USER nobody:nobody

ENTRYPOINT ["/busybox/sh"]
