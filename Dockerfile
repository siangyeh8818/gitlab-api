FROM golang:1.16.4-stretch as builder

COPY go.mod /go/src/gitlab.com/st-opger/gitlab-implement/go.mod
COPY go.sum /go/src/gitlab.com/st-opger/gitlab-implement/go.sum

# Run golang at any directory, not neccessary $GOROOT, $GOPATH
ENV GO111MODULE=on
WORKDIR /go/src/gitlab.com/st-opger/gitlab-implement

# RUN go mod init github.com/pnetwork/sre.monitor.metrics
RUN go mod download
COPY pkg /go/src/gitlab.com/st-opger/gitlab-implement/pkg
COPY internal /go/src/gitlab.com/st-opger/gitlab-implement/internal
COPY main.go /go/src/gitlab.com/st-opger/gitlab-implement

# Build the Go app
RUN env GOOS=linux GOARCH=amd64 go build -o gitlab-implement -v -ldflags "-s" gitlab.com/st-opger/gitlab-implement

##### To reduce the final image size, start a new stage with alpine from scratch #####
FROM alpine:3.13
RUN apk --no-cache add ca-certificates libc6-compat wget
#RUN wget https://github.com/etcd-io/etcd/releases/download/v3.4.5/etcd-v3.4.5-linux-amd64.tar.gz \
#    && tar zxvf etcd-v3.4.5-linux-amd64.tar.gz \
#    && mv etcd-v3.4.5-linux-amd64/etcdctl /usr/bin \
#    && chmod +x /usr/bin/etcdctl
# Run as root
WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /go/src/gitlab.com/st-opger/gitlab-implement/gitlab-implement /usr/local/bin/gitlab-implement
COPY conf/app.ini /root/conf

# EXPOSE 8081

ENTRYPOINT [ "gitlab-implement" ] 