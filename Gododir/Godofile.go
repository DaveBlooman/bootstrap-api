package main

import . "gopkg.in/godo.v1"

func tasks(p *Project) {
	p.Task("ci-build-binary", func() error {
		Run("CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags=\"-s\" -o bootstrap-api")
		return Run("mv ./bootstrap-api /SOURCES")
	})

	p.Task("default", func() error {
		return Run("go run main.go")
	})

	p.Task("save-deps", func() error {
		return Run("godep save")
	})

	p.Task("build", D{"save-deps"}, func() error {
		return Run("go build -o ../bin/bootstrap-api")
	})

	p.Task("tests", func() error {
		return Run("go test ./...")
	})

	p.Task("generate-docs", func() error {
		return Run(`swagger -apiPackage "github.com/DaveBlooman/bootstrap-api" -ignore "[^bootstrap\\-api]" -format "swagger" -output "docs/swagger/"`)
	})
}

func main() {
	Godo(tasks)
}
