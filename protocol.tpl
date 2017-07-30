/*
    此文件为go-protocol的模板文件
*/

{{range .PackComments}}{{.}}{{end}}
package {{.PackName}}

{{range .Messages}}
{{range .MessageComments}}{{.}}{{end}}
type {{.MessageName}} struct {
    {{range .MessageVariables}}
    {{.VariableName}} {{.VariableArray}}{{.VariableType}} {{.VariableLineComment}}
    {{end}}
}
func (msg *{{.MessageName}}) ReadMsg(buffer *BinaryBuffer) {
    {{range .MessageVariables}}
    msg.{{.VariableName}} = buffer.{{readVariable .VariableType}}
    {{end}}
}
func (msg *{{.MessageName}}) WriteMsg(buffer *BinaryBuffer) {
    {{range .MessageVariables}}
    {{writeVariable .VariableType}}(msg.{{.VariableName}})
    {{end}}
}
{{end}}