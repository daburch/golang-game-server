pipeline {
  
  environment {
    registry = '172.16.16.100:5000'
    imageName = 'daburch/golang-game-server'
    dockerImage = '' 
  }  

  agent any

  stages {

    stage('Checkout Source') {
      steps {
        git 'https://github.com/daburch/golang-game-server.git'
      }
    }
    
    stage('Build Image') {
      steps{
        script {
          dockerImage = docker.build "$registry/$imageName:$BUILD_NUMBER"
        }
      }
    }
    
    stage ('Push Image') {
      steps {
        script {
          docker.withRegistry("") {
            dockerImage.push()
          }
        }
      }
    }
    
  }
}
