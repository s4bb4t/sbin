FROM alpine:3.18.3

WORKDIR /

COPY bin/app /

EXPOSE 8011

ENTRYPOINT ["/app"]
