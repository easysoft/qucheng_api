package constant

type AppStatusType int

const (
	AppStatusUnknown    AppStatusType = iota // 未知
	AppStatusAbnormal                        // 异常
	AppStatusInit                            // 初始化
	AppStatusStoping                         // 关闭中
	AppStatusStoped                          // 停止
	AppStatusStarting                        // 启动中
	AppStatusSuspending                      // 暂停中
	AppStatusSuspended                       // 暂停
	AppStatusRunning                         // 运行中
)

var AppStatusMap = map[AppStatusType]string{
	AppStatusUnknown:    "unknown",
	AppStatusAbnormal:   "abnormal",
	AppStatusInit:       "initializing",
	AppStatusStoping:    "stoping",
	AppStatusStoped:     "stoped",
	AppStatusStarting:   "starting",
	AppStatusSuspending: "suspending",
	AppStatusSuspended:  "suspended",
	AppStatusRunning:    "running",
}
