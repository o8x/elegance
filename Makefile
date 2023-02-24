all:
	@./build.sh

clean:
	@sh -c '[[ -d build ]] && rm -rf build'

generate:
	@sh -c "cd app ; sqlc generate"
