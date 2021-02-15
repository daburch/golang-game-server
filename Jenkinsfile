pipeline {
    agent { kubernetes { image 'golang' } }

    stages {
        stage('build') {
            steps {
                sh 'go version'
            }
        }
    }
}
