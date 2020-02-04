pipeline {
        agent any
        stages {
            stage('Test') {
                steps {
                    withCredentials([file(credentialsId: '83ced48e-9ec7-4381-a219-f9d6c11289cc', variable: 'SERVICE_ACCOUNT_KEY')]) {
                      sh("gcloud auth activate-service-account --key-file $SERVICE_ACCOUNT_KEY")
                      sh("gcloud builds submit --config cloudbuild.yaml .")
                      sh("export")
                    }
                }
            }
            stage('Build') {
                steps {
                    sh("""
                    gcloud run deploy simple-http --namespace=default \
                    --image=gcr.io/jjzeng-knative-dev/simple-http:latest \
                    --platform=gke --cluster=jj-knative-cluster \
                    --connectivity=external \
                    --cluster-location=us-west1-a \
                    --set-env-vars=TARGET=Jenkins
                    """)
                }
            }
        }
    }
