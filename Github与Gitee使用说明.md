#### 本地插件安装
VSCode中使用markdown大部分是为了编写Github的README.md文档

安装插件：

Markdown Preview Github Styling

#### 公钥配置
以gitee为例：

账号 >> 设置 >> 安全设置 >> SSH公钥

添加本地公钥的内容：

ls -lh ~/.ssh

复制id_rsa.pub中的内容填写即可

#### 代码同步

同意使用使用ssh方式克隆仓库

```
添加github仓库关联：
git remote add 远程库名 远程库地址
git remote add github git@github.com:shenpeiliang/blog.git

查看：
git remote -v


推送：
git push 远程库名 分支名
git push github master

拉取：
git pull 远程库名 分支名
git pull github master
```

参考：

[将Gitee代码同步到Github](https://blog.csdn.net/icansoicrazy/article/details/116454389?utm_medium=distribute.pc_aggpage_search_result.none-task-blog-2~aggregatepage~first_rank_ecpm_v1~rank_v31_ecpm-1-116454389.pc_agg_new_rank&utm_term=gitee%E4%BB%93%E5%BA%93%E8%87%AA%E5%8A%A8%E5%90%8C%E6%AD%A5%E5%88%B0github&spm=1000.2123.3001.4430)