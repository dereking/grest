package mvc

import (
	"bufio"
	"os"
	"path"
)

type FileResult struct {
	ActionResult
}

func NewFileResult(code int, filePath string) *FileResult {
	ret := &FileResult{}
	ret.Header = make(map[string]string)
	ret.HttpCode = code
	ret.Message = []byte(filePath)
	return ret
}

func (ar *FileResult) ExecuteResult(c IController) {

	for k, v := range ar.Header {
		c.GetResponse().Header().Set(k, v)
	}

	//c.GetResponse().Header().Set("Content-Type", "application/bin; charset=utf-8")
	c.GetResponse().Header().Set("Content-Disposition", "attachment; filename="+path.Base(string(ar.Message)))

        c.GetResponse().Header().Set("Content-Type", "application/octet-stream")

	c.GetResponse().Header().Set("Accept-Ranges", "bytes")

	/*
	   header("Content-type: text/plain");
	           header("Accept-Ranges: bytes");
	           header("Content-Disposition: attachment; filename=".$filename);
	           header("Cache-Control: must-revalidate, post-check=0, pre-check=0" );
	           header("Pragma: no-cache" );
	           header("Expires: 0" );
	*/

	c.GetResponse().WriteHeader(ar.HttpCode)

	//fmt.Println("===========", ar.HttpCode, c.GetResponse().Header())
	f, err := os.Open(string(ar.Message))
	defer f.Close()
	if err != nil {
		c.GetResponse().Write([]byte(err.Error()))
	} else {
		r := bufio.NewReader(f)
		buf := make([]byte, 1024)
		for {
			n, err := r.Read(buf)
			if n == 0 || err != nil {
				break
			}
			c.GetResponse().Write(buf[:n])
		}
	}
}
