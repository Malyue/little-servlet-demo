FROM golang:alpine AS builder

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64\
    GOPROXY=https://goproxy.cn,direct

# 移动到工作目录：/build
WORKDIR /build

# 复制项目中的 go.mod 和 go.sum文件并下载依赖信息
COPY go.mod .
COPY go.sum .
RUN go mod download

# 将代码复制到容器中
COPY . .

# 将我们的代码编译成二进制可执行文件 main
RUN go build -o main .

ENTRYPOINT ["main", ".confiugre.json"]

###################
# 接下来创建一个小镜像
###################
#FROM scratch
#
##COPY ./conf /conf
#COPY ./config.yaml /
#COPY ./WW_verify_CSaj5WtuBq3AemhM.txt /
#ADD ca-bundle.crt /etc/ssl/certs/
#
## 从builder镜像中把/build/main 拷贝到当前目录
#COPY --from=builder /build/main /
#
## 需要运行的命令
#ENTRYPOINT ["/main", "/config.yaml"]