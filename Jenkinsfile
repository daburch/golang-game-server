pipeline {
  
  environment {
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
        sh "docker build ."
      }
    }
    
  }
}
