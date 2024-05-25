ACCOUNT_NUMBER="131791614471"
ECR_NAME="telegram_bot"
IMAGE_TAG="latest"

# add the parameter --progress=plain to view unaggregated logs
docker build -t "${ACCOUNT_NUMBER}.dkr.ecr.ap-southeast-1.amazonaws.com/${ECR_NAME}:${IMAGE_TAG}" . 
docker push "${ACCOUNT_NUMBER}.dkr.ecr.ap-southeast-1.amazonaws.com/${ECR_NAME}:${IMAGE_TAG}"