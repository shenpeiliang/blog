### 基础相关
#### deepin常用

```
进系统后修改超管账号密码：
sudo passwd

deepin相关快捷键：
查看应用目录：
ls -lh /usr/share/applications

直接打开对应的程序名，例如打开文本编辑器：
deepin-editor.desktop则打开命令：
deepin-editor

deepin打开文件夹窗口命令，例如打开当前文件夹窗口：
dde-file-manager .

打开终端快捷命令：
ctrl+alt+t

```

#### 跳转目录快捷键

```

设置别名快速跳转到指定目录：
vim ~/.bashrc

追加内容：
alias skc='cd /data/htdocs/sk_console'

function sk {
        if [ "$1" == "" ]; then
                DIR="src"
        else
                DIR="src_$1"
        fi  
        cd "/data/htdocs/$DIR"
}


alias cd-docker='cd /media/shikee/ntfs'
alias cd-www='cd /media/shikee/ntfs/docker/html/www'

立即生效：
source ~/.bashrc

之后可以通过别名跳转到指定目录，例如：
cd-www
```

#### 常见问题

vim复制粘贴不了：

touch /etc/vim/vimrc.local 

vim /etc/vim/vimrc.local 

填入以下内容，然后重启终端
```
source $VIMRUNTIME/defaults.vim

let skip_defaults_vim = 1

if has('mouse')

    set mouse=r

endif
```

#### 其他

查看所有用户：
vi /etc/passwd

查看是否安装：
dpkg -l | grep phpstorm

得到完整名称后查看软件安装位置：
dpkg -L com.jetbrains.phpstorm


进程查看：
ps -ef | grep redis
netstat -lntp|grep 80

pkill nginx
kill -9 123456

