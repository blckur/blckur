FROM pritunl/archlinux:2015-04-18
MAINTAINER Blckur <contact@blckur.com>

RUN pacman -S --noconfirm go git bzr wget
RUN wget https://s3.amazonaws.com/blckur/beanstalkd-1.10-1-x86_64.pkg.tar.xz
RUN pacman -U --noconfirm beanstalkd-1.10-1-x86_64.pkg.tar.xz

ENV GOPATH /go
ENV PATH $PATH:/go/bin

RUN go get github.com/blckur/blckur # 171acc07ff6c85ac6e0ee32118ad628b0e3fb100

CMD ["blckur", "queue"]
