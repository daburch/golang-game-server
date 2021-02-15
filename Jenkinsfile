pipeline {

  environment {
    dockerImage = ""
  }

  agent any

  stages {

    stage('Checkout Source') {
      steps {
        git 'https://github.com/daburch/golang-game-server.git'
      }
    }

    stage('Build image') {
      steps{
        docker build .
      }
    }
  }
}
