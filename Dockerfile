FROM alpine:3.18.3

WORKDIR /

CMD go build -o bin/app cmd/main.go
COPY bin/app /

EXPOSE 8011

ENTRYPOINT ["/app"]
