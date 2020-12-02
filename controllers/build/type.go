package build

const (
	DefaultCurlImage = "pstauffer/curl"
	DefaultImage     = "registry.cn-hangzhou.aliyuncs.com/tangcong/typhoon_monitor:v1"

	LabelOwnerKey   = "owner"
	LabelProjectKey = "project"
	LabelServiceKey = "service"
	LabelCreateTime = "create-time"
	LabelJobName    = "job-name"

	StopBuildFinalizer = "stop-build"

	DroneServer = "http://114.212.82.229:30170"

	AliyunLogsPrefix = "aliyun_logs"
	STDOUT           = "stdout"
)
