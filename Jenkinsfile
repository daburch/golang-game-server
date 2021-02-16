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
        git url: 'https://github.com/daburch/golang-game-server.git', branch: "$BRANCH_NAME"
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
      when { expression { return "$BRANCH_NAME" == "master" } }
      
      steps {
        script {
          docker.withRegistry("") {
            dockerImage.push()
          }
        }
      }
    }
    
    stage('Deploy App') {
      when { expression { return "$BRANCH_NAME" == "master" } }
      
      steps {
        script {
          kubernetesDeploy(configs: "deploy/deploy.yaml", kubeconfigId: "kconfig")
        }
      }
    }    
    
  }
}
