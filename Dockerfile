ARG GOARCH=amd64
ARG GOOS=linux
FROM --platform=${GOARCH} golang:1.19 as builder
WORKDIR /workspace

COPY . /workspace/

RUN go mod download
RUN CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -a -o controller main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder --chown=65532:65532 /workspace/controller /controller
USER 65532:65532
ENTRYPOINT ["/controller"]