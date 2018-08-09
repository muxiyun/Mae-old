package model


type Env struct{
	EnvKey string
	EnvVal string
}

type Volume struct {
	HostPath string
	TargetPath string
}

type Label struct {
	LabelKey string
	LabelVal string
}

type PortPair struct{
	ImagePort uint
	HostPort uint
}

type Container struct {
	ImageURL string
	StartCmd string
	Envs []Env
	Volume []Volume
	Ports []PortPair
}




type MAEDeployment struct{
	DeployName string
	NameSapce string
	Replicas int
	Labels []Label
	PodLabels []Label
	Containers []Container

}

type MAEService struct{
	Selector []Label
}

type ReqVersion struct {
	Deploy MAEDeployment
	SVC    MAEService
}

