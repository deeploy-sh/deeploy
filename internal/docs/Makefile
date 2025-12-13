templ:
	templ generate --watch --proxy="http://localhost:8090" --open-browser=false

server:
	air \
	--build.cmd "go build -o tmp/main ./cmd/main.go" \
	--build.bin "tmp/main" \
	--build.delay "100" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

tailwind-clean:
	tailwindcss -i ./assets/css/input.css -o ./assets/css/output.css --minify

tailwind-watch:
	tailwindcss -i ./assets/css/input.css -o ./assets/css/output.css --watch

dev:
	make tailwind-clean
	make -j3 tailwind-watch templ server
