package scheduler

import (
	"klog-backend/internal/config"
	"klog-backend/internal/utils"

	"github.com/robfig/cron/v3"
)

// Scheduler 定时任务调度器
type Scheduler struct {
	cron        *cron.Cron
	cleanupTask *CleanupTask
}

// NewScheduler 创建定时任务调度器
// @cleanupTask 清理任务
// @return 调度器
func NewScheduler(cleanupTask *CleanupTask) *Scheduler {
	// 创建cron实例，使用秒级精度
	c := cron.New(cron.WithSeconds())

	return &Scheduler{
		cron:        c,
		cleanupTask: cleanupTask,
	}
}

// Start 启动调度器
func (s *Scheduler) Start() {
	if !config.Cfg.Scheduler.Enabled {
		utils.SugarLogger.Info("定时任务调度器已禁用")
		return
	}

	// 注册清理任务
	cronExpr := config.Cfg.Scheduler.CleanupCron
	_, err := s.cron.AddFunc(cronExpr, func() {
		utils.SugarLogger.Info("开始执行定时清理任务")
		if err := s.cleanupTask.Execute(); err != nil {
			utils.SugarLogger.Errorf("定时清理任务执行失败: %v", err)
		} else {
			utils.SugarLogger.Info("定时清理任务执行完成")
		}
	})

	if err != nil {
		utils.SugarLogger.Errorf("注册清理任务失败: %v", err)
		return
	}

	// 启动调度器
	s.cron.Start()
	utils.SugarLogger.Infof("定时任务调度器已启动，清理任务Cron表达式: %s", cronExpr)
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	if s.cron != nil {
		ctx := s.cron.Stop()
		<-ctx.Done()
		utils.SugarLogger.Info("定时任务调度器已停止")
	}
}
