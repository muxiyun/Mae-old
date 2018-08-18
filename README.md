# [muxi application engine(Mae)](https://github.com/muxiyun/Mae/tree/master)

```
                 _      _      _
              __(.)< __(.)> __(.)=
              \___)  \___)  \___)
                     _      _      _
                  __(.)< __(.)> __(.)=
                  \___)  \___)  \___)
                 _      _      _
              __(.)< __(.)> __(.)=
              \___)  \___)  \___)
                     _      _      _
                  __(.)< __(.)> __(.)=
                  \___)  \___)  \___)

```

PaaS of Muxi-Studio. An easier way to manage Kubernetes cluser.

Click [http://zxc0328.github.io/2017/05/27/mae/](http://zxc0328.github.io/2017/05/27/mae/) to view details.


## Feature:
- [x] fast version switch and management
- [x] application management for microservices
- [x] casbin access control(RBAC with domains/tenants)
- [x] log query of sepcific container
- [x] web terminal of specific container
- [x] email confirm and email notification

## Build
Clone the source code and cd to the root dir of this project and execute the command below to install the dependencies.
``` bash
$ glide install
```
Users of mainland China may encounter some problem here. This project uses some dependencies that are blocked by the GWF. So you have to do terminal proxy configuration. How to config it? You can refer to this article:[https://andrewpqc.github.io/2018/04/30/let-the-terminal-penetrate-the-firewall.](https://andrewpqc.github.io/2018/04/30/let-the-terminal-penetrate-the-firewall)
you can also refer to [`glide mirror`](https://glide.readthedocs.io/en/latest/commands/#glide-mirror) to resove the problem.

Before you run Mae, you firstly have to config it. How to config? 

Firstly, you have to get the admin's kubeconfig file, and make sure the name of this file is `admin.kubeconfig`(if not so, you may have to rename it). Then put the `admin.kubeconfig` file in the `conf` folder of this project. the `admin.kubeconfig` is the link between this program and the kubernets cluster. So it's really important.

Secondly, You have to edit `conf/config.yaml` to config the mysql database connection information, the listen address and other configurable options. There are a lot of annotations in the config file, so you can view that to know more.

After you have finished the config part, you can build and run it by typing the following command in you shell.
``` bash
$ go build && ./Mae
```
Then, you can check `/api/v1.0/sd/health` to see whether it work properly or not.

## Test
After you have finished the config part, you can just to use the following command to run the whole test.
``` bash
$ go test -v -cover=true
```
But we don't to suggest you to do so. In order to prevent the rapid consumption of cluster resources and the interaction between test cases, we recommend that you run the test cases one by one in you integrated development environment(suggest [Goland](https://www.jetbrains.com/go/)).

## Next:

1. Grayscale release
2. Optimize the use of int types
3. Optimize the organization of the code
4. Further improve the documentation

## Thanks
Thanks for developers of kubernetes,client-go,iris and mysql.