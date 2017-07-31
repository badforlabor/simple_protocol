package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"unicode"
)

/*
	解析proto文件：

	支持文件
	根据io.Read流
		读取//，直接将此行内容当成注释
		读取/*，并查找*\/，找不到报错，找到了，那么把此范围内的当成日志
		读取token（字母开头或者下划线开头的，只包含字母数组和下划线的单词）
			检查是否是关键字
				package: 那么读取xxx的内容
				message：那么读取 xxx
				repeated：记录
				int,float等，作为变量类型
				剩下的作为变量。
	parse过程中，需要记录：
		包名，类名，类的注释，变量名，变量类型，变量注释。
		package
			-- class1
				-- variable1
				-- variable2
			-- class2
				-- variable1
				-- variable2

	如果解析的过程中，出现非法的定义，那么抛异常，并显示第N行出错，并且详细说明：
			包名错误
			消息名错误
			缺少{}
			缺少类型
			变量名错误
			语法错误，比如 int int a;
			包含了多个包名

	伪代码：
		工具函数：
			跳过空格
			跳过回车
			跳过；号
			读取注释：
				字符串以//开头，或者/*开头的。
				空白行，也作为注释；这样输出文件就非常好看了。
			读取token：
				以字母或者下划线开头的单词。
			抛异常：
				记录第N行
			解析完毕后：
				输出json
		数据结构：
			一个proto文件中，只能定义一个package。此package中可以包含N个class，每个class中可以包含N个variable
			// 变量
			class variable
			{
				string name;
				bool array;
				string type;
				string[] comments;
				string linecomment;		// 变量后面允许有单行注释
				bool parsed;			// 是否解析完毕
			}
			// 类
			class Tclass
			{
				string name
				variable[] v;
				string[] comments;
				bool parsed;			// 是否解析完毕
			}
			// 包
			class package
			{
				string name;
				Tclass[] c;
				string[] comments;
			}


		整个流程：
		while true
			while break
				跳过空格、回车（把;当成换行）
				读取注释
			读取token
				如果不是token，
					那么提示错误：非法
				如果是token，那么switch
					如果是package
						如果已经设置过了包名，那么提示：包含了多个包名
						while break
							读取注释
							跳过空格、回车
						读取token
							如果不是token，那么提示：包名错误
							如果是token，那么设置包名；继续读取字符并跳过空格
								如果读到；或者换行，那么语句结束，over
								如果是其他，那么报错：包名后需要有换行
					如果是message
						while break
							读取注释
							跳过空格、回车
						读取token
							如果不是token，那么提示类名错误
							如果是token，那么设置类名，
								继续读取
								while break
									读取注释
									跳过空格、回车
								如果读到{，那么跳转到process message函数
								如果读到其他，那么报错，提示：缺少{


		process message 流程：
		while true
			while break
				读取注释
				跳过空格、回车
			如果是}
				检查变量定义是否都成功了，如果没有，那么提示：语法错误。
				那么标记消息定义结束。
			读取token
				如果不是token，那么，提示非法
				如果是token，那么switch
					如果是repeated
						检查数据类型是否设置，如果设置了，那么提示语法错误。
						检查是否已经是数组，如果是，那么提示语法错误。
						标记自身是数组
					如果是int,float,string
						检查数据类型是否设置，如果设置了，那么提示语法错误。
						标记数据类型
					如果是其他
						检查类型是否设置
							如果没设置，那么把此token当成类型
							如果设置了，那么：
								检查变量是否设置了，如果设置了，那么提示语法错误
								检查数据类型是否设置了，如果设置了，那么提示语法错误。
								设置为是变量。
								while break
									跳过空格
								读取到回车或者分号时，设置该变量定义完整
									如果是分号，那么继续读取此行后面的内容，如果是注释，那么认为是linecomment。
									然后，开启下一个变量定义

*/

const (
	errComemnt        = "comment syntax error."
	errPackage        = "package name is invalid."
	errInvalidLine    = "line end with \r, not \n"
	errInvalidLineEnd = "line end with invalid character"
	errClass          = "class name is invalid."
	errClassBody      = "class body define failed"
	errVariableBody   = "variable body is invalid"

	runeSpace      = ' '
	runeTAB        = '\t'
	runeReturn     = '\r'
	runeReturnLine = '\n'
	runeSemicolon  = ';'
	runeUnderline  = '_'

	// 关键字
	keyPackage  = "package"
	keyClass    = "message"
	keyInt      = "int"
	keyString   = "string"
	keyFloat    = "float"
	keyRepeated = "repeated"
)

type Parser struct {
	scanner *strings.Reader
	pack    *messagePackage
}

type variable struct {
	Name         string   `json:"name"`
	VariableType string   `json:"variableType"`
	Array        bool     `json:"array"`
	Comments     []string `json:"comments"`
	LineComment  string   `json:"lineComment"`
	parsed       bool
}

type variableClass struct {
	Name      string     `json:"name"`
	Variables []variable `json:"variables"`
	Comments  []string   `json:"comments"`
	parsed    bool
}
type messagePackage struct {
	Name     string          `json:"name"`
	Classes  []variableClass `json:"classes"`
	Comments []string        `json:"comments"`
}

