FROM pritunl/archlinux:2015-04-18
MAINTAINER Blckur <contact@blckur.com>

RUN pacman -S --noconfirm go git bzr

ENV GOPATH /go
ENV PATH $PATH:/go/bin

RUN go get github.com/blckur/blckur # 8a355d9d4731eeddfb06dffaeb811edf61ec23ac

CMD ["blckur", "scheduler"]
