#### 环境配置

- 开始 > 设置 > 更新和安全 > 开发者选项 > 开发人员模式

- 开始 > 设置 > 应用 > 程序和功能 > 启用或关闭 Windows 功能，勾选 “适用于 Linux 的 Windows 子系统” 

#### 安装 Debian

在应用商店Microsoft Store中搜索 Debian 直接安装即可，在开始菜单就能看到安装好的程序，可以像普通的软件直接打开

#### Xshell连接

- 修改 root 管理员密码
```
sudo passwd root
```

- 更改软件源

替换成国内的阿里云仓库
```
sed -i 's/deb.debian.org/mirrors.aliyun.com/g' /etc/apt/sources.list
```

或直接修改文件
```
cp /etc/apt/sources.list /etc/apt/sources.list.old

vim /etc/apt/sources.list

deb http://mirrors.aliyun.com/debian/ stretch main non-free contrib
deb-src http://mirrors.aliyun.com/debian/ stretch main non-free contrib
deb http://mirrors.aliyun.com/debian-security stretch/updates main
deb-src http://mirrors.aliyun.com/debian-security stretch/updates main
deb http://mirrors.aliyun.com/debian/ stretch-updates main non-free contrib
deb-src http://mirrors.aliyun.com/debian/ stretch-updates main non-free contrib
deb http://mirrors.aliyun.com/debian/ stretch-backports main non-free contrib
deb-src http://mirrors.aliyun.com/debian/ stretch-backports main non-free contrib
```
- 更新程序

```
apt-get update
```

- 安装需要的软件

```
apt install vim

# 安装扩展使用 ifconfig 命令
apt install net-tools 

# 使用 SSH
apt install openssh-server
```
- 配置 SSH

```
vim /etc/ssh/sshd_config

PermitRootLogin yes
PubkeyAuthentication no
PasswordAuthentication yes
```

重启：
```
/etc/init.d/ssh restart
```

查看状态：
```
/etc/init.d/ssh status
```

添加开机自启动:
```
update-rc.d ssh enable
```

xshell连接：

主机：127.0.0.1
端口号：22

