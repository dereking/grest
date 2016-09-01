package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
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
	"io"
	"io/ioutil"
	"errors"
	"os"
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

	err = ioutil.WriteFile(projectDir+"/static.zip", []byte(dat), 0777)
	if err != nil {
		return err
	} 
	defer os.Remove(projectDir+"/static.zip")

	cf, err := zip.OpenReader(projectDir + "/static.zip")

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
			os.MkdirAll(projectDir+"/static/"+file.Name, 0777)
		} else {
			f, err := os.Create(projectDir + "/static/" + file.Name)
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
	ioutil.WriteFile("../writeStatic.go", []byte(src), 0777)

	fmt.Println("../writeStatic.go has been wroten.")
}
