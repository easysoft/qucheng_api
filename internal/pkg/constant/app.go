package constant

type AppStatusType string

const (
	AppStatusInit     AppStatusType = "init"     // 初始化
	AppStatusStop     AppStatusType = "stop"     // 停止
	AppStatusStartup  AppStatusType = "startup"  // 启动
	AppStatusSuspend  AppStatusType = "suspend"  // 暂停
	AppStatusDestroy  AppStatusType = "destroy"  // 销毁
	AppStatusRunning  AppStatusType = "running"  // 运行中
	AppStatusAbnormal AppStatusType = "abnormal" // 异常
	AppStatusUnknown 	AppStatusType = "unknown"  // 未知
)

func (a AppStatusType) String() string {
	return string(a)
}
