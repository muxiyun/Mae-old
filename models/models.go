package models

import (
	"time"
)

type VersionCreate struct {

	// App名字, 不需要用户填写，由前端自己传
	Appname string `json:"appname"`

	// App 版本号
	Versionname string `json:"versionname"`

	// 副本数量
	Replicas int32 `json:"replicas,omitempty"`

	// 容器信息
	Containers []Container `json:"containers"`
}

type Container struct {

	// container名字
	Containername string `json:"containername"`

	// 镜像地址
	Imageurl string `json:"imageurl"`

	// 镜像版本号
	Imagetag string `json:"imagetag,omitempty"`

	// 容器运行命令
	Runcmd string `json:"runcmd,omitempty"`

	// 环境变量键值对
	Envs []EnVpair `json:"envs,omitempty"`

	// 容器开放的端口
	Port int32 `json:"port,omitempty"`

	// 挂载点键值对
	Volumns []VolumnPair `json:"volumns,omitempty"`
}

type Log struct {

	// log
	Log string `json:"log"`
}

type EnVpair struct {
	Envname string `json:"envname"`

	Envvalue string `json:"envvalue"`
}

type Logininfo struct {
	Password string `json:"password"`

	Username string `json:"username"`
}

type NginxLocation struct {

	// 匹配规则，如\"/\" 或 \"^~ /api/\" 等
	Rule string `json:"rule"`

	// 导向的App, 如\"http://consume-fe.consume-fe.svc.cluster.local:3000\"
	ProxyPass string `json:"proxy_pass"`
}

type NginxServer struct {

	// 监听端口
	Listen int32 `json:"listen"`

	// 访问该应用使用的域名或IP
	Servername string `json:"servername"`

	// URL匹配规则
	Location []NginxLocation `json:"location"`
}

type Svcinfo struct {
	Svcname string `json:"svcname"`

	Port int32 `json:"port"`
}

type Version struct {

	// app名字
	Appname string `json:"appname"`

	// svc名字
	Svcname string `json:"svcname"`

	// 版本名
	Version string `json:"version"`

	// 状态:是否正在使用
	Status string `json:"status"`
}

type VolumnPair struct {

	// 在容器内的挂载点
	Innerlocation string `json:"innerlocation"`

	// 在宿主机上挂载点
	Outlocation string `json:"outlocation"`
}
