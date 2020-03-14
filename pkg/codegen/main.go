package main

import (
	"github.com/brandonkal/secret-replicator/pkg/apis/replicator.kite.run/v1alpha1"
	controllergen "github.com/rancher/wrangler/pkg/controller-gen"
	"github.com/rancher/wrangler/pkg/controller-gen/args"
)

func main() {
	controllergen.Run(args.Options{
		OutputPackage: "github.com/brandonkal/secret-replicator/pkg/generated",
		Boilerplate:   "pkg/codegen/boilerplate.go.txt",
		Groups: map[string]args.Group{
			"replicator.kite.run": {
				Types: []interface{}{
					v1alpha1.SecretExport{},
					v1alpha1.SecretRequest{},
				},
				GenerateTypes: true,
			},
			// wrangler-api project provides pre-generated k8s APIs which
		},
	})
}
