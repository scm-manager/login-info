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
          image 'golang:1.18.2'
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
          image 'lachlanevenson/k8s-helm:v3.2.1'
          args  '--entrypoint=""'
        }
      }
      steps {
        withCredentials([file(credentialsId: 'helm-client-scm-manager', variable: 'KUBECONFIG')]) {
          sh "helm upgrade --install --set image.tag=${version} login-info helm/login-info"
        }
      }
    }
    
    stage('Update GitHub') {
      when {
        branch pattern: 'master', comparator: 'GLOB'
	    expression { return isBuildSuccess() }
      }
      steps {
        sh 'git checkout master'
        
        // push changes to GitHub
        authGit 'cesmarvin', "push -f https://github.com/scm-manager/login-info master --tags"
      }
    }
  }

  post {
    failure {
      mail to: "scm-team@cloudogu.com",
        subject: "${JOB_NAME} - Build #${BUILD_NUMBER} - ${currentBuild.currentResult}!",
        body: "Check console output at ${BUILD_URL} to view the results."
    }
  }
}

void authGit(String credentials, String command) {
  withCredentials([
    usernamePassword(credentialsId: credentials, usernameVariable: 'AUTH_USR', passwordVariable: 'AUTH_PSW')
  ]) {
    sh "git -c credential.helper=\"!f() { echo username='\$AUTH_USR'; echo password='\$AUTH_PSW'; }; f\" ${command}"
  }
}


boolean isBuildSuccess() {
  return currentBuild.result == null || currentBuild.result == 'SUCCESS'
}