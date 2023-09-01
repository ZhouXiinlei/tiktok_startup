package tikcos

import (
	"context"
	"github.com/tencentyun/cos-go-sdk-v5"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
)

type TengxunyunCfg struct {
	Endpoint  string
	SecretID  string
	SecretKey string
}

func TengxunyunInit(TengxunyunCfg TengxunyunCfg) *cos.Client {
	u, _ := url.Parse(TengxunyunCfg.Endpoint)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{Transport: &cos.AuthorizationTransport{
		//填写用户账号密钥信息，也可以设置为环境变量
		SecretID:  os.Getenv(TengxunyunCfg.SecretID),
		SecretKey: os.Getenv(TengxunyunCfg.SecretKey),
	}})
	return client
}
func UploadVideoToCos(c *cos.Client, objectName string, reader multipart.File) (bool, error) {

	_, err := c.Object.Put(context.Background(), objectName, reader, nil)

	if err != nil {
		return true, nil
	} else {
		return false, err
	}
}
