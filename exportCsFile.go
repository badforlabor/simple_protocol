package main

import (
	"fmt"
	"text/template"
	"os"
	"bytes"
	"io/ioutil"
)

/*
	导出csharp格式：
		type MessageA struct {
			a int
			b float
			c string
		}
		func (msg *MessageA) Read(r io.Reader) {
			a = r.ReadInt()
			b = r.ReadFloat()
			c = r.ReadString()
		}
		func (msg* MessageA) Write(w io.Writer) {
			r.WriteInt(a)
			r.WriteFloat(b)
			r.WriteString(c)
		}


*/


func exportCsFile(pack *messagePackage) {

	t := template.New("protocol-cs.tpl")
	t = t.Funcs(template.FuncMap{"readVariable": ReadVariable})
	t = t.Funcs(template.FuncMap{"writeVariable": WriteVariable})
	t = t.Funcs(template.FuncMap{"isClsType": isClsType})
	t = t.Funcs(template.FuncMap{"isArray": isArray})

	t = template.Must(t.ParseFiles("protocol-cs.tpl"))

	pack1 := goPackage{}

	if pack != nil {

		pack1.PackName = pack.Name
		pack1.PackComments = pack.Comments

		for _, v := range pack.Classes {
			msg := goMessage{MessageName: v.Name, MessageComments: v.Comments}
			for _, vv := range v.Variables {
				variable := goVariable{vv.Name, getCharpTypeString(vv.VariableType), getArrayString(vv.Array), vv.LineComment}
				msg.MessageVariables = append(msg.MessageVariables, variable)
			}
			pack1.Messages = append(pack1.Messages, msg)
		}
		buffer := bytes.NewBufferString("")
		t.Execute(buffer, pack1)
		ioutil.WriteFile("binary_proto.cs", buffer.Bytes(), os.ModePerm)
	} else {

		// 调试
		v1 := goVariable{"i1", "int32", "", "//comment i1"}
		v2 := goVariable{"af2", "float32", "[]", ""}

		msg1 := goMessage{MessageComments: []string{"/* message comments", "cm1 */"}}
		msg1.MessageVariables = []goVariable{v1, v2}
		msg1.MessageName = "msg1"

		pack1 = goPackage{"test", []goMessage{msg1}, []string{"//1"}}

		t.Execute(os.Stdout, pack1)
	}


	fmt.Println("finished.")
}
func getCharpTypeString(variableType string) string {
	switch variableType {
	case "int":
		fallthrough
	case "int32":
		return "int"
	case "float":
		fallthrough
	case "float32":
		return "float"
	}
	return variableType
}