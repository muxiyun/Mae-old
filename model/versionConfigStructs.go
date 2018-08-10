//为了解析前端传过来的json而创建的type

package model


type Env struct{
	EnvKey string `json:"env_key"`
	EnvVal string  `json:"env_val"`
}

type Volume struct {
	HostPath string `json:"host_path"`
	TargetPath string `json:"target_path"`
}



type PortPair struct{
	ImagePort uint `json:"image_port"`
	TargetPort uint `json:"target_port"`
}

type Container struct {
	CtrName string `json:"ctr_name"`
	ImageURL string `json:"image_url"`
	StartCmd string `json:"start_cmd"`
	Envs []Env `json:"envs"`
	Volumes []Volume `json:"volumes"` 
	Ports []PortPair `json:"ports"`
}


type MAEDeployment struct{
	DeployName string `json:"deploy_name"`
	NameSapce string `json:"name_space"`
	Replicas int `json:"replicas"`
	Labels map[string]string `json:"labels"`
	PodLabels map[string]string `json:"pod_labels"`
	Containers []Container `json:"containers"`
}

type MAEService struct{
	SvcName string `json:"svc_name"`
	Selector map[string]string `json:"selector"`
	SvcType string `json:"svc_type"` //clusterip,nodeport,loadbalancer 
}

type VersionConfig struct{
	Deployment MAEDeployment `json:"deployment"`
	Svc    MAEService `json:"svc"`  //此处的Svc指的是k8s中的Service的概念
}
type ReqVersion struct {
	ServiceID uint `json:"svc_id"`
	VersionName string `json:"version_name"`
	VersionDesc string `json:"version_desc"`
	VersionConf VersionConfig `json:"version_conf"`
}


