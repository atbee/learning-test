#!/bin/sh

rm -rf /var/lib/cloud/*
sudo apt-get update -y
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh
sudo usermod -aG docker ubuntu
sudo curl -L "https://github.com/docker/compose/releases/download/1.25.5/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
sudo apt install -y awscli

# Download GIT_LAST_COMMIT and docker-compose.yml from S3
aws s3 cp s3://morchana-static-qr-code-deployment/api/production/GIT_LAST_COMMIT GIT_LAST_COMMIT
aws s3 cp s3://morchana-static-qr-code-deployment/api/production/docker-compose.yml docker-compose.yml

# Pull Docker images from ECR
$(aws ecr get-login --no-include-email --region ap-southeast-1)
GIT_LAST_COMMIT=$(cat GIT_LAST_COMMIT) docker-compose pull

# Run application
GIT_LAST_COMMIT=$(cat GIT_LAST_COMMIT) docker-compose up -d
