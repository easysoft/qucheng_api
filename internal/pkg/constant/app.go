package constant

type AppStatusType int

const (
	AppStatusUnknown  AppStatusType = iota // 初始化
	AppStatusAbnormal                      // 异常
	AppStatusInit                          // 初始化
	AppStatusStop                          // 停止
	AppStatusStartup                       // 启动
	AppStatusDestroy                       // 销毁
	AppStatusSuspend                       // 暂停
	AppStatusRunning                       // 运行中
)

var AppStatusMap = map[AppStatusType]string{
	AppStatusUnknown:  "unknown",
	AppStatusAbnormal: "abnormal",
	AppStatusInit:     "init",
	AppStatusStop:     "stop",
	AppStatusStartup:  "startup",
	AppStatusDestroy:  "destroy",
	AppStatusSuspend:  "suspend",
	AppStatusRunning:  "running",
}
