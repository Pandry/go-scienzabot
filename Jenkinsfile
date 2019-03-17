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
        sh '''go get "github.com/go-telegram-bot-api/telegram-bot-api"


'''
        sh 'apk add --update --no-cache alpine-sdk'
        sh 'go get "github.com/mattn/go-sqlite3"'
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