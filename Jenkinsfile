pipeline {
  agent any
  environment {
    IMAGE_URL="gcr.io/jjzeng-knative-dev/simple-http:${BUILD_TAG}"
  }
  stages {
    stage('Prepare') {
      steps {
        // This is the service account key we upload as a secret file.
        withCredentials([file(credentialsId: '83ced48e-9ec7-4381-a219-f9d6c11289cc', variable: 'SERVICE_ACCOUNT_KEY')]) {
          // Authenticate gcloud
          sh("gcloud auth activate-service-account --key-file $SERVICE_ACCOUNT_KEY")

          // Associate kubectl with the cluster.
          sh("gcloud container clusters get-credentials jj-knative-cluster --zone=us-west1-a")

          // These settings make subsequent "gcloud run" commands shorter.
          sh("gcloud config set run/platform gke")
          sh("gcloud config set run/cluster jj-knative-cluster")
          sh("gcloud config set run/cluster_location us-west1-a")

          // take a snapshot of the service, in case we need to rollback
          sh("kubectl get ksvc simple-http -oyaml | grep -v resourceVersion > /tmp/${BUILD_TAG}.yaml")
        }
      }
    }

    stage('Cloud Build') {
      steps {
        git url: 'https://github.com/jjzeng-seattle/simple-http.git'
        sh("cat cloudbuild.yaml")
        sh("gcloud builds submit --config cloudbuild.yaml --substitutions=_IMAGE_TAG=${BUILD_TAG}  .")
      }
    }

    stage('Deploy with no traffic') {
      steps {
        sh("gcloud alpha run deploy simple-http --namespace=default --image=${IMAGE_URL} --connectivity=external --set-env-vars=BUILD=${BUILD_TAG} --no-traffic --revision-suffix=${BUILD_TAG}")
      }
    }

    stage("Wait for revision ready") {
      steps {
        // TODO: We could query the status of the newly created revision.
        sh("sleep 10")
      }
    }

    stage('Verify revision') {
      steps {
        sh("kubectl exec -it sleep-7b9758757b-7578x -- curl -f simple-http-${BUILD_TAG}-private/healthcheck")
      }
    }
    stage('Add 50% traffic') {
      steps {
        sh(" gcloud alpha run services update-traffic simple-http --to-revisions simple-http-${BUILD_TAG}=50")
        }
    }
    stage('50% Rollout tests') {
      steps {
        sh("curl -f http://simple-http.default.35.185.251.139.xip.io/healthcheck?status=s")
      }
    }
    stage('Add 100% traffic') {
      steps {
        sh(" gcloud alpha run services update-traffic simple-http --to-revisions simple-http-${BUILD_TAG}=100")
        }
    }
    stage('100% Rollout tests') {
      steps {
        sh("curl -f http://simple-http.default.35.185.251.139.xip.io/healthcheck?status=s")
      }
    }
  }
  post {
    success {
      echo "success"
    }
    failure {
      sh("sleep 10")
      sh("kubectl apply -f /tmp/${BUILD_TAG}.yaml")
      echo "failure and rollback"
    }
    always {
      cleanWs()
    }
  }
}

