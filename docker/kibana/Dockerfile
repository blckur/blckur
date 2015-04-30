FROM pritunl/archlinux:2015-04-25
MAINTAINER Blckur <contact@blckur.com>

RUN pacman -S --noconfirm wget

RUN wget https://download.elastic.co/kibana/kibana/kibana-4.0.2-linux-x64.tar.gz
RUN tar xfz kibana-4.0.2-linux-x64.tar.gz
WORKDIR kibana-4.0.2-linux-x64

EXPOSE 5601

CMD bin/kibana -e $ELASTIC
