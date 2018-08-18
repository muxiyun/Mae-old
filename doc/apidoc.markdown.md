## Mae维护文档

### 关于用户认证
Mae认证采取Basic Auth认证的方式,首先携带 `Authorization:base64.encode(username:password)`请求头到/api/v1.0/token取到token。之后需要登录的操作携带` Authorization:base64.encode(token:)`请求头即可以。

用户注册要求提供的邮箱真实有效，注册完毕之后，系统会发送一封邮件验证邮箱是否真实。验证的handler还会赋予该用户对应的权限(user/admin)。就是说，如果用户没有点击验证邮件中的链接，该用户的绝大部分操作都会被Forbidden.

### 关于权限管理
Mae采用casbin的RBAC with domains/tenants 的访问控制模型。系统中所有的操作被分为三类:`roleAdmin`,`roleUser`,`roleAnonymous`。这三种类型也被称为三种角色(role)或者三个组(group).
一个操作是这样定义的`username,`


### 关于错误码管理

### 创建Version注意事项

### 关于app,service的删除
系统中实体之间的逻辑关系大致是这样的。系统中可以创建多个app(应用)，一个应用之下有一个或多个service(服务)，一个服务则有一个或者多个version(版本)。每一个版本其实就是一个在数据库中的用来在集群中创建deployment和service的配置文件的记录(序列化之后存成一个字段)。每一个版本有两个状态activa和unactive，active表明当前版本在集群中对应有资源，unactive则表明当前版本在集群中不存在资源。同时灭一个service也会有一个字段用来记录当前service对应的active版本是哪一个(一般来讲，一个service只有一个版本是active的)。

由于app,service,version这些抽象实体之间是一种树状的结构，所以在删除时采取的是级联删除的方式
### 关于邮件通知系统
当系统中发生向删除app,service对象这样的操作时，会发送邮件通知所有的管理员用户

### 关于namespace,deployment name,service naem的名称问题
namespace,deloyment name,service name这些对象的命名必须采取与Kubernetes对象的命名相同的规则