FROM jfloff/alpine-python:2.7-slim

LABEL Name "ealexhudson/stowage"
LABEL Description "stowage for Docker cli distribution management" 
LABEL Vendor "Alex Hudson"
LABEL Version "0.1"

WORKDIR /root

ENTRYPOINT ["/usr/sbin/stowage"]

CMD ["-h"]

COPY stowage /usr/sbin/stowage
