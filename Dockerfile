FROM golang:alpine AS build-env
MAINTAINER "Rajat Gupta <rajatgpt152@gmail.com>"
WORKDIR /cachingSystem
COPY . /cachingSystem
RUN apk update && apk upgrade && apk add git
RUN cd /cachingSystem && go build -o cachesystem . && ls

FROM alpine
RUN mkdir -p /app/service/configs
WORKDIR /app
COPY --from=build-env /cachingSystem/cachesystem /app/cachesystem
COPY --from=build-env /cachingSystem/service/configs/db_config.yaml /app/cachesystem
EXPOSE 8000
ENTRYPOINT [ "./cachesystem" ]