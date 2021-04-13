package basic

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const maxSectors = 1000

/* Universal Functions Start */

// CreateFile function to create file and write string into it
func CreateFile(fileName string, fileWrite string) {
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Errorf("error creating file: %v", err)
	}
    defer f.Close()
    
	_, err = f.WriteString(fileWrite)
	if err != nil {
		fmt.Errorf("error writing string: %v", err)
	}

	fmt.Println("written to", fileName)
}

// OpenFile function opens the file to read
func OpenFile(fileName string) {
	_, err := os.Open(fileName)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// CheckError function to check error
func CheckError(e error) {
	if e != nil {
		panic(e)
	}
}

// CreateDirIfNotExist function to create directory if not exist
func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		CheckError(err)
	}
}

// Createtmp function to create tmp folder and subsequent inner directories
func Createtmp() {
	CreateDirIfNotExist("./tmp")
	CreateDirIfNotExist("./tmp/contractInfo")
}

// FileCrWr function to create file 's' and write 'b' bytes in it
func FileCrWr(s string, b []byte) {
	f, err := os.OpenFile(s, os.O_RDWR|os.O_CREATE, 0755)
	CheckError(err)
	_, err = f.Write(b)
	CheckError(err)
	defer f.Close()
}

// Connect function to connect and post new request using http
func Connect(dest string, content string, contentType string) string {
	client := &http.Client{}
	r, err := http.NewRequest("POST", dest, strings.NewReader(content))
	CheckError(err)
	r.Header.Add("Content-Type", contentType)
	r.Header.Add("Content-Length", strconv.Itoa(len(content)))

	resp, err := client.Do(r)
	CheckError(err)
	defer resp.Body.Close()

	var message []byte
	if resp.StatusCode == http.StatusOK {
		message, err = ioutil.ReadAll(resp.Body)
		CheckError(err)
	}
	return string(message)
}

// GethPathAndKey function to get path of Geth node and keystore key
func GethPathAndKey() (string, string) {
	dir, err := os.Getwd()
	CheckError(err)

	pathParent := path.Dir(dir)
	pathParent += "/node-data/node2/"
	gethPath := pathParent + "geth.ipc"

	//var fileName string
	var filePath string
	keyStorePath := pathParent + "keystore/"

	err = filepath.Walk(keyStorePath, func(path string, info os.FileInfo, err error) error {
		CheckError(err)
		if !info.IsDir() {
			fmt.Printf("visited file or dir: %q\n", path)
			//fileName = info.Name()
			filePath = path
		}
		return nil
	})
	CheckError(err)

	var key []byte
	keyRead := func() {
		var err error
		key, err = ioutil.ReadFile(filePath)
		CheckError(err)
	}
	keyRead()
	return gethPath, string(key)
}

// TimeTrack function to check timings
func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("function %s took %s", name, elapsed)
}

/* Universal Functions End */
