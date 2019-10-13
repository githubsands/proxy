# build image 
FROM alpine:latest as builder
LABEL author="Ryan"
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh ca-certificates
RUN mkdir -p etc/proxy
VOLUME /etc/proxy

# install go
FROM golang:1.13 
ADD . /go/src/github.com/githubsands/proxy/
WORKDIR /go/src/github.com/githubsands/proxy/
RUN GO111MODULE=on go build -o proxy_linux
RUN ls && pwd
#ENV PATH=$PATH:/usr/local/go/bin
#ENV GOPATH=/go 

# RUN GO111MODULE=on go build -o proxy_linux

# final image
FROM alpine:latest

# add the binary file generated in the builder 
RUN cd usr/local && ls && ls
COPY --from=builder /go/src/github.com/githubsands/proxy/proxy_linux /usr/local/bin/proxy_linux

# create nonroot user proxy
ENV USER_GROUP=proxy
RUN groupadd -r ${USER_GROUP} && \ 
    useradd --no-create-home -g ${USER_GROUP} ${USER_GROUP}
USER ${USER_GROUP}

#TODO: make this configurable
EXPOSE 80

# copies from just the built artifact from the previous stage into this new stage
COPY config.json / 
COPY deploy/Dockerfile /
ENTRYPOINT ["/proxy_linux_amd64"]
