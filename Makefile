default: run

run:
	@echo "Running..."
	gow run .

swag:
	@echo "Generating swagger..."
	swag init