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

func main() {

	if !initFlag() {
		return
	}

	switch cmd {
	case "new":
		if len(cmdArgs) > 0 {
			newProject(cmdArgs[0])
		}
	default:
		usage()
	}

}

func isDirExists(path string) bool {
	_, err := os.Stat(path)

	if err != nil {
		return !os.IsNotExist(err)
	} else {
		return true
	}
}

func checkGopathAndProjectDir(name string) (string, error) {

	gopath := os.Getenv("GOPATH")
	if len(gopath) == 0 {
		return gopath, errors.New("未设置 GOPTH 环境变量")
	}

	ps := strings.Split(gopath, ";")
	if len(ps) > 1 {
		gopath = ps[0]
	}
	if !isDirExists(gopath) {
		return gopath, errors.New("GoPath:" + gopath + "note found")
	}

	pth := fmt.Sprintf("%s%csrc%c%s", gopath, os.PathSeparator, os.PathSeparator, name)

	if isDirExists(pth) {
		return gopath, errors.New("该项目已经存在:" + pth)
	}
	return gopath, nil
}

func newProject(name string) {

	gopath, err := checkGopathAndProjectDir(name)
	if err != nil {
		panic(err)
		return
	}

	basePath := fmt.Sprintf("%s%c%s%c%s%c", gopath, os.PathSeparator, "src", os.PathSeparator, name, os.PathSeparator)

	fmt.Println("Creating project:", basePath)

	os.MkdirAll(basePath, 0777)
	os.MkdirAll(basePath+"controllers", 0777)
	os.MkdirAll(basePath+"models", 0777)
	os.MkdirAll(basePath+"doc", 0777)
	os.MkdirAll(basePath+"views/Home", 0777)
	os.MkdirAll(basePath+"views/Shared", 0777)

	fmt.Println("create main.go      status:", writeMain(basePath, name))
	fmt.Println("create controller   status:", writeController(basePath))
	fmt.Println("create views        status:", writeViewHome(basePath))
	fmt.Println("create Model        status:", writeModel(basePath))
	fmt.Println("create app.conf     status:", writeConf(basePath))
	fmt.Println("create readme.md    status:", writeReadme(basePath))
	fmt.Println("create static files status:", writeStatic(basePath))

	fmt.Println("Project Created." + basePath)
}

func usage() {
	fmt.Println(
		`usage:
  grest subcmd args

e.g.
  create a new GREST project in $GOPATH:
  > grest new projectName
 
`)
}

func initFlag() bool {
	//flag.StringVar(&cmd, "cmd", "new", "create a new Grest project in $GOPATH dir.")
	//flag.StringVar(&projectName, "n", "testGrest", "project name")

	flag.Parse()

	args := flag.Args()
	fmt.Println(args)

	// 要求必须有子命令 以及参数。
	//格式： grest 子命令 参数
	if len(args) < 2 {
		usage()
		return false
	} else {
		cmd = args[0]
		cmdArgs = args[1:]
	}
	return true
}
