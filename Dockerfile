FROM alpine:3.12.0 AS builder

# Update/Upgrade/Add packages for building

RUN apk add --no-cache bash git go build-base

# Build scalingo_deployer

WORKDIR /build/scalingo_deployer

ADD . .

ENV GOPATH=/build/scalingo_deployer/gospace

RUN make clobber setup all

FROM alpine:3.12.0 AS runner

# Update/Upgrade/Add packages

RUN apk add --no-cache bash ca-certificates tzdata

RUN cp /usr/share/zoneinfo/Europe/Berlin /etc/localtime && \
  echo Europe/Berlin >/etc/timezone

WORKDIR /

COPY --from=builder /build/scalingo_deployer/scalingo_deployer /scalingo_deployer

ENTRYPOINT [ "/scalingo_deployer" ]
