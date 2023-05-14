FROM golang:1.20-alpine  as builder
ENV GOPROXY=https://goproxy.cn
#安装编译需要的环境gcc等
RUN apk add build-base git

WORKDIR /code

COPY . /code/src
WORKDIR /code/src

RUN make linux  

#编译
FROM alpine
RUN apk --no-cache add tzdata ca-certificates libc6-compat libgcc libstdc++
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo "Asia/Shanghai" > /etc/timezone 
COPY --from=builder  /code/src/dist/linux_amd64/chatserver-api  /app/chatserver-api/chatserver-api
COPY --from=builder  /code/src/dict  /app/chatserver-api/dict

WORKDIR /app/chatserver-api

EXPOSE 18080

VOLUME ["/app/chatserver-api/configs","/app/chatserver-api/logs","/app/chatserver-api/head_photo","/app/chatserver-api/uploadfile"]

ENTRYPOINT  ["/app/chatserver-api/chatserver-api"]