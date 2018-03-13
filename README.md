# MAE:MuxiAppEngine
An easier way to manipulate Kubernetes cluser.

The Paas of Muxi-Studio, Server Part of [Project MAE](http://zxc0328.github.io/2017/05/27/mae/)

# Note:

+ 对使用者隐藏Namespace, Deployment, Svc等概念，只需提供image镜像以及必须的参数即可部署服务.
+ 提供版本回滚：记录时间、版本配置信息、操作者等，可以回滚到服务的某个版本.

# Todo:

+ 查询服务状态 & log信息
+ 构建新服务
+ 更新服务
+ 删除服务
+ 用户权限管理
+ 数据库记录服务配置信息，保证在服务因各种原因挂掉后可以及时重启
+ 为服务提供相应Nginx配置


# 语言
Golang  https://github.com/kubernetes/client-go
