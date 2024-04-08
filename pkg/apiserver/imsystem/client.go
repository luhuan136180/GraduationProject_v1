package imsystem

import (
	"errors"
	"fmt"
	"io"
	"net"
)

var ClientMap map[string]*Client // uid:client{}

type ChatFlag int

const (
	ChatFlagPrivate ChatFlag = 1 // 单对单
	ChatFlagGroup   ChatFlag = 2
	ChatFlagPublic  ChatFlag = 3
)

type Client struct {
	ServerIp   string
	ServerPort int
	Adr        string

	UserUID string

	conn net.Conn
	flag int // 当前client的模式
}

func NewClient(serverIp string, serverPort int, uid string) (*Client, error) {
	// 创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       1,
		UserUID:    uid,
	}

	// 链接server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil, err
	}

	client.conn = conn
	client.Adr = conn.RemoteAddr().String()

	sendMsg := "uid:" + client.UserUID
	_, err = client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn Write err:", err)
	}

	// 返回对象
	return client, err
}

// 用于新增用户client
func AddChatClient(serverIp string, serverPort int, uid string) error {
	client, err := NewClient(serverIp, serverPort, uid)
	if client == nil {
		return err
	}

	// 监听对端传来的消息
	go client.listenHandler()

	ClientMap[uid] = client
	return nil
}

// 监听
func (client *Client) listenHandler() {
	buf := make([]byte, 4096)
	for {
		n, err := client.conn.Read(buf)

		if err != nil && err != io.EOF {
			fmt.Println("Conn Read err:", err)
			return
		}
		// 提取用户的消息(去除'\n')
		msg := string(buf[:n-1])

		// 用户针对msg进行消息处理
		client.Domessage(msg)
	}
}

func (client *Client) Domessage(msg string) {
	fmt.Println(client.UserUID + " get a message:" + msg) // 测试
}

// 删除client --- 仅在退登时使用
func DeleteClient(uid string) {
	client, ok := ClientMap[uid]
	if ok {
		sendMsg := "quit"
		_, err := client.conn.Write([]byte(sendMsg))
		if err != nil {
			fmt.Println("conn Write err:", err)
		}
	}

	delete(ClientMap, uid)

}

// 初始化
func InitChatClient() {
	ClientMap = map[string]*Client{}
}

func SendMsg(sendUID string, GetUID string, msg string) error {
	sendClient, ok := ClientMap[sendUID]
	if !ok {
		return errors.New("not get client")
	}
	sendClient.sendMsg(GetUID, msg)
	return nil
}

func (client *Client) sendMsg(getUID, msg string) {
	sendMsg := "to|" + getUID + "|" + msg + "\n\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn Write err:", err)
		return
	}
}
