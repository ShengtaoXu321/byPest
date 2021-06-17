# 获取golang简单的一个版本
FROM golang AS build-env

# 设置工作目录
WORKDIR /home/byserver/env

# go proxy代理
ENV GOPROXY https://goproxy.cn
ENV GO111MODULE on


# 将所有代码拷贝至工作目录
COPY .  /home/byserver/env

# 下载需要的库并编译
RUN go mod download && CGO_ENABLED=0 GOOS=linux go build -o byPest

# 第二阶段：拉取运行环境
FROM ubuntu:21.04

WORKDIR /home/byserverbyserver


COPY --from=build-env /home/byserver/env  /home/byserverbyserver
# 暴露端口
EXPOSE 18088

# 最终执行docker的命令
ENTRYPOINT [ "./byPest" ]



