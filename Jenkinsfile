pipeline {
  agent any
  stages {
    stage('Build') {
      agent {
        docker {
          image 'pandry/goanalysis'
        }

      }
      steps {
        echo 'Building..'
        sh 'go get ./...'
        sh 'go build ./...'
      }
    }
    stage('Static Analysis') {
      agent {
        docker {
          image 'pandry/goanalysis'
        }

      }
      steps {
        echo 'Testing..'
        sh 'golint .'
        sh 'go vet .'
        sh 'maligned ./...'
        sh 'dingo-hunter migo ./...'
      }
    }
  }
}