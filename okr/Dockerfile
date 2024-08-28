FROM golang:1.17-alpine3.13 AS builder
WORKDIR /usr/src/app
ENV GO111MODULE=on
#RUN adduser -u 10001 -D app-runner

ENV GOPROXY https://goproxy.cn
COPY ./pengcainiao  .
RUN go env
#RUN go mod download
COPY internal/v1/services .
#CMD ["tail", "-f", "/dev/null"]
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -o server ./okr/cmd
FROM alpine:3.12 AS final

WORKDIR /app
#COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/src/app/server /opt/app/
#USER app-runner
#CMD ["/opt/app/server"]