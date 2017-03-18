package templateManager

import (
	"fmt"
	"html/template"
	"time"
)

func templateFunc_html(x string) interface{} {
	return template.HTML(x)
}

func templateFunc_FileSize(s int64) string {
	if s < 1024 {
		return fmt.Sprintf("%d b", s)
	} else if s < 1024*1024 {
		return fmt.Sprintf("%.2f Kb", float64(s)/1024.0)
	} else if s < 1024*1024*1024 {
		return fmt.Sprintf("%.2f Mb", float64(s)/1024/1024)
	} else if s < 1024*1024*1024*1024 {
		return fmt.Sprintf("%.2f Gb", float64(s)/1024/1024/1024.0)
	} else {
		return fmt.Sprintf("%.2f Tb", float64(s)/1024/1024/1024/1024.0)
	}
}

func templateFunc_DateTime(s time.Time) string {
	return s.Format("2006-01-02 15:04:05")
}

func templateFunc_add(n int) int {
	return n + 1
}
