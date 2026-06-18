package domain

import "io"

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

// File 文件实体
type File struct {
	ID              uint       // 文件ID
	SchoolID        uint       // 所属学校ID
	Type            FileType   // 文件类型
	ContentType     string     // MIME类型
	Name            string     // 文件名
	Size            int64      // 文件大小(字节)
	Status          FileStatus // 上传状态
	UploadStartedAt int64      // 上传开始时间（Unix 时间戳）
	URL             string     // 文件访问URL
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
