pipeline {
    agent any

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Build and Deploy') {
            steps {
                sh '''
                    docker compose down || true
                    docker compose up -d --build
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
