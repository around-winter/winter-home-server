// Package handlers 包含所有 HTTP 请求处理器
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"winter-home-server/internal/models"
)

// 服务器启动时间，用于健康检查
var (
	startTime = time.Now()
	mu        sync.Mutex
)

// ==================== HTTP 处理器 ====================

// HealthHandler 健康检查处理器
// GET /health
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	// 确保只允许 GET 请求
	if r.Method != http.MethodGet {
		http.Error(w, "只允许 GET 请求", http.StatusMethodNotAllowed)
		return
	}

	// 计算服务器运行时间
	uptime := time.Since(startTime)

	// 构建响应
	resp := models.HealthResponse{
		Status:  "healthy",
		Uptime:  uptime.String(),
		Version: "1.0.0",
	}

	// 返回 JSON 响应
	respondJSON(w, http.StatusOK, resp)
}

// TestHandler 测试 API 处理器
// POST /api/test
func TestHandler(w http.ResponseWriter, r *http.Request) {
	// 确保只允许 POST 请求
	if r.Method != http.MethodPost {
		http.Error(w, "只允许 POST 请求", http.StatusMethodNotAllowed)
		return
	}

	// 解析请求体
	var req models.TestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "请求格式错误", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// 构建响应
	resp := models.TestResponse{
		Echo:      fmt.Sprintf("收到消息: %s", req.Message),
		Timestamp: time.Now(),
		Success:   true,
	}

	// 返回 JSON 响应
	respondJSON(w, http.StatusOK, resp)
}

// RootHandler 根路径处理器
// GET /
func RootHandler(w http.ResponseWriter, r *http.Request) {
	// 只响应根路径
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "🚀 Winter Home Server 运行中!\n")
	fmt.Fprintf(w, "可用接口:\n")
	fmt.Fprintf(w, "  - GET  /health     健康检查\n")
	fmt.Fprintf(w, "  - POST /api/test   测试 API\n")
}

// ==================== 辅助函数 ====================

// respondJSON 辅助函数，用于返回 JSON 响应
func respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "响应编码失败", http.StatusInternalServerError)
	}
}
