/*
    此文件为自动生成的。
*/

{{range .PackComments}}{{.}}{{end}}
package main//{{.PackName}}

type BinaryProtocol interface {
    ReadMsg(buffer * BinaryBuffer)
    WriteMsg(buffer * BinaryBuffer)
}

// 协议映射关系
const (
    {{range $i, $v := .Messages}}
    PID_{{$v.MessageName}} = 1 + {{$i}}
    {{end}}
)
func NewProtocol(msgid uint32) (interface{}) {
    switch(msgid) {
    {{range $i, $v := .Messages}}
    case PID_{{$v.MessageName}}:
        return &{{$v.MessageName}}{}
    {{end -}}
    }
    return nil
}
func GetProtocolID(proto interface{}) uint32 {

    switch proto.(type) {
    {{range $i, $v := .Messages}}
    case *{{$v.MessageName}}:
        return PID_{{$v.MessageName}}
    {{end -}}
    }

    return 0
}


{{range .Messages}}
{{range .MessageComments -}}
{{.}}
{{end -}}
type {{.MessageName}} struct {
    {{range .MessageVariables -}}
    {{.VariableName}} {{.VariableArray}}{{.VariableType}} {{.VariableLineComment}}
    {{end}}
}
func (msg *{{.MessageName}}) ReadMsg(buffer *BinaryBuffer) {
    {{range .MessageVariables -}}
        {{$bCls := isClsType .VariableType -}}
        {{$barray := isArray .VariableArray -}}
        {{if $bCls -}}
            {{if $barray}}
    {
        var size int32 = buffer.ReadInt()
        msg.{{.VariableName}} = make([]{{.VariableType}}, size)
        for i:=int32(0); i < size; i++ {
            msg.{{.VariableName}}[i].ReadMsg(buffer)
        }
    }
            {{else -}}
                msg.{{.VariableName}}.ReadMsg(buffer)
            {{- end}}
        {{else -}}
            msg.{{.VariableName}} = buffer.{{readVariable .VariableType .VariableArray}}
        {{- end}}
    {{end}}
}
func (msg *{{.MessageName}}) WriteMsg(buffer *BinaryBuffer) {
    {{range .MessageVariables -}}
        {{$bCls := isClsType .VariableType -}}
        {{$barray := isArray .VariableArray -}}
        {{if $bCls -}}
            {{if $barray}}
    {
        var size int32 = int32(len(msg.{{.VariableName}}))
        buffer.WriteInt(size)
        for i:=int32(0); i < size; i++ {
            msg.{{.VariableName}}[i].WriteMsg(buffer)
        }
    }
            {{else -}}
                msg.{{.VariableName}}.WriteMsg(buffer)
            {{- end}}
        {{else -}}
            buffer.{{writeVariable .VariableType .VariableArray}}(msg.{{.VariableName}})
        {{- end}}
    {{end}}
}
{{end}}