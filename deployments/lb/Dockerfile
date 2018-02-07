FROM nginx:latest
ENV CONSUL_HTTP_ADDR consul:8500
RUN apt-get update && apt-get install -y unzip wget
RUN mkdir -p /opt/ct
WORKDIR /opt/ct
ADD consul-template_0.19.4_linux_amd64.zip .
RUN unzip consul-template_0.19.4_linux_amd64.zip
ADD default.ctmpl /opt/ct
ADD startup.sh /opt/ct
RUN chmod +777 /opt/ct/startup.sh
RUN rm /etc/nginx/conf.d/default.conf
EXPOSE 80
ENTRYPOINT ["/opt/ct/startup.sh"]