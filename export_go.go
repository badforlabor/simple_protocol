package main

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

func exportGoFile(pack *messagePackage) {

}