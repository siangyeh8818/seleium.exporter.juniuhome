FROM golang:1.12.7-stretch as builder

COPY go.mod /go/src/github.com/siangyeh8818/seleium.exporter.juniuhome/go.mod
COPY go.sum /go/src/github.com/siangyeh8818/seleium.exporter.juniuhome/go.sum

# Run golang at any directory, not neccessary $GOROOT, $GOPATH
ENV GO111MODULE=on
WORKDIR /go/src/github.com/siangyeh8818/seleium.exporter.juniuhome

# RUN go mod init github.com/pnetwork/sre.monitor.metrics
RUN go mod download
COPY cmd /go/src/github.com/siangyeh8818/seleium.exporter.juniuhome/cmd
COPY internal /go/src/github.com/siangyeh8818/seleium.exporter.juniuhome/internal
#COPY pkg /go/src/github.com/pnetwork/sre.monitor.metrics/pkg

# Build the Go app
RUN env GOOS=linux GOARCH=amd64 go build -o juniuhome-exporter -v -ldflags "-s" github.com/siangyeh8818/seleium.exporter.juniuhome/cmd/

##### To reduce the final image size, start a new stage with alpine from scratch #####
FROM centos:7

COPY google.repo /etc/yum.repos.d/
RUN yum update & yum install -y google-chrome-stable git epel-release

# Run as root
WORKDIR /root/

COPY chromedriver /usr/local/bin/
RUN chmod +x /usr/local/bin/chromedriver
# Copy the pre-built binary file from the previous stage
COPY --from=builder /go/src/github.com/siangyeh8818/seleium.exporter.juniuhome/juniuhome-exporter /usr/local/bin/juniuhome-exporter

# EXPOSE 8081

ENTRYPOINT [ "juniuhome-exporter" ] 
