FROM golang:1.9.2

RUN mkdir -p /app/configs

WORKDIR /app

ADD ./finch .
ADD ./finch.yml ./configs/

ENV CONSUL_HTTP_ADDR consul:8500

CMD ["./finch" , "--config=configs/finch.yml"]