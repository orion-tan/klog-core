package scheduler

import (
	"fmt"
	"klog-backend/internal/config"
	"klog-backend/internal/repository"
	"klog-backend/internal/utils"
	"os"
	"path/filepath"
)

// CleanupTask 孤儿文件清理任务
type CleanupTask struct {
	mediaRepo *repository.MediaRepository
}

// NewCleanupTask 创建清理任务
// @mediaRepo 媒体repo
// @return 清理任务
func NewCleanupTask(mediaRepo *repository.MediaRepository) *CleanupTask {
	return &CleanupTask{
		mediaRepo: mediaRepo,
	}
}

// Execute 执行清理任务
// @return 错误
func (t *CleanupTask) Execute() error {
	utils.SugarLogger.Info("开始扫描孤儿文件...")

	// 1. 扫描本地文件
	localFiles, err := t.scanLocalFiles(config.Cfg.Media.MediaDir)
	if err != nil {
		return fmt.Errorf("扫描本地文件失败: %w", err)
	}
	utils.SugarLogger.Infof("扫描到本地文件数: %d", len(localFiles))

	// 2. 查询数据库中的文件路径
	dbFilePaths, err := t.mediaRepo.GetAllFilePaths()
	if err != nil {
		return fmt.Errorf("查询数据库文件路径失败: %w", err)
	}
	utils.SugarLogger.Infof("数据库记录文件数: %d", len(dbFilePaths))

	// 3. 构建数据库文件路径集合（用于快速查询）
	dbFileSet := make(map[string]bool)
	for _, filePath := range dbFilePaths {
		dbFileSet[filePath] = true
	}

	// 4. 找出孤儿文件并删除
	orphanCount := 0
	deletedCount := 0
	failedCount := 0

	for _, localFile := range localFiles {
		// 提取相对路径（相对于uploads目录）
		relPath, err := filepath.Rel(config.Cfg.Media.MediaDir, localFile)
		if err != nil {
			utils.SugarLogger.Warnf("获取相对路径失败: %s, 错误: %v", localFile, err)
			continue
		}

		// 检查文件是否在数据库中
		if !dbFileSet[relPath] {
			orphanCount++
			utils.SugarLogger.Infof("发现孤儿文件: %s", relPath)

			// 删除孤儿文件
			if err := os.Remove(localFile); err != nil {
				utils.SugarLogger.Errorf("删除孤儿文件失败: %s, 错误: %v", localFile, err)
				failedCount++
			} else {
				utils.SugarLogger.Infof("成功删除孤儿文件: %s", relPath)
				deletedCount++
			}
		}
	}

	// 5. 输出统计信息
	utils.SugarLogger.Infof("清理任务完成 - 扫描: %d, 孤儿文件: %d, 删除成功: %d, 删除失败: %d",
		len(localFiles), orphanCount, deletedCount, failedCount)

	return nil
}

// scanLocalFiles 扫描本地文件目录
// @dir 目录路径
// @return 文件路径列表, 错误
func (t *CleanupTask) scanLocalFiles(dir string) ([]string, error) {
	var files []string

	// 检查目录是否存在
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		utils.SugarLogger.Warnf("上传目录不存在: %s", dir)
		return files, nil
	}

	// 递归遍历目录
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			utils.SugarLogger.Warnf("访问文件失败: %s, 错误: %v", path, err)
			return nil // 继续遍历其他文件
		}

		// 只处理文件，跳过目录
		if !info.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}
