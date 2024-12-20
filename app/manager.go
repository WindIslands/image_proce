package app

import (
	"context"
)

// AppManager 管理所有应用模块
type AppManager struct {
	ctx       context.Context
	ImageApp  *ImageApp
}

func NewAppManager() *AppManager {
	return &AppManager{
		ImageApp: NewImageApp(),
	}
}

// Startup 统一处理所有模块的启动
func (m *AppManager) Startup(ctx context.Context) {
	m.ctx = ctx
	// 为所有模块设置context
	m.ImageApp.Startup(ctx)
} 