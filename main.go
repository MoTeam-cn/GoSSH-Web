package main

import (
	"embed"
	"gossh-web/ssh"
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

//go:embed templates/index.html
var embedFS embed.FS

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 在生产环境中应该更严格地检查origin
	},
}

func main() {
	r := gin.Default()

	// 从嵌入的文件系统加载模板
	templ := template.Must(template.ParseFS(embedFS, "index.html"))
	r.SetHTMLTemplate(templ)

	// 路由设置
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "SSH Web Client",
		})
	})

	// SSH相关API路由
	api := r.Group("/api")
	{
		api.POST("/auth", ssh.HandleSSHAuth) // 认证接口
		api.GET("/ws/:id", handleWebSocket)  // WebSocket接口
		api.POST("/disconnect", ssh.HandleSSHDisconnect)
	}

	log.Println("Server starting on :8080...")
	r.Run(":8080")
}

func handleWebSocket(c *gin.Context) {
	id := c.Param("id")
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket升级失败: %v", err)
		return
	}
	defer ws.Close()

	ssh.HandleSSHSession(id, ws)
}
