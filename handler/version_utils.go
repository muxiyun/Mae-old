package handler


import (
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/typed/apps/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/muxiyun/Mae/model"
	"fmt"
)

// set all the volume
func bindVolumeSource(version_config model.VersionConfig)[]apiv1.Volume{
	var volumes []apiv1.Volume
	for _,ctx:=range version_config.Deployment.Containers{
		for _,vol:=range ctx.Volumes{
			 var volume  apiv1.Volume //初始化一下
			 var hostpath apiv1.HostPathVolumeSource
			hostpath.Path=vol.HostPath
			hostpath.Type=&vol.HostPathType
			volume.Name=vol.VolumeName
			volume.HostPath=&hostpath

			volumes=append(volumes,volume)
		}
	}
	return volumes
}

// set all envVar
func bindEnvVar(container *apiv1.Container, version_config model.VersionConfig,index int){
	for _,env:=range version_config.Deployment.Containers[index].Envs{
		var envVar apiv1.EnvVar
		envVar.Name=env.EnvKey
		envVar.Value=env.EnvVal
		container.Env=append(container.Env,envVar)
	}
}

func bindVolumeTarget(container *apiv1.Container, version_config model.VersionConfig,index int){
	for _,volume:=range version_config.Deployment.Containers[index].Volumes{
		var volu apiv1.VolumeMount
		volu.MountPath=volume.TargetPath
		volu.ReadOnly=volume.ReadOnly
		volu.Name=volume.VolumeName
		container.VolumeMounts=append(container.VolumeMounts,volu)
	}
}


func bindPart(container *apiv1.Container, version_config model.VersionConfig,index int){
	for _,portpair:=range version_config.Deployment.Containers[index].Ports{
		var ports apiv1.ContainerPort
		ports.ContainerPort=portpair.ImagePort
		ports.Name=portpair.PortName
		ports.HostPort=portpair.TargetPort
		ports.Protocol=portpair.Protocol  //"TCP" or "UDP", default is TCP
		container.Ports=append(container.Ports,ports)
	}
}

func bindContainers( version_config model.VersionConfig)([]apiv1.Container){
	var containers []apiv1.Container
	for index,ctx :=range version_config.Deployment.Containers{
		var container apiv1.Container
		container.Name=ctx.CtrName
		container.Image=ctx.ImageURL
		container.Command=ctx.StartCmd
		bindPart(&container,version_config,index)
		bindEnvVar(&container,version_config,index)
		bindVolumeTarget(&container,version_config,index)
		containers=append(containers,container)
	}
	return containers
}



//监听Deployment变化
func startWatchDeployment(deploymentsClient v1beta1.DeploymentInterface) {
	w, _ := deploymentsClient.Watch(metav1.ListOptions{})
	for {
		select {
		case e, _ := <-w.ResultChan():
			fmt.Println(e.Type, e.Object)
		}
	}
}

