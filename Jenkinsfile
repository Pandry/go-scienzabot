pipeline {
  agent any
  stages {
    stage('Build') {
      agent {
        dockerfile {
          filename 'CIDockerfile'
        }

      }
      steps {
        sh 'cp -r $(pwd)  $GOPATH/src/scienzabot/'
        echo 'Building..'
        sh 'go build .'
      }
    }
    stage('Static Analysis') {
      agent {
        dockerfile {
          filename 'CIDockerfile'
        }

      }
      steps {
        sh 'cp -r $(pwd)  $GOPATH/src/scienzabot/'
        echo 'Testing..'
        sh 'golint .'
        sh 'go vet .'
        sh 'maligned ./...'
        sh 'dingo-hunter migo ./...'
      }
    }
  }
}