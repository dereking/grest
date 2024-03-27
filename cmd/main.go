package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	cmd string
	//projectName string
	cmdArgs []string
)

var _version_ = "1.0.0.0"

func main() {

	if !initFlag() {
		return
	}

	switch cmd {
	// case "new":
	// 	if len(cmdArgs) > 0 {
	// 		newProject(cmdArgs[0])
	// 	}
	case "new":
		if len(cmdArgs) > 0 {
			newProject_gomod(cmdArgs[0])
		}else{
			usage()
		}
	case "version":
		version()
	default:
		usage()
	}

}

func version() {
	fmt.Printf("Version: %s", _version_)
}

func isDirExists(path string) bool {
	_, err := os.Stat(path)

	if err != nil {
		return !os.IsNotExist(err)
	} else {
		return true
	}
}

func checkGopathAndProjectDir(name string) (gopath, projname string, er error) {

	gopath = os.Getenv("GOPATH")
	if len(gopath) == 0 {
		return gopath, name, errors.New("未设置 GOPTH 环境变量")
	}

	ps := strings.Split(gopath, ";")
	if len(ps) > 1 {
		gopath = ps[0]
	}
	if !isDirExists(gopath) {
		return gopath, name, errors.New("GoPath:" + gopath + "note found")
	}

	if os.PathSeparator == '/' {
		projname = strings.Replace(name, "\\", "/", -1)
	} else if os.PathSeparator == '\\' {
		projname = strings.Replace(name, "/", "\\", -1)
	}

	pth := fmt.Sprintf("%s%csrc%c%s", gopath, os.PathSeparator, os.PathSeparator, name)

	if isDirExists(pth) {
		return gopath, projname, errors.New("该项目已经存在:" + pth)
	}
	return gopath, projname, nil
}

func newProject_deprecated(name string) {

	gopath, projName, err := checkGopathAndProjectDir(name)
	if err != nil {
		panic(err)
		return
	}

	basePath := fmt.Sprintf("%s%c%s%c%s%c", gopath,
		os.PathSeparator, "src", os.PathSeparator, projName, os.PathSeparator)

	fmt.Println("Creating project:", basePath)

	os.MkdirAll(basePath, 0777)
	os.MkdirAll(basePath+"controllers", 0777)
	os.MkdirAll(basePath+"models", 0777)
	os.MkdirAll(basePath+"doc", 0777)
	os.MkdirAll(fmt.Sprintf("%s%s%c%s", basePath, "views", os.PathSeparator, "Home"), 0777)
	os.MkdirAll(fmt.Sprintf("%s%s%c%s", basePath, "views", os.PathSeparator, "Shared"), 0777)

	fmt.Println("create readme.md    status:", writeReadme(basePath))
	fmt.Println("create app.conf     status:", writeConf(basePath))

	fmt.Println("create main.go      status:", writeMain(basePath, name))
	fmt.Println("create controller   status:", writeController(basePath))
	fmt.Println("create Model        status:", writeModel(basePath))

	fmt.Println("create views        status:", writeViewHome(basePath))
	fmt.Println("create views        status:", writeViewShared(basePath))

	fmt.Println("create static files status:", writeStatic(basePath))

	fmt.Println("Project Created." + basePath)
}

func newProject_gomod(projName string) {

	basePath := fmt.Sprintf(".%c%s%c", os.PathSeparator, projName, os.PathSeparator)

	fmt.Println("Creating go mod project:", basePath)

	os.MkdirAll(basePath, 0777)
	os.MkdirAll(basePath+"controllers", 0777)
	os.MkdirAll(basePath+"models", 0777)
	os.MkdirAll(basePath+"doc", 0777)
	os.MkdirAll(fmt.Sprintf("%s%s%c%s", basePath, "views", os.PathSeparator, "Home"), 0777)
	os.MkdirAll(fmt.Sprintf("%s%s%c%s", basePath, "views", os.PathSeparator, "Shared"), 0777)

	fmt.Println("create readme.md    status:", writeReadme(basePath))
	fmt.Println("create app.conf     status:", writeConf(basePath))

	fmt.Println("create main.go      status:", writeMain(basePath, projName))
	fmt.Println("create controller   status:", writeController(basePath))
	fmt.Println("create Model        status:", writeModel(basePath))

	fmt.Println("create views        status:", writeViewHome(basePath))
	fmt.Println("create views        status:", writeViewShared(basePath))

	fmt.Println("create static files status:", writeStatic(basePath))

	fmt.Println("create go mod file status:", writeGoMod(basePath, projName))

	fmt.Println("Project Created." + basePath)
}

func usage() {
	fmt.Println(
		`usage:
  grest subcmd args

subcmd : 
1. help
	show this help
2. new
	create a new GREST project in $GOPATH/subdir/projectName:
	> grest new subdir/projectName
3. new2
	create a new GREST project in current dir and use go mod, ${projectName}:
	> grest new2 projectName

4. veriosn
	show the version
`)
}

func initFlag() bool {
	//flag.StringVar(&cmd, "cmd", "new", "create a new Grest project in $GOPATH dir.")
	//flag.StringVar(&projectName, "n", "testGrest", "project name")

	flag.Parse()

	args := flag.Args()

	// 要求必须有子命令 以及参数。
	//格式： grest 子命令 参数
	if len(args) >= 1 {
		cmd = args[0]
	}

	if len(args) >= 2 {
		cmdArgs = args[1:]
	}
	return true
}
