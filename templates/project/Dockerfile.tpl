FROM docker.hub.ipao.vip/alpine:3.20

# Set timezone
RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    apk del tzdata

COPY backend/build/app /app/app
COPY backend/config.toml /app/config.toml
COPY frontend/dist /app/dist

WORKDIR /app

ENTRYPOINT  ["/app/app"]

CMD [ "serve" ]
