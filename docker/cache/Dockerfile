FROM pritunl/archlinux:2015-04-18
MAINTAINER Blckur <contact@blckur.com>

RUN pacman -S --noconfirm go git bzr redis

ENV GOPATH /go
ENV PATH $PATH:/go/bin

RUN go get github.com/blckur/blckur # 247258264ced2c8234547fdc5a0ce7c38635f0bb

CMD ["blckur", "cache"]
