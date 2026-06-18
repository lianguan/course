package domain

import (
	"io"
	"time"
)

type (
	FileStatus int
	FileType   string
)

const (
	ClientUploadInProgress FileStatus = iota // 客户端上传中
	UploadedByClient                         // 客户端上传完成
	ClientUploadError                        // 客户端上传失败
	StorageUploadInProgress                   // 存储上传中
	UploadedToStorage                         // 存储上传完成
	StorageUploadError                        // 存储上传失败
)

const (
	Image FileType = "image" // 图片
	Video FileType = "video" // 视频
	Other FileType = "other" // 其他
)

type File struct {
	ID              uint       `gorm:"primaryKey;autoIncrement" json:"id"`              // 文件ID
	SchoolID        uint       `gorm:"not null;index" json:"schoolId"`                  // 所属学校ID
	Type            FileType   `gorm:"size:50;not null;index" json:"type"`              // 文件类型
	ContentType     string     `gorm:"size:100" json:"contentType"`                     // MIME类型
	Name            string     `gorm:"size:255;not null" json:"name"`                   // 文件名
	Size            int64      `gorm:"not null" json:"size"`                            // 文件大小(字节)
	Status          FileStatus `gorm:"not null;default:0;index" json:"status"`          // 上传状态
	UploadStartedAt time.Time  `gorm:"not null" json:"uploadStartedAt"`                 // 上传开始时间
	URL             string     `gorm:"size:500" json:"url"`                             // 文件访问URL
}

// UploadInput 文件上传输入（Service 层使用）
type UploadInput struct {
	File        io.Reader
	Filename    string
	Size        int64
	ContentType string
	SchoolID    uint
	Type        FileType
}
