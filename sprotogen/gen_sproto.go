package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/davyxu/gosproto/meta"
)

const spCodeTemplate = `// Generated by github.com/davyxu/gosproto/sprotogen
// DO NOT EDIT!{{range .Files}}
// Source: {{.FileName}}{{end}}

{{range .Files}}
{{range .Structs}}
.{{.Name}} {
	{{range .Fields}}	
	{{.Name}} {{.Tag}} : {{.CompatibleTypeString}}
	{{end}}
}
{{end}}{{end}}

`

type spFileModel struct {
	Files []*meta.FileDescriptor
}

func gen_sproto(pool []*meta.FileDescriptor, filename string) {

	tpl, err := template.New("sproto_go").Parse(spCodeTemplate)
	if err != nil {
		fmt.Println("template error ", err.Error())
		os.Exit(1)
	}

	var bf bytes.Buffer

	err = tpl.Execute(&bf, &spFileModel{
		Files: pool,
	})
	if err != nil {
		fmt.Println("template error ", err.Error())
		os.Exit(1)
	}

	if fileErr := ioutil.WriteFile(filename, bf.Bytes(), 666); fileErr != nil {
		fmt.Println("write file error ", fileErr.Error())
		os.Exit(1)
	}
}
