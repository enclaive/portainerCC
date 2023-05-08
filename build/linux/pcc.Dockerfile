FROM sgxdcaprastuff/gramine-os:latest

COPY ./dist /
COPY ./coordinator /coordinator
COPY build/docker-extension /
COPY ./confidential-templates.json /

VOLUME /data
WORKDIR /

EXPOSE 9000
EXPOSE 9443
EXPOSE 8000

ENTRYPOINT ["/portainer"]
