// Package models 定义了 API 请求和响应的数据结构
package models

import "time"

// ==================== 数据模型 ====================

// HealthResponse 健康检查响应结构
type HealthResponse struct {
	Status  string `json:"status"`  // 服务器状态
	Uptime  string `json:"uptime"`  // 运行时间
	Version string `json:"version"` // 版本号
}

// TestRequest 测试 API 请求结构
type TestRequest struct {
	Message string `json:"message"` // 测试消息
}

// TestResponse 测试 API 响应结构
type TestResponse struct {
	Echo      string    `json:"echo"`      // 回声消息
	Timestamp time.Time `json:"timestamp"` // 响应时间戳
	Success   bool      `json:"success"`   // 是否成功
}
