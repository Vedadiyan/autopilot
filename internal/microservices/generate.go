package microservices

import _ "embed"

var (
	//go:embed templates/microservice.go.tmpl
	_microservice string
	//go:embed templates/microservice.test.go.tmpl
	_test string
	//go:embed templates/microservice.init.go.tmpl
	_init string
	//go:embed templates/microservice.implement.go.tmpl
	_implement string
)

type Microservice struct {
	Create    Create    `long:"create" short:"" help:"Creates a microservice"`
	Init      Init      `long:"init" short:"" help:"Initiates a microservice"`
	Implement Implement `long:"implement" short:"" help:"Implements a microservice"`
}
