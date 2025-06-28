FROM caddy:2.4.6-alpine
RUN apk add --no-cache curl
COPY Caddyfile /etc/caddy/Caddyfile
