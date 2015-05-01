FROM pritunl/archlinux:2015-04-18
MAINTAINER Blckur <contact@blckur.com>

RUN pacman -S --noconfirm go git bzr wget
RUN wget https://s3.amazonaws.com/blckur/beanstalkd-1.10-1-x86_64.pkg.tar.xz
RUN pacman -U --noconfirm beanstalkd-1.10-1-x86_64.pkg.tar.xz

ENV GOPATH /go
ENV PATH $PATH:/go/bin

RUN go get github.com/blckur/blckur # 262b098c3b9b99fb16c4ff9e6ab33f86cffadd9e

CMD ["blckur", "queue"]
