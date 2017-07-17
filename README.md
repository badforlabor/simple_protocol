# simple_protocol
simple protocol

# 简单协议
- 目标是：简单的定义，应对复杂情况



### 协议文件
``` bash
package gamedef;

// ==========================================================
// 系统消息
// ==========================================================

// 一个连接接入
message SessionAccepted
{
	
}

// 已连接
message SessionConnected
{
	
}

// ==========================================================
// 测试用消息
// ==========================================================
message TestEchoACK
{
	string Content;
	string Content2;
}

message NestedMessage
{
 TestEchoACK Acks;
 int id;
 float f;
 string str;
}

message ArrayMessage
{
 repeated int datas;
 repeated string msgs;
 repeated NestedMessage p;
 int id;
 float f;
 string str;
}



```
