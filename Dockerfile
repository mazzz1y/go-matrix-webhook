FROM golang:1.20-alpine as build

ENV CGO_ENABLED 0
ARG VERSION

COPY . /src
RUN cd /src && \
  go build -ldflags="-s -w -X main.version=$VERSION" -trimpath -o /go-matrix-webhook ./cmd/go-matrix-webhook

FROM alpine:3.17 
COPY --from=build /go-matrix-webhook /
USER 1337
CMD /go-matrix-webhook
