更新yum仓库：

```
yum update

```

查看是否已经安装：

```
rpm -qa|grep docker
```

移除旧安装的

```
yum remove docker \
                  docker-client \
                  docker-client-latest \
                  docker-common \
                  docker-latest \
                  docker-latest-logrotate \
                  docker-logrotate \
                  docker-selinux \
                  docker-engine-selinux \
                  docker-engine
```				  

安装环境需要的插件

```
yum install -y yum-utils \
  device-mapper-persistent-data \
  lvm2
```

添加软件仓库

```
yum-config-manager \
    --add-repo \
    https://download.docker.com/linux/centos/docker-ce.repo
```

或

```
yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo
```

更新 yum 缓存：

```
yum makecache fast
```

说明：

```
yum clean all     # 清除系统所有的yum缓存
yum makecache     # 生成yum缓存
```

可以查看所有的软件仓库

```
yum repolist all
```

启用edge、test仓库（默认是不启用的）	

```
yum-config-manager --enable docker-ce-edge
yum-config-manager --enable docker-ce-test
```

如果需要停止使用，禁用的命令行：

```
yum-config-manager --disable docker-ce-edge
```

查看仓库中所有的docker版本

```
yum list docker-ce --showduplicates | sort -r
```

结果：stable稳定版

```
docker-ce.x86_64            18.03.0.ce-1.el7.centos             docker-ce-stable
```

安装docker-ce ，默认会安装最新的版本

```
yum install docker-ce
```

或指定版本：

```
yum install docker-ce-<VERSION STRING>
```

例如：

```
yum install docker-ce-18.03.0.ce
yum install docker-ce-18.06.1.ce
```

验证安装是否成功

```
docker version
```
