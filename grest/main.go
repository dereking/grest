package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

var (
	cmd         string
	projectName string
)

func main() {

	initFlag()

	switch cmd {
	case "new":
		if len(projectName) > 0 {
			newProject(projectName)
		} else {
			flag.PrintDefaults()
		}
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
	os.MkdirAll(basePath, 0777)
	os.MkdirAll(basePath+"controllers", 0777)
	os.MkdirAll(basePath+"models", 0777)
	os.MkdirAll(basePath+"doc", 0777)
	os.MkdirAll(basePath+"views/Home", 0777)
	os.MkdirAll(basePath+"views/Shared", 0777)

	fmt.Println("create main.go      status:", writeMain(basePath, projectName))
	fmt.Println("create controller   status:", writeController(basePath))
	fmt.Println("create views        status:", writeViewHome(basePath))
	fmt.Println("create Model        status:", writeModel(basePath))
	fmt.Println("create app.conf     status:", writeConf(basePath))
	fmt.Println("create readme.md    status:", writeReadme(basePath))
	fmt.Println("create static files status:", writeStatic(basePath))

	fmt.Println("Project Created.")
}

func initFlag() {
	flag.StringVar(&cmd, "cmd", "new", "create a new Grest project in $GOPATH dir.")
	flag.StringVar(&projectName, "n", "testGrest", "project name")

	flag.Parse()
}
