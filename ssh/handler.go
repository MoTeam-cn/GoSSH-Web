package ssh

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

// SSHClient 结构体用于存储SSH会话信息
type SSHClient struct {
	Client  *ssh.Client
	Session *ssh.Session
	Stdin   io.WriteCloser
	Stdout  io.Reader
	Stderr  io.Reader
	WsConn  *websocket.Conn
}

var (
	// 用于存储活动的SSH连接
	activeConnections = make(map[string]*SSHClient)
	connectionMutex   sync.RWMutex
)

// AuthRequest SSH认证请求结构
type AuthRequest struct {
	Host     string `json:"host" binding:"required"`
	Port     int    `json:"port" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// HandleSSHAuth 处理SSH认证
func HandleSSHAuth(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数: " + err.Error()})
		return
	}

	config := &ssh.ClientConfig{
		User: req.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(req.Password),
			ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) ([]string, error) {
				answers := make([]string, len(questions))
				for i := range questions {
					answers[i] = req.Password
				}
				return answers, nil
			}),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         15 * time.Second,
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", req.Host, req.Port), config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "SSH连接失败: " + err.Error()})
		return
	}

	// 生成唯一的会话ID
	sessionID := uuid.New().String()

	connectionMutex.Lock()
	activeConnections[sessionID] = &SSHClient{
		Client: client,
	}
	connectionMutex.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"message": "认证成功",
		"id":      sessionID,
	})
}

// HandleSSHSession 处理SSH WebSocket会话
func HandleSSHSession(id string, ws *websocket.Conn) {
	connectionMutex.RLock()
	sshClient, exists := activeConnections[id]
	connectionMutex.RUnlock()

	if !exists {
		ws.WriteJSON(gin.H{"error": "会话未找到"})
		return
	}

	session, err := sshClient.Client.NewSession()
	if err != nil {
		ws.WriteJSON(gin.H{"error": "创建会话失败: " + err.Error()})
		return
	}
	defer session.Close()

	// 请求伪终端，设置更合适的终端类型和大小
	if err = session.RequestPty("xterm-256color", 40, 80, ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
		ssh.ICRNL:         1,
		ssh.IEXTEN:        1,
		ssh.IGNCR:         0,
		ssh.IGNPAR:        0,
		ssh.IMAXBEL:       1,
		ssh.INLCR:         0,
		ssh.ISTRIP:        0,
		ssh.IUCLC:         0,
		ssh.IXANY:         1,
		ssh.IXOFF:         0,
		ssh.IXON:          1,
		ssh.OCRNL:         0,
		ssh.OLCUC:         0,
		ssh.ONLCR:         1,
		ssh.ONLRET:        0,
		ssh.ONOCR:         0,
		ssh.OPOST:         1,
		ssh.VINTR:         3,   // Ctrl-C
		ssh.VQUIT:         28,  // Ctrl-\
		ssh.VERASE:        127, // Backspace
		ssh.VKILL:         21,  // Ctrl-U
		ssh.VEOF:          4,   // Ctrl-D
		ssh.VEOL:          0,
		ssh.VEOL2:         0,
		ssh.VSTART:        17, // Ctrl-Q
		ssh.VSTOP:         19, // Ctrl-S
		ssh.VSUSP:         26, // Ctrl-Z
		ssh.VWERASE:       23, // Ctrl-W
	}); err != nil {
		ws.WriteJSON(gin.H{"error": "请求PTY失败: " + err.Error()})
		return
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		ws.WriteJSON(gin.H{"error": "获取STDIN失败: " + err.Error()})
		return
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		ws.WriteJSON(gin.H{"error": "获取STDOUT失败: " + err.Error()})
		return
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		ws.WriteJSON(gin.H{"error": "获取STDERR失败: " + err.Error()})
		return
	}

	// 启动远程shell，使用默认的shell或bash
	if err = session.Shell(); err != nil {
		ws.WriteJSON(gin.H{"error": "启动Shell失败: " + err.Error()})
		return
	}

	// 更新会话信息
	sshClient.Session = session
	sshClient.Stdin = stdin
	sshClient.Stdout = stdout
	sshClient.Stderr = stderr
	sshClient.WsConn = ws

	// 处理WebSocket消息
	go handleWebSocketInput(sshClient)
	go handleSSHOutput(sshClient)

	// 等待会话结束
	if err = session.Wait(); err != nil {
		log.Printf("会话结束: %v", err)
	}
}

// handleWebSocketInput 处理WebSocket输入
func handleWebSocketInput(client *SSHClient) {
	for {
		_, message, err := client.WsConn.ReadMessage()
		if err != nil {
			log.Printf("读取WebSocket消息失败: %v", err)
			return
		}

		var msg struct {
			Type    string `json:"type"`
			Command string `json:"command"`
		}
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("解析消息失败: %v", err)
			continue
		}

		if msg.Type == "command" {
			_, err = client.Stdin.Write([]byte(msg.Command))
			if err != nil {
				log.Printf("写入命令失败: %v", err)
				return
			}
		}
	}
}

// handleSSHOutput 处理SSH输出
func handleSSHOutput(client *SSHClient) {
	// 处理标准输出
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := client.Stdout.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Printf("读取STDOUT失败: %v", err)
				}
				return
			}
			if n > 0 {
				message := gin.H{
					"type": "output",
					"data": string(buf[:n]),
				}
				if err := client.WsConn.WriteJSON(message); err != nil {
					log.Printf("发送STDOUT失败: %v", err)
					return
				}
			}
		}
	}()

	// 处理标准错误
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := client.Stderr.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Printf("读取STDERR失败: %v", err)
				}
				return
			}
			if n > 0 {
				message := gin.H{
					"type": "error",
					"data": string(buf[:n]),
				}
				if err := client.WsConn.WriteJSON(message); err != nil {
					log.Printf("发送STDERR失败: %v", err)
					return
				}
			}
		}
	}()
}

// HandleSSHDisconnect 处理SSH断开连接
func HandleSSHDisconnect(c *gin.Context) {
	var req struct {
		ID string `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数: " + err.Error()})
		return
	}

	connectionMutex.Lock()
	defer connectionMutex.Unlock()

	if sshClient, exists := activeConnections[req.ID]; exists {
		if sshClient.Session != nil {
			sshClient.Session.Close()
		}
		if sshClient.Client != nil {
			sshClient.Client.Close()
		}
		if sshClient.WsConn != nil {
			sshClient.WsConn.Close()
		}
		delete(activeConnections, req.ID)
		c.JSON(http.StatusOK, gin.H{"message": "断开连接成功"})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "会话未找到"})
	}
}
