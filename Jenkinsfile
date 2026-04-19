pipeline {
    agent any

    parameters {
        string(name: 'APP_VERSION', defaultValue: '1.0', description: 'Docker image version')
    }

    environment {
        IMAGE_NAME = 'docferuza2024/orders-management'
        VERSION = "${params.APP_VERSION ?: '1.0'}"
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
                        set -e
                        echo "Using version: ${VERSION}"
                        echo "$DOCKER_PASS" | docker login -u "$DOCKER_USER" --password-stdin
                        docker build -t ${IMAGE_NAME}:${VERSION} .
                        docker push ${IMAGE_NAME}:${VERSION}
                    '''
                }
            }
        }

        stage('Deploy') {
            steps {
                sh '''
                    set -e
                    export APP_VERSION=${VERSION}
                    echo "Deploying version: ${APP_VERSION}"
                    docker compose down || true
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
        always {
            sh 'docker compose ps || true'
        }
        failure {
            sh 'docker compose logs || true'
        }
    }
}
