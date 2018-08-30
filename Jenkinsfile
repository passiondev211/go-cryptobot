node {
  def frontend
  def backend

  stage('Clone repository') {
    checkout scm
  }
  
  stage('Build frontend npm packages') {
    def node_home = tool name: 'Cryptobot Nodejs', type: 'nodejs'
        
    env.NODEJS_HOME = tool name: 'Cryptobot Nodejs', type: 'nodejs'
    env.PATH="${env.NODEJS_HOME}/bin:${env.PATH}"

    dir('fe/') {
      sh "npm install"
      sh "npm run build"
    }
  }

  stage('Build frontend image') {
    frontend = docker.build("registry.profitcoins.io/${ENV}-frontend", "-f Dockerfile-frontend .")
  }

  stage('Push frontend image') {
    withDockerRegistry([credentialsId: 'c407787d-9c7b-4d0f-9412-faacda541cf5', url: 'https://registry.profitcoins.io']) {
      frontend.push("${env.BUILD_NUMBER}")
      frontend.push("latest")
    }
  }

  stage('Remove image from local docker') {
    sh "docker rmi registry.profitcoins.io/${ENV}-frontend:latest"
    sh "docker rmi registry.profitcoins.io/${ENV}-frontend:${env.BUILD_NUMBER}"
  }

  stage('Build backend image') {
    backend = docker.build("registry.profitcoins.io/${ENV}-backend", "-f Dockerfile-backend .")
  }

  stage('Push backend image') {
    withDockerRegistry([credentialsId: 'c407787d-9c7b-4d0f-9412-faacda541cf5', url: 'https://registry.profitcoins.io']) {
      backend.push("${env.BUILD_NUMBER}")
      backend.push("latest")
    }
  }

  stage('Remove image from local docker') {
    sh "docker rmi registry.profitcoins.io/${ENV}-backend:latest"
    sh "docker rmi registry.profitcoins.io/${ENV}-backend:${env.BUILD_NUMBER}"
  }

  stage('Deploy') {
    withCredentials([sshUserPrivateKey(credentialsId: 'eed38b79-9505-4ed0-9777-d29d0c34d4ea', keyFileVariable: 'KEYFILE', passphraseVariable: '', usernameVariable: 'USER')]) {
      withCredentials([usernamePassword(credentialsId: 'c407787d-9c7b-4d0f-9412-faacda541cf5', passwordVariable: 'RegistryPassword', usernameVariable: 'RegistryUser')]) {
         sh 'ssh -o StrictHostKeyChecking=no -i  ${KEYFILE} ${USER}@${ENV}.profitcoins.io "echo ${RegistryPassword} | sudo docker login -u ${RegistryUser} --password-stdin https://registry.profitcoins.io"'
      }   
      sh 'ssh -o StrictHostKeyChecking=no -i ${KEYFILE} ${USER}@${ENV}.profitcoins.io "cd infrastructure && sudo docker-compose pull cryptobot-frontend cryptobot-backend"'
      sh 'ssh -o StrictHostKeyChecking=no -i  ${KEYFILE} ${USER}@${ENV}.profitcoins.io "cd infrastructure && sudo docker-compose rm -sf cryptobot-frontend cryptobot-backend"'
      sh 'ssh -o StrictHostKeyChecking=no -i  ${KEYFILE} ${USER}@${ENV}.profitcoins.io "cd infrastructure && sudo docker-compose up -d cryptobot-frontend cryptobot-backend"'
      sh 'ssh -o StrictHostKeyChecking=no -i  ${KEYFILE} ${USER}@${ENV}.profitcoins.io "sudo docker logout https://registry.profitcoins.io"'
    }
  }
}
