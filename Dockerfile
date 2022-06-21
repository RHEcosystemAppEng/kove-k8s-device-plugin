FROM scratch
LABEL maintainer="Red Hat Inc."
ARG GOARCH
COPY bin/amd64/kove-k8s-device-plugin /kove-k8s-device-plugin
ENTRYPOINT ["/kove-k8s-device-plugin"]
