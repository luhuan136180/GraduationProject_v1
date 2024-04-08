package imsystem

import (
	"fmt"
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	UID  string

	C    chan string
	conn net.Conn // 对应的链接

	server *Server
}

func NewUser(conn net.Conn, server *Server, userUID string) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		UID:  userUID,

		C:    make(chan string),
		conn: conn,

		server: server,
	}

	// 启动监听当前user channel消息的goroutine
	go user.ListenMessage()

	return user
}

// 用户的上线业务
func (this *User) Online() {

	// 用户上线,将用户加入到onlineMap中
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

}

// 用户的下线业务
func (this *User) Offline() {

	// 用户下线,将用户从onlineMap中删除
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()
}

// 给当前User对应的客户端发送消息
func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
	fmt.Println("to " + this.UID + "---" + msg) // 测试
}

// 监听当前User channel的 方法,一旦有消息，就直接发送给对端客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C

		this.conn.Write([]byte(msg + "\n"))
	}
}

func (this *User) Domessage(msg string) {
	if msg == "quit" {
		this.Offline()
	} else if len(msg) > 4 && msg[:3] == "to|" {
		// 消息格式:  to|uid|消息内容

		remoteUID := strings.Split(msg, "|")[1]
		if remoteUID == "" {
			this.SendMsg("消息格式不正确，请使用 \"to|张三|你好啊\"格式。\n")
			return
		}
		// 2 根据用户名 得到对方User对象
		remoteUser, ok := this.server.OnlineMap[remoteUID]
		if !ok {
			this.SendMsg("该用户名不不存在\n")
			return
		}
		// 3 获取消息内容，通过对方的User对象将消息内容发送过去
		content := strings.Split(msg, "|")[2]
		if content == "" {
			this.SendMsg("无消息内容，请重发\n")
			return
		}
		remoteUser.SendMsg(this.UID + "对您说:" + content)
	}
}
