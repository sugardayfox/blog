# TODO: switch to using scratch
# https://github.com/OpenBazaar/openbazaar-go/blob/98c9ac8ff6ad4b84674134398e890d3dc8912bb6/Dockerfile
FROM alpine:3.4

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY blog_linux /app/blog
COPY www /app/www/
COPY tmpl /app/tmpl/
COPY articles /app/articles/

EXPOSE 80 443

CMD ["/app/blog", "-production"]
