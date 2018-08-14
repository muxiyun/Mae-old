# [muxi application engine(Mae)](https://github.com/muxiyun/Mae/tree/master)

PaaS of Muxi-Studio. An easier way to manage Kubernetes cluser.

Click [http://zxc0328.github.io/2017/05/27/mae/](http://zxc0328.github.io/2017/05/27/mae/) to view details.


TODO:
- [x] api design
- [x] domain UML & database UML
- [x] user system
- [x] casbin access control
- [x] application (abstract entity)
- [x] service (abstract entity)
- [x] version (abstract entity)
- [x] log query
- [x] web terminal

NEXT:
1.在业务逻辑中实现app,service的级联删除. 即删除一个app或者service时，删除该app或service下的所有资源

2.优化int类型的使用，重构部分代码

3.文档