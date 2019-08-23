#!/usr/bin/env groovy
pipeline {
  agent any
  stages {

    stage('Build') {
      agent {
        docker {
          image 'golang:1.12.9'
        }
      }
      steps {
        sh 'go build -a -tags netgo -ldflags "-w -extldflags \'-static\'" -o login-info *.go'
      }
    }

    stage('Sonarqube') {
      environment {
        scannerHome = tool 'sonar-scanner'
      }
      steps {
        withSonarQubeEnv('sonarqube') {
          sh "${scannerHome}/bin/sonar-scanner"
        }
        timeout(time: 10, unit: 'MINUTES') {
          waitForQualityGate abortPipeline: true
        }
      }
    }

  }
}
