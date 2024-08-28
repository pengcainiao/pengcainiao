pipeline {
    agent any

    environment {
        // 设置环境变量，如Docker镜像名、GitLab仓库地址等
        DOCKER_IMAGE = 'your-docker-repo/your-image-name'
        GITLAB_REPO = 'https://gitlab.com/a16624741591/server.git'
        // 从Jenkins凭证中获取GitLab的凭证ID
        GIT_CREDENTIALS_ID = '11326794'
    }

    stages {

        stage('Clone Repository') {
            steps {
                git branch: 'main', credentialsId: env.GIT_CREDENTIALS_ID, url: env.GITLAB_REPO
            }
        }

        stage('Build Docker Image') {
            steps {
                script {
                    // 构建Docker镜像
                    docker.build(env.DOCKER_IMAGE, '-f ci/Dockerfile .')
                }
            }
        }

        stage('Run Tests') {
            steps {
                script {
                    // 假设在Dockerfile中有一个指定默认的ENTRYPOINT，执行测试
                    docker.image(env.DOCKER_IMAGE).inside {
                        sh 'go test ./...'
                    }
                }
            }
        }

        stage('Push Docker Image') {
            when {
                branch 'main'
            }
            steps {
                script {
                    withDockerRegistry([credentialsId: 'docker-credentials-id', url: '']) {
                        docker.image(env.DOCKER_IMAGE).push()
                    }
                }
            }
        }

        // 这里可以添加更多的stage，例如部署到Kubernetes，通知等
    }

    post {
        always {
            cleanWs()
        }
        success {
            echo 'Pipeline completed successfully'
        }
        failure {
            echo 'Pipeline failed'
        }
    }
}
