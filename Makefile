
REPO?=seb-demo

all: build auth push

auth:
	aws ecr-public get-login-password --region us-east-1 | docker login --username AWS --password-stdin public.ecr.aws

build:
	docker build -t public.ecr.aws/$(REPO)/eks-example-go:cluster-name .

push:
	docker push public.ecr.aws/$(REPO)/eks-example-go:cluster-name