# Build go binaries
build-backend:
	@echo "Building go binaries..."
	@go build -o bin/bucketfwd/main cmd/bucketfwd/main.go

# Deploy the service
deploy:
	$(MAKE) build
	@echo "Deploying service..."
	@cd deployments && cdk synth && cdk deploy

# docker build -t gcr.io/group24ece404/bucketfwd:0.1.1 --platform=linux/amd64 -f build/package/Dockerfile.bucketfwd .
# docker build -t gcr.io/group24ece404/frontend:0.1 --platform=linux/amd64 -f build/package/Dockerfile.frontend .
