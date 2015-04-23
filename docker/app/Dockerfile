FROM pritunl/archlinux:2015-04-18
MAINTAINER Blckur <contact@blckur.com>

RUN pacman -S --noconfirm go git bzr

ENV GOPATH /go
ENV PATH $PATH:/go/bin

RUN go get github.com/blckur/blckur # 3cf556542b46519f4a52aaff35d4bed21790eb3e26de392d72e0f6fe44af6b73ca184c9830bf16a26de392d72e0f6fe44af6b73ca184c9830bf16a

CMD ["blckur", "app"]
