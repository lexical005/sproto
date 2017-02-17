package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/davyxu/pbmeta"
	pbprotos "github.com/davyxu/pbmeta/proto"
)

const codeTemplate = `# Generated by github.com/davyxu/gosproto/pb2sproto
# DO NOT EDIT!
# Source: {{.FileName}}

{{range .Structs}}
{{.Comment}}
.{{.Name}} {
	{{range .Fields}}	
	{{.Name}} {{.Tag}} : {{.TypeString}} {{.Comment}}
	{{end}}
}
{{end}}
`

type fieldModel struct {
	*pbmeta.FieldDescriptor
}

func (self *fieldModel) Tag() int32 {
	return self.Define.GetNumber()
}

func (self *fieldModel) TypeString() (ret string) {

	if self.IsRepeated() {
		ret = "*"
	}

	switch self.Type() {
	case pbprotos.FieldDescriptorProto_TYPE_INT64:
		ret += "int64"
	case pbprotos.FieldDescriptorProto_TYPE_UINT64:
		ret += "uint64"
	case pbprotos.FieldDescriptorProto_TYPE_INT32,
		pbprotos.FieldDescriptorProto_TYPE_FLOAT, // 浮点数默认转为int32
		pbprotos.FieldDescriptorProto_TYPE_DOUBLE:
		ret += "int32"
	case pbprotos.FieldDescriptorProto_TYPE_UINT32:
		ret += "uint32"
	case pbprotos.FieldDescriptorProto_TYPE_BOOL:
		ret += "bool"
	case pbprotos.FieldDescriptorProto_TYPE_STRING:
		ret += "string"
	case pbprotos.FieldDescriptorProto_TYPE_MESSAGE:
		ret += self.MessageDesc().Name()
	case pbprotos.FieldDescriptorProto_TYPE_ENUM:
		ret += "int32"
	}

	return
}

func addCommentSignAtEachLine(sign, comment string) string {

	if comment == "" {
		return ""
	}
	var out bytes.Buffer

	scanner := bufio.NewScanner(strings.NewReader(comment))

	var index int
	for scanner.Scan() {

		if index > 0 {
			out.WriteString("\n")
		}

		out.WriteString(sign)
		out.WriteString(" ")
		out.WriteString(scanner.Text())

		index++
	}

	return out.String()
}

func (self *fieldModel) Comment() string {

	return addCommentSignAtEachLine("#", self.CommentMeta.TrailingComment())

}

type structModel struct {
	*pbmeta.Descriptor

	Fields []fieldModel
}

func (self *structModel) Comment() string {

	return addCommentSignAtEachLine("#", self.CommentMeta.LeadingComment())
}

type protoFileModel struct {
	*pbmeta.FileDescriptor

	Structs []*structModel
}

func gen_proto(fileD *pbmeta.FileDescriptor, outputDir string) {

	tpl, err := template.New("sproto_to_go").Parse(codeTemplate)
	if err != nil {
		fmt.Println("template error ", err.Error())
		os.Exit(1)
	}

	pm := &protoFileModel{
		FileDescriptor: fileD,
	}

	for structIndex := 0; structIndex < fileD.MessageCount(); structIndex++ {
		st := fileD.Message(structIndex)

		stModel := &structModel{
			Descriptor: st,
		}

		for fieldIndex := 0; fieldIndex < st.FieldCount(); fieldIndex++ {
			fd := st.Field(fieldIndex)

			fdModel := fieldModel{
				FieldDescriptor: fd,
			}

			stModel.Fields = append(stModel.Fields, fdModel)
		}

		pm.Structs = append(pm.Structs, stModel)
	}

	var bf bytes.Buffer

	err = tpl.Execute(&bf, &pm)
	if err != nil {
		fmt.Println("template error ", err.Error())
		os.Exit(1)
	}

	if err != nil {
		fmt.Println("format error ", err.Error())
	}

	final := path.Join(outputDir, changeExt(fileD.FileName(), ".sp"))

	if fileErr := ioutil.WriteFile(final, bf.Bytes(), 666); fileErr != nil {
		fmt.Println("write file error ", fileErr.Error())
		os.Exit(1)
	}
}

// newExt = .xxx
func changeExt(name, newExt string) string {
	ext := path.Ext(name)
	name = name[0 : len(name)-len(ext)]
	return name + newExt
}
