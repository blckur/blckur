FROM pritunl/archlinux:2015-04-18
MAINTAINER Blckur <contact@blckur.com>

RUN pacman -S --noconfirm go git bzr

ENV GOPATH /go
ENV PATH $PATH:/go/bin

RUN go get github.com/blckur/blckur # f25a2cb084cf6f53f8f6c9c18080f69ac7f308b5

CMD ["blckur", "scheduler"]
