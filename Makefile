build: pkg/generated
	go build

pkg/generated:
	go generate
	rg "generated by main" -t go --files-with-matches | xargs sed -i 's/generated by main/generated by wrangler/'

run: build
	kapp deploy -a sr -f manifests/crd.yaml -f manifests/example-foo.yaml -y
	./secret-replicator

clean:
	rm -rf ~/go/src/github.com/brandonkal/pkg
	go run pkg/codegen/cleanup/main.go
	rm pkg/apis/replicator.kite.run/v1alpha1/doc.go

deploy:
	kapp deploy -a sr -f manifests/crd.yaml -f manifests/example-foo.yaml -y
