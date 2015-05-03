FROM pritunl/archlinux:2015-04-18
MAINTAINER Blckur <contact@blckur.com>

RUN pacman -S --noconfirm go git bzr redis

ENV GOPATH /go
ENV PATH $PATH:/go/bin

RUN go get github.com/blckur/blckur # fcc3d35f8497c2a64ff8e0ea0e6c5ff7c39ca247

CMD ["blckur", "cache"]