func NewParser(str string) *Parser {
	p := &Parser{scanner: strings.NewReader(str)}
	return p
}

func (p *Parser) DoParse() error {

	fmt.Println("do parse.")

	p.parsePackage()

	data, err := json.Marshal(p.pack)
	if err == nil {
		fmt.Println(string(data))
	}
	ioutil.WriteFile("out.json", data, os.ModePerm)

	exportGoFile(p.pack)

	return nil
}
func (p *Parser) unread(n int) {
	for n > 0 {
		p.scanner.UnreadRune()
		n--
	}
}
func (p *Parser) readRune() rune {
	c, _, err := p.scanner.ReadRune()

	// 如果发生错误了，那么应该是到文件末尾了。
	if err != nil {
		c = 0
	}
	return c
}

// 读取空格，直到不是空格为止
func (p *Parser) readAllSpace() int {
	n := 0
	for {
		s, _, err := p.scanner.ReadRune()
		if err != nil {
			break
		} else {
			if s == runeSpace || s == runeTAB {
				n++
			} else {
				p.scanner.UnreadRune()
				break
			}
		}
	}
	return n
}

// 读取换行，直到不是换行为止
func (p *Parser) readReturn() int {
	n := 0
	for {
		s, _, err := p.scanner.ReadRune()
		if err != nil {
			break
		} else {

			if s == runeReturn {
				s, _, err = p.scanner.ReadRune()
				if err != nil || s != runeReturnLine {
					p.scanner.UnreadRune()
					p.scanner.UnreadRune()
					break
				}
			}

			if s == runeReturnLine || s == runeSemicolon {
				n++
			} else if s == runeSpace || s == runeTAB {
				space := p.readAllSpace()
				if space > 0 {
					// 如果空格后面是换行符，那么说明是有效的空行
					s, _, err = p.scanner.ReadRune()

					if s == runeReturn {
						s, _, err = p.scanner.ReadRune()
						if err != nil || s != runeReturnLine {
							// 无效的空行，所以得回退咯
							p.scanner.UnreadRune()
							p.unread(space + 1)
							break
						}
					}
					if s == runeReturnLine || s == runeSemicolon {
						n++
					} else {
						// 无效的空行，所以得回退咯
						p.unread(space + 1)
						break
					}
				} else {
					break
				}
			} else {
				p.scanner.UnreadRune()
				break
			}
		}
	}
	return n
}

func (p *Parser) readComment() string {

	s, _, err := p.scanner.ReadRune()
	if err != nil {
		return ""
	}
	if s != '/' {
		p.scanner.UnreadRune()
		return ""
	}

	s, _, err = p.scanner.ReadRune()
	if err != nil {
		p.scanner.UnreadRune()
		return ""
	}

	if s == '/' {
		// 整行内容为注释
		comment := "//"
		for {
			s, _, err = p.scanner.ReadRune()
			if err != nil {
				break
			}
			if s == runeReturnLine {
				break
			}
			comment = strings.Join([]string{comment, string(s)}, "")
		}

		return comment

	} else if s == '*' {
		// 接下来的内容为注释
		comment := "/*"
		last := ' '
		parsed := false
		for {
			s, _, err = p.scanner.ReadRune()
			if err != nil {
				break
			}

			comment = strings.Join([]string{comment, string(s)}, "")

			if last == '*' && s == '/' {
				parsed = true
				break
			}
			last = s
		}

		// 没有找到对应的‘*/’，所以直接抛异常
		if !parsed {
			panic(errComemnt)
		}

		return comment
	} else {
		p.scanner.UnreadRune()
		p.scanner.UnreadRune()
		return ""
	}
}

// 读取行注释，形如这种的：  TestEchoACK Acks;      // 嵌套语句
func (p *Parser) readLineComment() string {

	p.readAllSpace()

	c := p.readRune()
	if c == runeSemicolon {
		c = p.readRune()
	}

	if c == runeReturnLine {
		return ""
	}

	if c == runeReturn {
		c = p.readRune()
		if c == runeReturnLine {
			return ""
		}
		panic(errInvalidLine)
	}

	if c == '/' {
		c = p.readRune()
		if c == '/' {
			comment := "//"

			for {
				c = p.readRune()

				if c == runeReturn {
					c = p.readRune()
					if c == runeReturnLine {
						break
					} else {
						panic(errInvalidLineEnd)
					}
				}

				if c == runeReturnLine {
					break
				}

				comment += string(c)
			}

			return comment

		} else if c == '*' {
			// 一行结束后，不支持/**/注释，只支持//
		}
	}

	// 如果读到了文件末尾，那么正常返回
	if c == 0 {
		return ""
	}

	// 一行结束，要么是';'，要么是'\n\，要么是\\注释，其他的，均非法
	panic(errInvalidLineEnd)
}