[文件查找](https://www.cnblogs.com/jiftle/p/9707518.html)

find /home -name "*.txt"


批量修改文件名的某一部分：

sudo apt install rename

例如在当前目录把user.sk0.com.conf文件名修改为user.taotaofa1.cn.conf

rename 's/sk0\.com/taotaofa0\.cn/' ./*


[vim替换内容](https://www.cnblogs.com/GODYCA/archive/2013/02/22/2922840.html)


:%s#/home/szs/stb/stsdk/A36/rpmbuild/BUILD#/home/yinjiabin/qt#g
 
解释：

将/home/szs/stb/stsdk/A36/rpmbuild/BUILD替换为/home/yhinjiabin/qt

gg=G 格式化
ggvG 全选复制

### 环境搭建
#### git安装

```
sudo apt install git
git --version
git config --global user.name "shenpeiliang"
git config --global user.email 2172592393@qq.com

生成公钥：
ssh-keygen -t rsa -C "shenpeiliang"

提供给仓库管理员公钥：id_rsa.pub
ls -lh ~/.ssh

中文乱码处理：
git config --global core.quotepath false
```

#### docker安装

```
安装：
curl -sSL https://get.daocloud.io/docker | sh

查看版本：
docker version

查看状态：
systemctl status docker 

设置开机启动：
systemctl daemon-reload
systemctl restart docker
```

#### docker-compose安装

[docker相关文件](docker)

```
需要管理员安装：
curl -L https://get.daocloud.io/docker/compose/releases/download/1.27.4/docker-compose-`uname -s`-`uname -m` > /usr/local/bin/docker-compose

授权
chmod +x /usr/local/bin/docker-compose

版本查看：
docker-compose --version


docker-compose.yml所属目录构建服务：
docker-compose up -d --build

查看服务状态：
docker-compose ps


启动服务：
docker-compose up -d 

```

#### mysql相关

```
进入容器：
docker-compose exec mysql-master bash

查看数据库：
mysql -uroot -p

创建数据库：
CREATE DATABASE IF NOT EXISTS my_db DEFAULT CHARSET utf8 COLLATE utf8_general_ci;

还原数据库：
use my_db;
source /var/backup/my_db.sql

/var/backup是docker-compose.yml配置文件里的路径映射


查看配置文件读取顺序：
mysql --help |grep '\.cnf'

问题：
mysql: [Warning] World-writable config file '/etc/mysql/mysql.conf.d/mysqld.cnf' is ignored.
解决：
这个时候需要将mysql.conf.d文件通过chmod进行权限降级
sudo chmod 644 mysql.conf.d

```

#### nginx站点配置

```
域名解析：
vim /etc/hosts

站点配置目录：docker/nginx/conf.d

修改配置之后进容器重新加载：
docker-compose exec nginx sh

检查配置是否有错：
nginx -t

重新加载配置：
nginx -s reload

```

#### php配置

- 安装xdebug扩展
```
配置：
upload_max_filesize=100M
post_max_size=108M
memory_limit=1024M
date.timezone=Asia/Shanghai

display_errors = On
display_startup_errors = On

[XDebug]
xdebug.remote_enable = 1
xdebug.remote_autostart = 1
xdebug.log  = /usr/local/php/xdebug.log
xdebug.mode = debug
xdebug.client_port = 9003
xdebug.client_host = 172.24.0.1
xdebug.idekey = PHPSTORM
xdebug.cli_color = 2
xdebug.var_display_max_depth = 15
xdebug.var_display_max_data  = 2048
xdebug.profiler_append = 0
xdebug.profiler_output_name = cachegrind.out.%p
xdebug.start_with_request = yes
xdebug.trigger_value = StartProfileForMe
```

相关参考：
https://www.cnblogs.com/jun1019/p/9735250.html

xdebug3配置:
https://www.cnblogs.com/feimoc/p/14684730.html

官方配置说明：
https://www.jetbrains.com/help/phpstorm/2021.2/configuring-xdebug.html#updatingPhpIni

配置实例：
https://blog.csdn.net/benpaodelulu_guajian/article/details/90574728

注意client_host，官方推荐xdebug.client_host = host.docker.internal，但在容器内ping不通
可以在/etc/hosts中添加：
172.24.0.1 docker.host.internal 

网络相关：
https://blog.csdn.net/weixin_2158/article/details/106481238


- 开启错误提示：

查看phpinfo:

http://localhost/phpinfo.php

如果需要修改php配置：

99-overrides.ini

例如添加：

display_errors = On
display_startup_errors = On

重启服务：

docker-compose restart php

- 其他问题

docker中php的执行权限需要添加主机的用户：
在docker-compose.yml中的php配置中加入：
user: 1000:1000

或者在dockerfile中加入命令，改变用户的uid：

RUN usermod -u 1000 www-data

- 文件上传权限不足

```
以管理员权限进到容器：
docker-compose  exec -u root php bash

mkdir -p /data/htdocs/img.sylm.com
chown -R 1000:1000 /data/htdocs

需要在docker-compose.yml做路径映射：（php模块）
- $PWD/html/www/img.shikee.com/:/data/htdocs/img.sylm.com/
- $PWD/html/www/img.shikee.com/:/home/shi/code/img.shikee.com/
```

### 开发工具
#### phpstorm汉化

settings >> Plugins：

搜索Chinese安装重启

主题下载：

Material Theme UI

安装实例：

https://www.cnblogs.com/Dong-Ge/articles/11248689.html

下载低版本：

https://www.jetbrains.com/phpstorm/download/other.html

手动安装：
```
sudo tar xf PhpStorm-2019.2.tar.gz
sudo mkdir /opt/phpstorm/
sudo mv PhpStorm-192.5728.108/* /opt/phpstorm/
sudo ln -s /opt/phpstorm/bin/phpstorm.sh /usr/local/bin/phpstorm

启动：
phpstorm
```

创建快捷桌面：
```
cd /usr/share/applications
vim phpstorm.desktop

[Desktop Entry]
Type=Application
Version=2019.2
Name=phpstorm
Comment=phpstorm ide
Exec=phpstorm
Icon=/usr/share/phpstorm/bin/phpstorm.svg
Terminal=false
Categories=Development;IDE;
StartupNotify=true
```

#### 数据库管理工具

mysql-workbench

navicat for mysql

HeidiSQL

dbeaver

个人习惯使用HeidiSQL，安装参考：

```
sudo apt install deepin-wine
deepin-wine --version
wget -c https://www.heidisql.com/installers/HeidiSQL_11.3.0.6295_Setup.exe

安装：
deepin-wine HeidiSQL_11.2.0.6213_Setup.exe 

启动：
deepin-wine ./HeidiSQL/heidisql.exe

可以设置别名：
vim ~/.bashrc
alias heidisql='deepin-wine /media/shikee/ntfs/HeidiSQL/heidisql.exe'

生效：
source ~/.bashrc
```

#### 网络抓包工具

```
安装：
wget -q -O - https://www.charlesproxy.com/packages/apt/PublicKey | sudo apt-key add -
sudo sh -c 'echo deb https://www.charlesproxy.com/packages/apt/ charles-proxy main > /etc/apt/sources.list.d/charles.list'
sudo apt-get update
sudo apt-get install charles-proxy

启动命令：
charles

生成证书：
help->ssl-proxying->install proxy root certificate
右击导出

到相应浏览器导入证书或设置手动代理
```

参考：
https://blog.csdn.net/pineapple_C/article/details/109168828

https://www.charlesproxy.com/documentation/installation/apt-repository/

