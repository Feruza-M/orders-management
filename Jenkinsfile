pipeline {
    agent any

    parameters {
        string(name: 'APP_VERSION', defaultValue: '1.0', description: 'Docker image version')
    }

    environment {
        IMAGE_NAME = 'filizmamedova/orders-management'
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Build and Push') {
            steps {
                withCredentials([usernamePassword(
                    credentialsId: 'dockerhub-creds',
                    usernameVariable: 'DOCKER_USER',
                    passwordVariable: 'DOCKER_PASS'
                )]) {
                    sh '''
                        echo "$DOCKER_PASS" | docker login -u "$DOCKER_USER" --password-stdin
                        docker build -t ${IMAGE_NAME}:${APP_VERSION} .
                        docker push ${IMAGE_NAME}:${APP_VERSION}
                    '''
                }
            }
        }

        stage('Deploy') {
            steps {
                sh '''
                    export APP_VERSION=${APP_VERSION}
                    docker compose pull
                    docker compose up -d
                '''
            }
        }

        stage('Health Check') {
            steps {
                sh '''
                    chmod +x check_orders.sh
                    ./check_orders.sh
                '''
            }
        }
    }

    post {
        failure {
            sh 'docker compose logs || true'
        }
    }
}
