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
            buildctl build \
              --frontend=dockerfile.v0 \
              --local context=. \
              --local dockerfile=. \
              --output type=image,name=forge.nhmk.de/khuong/kurz:${GIT_COMMIT},push=true
            '''
        }
      }
    }
    stage('Update GitOps repo') {
      when {
        branch 'main'
      }
      steps {
        sshagent(credentials: ['jenkins-github']) {
          sh '''
            export GIT_SSH_COMMAND="ssh -o StrictHostKeyChecking=no"
            git config --global user.name "Jenkins CI"
            git config --global user.email "khuong.nguyen@tutamail.com"

            git clone git@github.com:minhkhuong0/homelab-iac.git
            cd homelab-iac/kubernetes/deploy/kurz/dev

            sed -i "s/Tag:.*/Tag: ${GIT_COMMIT}/g" kustomization.yaml

            git add kustomization.yaml
            git commit -m "ci: update image tag to ${GIT_COMMIT}"
            git push origin main
          '''
        }
      }
    }
  }
}
