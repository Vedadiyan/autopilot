package microservices

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	funcs "github.com/vedadiyan/autopilot/internal/funcs"
)

type Create struct {
	Package          string `long:"--package" short:"-p" help:"Name of the package"`
	Service          string `long:"--service" short:"-s" help:"Name of the service"`
	ImportPath       string `long:"--import-path" short:"-i" help:"Import path"`
	Request          string `long:"--request" short:"" help:"Name of the request type"`
	Response         string `long:"--response" short:"" help:"Name of the response type"`
	IsApiIntegration bool   `long:"--api-integration" short:"" help:"Specifies if the microservice should integrate an API"`
}

func (c Create) Run() error {
	err := c.GenerateMicroservice()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	err = c.GenerateTest()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func (c Create) GenerateMicroservice() error {
	t, err := template.New("Microservice").Funcs(funcs.Funcs).Parse(_microservice)
	if err != nil {
		return err
	}
	var buffer bytes.Buffer
	err = t.Execute(&buffer, c)
	if err != nil {
		return err
	}
	target, err := createPath(c.Service, "definition")
	if err != nil {
		return err
	}
	fmt.Println(target)
	err = os.WriteFile(target, buffer.Bytes(), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (c Create) GenerateTest() error {
	t, err := template.New("Test").Funcs(funcs.Funcs).Parse(_test)
	if err != nil {
		return err
	}
	var buffer bytes.Buffer
	err = t.Execute(&buffer, c)
	if err != nil {
		return err
	}
	target, err := createPath(c.Service, "test")
	if err != nil {
		return err
	}
	fmt.Println(target)
	err = os.WriteFile(target, buffer.Bytes(), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func Exists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	var output bool
	if err == nil {
		output = true
		return output, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		output = false
		return output, nil
	}
	return false, err
}

func createPath(serviceName string, fileType string) (string, error) {
	path := fmt.Sprintf("%s/%s", os.Getenv("GOGEN_TARGET"), strings.ToLower(serviceName))
	exists, err := Exists(path)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	if !exists {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			fmt.Println(err.Error())
			return "", err
		}
	}
	target := fmt.Sprintf("%s/a.%s_%s.go", path, strings.ToLower(serviceName), fileType)
	return target, nil
}
