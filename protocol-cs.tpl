/*
    此文件为自动生成的。
*/

{{range .PackComments}}{{.}}{{end}}
namespace Net
{

    // 协议映射关系
    public enum EPID
    {
            PID_None = 0,

        {{range $i, $v := .Messages}}
        PID_{{$v.MessageName}},
        {{end}}

            PID_MAX = 65535,
    }

    public static class BinaryProtocol
    {
        public static GbProtocol NewProtocol(EPID msgid)
        {
            switch (msgid)
            {
                {{range $i, $v := .Messages}}
                case EPID.PID_{{$v.MessageName}}:
                    return new {{$v.MessageName}}();
                {{end -}}
            }

            return null;
        }
        public static uint GetProtocolID(GbProtocol protocol)
        {
            var t = protocol.GetType();


    {{range $i, $v := .Messages}}
            if (t == typeof({{$v.MessageName}}))
            {
                return (uint)EPID.PID_{{$v.MessageName}};
            }
    {{end -}}

            return 0;
        }
    }

    {{range .Messages}}
    {{range .MessageComments -}}
    {{.}}
    {{end -}}
    public class {{.MessageName}} : GbProtocol
    {
        {{range .MessageVariables -}}
        public {{.VariableType}}{{.VariableArray}} {{.VariableName}}; {{.VariableLineComment}}
        {{end}}

        public void ReadMsg(BinaryBuffer buffer)
        {
            {{range .MessageVariables -}}
                {{$bCls := isClsType .VariableType -}}
                {{$barray := isArray .VariableArray -}}
                {{if $bCls -}}
                    {{if $barray}}
            {
                int size = buffer.ReadInt();
                {{.VariableName}} = new {{.VariableType}}[size];
                for (int i=0; i < size; i++)
                {
                    {{.VariableName}}[i] = new {{.VariableType}}();
                    {{.VariableName}}[i].ReadMsg(buffer);
                }
            }
                    {{else -}}
                        {{.VariableName}} = new {{.VariableType}}();
                        {{.VariableName}}.ReadMsg(buffer);
                    {{- end}}
                {{else -}}
                    {{.VariableName}} = buffer.{{readVariable .VariableType .VariableArray}};
                {{- end}}
            {{end}}
        }
        public void WriteMsg(BinaryBuffer buffer)
        {
            {{range .MessageVariables -}}
                {{$bCls := isClsType .VariableType -}}
                {{$barray := isArray .VariableArray -}}
                {{if $bCls -}}
                    {{if $barray}}
            {
                int size = {{.VariableName}}.Length;
                buffer.WriteInt(size);
                for(int i = 0; i < size; i++)
                {
                    {{.VariableName}}[i].WriteMsg(buffer);
                }
            }
                    {{else -}}
                        {{.VariableName}}.WriteMsg(buffer);
                    {{- end}}
                {{else -}}
                    buffer.{{writeVariable .VariableType .VariableArray}}({{.VariableName}});
                {{- end}}
            {{end}}
        }
    }
    {{end}}
}

