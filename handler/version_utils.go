package handler


import (
	apiv1 "k8s.io/api/core/v1"
	//"k8s.io/client-go/kubernetes/typed/apps/v1beta1"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1b1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/muxiyun/Mae/model"
	"k8s.io/apimachinery/pkg/util/intstr"
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



func bindServicePort(version_config model.VersionConfig)([]apiv1.ServicePort){
	var ports []apiv1.ServicePort
	for _, ctr:=range version_config.Deployment.Containers{
		for _,port :=range ctr.Ports{
			var svcport apiv1.ServicePort
			svcport.Protocol=port.Protocol
			svcport.TargetPort=intstr.IntOrString{IntVal:int32(port.ImagePort), StrVal:string(port.ImagePort)}
			svcport.Name=port.PortName
			svcport.Port=port.TargetPort
			ports=append(ports,svcport)
		}
	}
	return ports
}

// config the Deployment with version_config struct and return the pointer
// of the deployment object
func ConfigDeployment(version_config model.VersionConfig)(*appsv1b1.Deployment){
	deployment := &appsv1b1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: version_config.Deployment.DeployName,
		},
		Spec: appsv1b1.DeploymentSpec{
			Replicas: int32Ptr(int32(version_config.Deployment.Replicas)),
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: version_config.Deployment.Labels,
				},
				Spec: apiv1.PodSpec{

					Volumes:    bindVolumeSource(version_config),
					Containers: bindContainers(version_config),
				},
			},
		},
	}
return deployment
}


// config the service with the version_config struct and return the pointer
// of the service object
func ConfigService(version_config model.VersionConfig)(*apiv1.Service){
	svc := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      version_config.Svc.SvcName,
			Namespace: version_config.Deployment.NameSapce,
			Labels:    version_config.Svc.Labels,
		},
		Spec: apiv1.ServiceSpec{
			Type:     version_config.Svc.SvcType,
			Ports:    bindServicePort(version_config),
			Selector: version_config.Svc.Selector,
		},
	}
	return svc
}
