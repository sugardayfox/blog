FROM alpine:3.4

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY config.json /app/
COPY blog_app_linux /app/blog
COPY scripts/entrypoint.sh /app/entrypoint.sh
COPY www /app/www/
COPY tmpl /app/tmpl/
COPY blog_posts /app/blog_posts/

EXPOSE 80

CMD ["./entrypoint.sh"]
