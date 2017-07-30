package main

import (
	"html/template"
	"fmt"
	"os"
)

/*
	导出go格式：
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

type goPackage struct {
	PackName string
	Messages []goMessage
	PackComments []string
}
type goMessage struct {
	MessageName string
	MessageComments []string
	MessageVariables []goVariable
}
type goVariable struct {
	VariableName string
	VariableType string
	VariableArray string
	VariableLineComment string
}

func exportGoFile(pack *messagePackage) {

	t := template.New("protocol.tpl")
	t = t.Funcs(template.FuncMap{"readVariable": ReadVariable})
	t = t.Funcs(template.FuncMap{"writeVariable": WriteVariable})

	t = template.Must(t.ParseFiles("protocol.tpl"))

	v1 := goVariable{"i1", "int32", "", "//comment i1"}
	v2 := goVariable{"af2", "float32", "[]", ""}

	msg1 := goMessage{MessageComments:[]string{"/* message comments", "cm1 */"}}
	msg1.MessageVariables = []goVariable{v1, v2}
	msg1.MessageName = "msg1"

	pack1 := goPackage{"test", []goMessage{msg1}, []string{"//1"}}

	t.Execute(os.Stdout, pack1)

	fmt.Println("finished.")
}
func ReadVariable(variableType string) string{
	switch variableType {
	case "int32":
		return "ReadInt()"
	case "float32":
		return "ReadFloat()"
	case "string":
		return "ReadString()"
	}
	return "error"
}
func WriteVariable(variableType string) string{
	switch variableType {
	case "int32":
		return "WriteInt"
	case "float32":
		return "WriteFloat"
	case "string":
		return "WriteString"
	}
	return "error"
}
