FROM rithium/smartstackgo

WORKDIR $GOPATH/src/

COPY . $GOPATH/src/

RUN apk update && apk upgrade

RUN apk add git
RUN apk add make

RUN git config --global http.https://gopkg.in.followRedirects true

RUN make xbuild

RUN rm -rf /var/cache/apk/*

ADD nerve.conf.json /etc/nerve/
ADD nerve.conf.json /etc/

ADD synapse.conf.json /etc/synapse/
ADD synapse.conf.json /etc/

RUN chmod +x run.sh
RUN chmod +x /opt/startNerve.sh
RUN chmod +x /opt/startSynapse.sh

ENTRYPOINT /go/src/run.sh