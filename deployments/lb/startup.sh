#!/bin/bash

service nginx start && 
/opt/ct/consul-template \
        -consul-addr=$CONSUL_HTTP_ADDR \
        -template="/opt/ct/default.ctmpl:/etc/nginx/conf.d/default.conf:service nginx reload"