FROM alpine:3.6

RUN adduser -D kubevirt-operator
USER kubevirt-operator

ADD tmp/_output/bin/kubevirt-operator /usr/local/bin/kubevirt-operator
