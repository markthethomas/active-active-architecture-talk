.PHONY: build clean deploy

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/api/health/check api/health/check.go &
	env GOOS=linux go build -ldflags="-s -w" -o bin/api/live/connect api/live/connect.go &
	env GOOS=linux go build -ldflags="-s -w" -o bin/api/live/disconnect api/live/disconnect.go &
	env GOOS=linux go build -ldflags="-s -w" -o bin/api/live/broadcast api/live/broadcast.go &
	env GOOS=linux go build -ldflags="-s -w" -o bin/api/people/create api/people/create.go &
	env GOOS=linux go build -ldflags="-s -w" -o bin/api/people/get api/people/get.go &

	wait
clean:
	rm -rf ./bin
deploy: clean build
	sls deploy -v --aws-s3-accelerate -s $(STAGE) -r $(REGION)
deploy-func: clean build
	sls deploy -v --aws-s3-accelerate -s $(STAGE) -f $(FUNC) -r $(REGION)
deploy-all: clean build
	# US	
	sls deploy --aws-s3-accelerate -s $(STAGE) -r us-west-1
	sls deploy --aws-s3-accelerate -s $(STAGE) -r us-west-2
	sls deploy --aws-s3-accelerate -s $(STAGE) -r us-east-1
	sls deploy --aws-s3-accelerate -s $(STAGE) -r us-east-2