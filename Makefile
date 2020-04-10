
GIT_COMMIT := $(shell git rev-list -1 HEAD)

build:
	env GIT_COMMIT=$(GIT_COMMIT) docker build --build-arg GIT_COMMIT -t arunsworld/go-vue-app:latest .

run:
	docker run -d -e SESSION_KEY=ba6c108c147b1472b9a73c06c0e25ca8390f47a0e439bf9252b37702a9afc529 -p 1000:1000 arunsworld/go-vue-app:latest app

gcp:
	docker run -d -e SESSION_KEY=ba6c108c147b1472b9a73c06c0e25ca8390f47a0e439bf9252b37702a9afc529 --name bootstrap --network apps-net arunsworld/go-vue-app:latest app