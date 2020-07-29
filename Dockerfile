FROM debian:stretch-slim

WORKDIR /

COPY _output/bin/scheduler-extender-demo /usr/local/bin

CMD ["scheduler-extender-demo"]