// 循环读取空格、回车、注释，直到遇到token为止
func (p *Parser) whileReadAllLineReturnAndComments() []string {
	comments := []string{}

	for {
		cnt := 0
		cnt += p.readAllSpace()
		cnt += p.readReturn()
		c := p.readComment()
		if cnt == 0 && c == "" {
			break
		}
		if c != "" {
			comments = append(comments, c)
		}
	}

	return comments
}

// 读取token
func (p *Parser) readToken() string {
	token := ""

	for {
		c, _, err := p.scanner.ReadRune()
		if err != nil {
			break
		}

		if c == '{' || c == '}' {
			p.scanner.UnreadRune()
			break
		}

		if c == runeSpace || c == runeTAB || c == runeReturn || c == runeReturnLine || c == runeSemicolon {
			break
		}

		token = token + string(c)
	}

	return token
}

func isNumber(c rune) bool {
	return c < unicode.MaxASCII && c >= '0' && c < '9'
}
func isLetter(c rune) bool {
	return c < unicode.MaxASCII && ((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z'))
}

// 变量名必须以字母或者下划线开头，且只能包含字母、下划线、数字
func isVariableName(s string) bool {
	if len(s) == 0 {
		return false
	}

	// 获取第一个字符
	var c rune
	for _, c = range s {
		break
	}

	if c == runeUnderline || isLetter(c) {

	} else {
		return false
	}

	for _, c = range s {
		if c == runeUnderline || isLetter(c) || isNumber(c) {

		} else {
			return false
		}
	}

	return true
}

// 包名必须
func isPackageName(s string) bool {
	return isVariableName(s)
}

// 类名
func isClassName(s string) bool {
	return isVariableName(s)
}

func (p *Parser) parsePackage() {

	comments := p.whileReadAllLineReturnAndComments()
	for _, s := range comments {
		fmt.Println(s)
	}

	// 第一个一定是包名，格式是： package pkgname;

	keyPkg := p.readToken()
	fmt.Println(keyPkg)

	if keyPkg != keyPackage {
		panic(errPackage)
	}

	p.readAllSpace()

	pkgname := p.readToken()
	if !isPackageName(pkgname) {
		panic(errPackage)
	}

	// package后面的行注释，直接忽略掉，不保留
	p.readLineComment()

	p.pack = &messagePackage{}
	p.pack.Comments = comments
	p.pack.Name = pkgname

	// 接下来就是一系列消息体
	p.parseAllMessage()
}

func (p *Parser) parseAllMessage() {

	comments := p.whileReadAllLineReturnAndComments()

	// 读到了文件末尾
	if p.readRune() == 0 {
		fmt.Println("parse end.")
		return
	}
	p.unread(1)

	keyclass := p.readToken()
	if keyclass != keyClass {
		panic(errClass)
	}

	p.readAllSpace()
	classname := p.readToken()
	if !isClassName(classname) {
		panic(errClass)
	}

	var cls variableClass
	cls.Comments = comments
	cls.Name = classname

	p.parseClassBody(&cls)

	if !cls.parsed {
		panic(errClassBody)
	}

	if p.pack.Classes == nil {
		p.pack.Classes = []variableClass{}
	}
	p.pack.Classes = append(p.pack.Classes, cls)

	// 继续读取消息体
	p.parseAllMessage()
}

// 解析整个类
func (p *Parser) parseClassBody(cls *variableClass) {

	p.whileReadAllLineReturnAndComments()

	// 以 { 开头
	c := p.readRune()
	if c != '{' {
		panic(errClassBody)
	}

	// 类定义体中，有没有变量，都可以
	p.parseVariables(cls)

	// 以 } 结尾
	c = p.readRune()
	if c != '}' {
		panic(errClassBody)
	}
	cls.parsed = true

	// } 后面要么是注释，要么是换行符
	p.readLineComment()
}

// 读取类中的所有行变量
func (p *Parser) parseVariables(cls *variableClass) {

	p.whileReadAllLineReturnAndComments()
	token := p.readToken()

	if token != "" {

		v := variable{}
		p.parseVariableBody(&v, token)
		if !v.parsed {
			panic(errVariableBody)
		}

		if cls.Variables == nil {
			cls.Variables = []variable{}
		}
		// 成功读取到一个变量
		cls.Variables = append(cls.Variables, v)

		// 继续读取下一个变量
		p.parseVariables(cls)
	}

}

// 读取一行变量
func (p *Parser) parseVariableBody(v *variable, token string) {

	switch token {
	case "repeated":
		if v.Array {
			panic(errVariableBody)
		}
		v.Array = true

	case "int":
		fallthrough
	case "string":
		fallthrough
	case "float":
		if v.VariableType != "" {
			panic(errVariableBody)
		}
		v.VariableType = token

	default:
		if v.VariableType == "" {
			v.VariableType = token
		} else if v.Name == "" {
			v.Name = token
			// 变量后面可能有注释
			linecomment := p.readLineComment()
			v.LineComment = linecomment
			v.parsed = true
		}
	}

	// 如果一个变量体没有定义完全，那么继续parse
	if !v.parsed {
		p.readAllSpace()
		token = p.readToken()
		p.parseVariableBody(v, token)
	}
}
