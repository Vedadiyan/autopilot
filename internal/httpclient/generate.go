package httpclient

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/vedadiyan/autopilot/internal/funcs"
)

var (
	//go:embed templates/httpclient.go.tmpl
	_httpclient string
)

type Client struct {
	Name        string
	ContentType string
	Method      string
	URL         string
}

type ClientContext struct {
	Package string
	Clients []Client
}

type HttpClient struct {
	FileName  string `long:"--filename" short:"-f" help:"Path to the Postman collection file"`
	OutputDir string `long:"--output-dir" short:"-o" help:"Output directory"`
}

func (h HttpClient) Run() error {
	file, err := os.ReadFile(h.FileName)
	if err != nil {
		return err
	}
	model := Postman{}
	err = json.Unmarshal(file, &model)
	if err != nil {
		return err
	}
	clients := make([]Client, 0)
	for _, item := range model.Item {
		var contentType string
		for _, header := range item.Request.Header {
			if strings.ToLower(header.Key) == "content-type" {
				contentType = header.Value
				break
			}
		}
		client := Client{
			Name:        item.Name,
			ContentType: contentType,
			Method:      item.Request.Method,
			URL:         item.Request.URL.Raw,
		}
		clients = append(clients, client)
	}
	t, err := template.New("httpclient").Funcs(funcs.Funcs).Parse(_httpclient)
	if err != nil {
		return err
	}
	clientContext := ClientContext{
		Package: "test",
		Clients: clients,
	}
	var output bytes.Buffer
	err = t.Execute(&output, clientContext)
	if err != nil {
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/a.httpclient.go", h.OutputDir), output.Bytes(), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
