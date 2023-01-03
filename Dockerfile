FROM golang:1.19 as build-env

ENV GO111MODULE on

WORKDIR /go/src/github.com/bfu4/mipscalls
ADD . /go/src/github.com/bfu4/mipscalls

RUN go get -d -v ./...

RUN make
RUN mv ./.env /go/bin/.env
RUN mv ./mipscalls /go/bin/mipscalls
FROM gcr.io/distroless/base
COPY --from=build-env /go/bin/mipscalls /
COPY --from=build-env /go/bin/.env /
CMD ["/mipscalls"]