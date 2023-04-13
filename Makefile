# Build go binaries
build-backend:
	@echo "Building go binaries..."
	@go build -o bin/bucketfwd/main cmd/bucketfwd/main.go

# Deploy the service
deploy:
	$(MAKE) build
	@echo "Deploying service..."
	@cd deployments && cdk synth && cdk deploy

copy-frontend:	
	@gsutil cp -r frontend/my-app/dist/* gs://ece461-dev.tonyni.ca 

docker-go-bp:
	@docker build -t gcr.io/group-24-ece461/api:0.1 --platform=linux/amd64 -f build/package/Dockerfile.api .
	@docker push gcr.io/group-24-ece461/api:0.1
