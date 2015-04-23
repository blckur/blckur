FROM pritunl/archlinux:2015-04-18
MAINTAINER Blckur <contact@blckur.com>

RUN pacman -S --noconfirm go git bzr wget
RUN wget https://s3.amazonaws.com/blckur/beanstalkd-1.10-1-x86_64.pkg.tar.xz
RUN pacman -U --noconfirm beanstalkd-1.10-1-x86_64.pkg.tar.xz

ENV GOPATH /go
ENV PATH $PATH:/go/bin

RUN go get github.com/blckur/blckur # 3cf556542b46519f4a52aaff35d4bed21790eb3e26de392d72e0f6fe44af6b73ca184c9830bf16a26de392d72e0f6fe44af6b73ca184c9830bf16a

CMD ["blckur", "queue"]
