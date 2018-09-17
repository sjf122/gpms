# gpm
golang多项目配置工具
==========

       1，安装golang到C:\Go目录
       2，project文件夹用于存放全部的项目，work存放当前工作的项目，每个项目包含与GOPATH一致的3个文件夹bin,pkg,src，额外包含一个project.txt文件，内容为项目名称
       3，设置系统环境变量GOPATH 为 D:\xxx\project\work
       4，设置系统环境变量GOBIN 为 D:\xxx\project\work\bin
----------
       gpm命令列表
       将xxx项目目录设置为work: use xxx项目名
       清空work目录: clean
       创建新项目: create xxx
       退出命令行: exit
       帮助: help
