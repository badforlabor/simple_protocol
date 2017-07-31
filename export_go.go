package main

import (
	"fmt"
	"html/template"
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
	PackName     string
	Messages     []goMessage
	PackComments []string
}
type goMessage struct {
	MessageName      string
	MessageComments  []string
	MessageVariables []goVariable
}
type goVariable struct {
	VariableName        string
	VariableType        string
	VariableArray       string
	VariableLineComment string
}

func exportGoFile(pack *messagePackage) {

	t := template.New("protocol.tpl")
	t = t.Funcs(template.FuncMap{"readVariable": ReadVariable})
	t = t.Funcs(template.FuncMap{"writeVariable": WriteVariable})
	t = t.Funcs(template.FuncMap{"isClsType": isClsType})

	t = template.Must(t.ParseFiles("protocol.tpl"))

	pack1 := goPackage{}

	if pack != nil {

		pack1.PackName = pack.Name
		pack1.PackComments = pack.Comments

		for _, v := range pack.Classes {
			msg := goMessage{MessageName: v.Name, MessageComments: v.Comments}
			for _, vv := range v.Variables {
				variable := goVariable{vv.Name, getTypeString(vv.VariableType), getArrayString(vv.Array), vv.LineComment}
				msg.MessageVariables = append(msg.MessageVariables, variable)
			}
			pack1.Messages = append(pack1.Messages, msg)
		}

	} else {

		// 调试
		v1 := goVariable{"i1", "int32", "", "//comment i1"}
		v2 := goVariable{"af2", "float32", "[]", ""}

		msg1 := goMessage{MessageComments: []string{"/* message comments", "cm1 */"}}
		msg1.MessageVariables = []goVariable{v1, v2}
		msg1.MessageName = "msg1"

		pack1 = goPackage{"test", []goMessage{msg1}, []string{"//1"}}
	}

	t.Execute(os.Stdout, pack1)

	fmt.Println("finished.")
}
func getArrayString(barray bool) string {
	if barray {
		return "[]"
	} else {
		return ""
	}
}
func getTypeString(variableType string) string {
	switch variableType {
	case "int":
		fallthrough
	case "int32":
		return "int32"
	case "float":
		fallthrough
	case "float32":
		return "float32"
	}
	return variableType
}
func ReadVariable(variableType string, variableArray string) string {

	barray := isArray(variableArray)

	switch variableType {
	case "int32":
		if barray {
			return "ReadIntArray()"
		} else {
			return "ReadInt()"
		}
	case "float32":
		if barray {
			return "ReadFloatArray()"
		} else {
			return "ReadFloat()"
		}
	case "string":
		if barray {
			panic("not support string array")
		}
		return "ReadString()"
	}
	panic("not support type")
}
func WriteVariable(variableType string, variableArray string) string {

	barray := isArray(variableArray)

	switch variableType {
	case "int32":
		if barray {
			return "WriteIntArray"
		} else {
			return "WriteInt"
		}
	case "float32":
		if barray {
			return "WriteFloatArray"
		} else {
			return "WriteFloat"
		}
	case "string":
		if barray {
			panic("not support string array")
		}
		return "WriteString"
	}
	panic("not support type")
}
func isArray(variableArray string) bool {
	return variableArray == "[]"
}
func isClsType(variableType string) bool {
	switch variableType {
	case "int32":
		return false
	case "float32":
		return false
	case "string":
		return false
	}
	return true
}
