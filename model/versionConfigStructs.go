//为了解析前端传过来的json而创建的type

package model

import (
	apiv1 "k8s.io/api/core/v1"
)

type Env struct {
	EnvKey string `json:"env_key"`
	EnvVal string `json:"env_val"`
}

type Volume struct {
	VolumeName string `json:"volume_name"`
	ReadOnly   bool   `json:"read_only"`
	HostPath   string `json:"host_path"`

	//https://kubernetes.io/docs/concepts/storage/volumes/#hostpath
	HostPathType apiv1.HostPathType `json:"host_path_type"`

	TargetPath string `json:"target_path"`
}

type PortPair struct {
	PortName   string         `json:"port_name"` // this name can be referred by service
	ImagePort  int32          `json:"image_port"`
	TargetPort int32          `json:"target_port"`
	Protocol   apiv1.Protocol `json:"protocol"` // TCP or UDP
}

type Container struct {
	CtrName  string     `json:"ctr_name"`
	ImageURL string     `json:"image_url"`
	StartCmd []string   `json:"start_cmd"`
	Envs     []Env      `json:"envs"`
	Volumes  []Volume   `json:"volumes"`
	Ports    []PortPair `json:"ports"`
}

type MAEDeployment struct {
	DeployName string            `json:"deploy_name"`
	NameSapce  string            `json:"name_space"`
	Replicas   int               `json:"replicas"`
	Labels     map[string]string `json:"labels"`
	Containers []Container       `json:"containers"`
}

type MAEService struct {
	SvcName  string            `json:"svc_name"`
	Selector map[string]string `json:"selector"`
	SvcType  apiv1.ServiceType `json:"svc_type"` //clusterip,nodeport,loadbalancer
	Labels   map[string]string `json:"labels"`
}

type VersionConfig struct {
	Deployment MAEDeployment `json:"deployment"`
	Svc        MAEService    `json:"svc"` //此处的Svc指的是k8s中的Service的概念
}

type ReqVersion struct {
	ServiceID   uint          `json:"svc_id"`
	VersionName string        `json:"version_name"`
	VersionDesc string        `json:"version_desc"`
	VersionConf VersionConfig `json:"version_conf"`
}
