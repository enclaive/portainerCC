FROM alpine:latest

COPY dist /app
COPY coordinator /

VOLUME /data
WORKDIR /
RUN ls /

EXPOSE 9000
EXPOSE 9443
EXPOSE 8000

ENTRYPOINT ["/app/portainer"]
