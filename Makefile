BINARY_NAME := kuma-waybar

.PHONY: build test install

build:
	@echo "Building $(BINARY_NAME)..."
	go build -ldflags="-s -w" -o $(BINARY_NAME) .

test:
	@echo "Running tests..."
	go test ./...

fmt:
	@echo "Formatting code..."
	go fmt ./...

install: build
	@# Check for the existence of ~/.config/waybar
	@if [ ! -d "$$HOME/.config/waybar" ]; then \
		echo "waybar not found"; \
		exit 1; \
	fi
	@echo "copying $(BINARY_NAME) to $$HOME/.local/bin/"
	@mkdir -p "$$HOME/.local/bin"
	@cp ./$(BINARY_NAME) "$$HOME/.local/bin/$(BINARY_NAME)"
	@echo "Installing at /usr/local/bin/$(BINARY_NAME) (CTRL+C to cancel system install)"
	@sudo cp ./$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	@echo "done."
	@echo ""
	@echo "Please add the following to your waybar config:"
	@echo "    \"custom/kuma-waybar\": {"
	@echo "        \"exec\": \"$(BINARY_NAME) --format=waybar --env=$$HOME/.config/waybar/kuma.env\","
	@echo "        \"interval\": 60,"
	@echo "        \"on-click\": \"$(BINARY_NAME) open --env=$$HOME/.config/waybar/kuma.env\","
	@echo "        \"format\": \"Kuma {}\","
	@echo "    },"
	@echo ""
	@echo "Please add your UPTIME_KUMA_API_KEY & UPTIME_KUMA_BASE_URL to $$HOME/.config/waybar/scripts/kuma-waybar.env"
	@echo "Optionally use --env=./second.env to allow for multiple Uptime Kuma"