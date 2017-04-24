FROM alpine:3.4

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY config.json /app/
COPY blog_linux /app/blog
COPY www /app/www/
COPY tmpl /app/tmpl/
COPY blog_posts /app/blog_posts/

EXPOSE 80 443

CMD ["/app/blog", "-production", "-addr=:80"]
