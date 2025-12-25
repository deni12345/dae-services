FROM node:20-alpine

RUN apk add --no-cache openjdk21-jre-headless bash curl \
  && npm i -g firebase-tools@13


# Thư mục persist dữ liệu emulator (dùng với --import/--export-on-exit)
VOLUME ["/data"]

# Set ownership
RUN mkdir -p /data && chown -R node:node /data

# Switch to non-root user
USER node
WORKDIR /data

# Mặc định mở cổng Firestore (8080) và UI (4000) nếu bạn bật UI trong firebase.json
EXPOSE 8080 4000

# Mặc định chạy firebase CLI; flags sẽ truyền ở docker run / docker-compose
ENTRYPOINT ["firebase"]  
CMD ["--version"]