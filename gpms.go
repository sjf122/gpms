// 项目管理工具

package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var cmd string
var projectname string

/*
* 项目名不能用clean，create，work
 */
func main() {
	help()
	for i := 0; i < 1; i++ {
		// 读取用户输入
		fmt.Scanln(&cmd, &projectname)
		if cmd != "" {
			result := handle(cmd, projectname)
			if result == false {
				i--
			}
		} else {
			// 没输入命令 继续循环
			fmt.Println("您的输入为空，请输入命令")
			i--
		}
		// cmd初始化
		cmd = ""
	}
}

// handle 处理
func handle(cmd string, projectname string) bool {
	// 获取当前目录
	newGoPATH, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
		return false
	}
	//
	switch cmd {
	// 退出
	case "exit":
		{
			fmt.Println("Goodbye")
			return true
		}
	// 列出当前的全部
	case "list":
		{
			list()
			return false
		}
	// 显示帮助
	case "help":
		{
			help()
			return false
		}
	// 清空work目录
	case "clean":
		{
			clean()
			return false
		}
	// 创建项目
	case "create":
		{
			create(projectname)
			return false
		}
	// 使用项目
	case "use":
		{
			if clean() == true {
				use(projectname)
			} else {
				fmt.Println("work目录清理失败，请手工修改work目录文件名")
			}
			return false
		}
	// GOPATH与GOBIN设置
	case "thisdir", "newdir":
		{
			// 判断是否是新目标目录
			if cmd == "newdir" {
				if projectname != "" {
					newGoPATH = projectname
				} else {
					fmt.Println("请输入新目录")
					return false
				}
			}
			//
			switch runtime.GOOS {
			case "linux":
				{
					linuxSet(newGoPATH)
				}
			case "windows":
				{
					windowsSet(newGoPATH)
				}
			case "darwin":
				{
					darwinSet(newGoPATH)
				}
			default:
				{
					fmt.Println("未知的操作系统")
				}
			}
			return false
		}
	// 没有此命令
	default:
		{
			fmt.Println("没有此命令")
			return false
		}
	}
	return false
}

// create 创建项目
func create(projectname string) bool {
	if projectname == "" {
		fmt.Println("新建项目名称不能为空")
		return false
	}
	// 获取当前目录
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
		return false
	}
	// 不能用work
	if projectname == "work" {
		fmt.Println("创建失败，不能使用此项目名", projectname)
		return false
	}
	// 读取work项目名称
	f, fileerr := os.Open(dir + "/work/project.txt")
	if fileerr == nil {
		workprojectname, _ := ioutil.ReadAll(f)
		if projectname == string(workprojectname) {
			fmt.Println("创建失败，此项目名正在使用", projectname)
			return false
		}
	}
	defer f.Close()
	// 判断是否已经存在projectname
	_, haveerr := os.Stat(dir + "/" + projectname)
	if haveerr == nil {
		fmt.Println("创建失败，项目已存在", haveerr)
		return false
	}
	// 创建目录 子目录 bin pkg src project.txt
	rooterr := os.Mkdir(dir+"/"+projectname, os.ModePerm)
	if rooterr != nil {
		fmt.Println(rooterr)
		return false
	}
	binerr := os.Mkdir(dir+"/"+projectname+"/bin", os.ModePerm)
	if binerr != nil {
		fmt.Println(binerr)
		return false
	}
	pkgerr := os.Mkdir(dir+"/"+projectname+"/pkg", os.ModePerm)
	if pkgerr != nil {
		fmt.Println(pkgerr)
		return false
	}
	srcerr := os.Mkdir(dir+"/"+projectname+"/src", os.ModePerm)
	if srcerr != nil {
		fmt.Println(srcerr)
		return false
	}
	// 新建文件
	fw, fwerr := os.Create(dir + "/" + projectname + "/project.txt")
	if fwerr != nil {
		fmt.Println(fwerr)
		return false
	}
	fwstrerr := ioutil.WriteFile(dir+"/"+projectname+"/project.txt", []byte(projectname), 0666)
	if fwstrerr != nil {
		fmt.Println(fwstrerr)
		return false
	}
	defer fw.Close()
	fmt.Println("成功创建项目", projectname)
	return true
}

// use 使用项目：将项目的名字改为work，操作之前先将work名字改为其他
func use(projectname string) bool {
	if projectname == "" {
		fmt.Println("要操作的项目名称不能为空")
		return false
	}
	// 获取当前目录
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
		return false
	}
	// 判断是否存在work，存在直接返回失败
	_, workerr := os.Stat(dir + "/work")
	if workerr == nil {
		fmt.Println("work目录已经存在，请先清理work", workerr)
		return false
	}
	// 判断是否存在目标文件夹
	objectdir := dir + "/" + projectname
	_, haveerr := os.Stat(objectdir)
	if haveerr != nil {
		fmt.Println("项目"+projectname+"不存在，请输入正确的项目名称", haveerr)
		return false
	}
	// 重命名文件夹
	actionerr := os.Rename(objectdir, dir+"/work")
	if actionerr != nil {
		fmt.Println(actionerr)
		return false
	}
	fmt.Println("切换工作区成功", projectname)
	return true
}

