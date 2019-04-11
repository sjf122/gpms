# gpms

##Golang多项目配置工具
  发现golang在同时进行好几个不相关的项目时候，很难管理，我希望我公司的项目和个人的项目严格分离。然而GOPATH里多个项目路径的设置导致引用的第三方包会在一个GOPATH路径里，管理很麻烦，所以写了这样一个小工具。初衷就是只设置一个GOPATH作为当前工作目录，通过重命名文件夹的方式将需要设置为当前工作目录的正开发项目路径改为GOPATH路径。不是当前正开发项目的改为其他路径。这样做到了每个项目都严格分离。

* 安装golang到C:\Go目录
* project文件夹用于存放全部的项目，work存放当前工作的项目，每个项目包含与GOPATH一致的3个文件夹bin,pkg,src，额外包含一个project.txt文件，内容为项目名称
* 设置系统环境变量GOPATH 为 D:\xxx\project\work
* 设置系统环境变量GOBIN 为 D:\xxx\project\work\bin
* 如果windows下出现 rename xxx xxx Access is denied. 现象，请重启windows资源管理器/关闭文件夹内打开的文件再进行操作
  
```javascript
  gpm命令列表
  列出全部项目: list
  将xxx项目目录设置为work: use xxx项目名
  清空work目录: clean
  创建新项目: create xxx
  退出命令行: exit
  帮助: help
```
