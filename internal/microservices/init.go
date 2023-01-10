package microservices

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"os/exec"
	"strings"
)

type Init struct {
	Service      string `long:"--service" short:"-s" help:"Name of the service"`
	ImportPath   string `long:"--import-path" short:"-i" help:"Import path"`
	FunctionName string `long:"--function-name" short:"-f" help:"Name of the microservice function"`
}

func (i Init) Run() error {
	wd := fmt.Sprintf("%s/cmd", os.Getenv("GOGEN_WD"))
	exists, err := Exists(wd)
	if err != nil {
		return err
	}
	if !exists {
		err := os.MkdirAll(wd, os.ModePerm)
		if err != nil {
			return err
		}
	}
	target := fmt.Sprintf("%s/a.microservice.go", wd)
	fmt.Println(i.ImportPath)
	exists, err = Exists(target)
	if err != nil {
		return err
	}
	if !exists {
		err := os.WriteFile(target, []byte(_init), os.ModeAppend)
		if err != nil {
			return err
		}
	}
	fset := token.FileSet{}
	file, err := parser.ParseFile(&fset, target, nil, parser.ParseComments)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	var initFn *ast.FuncDecl
	var pkg string
	for _, decl := range file.Decls {
		switch t := decl.(type) {
		case *ast.GenDecl:
			{
				if t.Tok == token.IMPORT {
					path := fmt.Sprintf(`"%s"`, i.ImportPath)
					isAvailable := false
					for _, value := range t.Specs {
						if v, ok := value.(*ast.ImportSpec); ok {
							if v.Path.Value == path {
								if v.Name == nil {
									return fmt.Errorf("package name should be always provided")
								}
								pkg = v.Name.Name
								isAvailable = true
								break
							}
						}
					}
					if isAvailable {
						continue
					}
					pkg = fmt.Sprintf("%spkg", strings.ToLower(i.Service))
					imprt := ast.ImportSpec{}
					imprt.Name = &ast.Ident{}
					imprt.Name.Name = pkg
					imprt.Path = &ast.BasicLit{}
					imprt.Path.Kind = token.STRING
					imprt.Path.Value = fmt.Sprintf(`"%s"`, i.ImportPath)
					t.Specs = append(t.Specs, &imprt)
				}
			}
		case *ast.FuncDecl:
			{
				if t.Name.Name == "init" {
					initFn = t
					break
				}
			}
		}
	}
	if initFn == nil {
		panic("init function not found")
	}
	var microserviceArray *ast.CompositeLit
	for _, expr := range initFn.Body.List {
		switch t := expr.(type) {
		case *ast.AssignStmt:
			{
				if len(t.Rhs) > 1 {
					continue
				}
				if value, ok := t.Rhs[0].(*ast.CompositeLit); ok {
					microserviceArray = value
					break
				}
			}
		}
	}
	for _, value := range microserviceArray.Elts {
		if v, ok := value.(*ast.SelectorExpr); ok {
			left := v.X.(*ast.Ident)
			if left.Name == pkg {
				if v.Sel.Name == i.FunctionName {
					return nil
				}
			}
		}
	}
	element := ast.SelectorExpr{}
	x := &ast.Ident{}
	x.Name = pkg
	element.X = x
	selector := ast.Ident{}
	selector.Name = i.FunctionName
	element.Sel = &selector
	microserviceArray.Elts = append(microserviceArray.Elts, &element)
	buffer := bytes.NewBuffer([]byte{})
	printer.Fprint(buffer, &fset, file)
	err = os.Remove(target)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = os.WriteFile(target, buffer.Bytes(), os.ModePerm)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	exec.Command("gofmt", "-w", target).Run()
	return nil
}
