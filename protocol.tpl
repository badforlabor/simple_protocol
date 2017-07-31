/*
    此文件为go-protocol的模板文件
*/

{{range .PackComments}}{{.}}{{end}}
package {{.PackName}}

{{range .Messages}}
{{range .MessageComments}}{{.}}{{end}}
type {{.MessageName}} struct {
    {{range .MessageVariables -}}
    {{.VariableName}} {{.VariableArray}}{{.VariableType}} {{.VariableLineComment}}
    {{end}}
}
func (msg *{{.MessageName}}) ReadMsg(buffer *BinaryBuffer) {
    {{range .MessageVariables -}}
        {{$bCls := isClsType .VariableType -}}
        {{if $bCls -}}
            msg.{{.VariableName}}.ReadMsg(buffer)
        {{else -}}
            msg.{{.VariableName}} = buffer.{{readVariable .VariableType .VariableArray}}
        {{end}}
    {{end}}
}
func (msg *{{.MessageName}}) WriteMsg(buffer *BinaryBuffer) {
    {{range .MessageVariables -}}
        {{$bCls := isClsType .VariableType -}}
        {{if $bCls -}}
            msg.{{.VariableName}}.WriteMsg(buffer)
        {{else -}}
            {{writeVariable .VariableType .VariableArray}}(msg.{{.VariableName}})
        {{- end}}
    {{end}}
}
{{end}}