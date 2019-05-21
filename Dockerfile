#源镜像
FROM golang:1.12.4-alpine

#作者
MAINTAINER Pharbers "pqian@pharbers.com"

#LABEL
LABEL NtmPods.version="0.0.1" maintainer="Alex"

# 安装git
RUN apk add --no-cache git gcc musl-dev

# 下载依赖
RUN git clone https://github.com/go-yaml/yaml $GOPATH/src/gopkg.in/yaml.v2 && \
    cd $GOPATH/src/gopkg.in/yaml.v2 && git checkout tags/v2.2.2 && \
    git clone https://github.com/go-mgo/mgo $GOPATH/src/gopkg.in/mgo.v2 && \
    cd $GOPATH/src/gopkg.in/mgo.v2 && git checkout -b v2 && \
    git clone https://github.com/PharbersDeveloper/PharbersWebPods.git $GOPATH/src/github.com/PharbersDeveloper/PharbersWebPods

#ADD  src/   $GOPATH/src/

# 设置工程配置文件的环境变量
ENV PHARBERS_WEB_HOME $GOPATH/src/github.com/PharbersDeveloper/PharbersWebPods/resources
ENV GO111MODULE on

# 构建可执行文件
RUN cd $GOPATH/src/github.com/PharbersDeveloper/PharbersWebPods && \
    go build && go install

# 暴露端口
EXPOSE 35415

# 设置工作目录
WORKDIR $GOPATH/bin

ENTRYPOINT ["Web"]
