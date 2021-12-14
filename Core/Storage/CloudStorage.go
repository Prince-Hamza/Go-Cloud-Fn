package CloudStorage

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	// "sync"

	// "io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	//"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
)

const (
	projectID  = "my-first-project-ce24e"             // FILL IN WITH YOURS
	bucketName = "my-first-project-ce24e.appspot.com" // FILL IN WITH YOURS
)

type ClientUploader struct {
	cl         *storage.Client
	projectID  string
	bucketName string
	uploadPath string
}

var uploader *ClientUploader

func init() {
	//os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "") // FILL IN WITH YOUR FILE PATH
	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	uploader = &ClientUploader{
		cl:         client,
		bucketName: bucketName,
		projectID:  projectID,
		uploadPath: "blogheroku.png",
	}

}

func Main() {
	//uploader.UploadFile("notes_test/abc.txt")
	r := gin.Default()
	r.POST("/upload", func(c *gin.Context) {

		f, _ := c.FormFile("file_input")
		blobFile, _ := f.Open()

		fmt.Println("From /Upload Api")

		err3 := uploader.UploadFileApi(blobFile, f.Filename)

		if err3 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err3.Error(),
			})
			return
		}

		// if err2 != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{
		// 		"error": err2.Error(),
		// 	})
		// 	return
		// }
		// if err1 != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{
		// 		"error": err1.Error(),
		// 	})
		// 	return
		// }

		c.JSON(200, gin.H{
			"message": "success",
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// UploadFile uploads an object
func (c ClientUploader) UploadFileApi(file multipart.File, object string) error {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	fmt.Println("upload Path : ", c.uploadPath)
	wc := c.cl.Bucket(c.bucketName).Object(c.uploadPath).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	return nil
}

func UploadImage(urlPath, method string) error {
	fmt.Println("Call")
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err1 := writer.CreateFormFile("file_input", filepath.Base("Images/blogheroku.png"))
	file, err2 := os.Open("blogheroku.png")
	_, err3 := io.Copy(fw, file)

	writer.Close()

	req, err4 := http.NewRequest("POST", urlPath, bytes.NewReader(body.Bytes()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rsp, _ := client.Do(req)

	if rsp.StatusCode != http.StatusOK {
		log.Printf("Request failed with response code: %d", rsp.StatusCode)
	}

	errz := [4]error{err1, err2, err3, err4}
	onError(errz)

	fmt.Println(rsp)
	return nil
}

func onError(errz [4]error) {
	for _, item := range errz {
		if item != nil {
			fmt.Println(item)
		}
	}

}

func Upload2(res http.ResponseWriter, req *http.Request) {

	file, _, err := req.FormFile("file_input")
	// file, _ := os.Open(f)

	client, err := storage.NewClient(context.Background())

	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	wc := client.Bucket("my-first-project-ce24e.appspot.com").Object("basePathimg").NewWriter(ctx)

	// Execute Upload

	if _, err := io.Copy(wc, file); err != nil {
		fmt.Println("io.Copy")
	}

	if err := wc.Close(); err != nil {
		fmt.Println("Writer.Close")
	}

}

func toBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func UploadFromRemoteUrl(fileName string, remoteUrl string) string {

	fmt.Println("uploading image to cloud storage")
	resp, err := http.Get(remoteUrl)
	remoteFile, err := ioutil.ReadAll(resp.Body)

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		fmt.Println("err", err)
	}

	defer client.Close()
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	wc := client.Bucket("my-first-project-ce24e.appspot.com").Object(fileName).NewWriter(ctx)
	if _, err = io.Copy(wc, bytes.NewReader(remoteFile)); err != nil {
		fmt.Println("err", err)
	}
	
	if err := wc.Close(); err != nil {
		fmt.Println("err", err)
	}

	fmt.Println("Uploaded File Successfully")
	// wait.Done()
	respUrl := "https://" + "storage.cloud.google.com/my-first-project-ce24e.appspot.com/" + fileName
	return respUrl
}
