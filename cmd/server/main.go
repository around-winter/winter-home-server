// Package main 是服务器的入口点
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"winter-home-server/internal/handlers"
	"winter-home-server/internal/middleware"
)

// ==================== 服务器主函数 ====================

func main() {
	startTime := time.Now()

	fmt.Println("========================================")
	fmt.Println("🚀  Winter Home Server")
	fmt.Println("========================================")
	fmt.Printf("📡  启动时间: %s\n", startTime.Format("2006-01-02 15:04:05"))
	fmt.Println("========================================")

	// 注册路由
	setupRoutes()

	// 启动服务器
	addr := ":8080"
	fmt.Printf("\n🌐  服务器启动中，监听地址: http://localhost%s\n", addr)
	fmt.Println("   按 Ctrl+C 停止服务器")
	fmt.Println()

	// 启动 HTTP 服务器
	err := http.ListenAndServe(addr, nil)
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("❌ 服务器启动失败: %v", err)
	}

	fmt.Println("\n👋 服务器已停止")
}

// setupRoutes 注册所有 HTTP 路由
func setupRoutes() {
	// 公开接口（无需认证）
	http.HandleFunc("/", handlers.RootHandler)
	http.HandleFunc("/health", handlers.HealthHandler)

	// 用户认证接口（无需认证）
	http.HandleFunc("/api/register", handlers.RegisterHandler)
	http.HandleFunc("/api/login", handlers.LoginHandler)

	// 需要认证的接口
	authHandler := middleware.JWTAuthMiddleware(http.HandlerFunc(handlers.TestHandler))
	http.Handle("/api/test", authHandler)
}
