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
          - name: buildkitd
            image: moby/buildkit:master-rootless
            args:
              - --oci-worker-no-process-sandbox
            securityContext:
              seccompProfile:
                type: Unconfined
              appArmorProfile:
                type: Unconfined
              runAsUser: 1000
              runAsGroup: 1000
            volumeMounts:
              - mountPath: /home/user/.local/share/buildkit
                name: buildkitd
          volumes:
            - name: buildkitd
              emptyDir: {}
          '''
    }
  }

  environment {
    DOCKER_CONFIG = "${WORKSPACE}"
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
              echo Run "gofmt -w ." before pushing.
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
          sh 'apk add --no-cache make'
          sh 'make api'
        }
      }
    }
    stage('Build and push Docker image') {
      environment {
        FORGEJO_TOKEN = credentials('khuong-forge-pat')
      }
      steps {
        script {
          if (env.TAG_NAME) {
            env.IMAGE_TAG = env.TAG_NAME
          } else {
            env.IMAGE_TAG = env.GIT_COMMIT
          }
        }
        container('buildkitd') {
          sh '''
            AUTH=$(echo -n "khuong:$FORGEJO_TOKEN" | base64 -w 0)

            printf '{
  "auths": {
    "forge.nhmk.de": {
      "auth": "%s"
    }
  }
}
' "$AUTH" > "$DOCKER_CONFIG/config.json"
            '''
          sh '''
            echo ${IMAGE_TAG}
            echo ${TAG_NAME}
            buildctl build \
              --frontend=dockerfile.v0 \
              --local context=. \
              --local dockerfile=. \
              --output type=image,name=forge.nhmk.de/khuong/kurz:${IMAGE_TAG},push=true
            '''
        }
      }
    }
    stage('Deploy to dev') {
      when {
        branch 'main'
      }
      steps {
        updateGitOps('dev')
      }
    }
    stage('Deploy to prod') {
      when {
        buildingTag()
      }
      steps {
        updateGitOps('prod')
      }
    }
  }
}

def updateGitOps(String overlayPath){
  dir('gitops-tmp') {
    deleteDir()
    sshagent(credentials: ['jenkins-github']) {
      sh '''
        export GIT_SSH_COMMAND="ssh -o StrictHostKeyChecking=no"
        git clone --depth 1 git@github.com:minhkhuong0/homelab-iac.git
      '''
      dir('homelab-iac/kubernetes/deploy/kurz') {
        dir(overlayPath) {
          sh 'sed -i "s/Tag:.*/Tag: ${IMAGE_TAG}/g" kustomization.yaml'
        }
      }
      sh '''
        cd homelab-iac
        echo "${overlayPath}"
        git config --global user.name "Jenkins CI"
        git config --global user.email "khuong.nguyen@tutamail.com"

        git add .
        git commit -m "ci: update kurz image tag to ${IMAGE_TAG}"
        git push origin main
      '''
    }
  }
}
