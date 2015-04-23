FROM pritunl/archlinux:2015-04-18
MAINTAINER Blckur <contact@blckur.com>

RUN pacman -S --noconfirm go git bzr

ENV GOPATH /go
ENV PATH $PATH:/go/bin

RUN go get github.com/blckur/blckur # 526de392d72e0f6fe44af6b73ca184c9830bf16a26de392d72e0f6fe44af6b73ca184c9830bf16a

CMD ["blckur", "worker"]
