# 构建环境
FROM golang:1.23.0 as builder
WORKDIR /build
ENV GOPROXY https://goproxy.cn
COPY . .
RUN CGO_ENABLED=0  GOOS=linux  GOARCH=amd64 go build -o lpdrive .

# 运行环境
FROM alpine:3.20.3 as runner
WORKDIR /app
COPY --from=builder /build/lpdrive /app/
RUN chmod u+x /app/lpdrive
EXPOSE 8080
CMD ["/app/lpdrive"]
