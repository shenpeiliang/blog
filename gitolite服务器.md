#### 安装
前提需要安装好git

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

创建用户
useradd -d /home/git -m git
su git 
cd ~

git clone https://github.com/sitaramc/gitolite.git

确保这个文件夹是空的：
rm -rf ~/.ssh/authorized_keys

安装：
mkdir $HOME/bin
gitolite/install -to $HOME/bin

生成公钥：
git config --global user.name "shen"
git config --global user.email "harry.shen90@hotmail.com"

ssh-keygen -t rsa -C "harry.shen90@hotmail.com"

查看秘钥：
ls -lh ~/.ssh

cp ~/.ssh/id_rsa.pub ~/gitolite-admin.pub

$HOME/bin/gitolite setup -pk $HOME/gitolite-admin.pub
```

#### 配置

```
cd ~/gitolite-admin

公钥目录，存放可访问用户的公钥文件
ls -lh ./keydir

配置用户组
vim conf/users.conf

# 管理员
@admin = harry.shen

# php组
@php = shenpeiliang

仓库配置
vim conf/gitolite.conf

include "users.conf"

repo gitolite-admin
    RW+     =   gitolite-admin

repo testing
    RW+     =   @php
```

配置好后提交生效
```
cd ~/gitolite-admin
git add .
git commit -m 'init'
git push
```

#### 其他

- 更改远程仓库指向 remote url

```
git remote set-url origin git@github.com:test/thinkphp.git

或进入文件修改地址

git config -e
```