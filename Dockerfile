FROM registry.fedoraproject.org/fedora-minimal
#FROM centos:7
MAINTAINER Matthias Wessendorf <matzew@apache.org>
ARG BINARY=./kes

COPY ${BINARY} /opt/kes
ENTRYPOINT ["/opt/kes"]
