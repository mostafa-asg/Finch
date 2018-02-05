#!/bin/bash

service nginx start && /opt/ct/consul-template/consul-template -consul=$CONSUL_HTTP_ADDR -template="default.ctmpl:/etc/nginx/conf.d/default.conf:service nginx reload"