// clean 清空work目录：读取work里的项目名称 将work文件夹改名为项目名称
func clean() bool {
	// 获取当前目录
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
		return false
	}
	// 判断是否存在work，不存在直接返回成功
	_, workerr := os.Stat(dir + "/work")
	if workerr != nil {
		fmt.Println("work已清理")
		return true
	}
	// 读取项目名称
	f, fileerr := os.Open(dir + "/work/project.txt")
	if fileerr != nil {
		fmt.Println(fileerr)
		return false
	}
	projectname, _ := ioutil.ReadAll(f)
	f.Close()
	// 判断是否存在目标文件夹
	objectdir := dir + "/" + string(projectname)
	_, haveerr := os.Stat(objectdir)
	if haveerr == nil {
		fmt.Println(haveerr)
		return false
	}
	// 重命名文件夹
	actionerr := os.Rename(dir+"/work", objectdir)
	if actionerr != nil {
		fmt.Println(actionerr)
		return false
	}
	fmt.Println("清理work成功")
	return true
}

// list 列出来所有的项目
func list() {
	var project string
	fmt.Println("----------")
	// 获取当前目录
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
		return
	}
	// 读取文件夹内的文件
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		return
	}
	var iswork string
	for _, f := range files {
		iswork = ""
		//fmt.Println(f.Name())
		//fmt.Println(dir + "/" + f.Name())
		ftype, err := os.Stat(dir + "/" + f.Name())
		if err != nil {
			fmt.Println(err)
			return
		}
		if ftype.IsDir() {
			// 读取项目名称
			//fmt.Println(dir + "/" + f.Name() + "/project.txt")
			filename, fileerr := os.Open(dir + "/" + f.Name() + "/project.txt")
			if fileerr != nil {
				//fmt.Println(fileerr)
				return
			}
			projectname, _ := ioutil.ReadAll(filename)
			filename.Close()
			if f.Name() == "work" {
				iswork = " [working]"
			}
			project = project + string(projectname) + iswork + "\r\n"
		}
	}
	fmt.Print(project)
	fmt.Println("----------")
}

// help 显示帮助文档
func help() {
	help := `
----------
命令列表
列出全部项目: list
将xxx项目目录设置为work: use xxx项目名
清空work目录: clean
创建新项目: create xxx
将当前目录的work文件夹设置为GOPATH: thisdir
将xxx目录的work文件夹设置为GOPATH: newdir xxx
退出命令行: exit
帮助: help
PS：如果windows下出现 rename xxx xxx Access is denied. 现象，请重启windows资源管理器/关闭文件夹内打开的文件再进行操作
----------`
	fmt.Println(help)
	fmt.Println("OS: " + runtime.GOOS + " >>> " + runtime.GOARCH)
	fmt.Println("GOPATH: " + os.Getenv("GOPATH"))
	fmt.Println("GOBIN: " + os.Getenv("GOBIN"))
	fmt.Println("----------")
	fmt.Println("当前项目列表: ")
	list()
	fmt.Println(`请输入命令:`)
}

// windows修改gopath
func windowsSet(newGoPATH string) bool {
	// 通过命令行调用setx命令实现 注意：必须加上.Output()才能操作成功
	exec.Command("CMD", "/C", " setx GOPATH "+newGoPATH+"/work /m").Output()
	exec.Command("CMD", "/C", " setx GOBIN "+newGoPATH+"/work/bin /m").Output()
	//
	exec.Command("CMD", "/C", " setx GOPATH "+newGoPATH+"/work").Output()
	exec.Command("CMD", "/C", " setx GOBIN "+newGoPATH+"/work/bin").Output()
	fmt.Println("设置成功, 请重启会话 reload session >>> " + newGoPATH)
	return true
}

// linux修改gopath
func linuxSet(newGoPATH string) bool {
	return darwinSet(newGoPATH)
}

// mac修改gopath
func darwinSet(newGoPATH string) bool {
	// 编辑 ~/.bash_profile文件实现环境变量的修改
	filename := os.Getenv("HOME") + "/.bash_profile"
	// 读取文件
	fi, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return false
	}
	defer fi.Close()
	// 按行查找
	br := bufio.NewReader(fi)
	var out = ""
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		l := string(a)
		if !strings.Contains(l, "export GOPATH") && !strings.Contains(l, "export GOBIN") && !strings.Contains(l, "export PATH=$PATH:$GOBIN") {
			if out != "" {
				out = out + "\n" + l
			} else {
				out = l
			}
		}
	}
	// 输出新文件
	out = out + "\n" + "export GOPATH=" + newGoPATH + "/work"
	out = out + "\n" + "export GOBIN=$GOPATH/bin"
	out = out + "\n" + "export PATH=$PATH:$GOBIN"
	//fmt.Println(out)
	var ws = []byte(out)
	err = ioutil.WriteFile(filename, ws, 0644)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return false
	}
	//
	fmt.Println("设置成功, 请重启会话 reload session >>> " + newGoPATH)
	return true
}
