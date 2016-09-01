package debug

import (
	"log"
	"runtime"
)

const DEBUG = true

func Debug(obj ...interface{}) {
	if DEBUG {
		pc, _, line, _ := runtime.Caller(1)
		f := runtime.FuncForPC(pc)

		args := make([]interface{}, 0)
		args = append(args, obj...)
		//args = append(args, file)
		args = append(args, "\t\t")
		args = append(args, f.Name())
		args = append(args, line)
		log.Println(args...)
	}
}
