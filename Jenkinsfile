node ('master'){
    def dockerFile = 'Dockerfile'
    def namespace = 'integrasi-mitra'
    def appName = "integrasi-mitra-grpc"
    def gitUrl = 'https://bitbucket.tapera.go.id/scm/int/integrasi.git'
    def imageTag = 'latest'
    def sonarSrc = 'grpc/v1'
    def sonarTest ='grpc/v1'
    def sonarCoverage = 'grpc/cover.out'


    def nexusPluginsRepoUrl = 'https://nexus.tapera.go.id/repository/maven-central/'
    def memLimit = '512Mi'
    def cpuLimit = '500m'
    def imagePullSecret = 'nexus-dev-repo'
    def serviceTypeGKE = 'LoadBalancer'
    def serviceTypeALI = 'NodePort'
    def gateway = 'istio-gateway'
    def host = 'hello.sotech.info'
    def gkeKubeConfig = '/jenkinsdev01/jenkins-home/kubeconfig/config-gke-dev'
    def aliKubeConfig = '/jenkinsdev01/jenkins-home/kubeconfig/config-ali-dev'
    def jmeterTestFileGCP = 'petclinic-gke.jmx'
    def jmeterTestFileALI = 'petclinic-ali.jmx'
    def jmeterNumberThreads = '20'
    def jmeterRampUp = '3'
    def jmeterLoopCount = '10'
    def jmeterErrorRateThresholdPercent = '99'
    def jmeterGitRepo = 'https://bitbucket.tapera.go.id/scm/sam/jmeter.git'
    def katalonGitRepo = 'https://bitbucket.tapera.go.id/scm/sam/katalon.git'
    def katalonBranch = 'master'
    def katalonProjectName = 'TestPlease.prj'
    def katalonTestSuiteName = 'TestSuiteTapera'
    def jiraIssueKey
    def jiraUrl = 'https://jira.tapera.go.id'
    def branch = 'master'

    def nexusGoCentral = 'nexus.tapera.go.id/repository/go-central'

    def nexusDockerDevRepoGCP = '10.172.24.50:8082'
    def nexusDockerDevRepoALI = '10.172.24.50:8082'

    environment {
        GO111MODULE = 'on'
    }

    stage('Checkout'){
        echo 'Checking out SCM'
        checkout scm: [$class: 'GitSCM', userRemoteConfigs: [[credentialsId: 'ci-cd', url: "${gitUrl}"]], branches: [[name: "${branch}"]]]
    }

    stage('Unit Test') {
        echo 'unit test'
        sh '''
        go version
        cd grpc

        rm -rf report.xml
        rm -rf cover.out
        rm -rf cover
        
        export PATH=$PATH:$(go env GOPATH)/bin

        go mod tidy -v

        go get -u golang.org/x/lint/golint
        golint ./v1/...
        
        go get -u github.com/jstemmer/go-junit-report
        go clean -testcache
        '''

        def sts = 1
        try {
            sts = sh (
                returnStatus: true, 
                script: '''
                    cd grpc
                    export PATH=$PATH:$(go env GOPATH)/bin
                    # CGO_ENABLED=0 go test ./v1/... -cover -v -coverprofile=cover.out 
                    CGO_ENABLED=0 go test ./v1/... -cover -v -coverprofile=cover.out 2>&1 | go-junit-report -set-exit-code > ./report.xml
                    echo $?
                '''
            )
            echo sts.toString()
        }
        finally{
            if (fileExists('./grpc/report.xml')) { 
                echo 'junit report'
                try{
                    junit '**/grpc/report.xml'
                }
                catch(e){
                }
            }
            if(sts == 1){
                error('unit testing Fail!!!!')
            }
        }
    }

    stage("SonarQube Analysis"){
        //scan
		withSonarQubeEnv(credentialsId: 'sonarqube-token', installationName: 'sonarqube') {
            // kalau tidak bisa pakai bat
			sh """
                sonar-scanner -Dsonar.projectKey=${appName} \
                            -Dsonar.sources=${sonarSrc} \
                            -Dsonar.exclusions="**/*.pb.go" \
                            -Dsonar.coverage.exclusions="**/*.pb.go" \
                            -Dsonar.qualitygate.wait=true \
                            -Dsonar.language=go \
                            -Dsonar.dynamicAnalysis=reuseReports \
                            -Dsonar.go.coverage.reportPaths=${sonarCoverage} \
                            -Dsonar.tests=${sonarTest} \
                            -Dsonar.test.inclusions="**/*_test.go" \
                            -Dsonar.test.exclusions="**/vendor/**,config/**,docs/**,resources/**, **/*.pb.go" \
            """
		}
		//quality gate
		timeout(time: 1, unit: 'HOURS') {
              def qg = waitForQualityGate()
              if (qg.status != 'OK') {
                  error "Pipeline aborted due to quality gate failure: ${qg.status}"
              }
        }
    }

    stage("Security Check"){
        echo 'security check'
    }

    stage('Build Image') {
        echo 'build image'
        sh """ 
            docker-compose build --force-rm
        """
    }

    stage('Publish to GKE'){
        echo 'publish to GKE'
        //def nexusUrl = nexusDockerDevRepoGCP.replace("http://","")
        withCredentials([usernamePassword(credentialsId: 'ci-cd', passwordVariable: 'nexusPassword', usernameVariable: 'nexusUsername')]) {
            sh """
                docker login -u=${nexusUsername} -p=${nexusPassword} ${nexusDockerDevRepoGCP}
                docker push ${nexusDockerDevRepoGCP}/${appName}:${imageTag}
                docker logout ${nexusDockerDevRepoGCP} 
            """
        }
    }

    /*stage('Publish to ALI'){
        echo 'publish to li'
        def nexusUrl = nexusDockerDevRepoALI.replace("http://","")
        withCredentials([usernamePassword(credentialsId: 'ci-cd', passwordVariable: 'nexusPassword', usernameVariable: 'nexusUsername')]) {
            sh """
                docker login -u=${nexusUsername} -p=${nexusPassword} ${nexusDockerDevRepoALI}
                docker push ${nexusUrl}/${appName}:${imageTag} 
                docker logout ${nexusDockerDevRepoALI}
            """
        }
    }*/

    stage("Delete Image"){
        echo 'delete image'
        sh """
            docker-compose rm -f
        """       
    }
        
    stage('Deploy to GKE'){
        echo 'deploy to GKE'
        
        sh """
		cat kubernetes/gcp-deployment-template.yaml | sed 's/{APP_NAME}/${appName}/g'  \
		| sed 's/{NEXUS_REPO}/${nexusDockerDevRepoGCP}/g' | sed 's/{IMG_TAG}/${imageTag}/g' \
		| sed 's/{MEMORY_LIMIT}/${memLimit}/g' | sed 's/{CPU_LIMIT}/${cpuLimit}/g' \
		| sed 's/{IMG_PULL_SECRET}/${imagePullSecret}/g' |  sed 's/{SERVICE_TYPE}/${serviceTypeGKE}/g' \
		| kubectl --kubeconfig='${gkeKubeConfig}' apply -n '${namespace}' -f -
		
		kubectl --kubeconfig='${gkeKubeConfig}' rollout status deployment/'${appName}' -n '${namespace}'
		"""
    }

    /*
    stage('Deploy to ALI'){
        echo 'deploy to ALI'

        sh """
		cat kubernetes/ali-deployment-template.yaml | sed 's/{APP_NAME}/${appName}/g'  \
		| sed 's/{NEXUS_REPO}/${nexusDockerDevRepoGCP}/g' | sed 's/{IMG_TAG}/${imageTag}/g' \
		| sed 's/{MEMORY_LIMIT}/${memLimit}/g' | sed 's/{CPU_LIMIT}/${cpuLimit}/g' \
		| sed 's/{IMG_PULL_SECRET}/${imagePullSecret}/g' |  sed 's/{SERVICE_TYPE}/${serviceTypeALI}/g' \
		| kubectl --kubeconfig='${aliKubeConfig}' apply -n '${namespace}' -f -
		
		cat kubernetes/ali-virtual-service-template.yaml | sed 's/{APP_NAME}/${appName}/g'  \
		| sed 's/{GATEWAY}/${gateway}/g' | sed 's/{HOST}/${host}/g' \
		| sed 's/{NAMESPACE}/${namespace}/g' \
		| kubectl --kubeconfig='${aliKubeConfig}' apply -f -
		
		kubectl --kubeconfig='${aliKubeConfig}' rollout status deployment/'${appName}' -n '${namespace}'
		"""
    }*/
    
    stage ("Regression Test") {
        echo 'regresion test'
        /*
        echo "Build"
		node ('kre-centos') {
			sleep 60
			cleanWs deleteDirs: true
			//git credentialsId: 'portal', branch: 'development', url: 'https://bitbucket.tapera.go.id/scm/PRT/portal-peserta-katalon-script.git'
            checkout scm: [$class: 'GitSCM', userRemoteConfigs: [[credentialsId: 'ci-cd', url: "${gitUrlKatalon}"]], branches: [[name: "${branch}"]]]
			
			withCredentials([string(credentialsId: 'katalon-api-key', variable: 'secret')]) {
				echo "workspace : ${workspace}"
				sh """
				pwd
				ls
				
				/katalon01/katalon-studio-engine/katalonc -noSplash -runMode=console \
				-projectPath='${workspace}/${katalonProjectName}' -retry=0 \
				-testSuitePath='Test Suites/${katalonTestSuiteName}' -executionProfile='default' \
				-browserType="Chrome (headless)" -apiKey='${secret}'
				"""
			}
		}*/
	}
        
    stage ("Jira Update Status") {
        echo 'jira update status'
    }
}


