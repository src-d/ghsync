FROM alpine

RUN apk update && \
  apk add ca-certificates && \
  rm -rf /var/cache/apk/*

COPY ./build/bin/ghsync /bin/ghsync

ENTRYPOINT ["/bin/ghsync"]
