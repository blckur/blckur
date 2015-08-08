FROM pritunl/archlinux:2015-04-18
MAINTAINER Blckur <contact@blckur.com>

RUN pacman -S --noconfirm go git bzr

ENV GOPATH /go
ENV PATH $PATH:/go/bin

RUN go get github.com/blckur/blckur # 2ee01eb0e902e1ba9467fa0937d5f3d1f58c3d1b

CMD ["blckur", "app"]
