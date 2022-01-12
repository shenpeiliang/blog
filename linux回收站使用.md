#### 普通提醒

意思是输入rm命令就相当于输入rm -i命令，即：每次删除都会询问

```
vim ~/.bashrc

添加别名：
alias rm='rm -i'
alias cp='cp -i'
alias mv='mv -i'

```

#### 自定义命令

```
当前用户：
~/.bashrc

全局：
/etc/bashrc

trash_path="~/.Trash"
  
# 目录是否存在
if [ ! -d "$trash_path" ]; then
   mkdir -p "$trash_path"
fi

alias rm=trash_put
alias trash-list='ls -lh ~/.Trash'
alias trash-restore=trash_restore
alias trash-empty=trash_clear
#撤回
trash_restore()  
{
        mv -i "$trash_path"/$@ ./
}
#放到回收站
trash_put()  
{
        mv $@ "$trash_path"/
}

trash_clear()  
{
        read -p "Clear trash?[n]" confirm
        [ $confirm == 'y' ] || [ $confirm == 'Y' ]  && /usr/bin/rm -rf "$trash_path"/*
}


区别：
1. 两个文件都是设置环境变量文件的，/etc/profile是永久性的环境变量,是全局变量，/etc/profile.d/设置所有用户生效

2. /etc/profile.d/比/etc/profile好维护，不想要什么变量直接删除/etc/profile.d/下对应的shell脚本即可，不用像/etc/profile需要改动此文件
```

#### trash-cli插件

```
安装依赖python

cd /usr/local/src
wget -c https://www.python.org/ftp/python/3.9.0/Python-3.9.0.tgz
tar -zxf Python-3.9.0.tgz
cd Python-3.9.0
./configure --prefix=/usr/local/python
make && make install

python3 -V

环境变量：
export PATH="$PATH:/usr/local/python/bin" 
或
echo 'export PATH=$PATH:/usr/local/python/bin' >> /etc/profile && source /etc/profile

查看是否安装成功：
python3 pip3


安装插件：
pip3 install trash-cli

查看安装成功的插件：
ls -lh /usr/local/bin | grep trash

trash-put           trash files and directories. 放入回收站
trash-empty         empty the trashcan(s). 清空
trash-list          list trashed files. 查看回收站
trash-restore       restore a trashed file. 恢复
trash-rm            remove individual files from the trashcan. 从回收站移除指定文件

重修设置别名，但其他脚本也要修改
alias rm='echo "This is not the command you are looking for."; false'

直接修改:vim /etc/profile.d/trash.sh
alias rm="trash-put"

回收站文件目录地址：
~/.local/share/Trash/
```