FROM pritunl/archlinux:2015-04-18
MAINTAINER Blckur <contact@blckur.com>

RUN pacman -S --noconfirm go git bzr redis

ENV GOPATH /go
ENV PATH $PATH:/go/bin

RUN go get github.com/blckur/blckur # 0e5366fac26b4d7513a39a9def06c259f66c8f3d

CMD ["blckur", "cache"]
