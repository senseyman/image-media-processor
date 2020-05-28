FROM golang:1.14-alpine AS build

WORKDIR /src/
COPY . /src/
RUN CGO_ENABLED=0 go build -o /bin/image-media-processor

FROM alpine
RUN apk update && apk add ca-certificates
COPY --from=build /bin/image-media-processor /image-media-processor
EXPOSE 8080
ENTRYPOINT ["/image-media-processor"]