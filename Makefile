build: pkg/generated
	go build

pkg/generated:
	go generate
	rg "generated by main" -t go --files-with-matches | xargs sed -i 's/generated by main/generated by wrangler/'

run: build
	kapp deploy -a sr -f manifests/crd.yaml -f manifests/example-foo.yaml -y
	./secret-replicator

clean:
	rm -rf pkg/generated
