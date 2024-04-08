package imsystem

import (
	"flag"
	"fmt"
	"io"
	"net"
	"sync"
)

var ChatServerIp string
var ChatServerPort int

func init() {
	flag.StringVar(&ChatServerIp, "ip", "127.0.0.1", "设置服务器IP地址(默认是127.0.0.1)")
	flag.IntVar(&ChatServerPort, "port", 8888, "设置服务器端口(默认是8888)")
}

type Server struct {
	Ip   string
	Port int

	// 在线用户的列表
	OnlineMap map[string]*User // addr: uid
	mapLock   sync.RWMutex

	// 消息广播的channel
	Message chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}

	return server
}

func (this *Server) Handler(conn net.Conn) {
	// ...当前链接的业务
	// fmt.Println("链接建立成功")
	startMsg := make([]byte, 4096)
	n, err := conn.Read(startMsg)
	if n == 0 || err != nil {
		return
	}
	startmsg := string(startMsg[:n])

	var UserUID string
	if len(startmsg) > 4 && startmsg[:4] == "uid:" {
		UserUID = startmsg[4:]
	} else {
		return
	}
	// 可以在服务端新增聊天用户
	user := NewUser(conn, this, UserUID)

	user.Online()
	wg := sync.WaitGroup{}
	// 接受客户端发送的消息
	wg.Add(1)
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}

			// 提取用户的消息(去除'\n')
			msg := string(buf[:n])

			// 用户针对msg进行消息处理
			user.Domessage(msg)
			if msg == "quit" {
				wg.Done()
			}
		}
	}()

	// 当前handler阻塞
	wg.Wait()
}

// 启动服务器的接口
func (this *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	// close listen socket
	defer listener.Close()

	for {
		// accept-- 监听链接
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}

		// do handler
		go this.Handler(conn)
	}
}

// 初始化时调用
func InitChatServer(ip string, port int) *Server {
	server := NewServer(ip, port)
	go server.Start()
	fmt.Println("启动聊天服务...")
	return server
}
