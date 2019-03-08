FROM golang:alpine AS builder

# install packages
RUN apk --no-cache add ca-certificates git

# install dep manager
ADD https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

WORKDIR $GOPATH/src/github.com/BlitzinBuffalo/liveness-prober

COPY . .
RUN dep ensure --vendor-only \
    && CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo -o /liveness-prober .

FROM scratch
LABEL Author="Michael K. Essandoh <mexcon.mike@gmail.com>"

ENV PROBE_INTERVAL=1000
ENV PROBE_HOSTS="www.google.com"

COPY --from=builder /liveness-prober .

ENTRYPOINT ["./liveness-prober"]
