FROM golang:1.22.0-alpine as build

RUN mkdir /blogging_platform

ADD . /blogging_platform

WORKDIR /blogging_platform

RUN go build -o blogging_platform ./cmd

FROM alpine:latest
COPY --from=build /blogging_platform /blogging_platform

WORKDIR /blogging_platform

CMD ["/blogging_platform/blogging_platform"]