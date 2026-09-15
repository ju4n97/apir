FROM gcr.io/distroless/static-debian12:nonroot
ARG TARGETPLATFORM
COPY ${TARGETPLATFORM}/esquema /esquema
USER nonroot:nonroot
ENTRYPOINT ["/esquema"]