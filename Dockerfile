FROM golang:latest

ENV CLUSTER_NAME ""
ENV ADDRESS 0.0.0.0:8888

RUN mkdir /app

COPY . /app

RUN cd /app && go install

CMD /go/bin/configsvc server -a ${ADDRESS}
