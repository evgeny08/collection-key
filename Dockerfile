# build binary
FROM golang:1.11-alpine3.8 AS build
WORKDIR /go/src/github.com/evgeny08/collection-key
COPY . /go/src/github.com/evgeny08/collection-key
RUN CGO_ENABLED=0 go build -o /out/collection-key github.com/evgeny08/collection-key/cmd/collection-key-d

# copy to alpine image
FROM alpine:3.8
WORKDIR /app
COPY --from=build /out/collection-key /app
CMD ["/app/collection-key"]
