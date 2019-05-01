FROM registry.access.redhat.com/ubi7-dev-preview/ubi-minimal:7.6@sha256:a572439e8d7ddced4565e2703c940b13a7a8087807db4b6fbbb9889dbfb731e8

ENV OPERATOR=/usr/local/bin/cdays-namespace-poc \
    USER_UID=1001 \
    USER_NAME=cdays-namespace-poc

# install operator binary
COPY build/_output/bin/cdays-namespace-poc ${OPERATOR}

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
