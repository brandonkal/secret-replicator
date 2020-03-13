build: pkg/generated
	go build

pkg/generated:
	go generate

run: build
	kapp deploy -a sr -f manifests/crd.yaml -f manifests/example-foo.yaml -y
	./secret-replicator

clean:
	rm -rf pkg/generated
