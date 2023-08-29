FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY controller /controller
USER 65532:65532
ENTRYPOINT ["/controller"]
