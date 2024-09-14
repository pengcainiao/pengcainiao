#!/usr/bin/env groovy
 
def registry = "http://119.29.5.54:1180/"
def project = "library"
def Branch = "main"
def app_name = "okr"
def image_name = "${registry}/${project}/${app_name}:${Branch}-${BUILD_NUMBER}"
def git_address = "https://gitlab.com/a16624741591/server.git"
def docker_registry_auth = "a4ef0fc8-fcad-4859-8a90-1d390a74d7da"
def git_auth = "d090072a-cfaa-4858-b192-197161849d18"

pipeline {
    agent any
        environment {
        _VERSION = sh(script: "echo `date '+%Y%m%d-%H%M%S'`", returnStdout: true).trim()
    }
    stages {
        stage('拉取代码'){
            steps {
                sh 'cd $WORKSPACE'
                sh 'mkdir -p server'
                sh 'cd server'
                checkout([$class: 'GitSCM', branches: [[name: '${SERVER}']], userRemoteConfigs: [[url: 'https://gitlab.com/a16624741591/server.git']]])
            }
        }



        stage('构建镜像'){
           steps {
                withCredentials([usernamePassword(credentialsId: "${docker_registry_auth}", passwordVariable: 'plh12345', usernameVariable: 'admin')]) {
                sh '''
                  echo '
                    FROM golang:1.17-alpine3.13 AS builder
                    WORKDIR /usr/src/app
                    ENV GO111MODULE=on

                    ENV GOPROXY https://goproxy.cn
                    COPY ..  .
                    RUN go env
                    COPY ../.. .
                    RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -o server ./okr/cmd
                    FROM alpine:3.12 AS final

                    WORKDIR /app
                    #COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
                    COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
                    COPY --from=builder /usr/src/app/server /opt/app/
                    #USER app-runner
                    CMD ["/opt/app/cmd"]
                  ' > Dockerfile
                '''
                sh 'docker build --platform linux/amd64 -t 119.29.5.54:1180/library/okr:t${_VERSION} .'
                sh 'docker login -u admin -p plh12345 http://119.29.5.54:1180/'
                sh 'docker push 119.29.5.54:1180/library/okr:t${_VERSION}'
                }
           }
        }

        stage('登录服务器pull'){
            steps{
                sh '''
                ssh -tt root@119.29.5.54 <<EOF
                   ./check_containers.sh "t${_VERSION}"
                   docker run -d -p 8085:8085 119.29.5.54:1180/library/okr:t${_VERSION}
                   exit
                '''
                echo 'pull ok'
            }
        }
    }
}