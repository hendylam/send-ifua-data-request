cd api
rm -rf report.xml

echo -preparing go lint
go get -u golang.org/x/lint/golint
echo -performing go lint
golint ./controller/...

echo -performing go tidy
go mod tidy

echo -preparing unit testing
go get -u github.com/jstemmer/go-junit-report
go clean -testcache

echo -performing unit testing
set -o pipefail
CGO_ENABLED=0 go test ./controller/... -cover -v 2>&1 | tee /dev/stderr | go-junit-report -set-exit-code > ./report.xml

echo -build app
#rm -rf bin 
mkdir -p bin
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin

cd ..


appName="pemanfaatan-dana-api"
imageTag="latest"
namespace="pemanfaatan-dana"
nexusUsername="raditya.umbas"
nexusPassword="Ecomindo@2o2o"
nexusDockerDevRepoGCP="10.172.24.50:8082"
nexusDockerDevRepoALI="10.172.24.50:8082"
memLimit="512Mi"
cpuLimit="500m"
imagePullSecret="nexus-dev-repo"
serviceTypeGKE="LoadBalancer"
serviceTypeALI="NodePort"
gkeKubeConfig="/jenkinsdev01/jenkins-home/kubeconfig/config-gke-dev"
aliKubeConfig="/jenkinsdev01/jenkins-home/kubeconfig/config-ali-dev"
gateway="istio-gateway"
host="hello.sotech.info"
imagePullSecret="nexus-dev-repo"

echo -build image
docker-compose build --force-rm 

echo -push image to GCP nexus
docker login -u=$nexusUsername -p=$nexusPassword $nexusDockerDevRepoGCP
docker push $nexusDockerDevRepoGCP/$appName:$imageTag
docker logout $nexusDockerDevRepoGCP

echo -deploy to GKE
cat kubernetes/gcp-deployment-template.yaml | sed 's/{APP_NAME}/$appName/g'  \
| sed 's/{NEXUS_REPO}/$nexusDockerDevRepoGCP/g' | sed 's/{IMG_TAG}/$imageTag/g' \
| sed 's/{MEMORY_LIMIT}/$memLimit/g' | sed 's/{CPU_LIMIT}/$cpuLimit/g' \
| sed 's/{IMG_PULL_SECRET}/$imagePullSecret/g' |  sed 's/{SERVICE_TYPE}/$serviceTypeGKE/g' \
| kubectl --kubeconfig='$gkeKubeConfig' apply -n '$namespace' -f -

kubectl --kubeconfig='$gkeKubeConfig' rollout status deployment/'$appName' -n '$namespace'


#echo -push image to ALI nexus
#docker login -u=$nexusUsername -p=$nexusPassword $nexusDockerDevRepoALI
#docker push $nexusDockerDevRepoALI/$appName:$imageTag 
#docker logout $nexusDockerDevRepoALI

#echo -deploy to ALI
#cat kubernetes/ali-deployment-template.yaml | sed 's/{APP_NAME}/$appName/g'  \
#| sed 's/{NEXUS_REPO}/$nexusDockerDevRepoGCP/g' | sed 's/{IMG_TAG}/$imageTag/g' \
#| sed 's/{MEMORY_LIMIT}/$memLimit/g' | sed 's/{CPU_LIMIT}/$cpuLimit/g' \
#| sed 's/{IMG_PULL_SECRET}/$imagePullSecret/g' |  sed 's/{SERVICE_TYPE}/$serviceTypeALI/g' \
#| kubectl --kubeconfig='$aliKubeConfig' apply -n '$namespace' -f -

#cat kubernetes/ali-virtual-service-template.yaml | sed 's/{APP_NAME}/$appName/g'  \
#| sed 's/{GATEWAY}/$gateway/g' | sed 's/{HOST}/$host/g' \
#| sed 's/{NAMESPACE}/$namespace/g' \
#| kubectl --kubeconfig='$aliKubeConfig' apply -f -

#kubectl --kubeconfig='$aliKubeConfig' rollout status deployment/'$appName' -n '$namespace'

echo -delete image
docker-compose rm -f
