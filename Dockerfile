FROM golang:1.6.2

#因为构建下载资源比较慢,此文件建立一个基础镜像
#docker build -t gobase .
MAINTAINER dxwsker@qq.com

RUN apt-get update \
    && apt-get install -y libldap2-dev \
    && rm -r /var/lib/apt/lists/* \
    && go get -d github.com/docker/distribution \
    && go get -d github.com/docker/libtrust \
    && go get -d github.com/go-sql-driver/mysql