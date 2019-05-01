FROM registry.access.redhat.com/ubi7-dev-preview/ubi-minimal:7.6@sha256:b50ba3aac799c34c9316442d4541f9f9682d7860cf9290c15dfcc459c68c4c3f

ENV OPERATOR=/usr/local/bin/cdays-namespace-poc \
    USER_UID=1001 \
    USER_NAME=cdays-namespace-poc

# install operator binary
COPY build/_output/bin/cdays-namespace-poc ${OPERATOR}

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
