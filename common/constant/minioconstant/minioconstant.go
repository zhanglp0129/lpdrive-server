package minioconstant

const (
	// AbortUploadLifecycleID 自动终止分片上传的生命周期规则
	AbortUploadLifecycleID = "abort-upload"
	// UploadExpireDays 上传过期天数
	UploadExpireDays = 1
	// MinPartSize 最小分片大小
	MinPartSize = 5 << 20
	// MaxPartSize 最大分片大小
	MaxPartSize = 20 << 20
)
