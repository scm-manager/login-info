#!/usr/bin/env groovy
def version = 'UNKNOWN'

pipeline {

  agent {
    node {
      label 'docker'
    }
  }

  stages {

    stage('Environment') {
      steps {
        script {
          def commitHashShort = sh(returnStdout: true, script: 'git rev-parse --short HEAD')
          version = "${new Date().format('yyyyMMddHHmm')}-${commitHashShort}".trim()
        }
      }
    }

    stage('Build') {
      agent {
        docker {
          image 'golang:1.12.9'
        }
      }
      environment {
        // change go cache location
        XDG_CACHE_HOME = "${WORKSPACE}/.cache"
      }
      steps {
        sh 'go build -a -tags netgo -ldflags "-w -extldflags \'-static\'" -o target/login-info *.go'
        stash name: 'target', includes: 'target/*'
      }
    }

    stage('Sonarqube') {
      environment {
        scannerHome = tool 'sonar-scanner'
      }
      steps {
        withSonarQubeEnv('sonarcloud.io-scm') {
          sh "${scannerHome}/bin/sonar-scanner"
        }
        timeout(time: 10, unit: 'MINUTES') {
          waitForQualityGate abortPipeline: true
        }
      }
    }

    stage('Docker') {
      steps {
        unstash 'target'
        script {
          docker.withRegistry('', 'hub.docker.com-cesmarvin') {
            def image = docker.build("scmmanager/login-info:${version}")
            image.push()
          }
        }
      }
    }

    stage('Deployment') {
      agent {
        docker {
          image 'lachlanevenson/k8s-helm:v2.14.2'
          args  '--entrypoint=""'
        }
      }
      steps {
        withCredentials([file(credentialsId: 'helm-client-scm-manager', variable: 'KUBECONFIG')]) {
          sh "helm upgrade --install --set image.tag=${version} login-info helm/login-info"
        }
      }
    }

  }
}
