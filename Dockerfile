FROM golang:1.16 as build-env

ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /go/src/heatilly
ADD . /go/src/heatilly

RUN go build -o /go/bin/heatilly

FROM gcr.io/distroless/base
COPY --from=build-env /go/bin/heatilly /
EXPOSE 8000
CMD ["/heatilly"]
