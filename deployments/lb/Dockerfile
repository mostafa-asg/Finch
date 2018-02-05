FROM nginx:latest
ENV CONSUL_HTTP_ADDR consul:8500
RUN apt-get update && apt-get install -y unzip wget
RUN mkdir -p /opt/ct/consul-template
WORKDIR /opt/ct/consul-template
RUN wget https://releases.hashicorp.com/consul-template/0.18.2/consul-template_0.19.4_linux_amd64.zip && unzip consul-template_0.18.2_linux_amd64.zip
ADD default.ctmpl /opt/ct/consul-template/
ADD startup.sh /opt/ct/consul-template/
RUN chmod +777 /opt/ct/consul-template/startup.sh
RUN rm /etc/nginx/conf.d/default.conf
EXPOSE 80
ENTRYPOINT ["/opt/ct/consul-template/startup.sh"]