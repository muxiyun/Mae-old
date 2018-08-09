/*************************************************************************8

目前的想法，一个Version就是一个部署的实例(deployment+service的组合),这里的service是k8s里面的service的概念
创建version需要一次性收集所有的创建一个deployment+service组合的配置信息，由于这些配置信息结构化不强，所以目前
的想法是存放在mongodb中，每一份配置信息就对应着一个版本。所有的版本都存放到mongodb中，用户就可以轻松的在版本之间
回滚了

需要存储到mongodb的一份配置信息与k8s的deployment+service的配置文件大致相同；
deployment:
	ns
	name
	replicas
	[]lables //deploy
	[]lables //pod
	[]containers
		name
		imageURL
		runcmd
		[]volumes
		[]ports
		[]envs
service:
	ns
	name
	[]labels
	[]selector
	[]port
*/

package handler

import (
	"github.com/kataras/iris"
	//"github.com/muxiyun/Mae/model"
)

func CreateVersion(ctx iris.Context) {

}

func GetVersion(ctx iris.Context) {

}

func UpdateVersion(ctx iris.Context) {

}

func DeleteVersion(ctx iris.Context) {

}

func GetVersionList(ctx iris.Context) {

}
