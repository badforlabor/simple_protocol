package main

import (
	"fmt"
	"strings"
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
	errComemnt = "comment syntax error."

	runeSpace      = ' '
	runeTAB        = '\t'
	runeReturn     = '\r'
	runeReturnLine = '\n'
	runeSemicolon  = ';'
)

type Parser struct {
	scanner *strings.Reader
}

type variable struct {
	name         string
	variableType string
	array        bool
	comments     []string
	lineComment  string
	parsed       bool
}

type variableClass struct {
	name      string
	variables []variable
	comments  []string
	parsed    bool
}
type messagePackage struct {
	name     string
	classes  []variableClass
	comments []string
}

func NewParser(str string) *Parser {
	p := &Parser{scanner: strings.NewReader(str)}
	return p
}

func (p *Parser) DoParse() error {

	fmt.Println("do parse.")

	s := p.readComment()
	fmt.Print(s)

	return nil
}
func (p *Parser) unread(n int) {
	for n > 0 {
		p.scanner.UnreadRune()
		n--
	}
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

	backup := p.scanner.Size() - int64(p.scanner.Len())

	defer func() {
		if backup >= 0 {
			cur := p.scanner.Size() - int64(p.scanner.Len())
			for cur > backup {
				p.scanner.UnreadByte()
				cur--
			}
		}
	}()

	p.readAllSpace()
	p.readReturn()

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
		backup = -1
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

		backup = -1
		return comment

	} else if s == '*' {
		backup = -1
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
