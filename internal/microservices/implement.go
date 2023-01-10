package microservices

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/vedadiyan/autopilot/internal/funcs"
)

type Implement struct {
	Package    string `long:"--package" short:"-p" help:"Name of the package"`
	Service    string `long:"--service" short:"-s" help:"Name of the service"`
	ImportPath string `long:"--import-path" short:"-i" help:"Import path"`
	Request    string `long:"--request" short:"" help:"Name of the request type"`
	Response   string `long:"--response" short:"" help:"Name of the response type"`
	ClientPath string
}

func (i Implement) Run() error {
	fmt.Println("STARTED")
	t, err := template.New("Handler").Funcs(funcs.Funcs).Parse(_implement)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	i.ClientPath = fmt.Sprintf("%s/client", strings.TrimSuffix(i.ImportPath, "/pb"))
	var buffer bytes.Buffer
	err = t.Execute(&buffer, i)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	target := os.Getenv("GOGEN_TARGET")
	err = os.WriteFile(fmt.Sprintf("%s/a.%s_handler.go", target, strings.ToLower(i.Service)), buffer.Bytes(), os.ModePerm)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
