FROM alpine:3.6

RUN wget -P /usr/bin http://storage.googleapis.com/kubernetes-release/release/v1.10.7/bin/linux/amd64/kubectl
RUN chmod 755 /usr/bin/kubectl

RUN mkdir /etc/kubevirt
COPY kubevirt.yaml /etc/kubevirt/

RUN adduser -D kubevirt-operator
RUN chown -R kubevirt-operator: /etc/kubevirt
USER kubevirt-operator

ADD tmp/_output/bin/kubevirt-operator /usr/local/bin/kubevirt-operator
