FROM golang:1.20-alpine  as builder
ENV GOPROXY=https://goproxy.cn
#安装编译需要的环境gcc等
RUN apk add build-base git

WORKDIR /code
#将上层整个文件夹拷贝到/go/release
COPY . /code/src
WORKDIR /code/src
#交叉编译，需要制定CGO_ENABLED=1，默认是关闭的
RUN make linux  

#编译
FROM alpine
RUN apk --no-cache add tzdata ca-certificates libc6-compat libgcc libstdc++
COPY --from=builder  /code/src/dist/linux_amd64/chatserver-api  /app/chatserver-api/chatserver-api

WORKDIR /app/chatserver-api
EXPOSE 18080
VOLUME ["/app/chatserver-api/configs","/app/chatserver-api/logs","/app/chatserver-api/head_photo","/app/chatserver-api/uploadfile"]
ENTRYPOINT  ["/app/chatserver-api/chatserver-api"]