pipeline {
  agent {
    kubernetes {
      yaml '''
        apiVersion: v1
        kind: Pod
        spec:
          containers:
          - name: golang
            image: golang:alpine
            command:
              - sleep
            args:
              - 99d
        '''
    }
  }

  stages {
    stage('Check go formatting') {
      steps {
        container('golang') {
          sh '''
            UNFORMATTED=$(gofmt -l .)

            if [ -n "$UNFORMATTED" ]; then
              echo "ERROR: go files not formatted correctly:"
              echo "$UNFORMATTED"
              echo Run "gofmt -w" . before pushing.
              exit 1
            else
              echo "Code is correctly formatted."
            fi
          '''
        }
      }
    }
    stage('Build api server') {
      steps {
        container('golang') {
          sh 'go build -o bin/server ./cmd/server'
        }
      }
    }
  }
}
