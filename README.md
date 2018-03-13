# MAE:MuxiAppEngine
An easier way to manipulate Kubernetes cluser.

The Server Part of [Projetc MAE](http://zxc0328.github.io/2017/05/27/mae/)

# Todo:
暂时不打算在MAE中增加删除选项，而是交给Devops去处理，避免出现不可控的错误.

## 查询
+ 获取所有namespaces
+ 根据namespaces获取相应deployment, service, pods以及其状态
+ 根据namespace/pod/container 获取 log(缺省设置--tail=100等参数)

## 创建
+ 创建namespaces
+ 创建namespaces/deployment
+ 创建namespaces/svc

## 更新
+ 重新构建镜像后重启deployment
+ 更改镜像地址重启deployment
+ 改变副本数量后重启deployment

# 语言
Python OR Go? <br>
K8s Python Client: https://github.com/kubernetes-client/python <br>
K8s Golang Client: https://github.com/kubernetes/client-go <br>
