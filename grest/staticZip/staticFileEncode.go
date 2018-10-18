package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {

	encodeFileToGoCode("static.zip")
}

func encodeFileToGoCode(fn string) {
	data, _ := ioutil.ReadFile(fn)

	sEnc := base64.StdEncoding.EncodeToString([]byte(data))
	//fmt.Println(sEnc)

	src := `package main
	
import(
	"archive/zip"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func writeStatic(projectDir string) error{
	if len(projectDir) ==0{
		return errors.New("projectDir can not be  empty.")
	}
	
	b64 :=""`

	const lineChar = 78
	lineCnt := len(sEnc) / lineChar
	for i := 0; i < lineCnt; i++ {
		src = src + "\n b64 = b64 +\"" + sEnc[i*lineChar:(i+1)*lineChar] + "\""
	}
	if lineCnt*lineChar < len(sEnc) {
		src = src + "\n b64 = b64 +\"" + sEnc[lineCnt*lineChar:] + "\""
	}

	src = src + ` 
	
	dat, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return err
	}
	zipFile := fmt.Sprintf("%s%c%s",projectDir,os.PathSeparator,"static.zip")
	err = ioutil.WriteFile(zipFile, []byte(dat), 0777)
	if err != nil {
		return err
	} 
	defer os.Remove(zipFile)

	cf, err := zip.OpenReader(zipFile)

	if err != nil {
		return err
	}
	defer cf.Close()

	for _, file := range cf.File {
		rc, err := file.Open()
		if err != nil {
			return err
		}

		if file.FileInfo().IsDir() {
			os.MkdirAll(
				fmt.Sprintf("%sstatic%c%s", projectDir,  os.PathSeparator, file.Name), 
				0777)
		} else {
			fshortName := file.Name
			if os.PathSeparator=='\\'{
				fshortName = strings.Replace(file.Name,"/","\\",-1) 
			}
			fn := fmt.Sprintf("%s%s%c%s", projectDir,"static",os.PathSeparator , fshortName)
			
			os.MkdirAll(filepath.Dir(fn))
			
			f, err := os.Create(fn)
			if err != nil {
				return err
			}
			defer f.Close()
			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}

	}
	return nil
}
	`

	srcFile := fmt.Sprintf("%s%c%s", "..", os.PathSeparator, "writeStatic.go")
	ioutil.WriteFile(srcFile, []byte(src), 0777)

	fmt.Println(srcFile + " has been wroten.")
}
