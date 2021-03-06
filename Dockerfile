FROM alpine:3.6

RUN apk update && \
    apk add curl git

RUN wget -P /usr/bin http://storage.googleapis.com/kubernetes-release/release/v1.10.7/bin/linux/amd64/kubectl
RUN chmod 755 /usr/bin/kubectl

RUN mkdir /etc/kubevirt
COPY download-templates /bin
RUN download-templates /etc/kubevirt

RUN adduser -D kubevirt-operator
RUN chown -R kubevirt-operator: /etc/kubevirt
USER kubevirt-operator

ADD kubevirt-operator /usr/local/bin/kubevirt-operator
