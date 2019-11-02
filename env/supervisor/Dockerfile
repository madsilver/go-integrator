FROM debian:latest

RUN apt-get update \
    && apt-get install -y supervisor \
    && apt-get clean; rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /usr/share/doc/*

COPY ./config/supervisord.conf /etc/supervisor
COPY ./config/listener.conf /etc/supervisor/conf.d

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/supervisord.conf"]