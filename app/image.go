package app

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/h2non/bimg"
)

type ImageApp struct {
	ctx context.Context
}

func NewImageApp() *ImageApp {
	return &ImageApp{}
}

// Startup 在应用启动时被调用
func (a *ImageApp) Startup(ctx context.Context) {
	a.ctx = ctx
}

// FileInfo 文件信息结构体
type FileInfo struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Size int64  `json:"size"`
	Path string `json:"path"`
}

// 打开文件选择对话框
func (a *ImageApp) OpenMultipleFilesDialog() ([]FileInfo, error) {
	fileInfos := make([]FileInfo, 0)
	filepaths, err := runtime.OpenMultipleFilesDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择图片文件",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "图片文件 (*.png;*.jpg;*.jpeg;*.webp)",
				Pattern:     "*.png;*.jpg;*.jpeg;*.webp",
			},
		},
	})
	for _, filepath := range filepaths {

		file, err := os.Stat(filepath)
		if err != nil {
			continue
		}
		fileInfo := FileInfo{
			Id:   uuid.New().String(),
			Name: file.Name(),
			Size: file.Size(),
			Path: filepath,
		}
		fileInfos = append(fileInfos, fileInfo)
	}

	return fileInfos, err
}

func (a *ImageApp) OpenDirectoryDialog() (string, error) {
	return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择文件夹",
	})
}

// ProcessImagesOptions 处理选项
type ProcessImagesOptions struct {
	Tasks      []ImageTask `json:"tasks"`
	OutputPath string      `json:"outputPath"`
	Workers    int        `json:"workers"`
}

// 修改ProcessImages方法签名
func (a *ImageApp) ProcessImages(options ProcessImagesOptions) error {
	tasks := options.Tasks
	outputPath := options.OutputPath
	total := len(tasks)
	processed := 0
	errorChan := make(chan error, total)
	doneChan := make(chan bool, total)

	// 使用传入的workers数量
	maxWorkers := options.Workers
	if maxWorkers <= 0 {
		maxWorkers = 1
	}
	if len(tasks) < maxWorkers {
		maxWorkers = len(tasks)
	}
	
	// 创建任务通道
	taskChan := make(chan ImageTask, total)
	
	// 启动工作协程
	for i := 0; i < maxWorkers; i++ {
		go func() {
			for task := range taskChan {
				// 发送进度事件
				runtime.EventsEmit(a.ctx, "process-progress", map[string]interface{}{
					"current": processed + 1,
					"total":   total,
					"file":    task.Name,
				})

				// 处理图片
				if err := BimgHandle(task, outputPath); err != nil {
					// 发送错误事件
					runtime.EventsEmit(a.ctx, "process-error", map[string]interface{}{
						"file":  task.Name,
						"error": err.Error(),
					})
					errorChan <- err
				} else {
					doneChan <- true
				}
			}
		}()
	}

	// 发送任务到通道
	go func() {
		for _, task := range tasks {
			taskChan <- task
		}
		close(taskChan)
	}()

	// 等待所有任务完成
	var errors []error
	for i := 0; i < total; i++ {
		select {
		case err := <-errorChan:
			errors = append(errors, err)
		case <-doneChan:
			processed++
		}
	}

	// 发送完成事件
	runtime.EventsEmit(a.ctx, "process-complete")

	// 如果有错误，返回第一个错误
	if len(errors) > 0 {
		return errors[0]
	}
	return nil
}

// ImageTask 表示一个图片处理任务
type ImageTask struct {
	Name             string `json:"name"`
	Path             string `json:"path"`
	OutputFormat     string `json:"outputFormat"`
	Quality          int    `json:"quality"`
	ConvertFormat    bool   `json:"convertFormat"`
	Width            int    `json:"width"`
	Height           int    `json:"height"`
	KeepOriginalSize bool   `json:"keepOriginalSize"`
}

func BimgHandle(task ImageTask, outputPath string) error {

	filename := task.Name
	// 如果需要转换格式
	if task.ConvertFormat {
		filename = filename[:len(filename)-len(filepath.Ext(filename))] + "." + task.OutputFormat
	}
	outputFilePath := filepath.Join(outputPath, filename)

	data, err := os.ReadFile(task.Path)
	if err != nil {
		log.Println(err)
		return err
	}

	buffer := bimg.NewImage(data)

	// 获取原始图片的尺寸
	originalSize, err := buffer.Size()
	if err != nil {
		log.Println(err)
		return err
	}

	// 设置输出选项
	options := bimg.Options{
		Quality: task.Quality,
	}

	// 如果需要转换格式
	if task.ConvertFormat {
		switch task.OutputFormat {
		case "jpeg", "jpg":
			options.Type = bimg.JPEG
		case "png":
			options.Type = bimg.PNG
		case "webp":
			options.Type = bimg.WEBP
		default:
			log.Println("unsupported output format: %s", task.OutputFormat)
		}
	}

	// 如果需要调整图片大小
	if !task.KeepOriginalSize {
		options.Width = task.Width
		options.Height = task.Height
		// 如果宽度和高度都为0，则保持原始尺寸
		if options.Width == 0 && options.Height == 0 {
			options.Width = originalSize.Width
			options.Height = originalSize.Height
		}
	}

	// 处理图片
	data, err = buffer.Process(options)
	if err != nil {
		log.Println(err)
		return err
	}

	// 将处理后的图片写入文件
	if err := bimg.Write(outputFilePath, data); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
