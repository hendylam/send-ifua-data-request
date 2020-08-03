@echo off 
cd api

echo -preparing go lint
go get -u golang.org/x/lint/golint
echo -performing go lint
golint -set_exit_status ./controller/...

echo -performing go tidy
go mod tidy

echo -preparing unit testing
go get -u github.com/jstemmer/go-junit-report
go clean -testcache

echo -performing unit testing
set CGO_ENABLED=0
go test ./controller/... -cover -v || exit /b

echo -build app
set curd=%~dp0%api\bin
if not exist bin (
    md bin
    if errorlevel 1 (
        exit %ERRORLEVEL%
    )
)

SET GOOS=linux
SET GOARCH=amd64

go build -o ./bin || exit /b

SETLOCAL ENABLEEXTENSIONS DISABLEDELAYEDEXPANSION


cd ..


echo -preparing image push
set appName=integrasi-mitra-api
set imageTag=latest
set namespace=integrasi-mitra
set nexusUsername=raditya.umbas
set nexusPassword=Ecomindo@2o2o
set nexusDockerDevRepoGCP=10.172.24.50:8082
set nexusDockerDevRepoALI=10.172.24.50:8082
set memLimit=512Mi
set cpuLimit=500m
set imagePullSecret=nexus-dev-repo
set serviceTypeGKE=LoadBalancer
set serviceTypeALI=NodePort
set gkeKubeConfig=/jenkinsdev01/jenkins-home/kubeconfig/config-gke-dev
set aliKubeConfig=/jenkinsdev01/jenkins-home/kubeconfig/config-ali-dev
set gateway=istio-gateway
set host=hello.sotech.info
set imagePullSecret=nexus-dev-repo

echo -build image
docker-compose build --force-rm || exit /b

echo -push image to GCP nexus %nexusDockerDevRepoGCP%
docker login -u=%nexusUsername% -p=%nexusPassword% %nexusDockerDevRepoGCP% || exit /b
docker push %nexusDockerDevRepoGCP%/%appName%:%imageTag% || exit /b
docker logout %nexusDockerDevRepoGCP% || exit /b

echo -deploy to GKE
cat kubernetes/gcp-deployment-template.yaml | sed 's/{APP_NAME}/%appName%/g' | sed 's/{NEXUS_REPO}/%nexusDockerDevRepoGCP%/g' | sed 's/{IMG_TAG}/%imageTag%/g' | sed 's/{MEMORY_LIMIT}/%memLimit%/g' | sed 's/{CPU_LIMIT}/%cpuLimit%/g' | sed 's/{IMG_PULL_SECRET}/%imagePullSecret%/g' |  sed 's/{SERVICE_TYPE}/%serviceTypeGKE%/g' | kubectl --kubeconfig='%gkeKubeConfig%' apply -n '%namespace%' -f -
if errorlevel 1 (
    echo error
    exit %ERRORLEVEL%
)
kubectl --kubeconfig='%gkeKubeConfig%' rollout status deployment/'%appName%' -n '%namespace%'
if errorlevel 1 (
    echo error
    exit %ERRORLEVEL%
)


rem echo -build image ALI
rem docker-compose build --build-arg FLAG=nexusDockerDevRepoALI || exit /b

rem echo -push image to ALI nexus %nexusDockerDevRepoALI% || exit /b
rem docker login -u=%nexusUsername% -p=%nexusPassword% %nexusDockerDevRepoALI% || exit /b
rem docker push %nexusDockerDevRepoALI%/%appName%:%imageTag% || exit /b
rem docker logout %nexusDockerDevRepoALI% || exit /b

rem echo -deploy to ALI
rem cat kubernetes/ali-deployment-template.yaml | sed 's/{APP_NAME}/$appName/g' | sed 's/{NEXUS_REPO}/$nexusDockerDevRepoGCP/g' | sed 's/{IMG_TAG}/$imageTag/g' | sed 's/{MEMORY_LIMIT}/$memLimit/g' | sed 's/{CPU_LIMIT}/$cpuLimit/g' | sed 's/{IMG_PULL_SECRET}/$imagePullSecret/g' |  sed 's/{SERVICE_TYPE}/$serviceTypeALI/g' | kubectl --kubeconfig='$aliKubeConfig' apply -n '$namespace' -f -
rem if errorlevel 1 (
rem     exit %ERRORLEVEL%
rem )
rem cat kubernetes/ali-virtual-service-template.yaml | sed 's/{APP_NAME}/$appName/g' | sed 's/{GATEWAY}/$gateway/g' | sed 's/{HOST}/$host/g' | sed 's/{NAMESPACE}/$namespace/g' | kubectl --kubeconfig='$aliKubeConfig' apply -f -
rem if errorlevel 1 (
rem     exit %ERRORLEVEL%
rem )
rem kubectl --kubeconfig='$aliKubeConfig' rollout status deployment/'$appName' -n '$namespace'
rem if errorlevel 1 (
rem     exit %ERRORLEVEL%
rem )

rem docker rmi -f %nexusDockerDevRepoALI%/$appName:$imageTag

echo -delete iamage
docker-compose rm -f

ENDLOCAL