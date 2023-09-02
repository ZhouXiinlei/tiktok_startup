package tikcos

import (
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
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
		SecretID:  TengxunyunCfg.SecretID,
		SecretKey: TengxunyunCfg.SecretKey,
	}})
	//_, err := client.Bucket.Put(context.Background(), nil)
	//if err != nil {
	//	panic(err)
	//}
	return client
}
