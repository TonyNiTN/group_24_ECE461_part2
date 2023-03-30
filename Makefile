# Build go binaries
build-backend:
	@echo "Building go binaries..."
	@go build -o bin/bucketfwd/main cmd/bucketfwd/main.go

# Deploy the service
deploy:
	$(MAKE) build
	@echo "Deploying service..."
	@cd deployments && cdk synth && cdk deploy