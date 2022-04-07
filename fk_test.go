package ishumei

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestNewBaseURL(t *testing.T) {
	businessName := "text"
	region := "bj"
	secure := false
	u, _ := NewBaseURL(businessName,region,secure)
	fmt.Println(u)
}


//数美文本送审
func TestSmService_TextAuditingJob(t *testing.T) {
    ctx := context.TODO()
	businessName := "text"
	region := "bj"
	secure := false
	u, _ := NewBaseURL(businessName,region,secure)
	url := &BaseURL{
		SMURL: u,
	}
	client := NewClient(url,&http.Client{
		Timeout: time.Second * 15,
	})
   p,r,e := client.SM.TextAuditingJob(ctx,&TextAuditingJobOptions{
   	   AccessKey: "*****************",
   	   AppID: "default",
	   EventID: "title",
	   Type: "DEFAULT",
	   BusinessType: "",
	   Data: &TextData {
   	   	 Text: "****",
   	   	 TokenID: "111111",
   	   	 Lang: "auto",
	   },
   })
   fmt.Println(p,r,e)
}

//数美同步单张图片送审
func TestSmService_ImageSingleAuditing(t *testing.T) {
	ctx := context.TODO()
	businessName := "img"
	region := "bj"
	secure := false
	u, _ := NewBaseURL(businessName,region,secure)
	url := &BaseURL{
		SMURL: u,
	}
	client := NewClient(url,&http.Client{
		Timeout: time.Second * 15,
	})
	p,r,e := client.SM.ImageSingleAuditing(ctx,&ImageRecognitionOptions{
		AccessKey: "*****************",
		AppID: "default",
		EventID: "default",
		Type: "POLITICS_PORN_BAN_VIOLENCE",
		Data: &ImageData {
			Img: "https://gimg2.baidu.com/image_search/src=http%3A%2F%2Fimages.china.cn%2Fattachement%2Fjpg%2Fsite1000%2F20150206%2F002564ba9cf7163eaa825a.jpg&refer=http%3A%2F%2Fimages.china.cn&app=2002&size=f9999,10000&q=a80&n=0&g=0n&fmt=auto?sec=1651915033&t=a13b3db6f448b2174f64b6efe5a015cb",
		},
	})
	fmt.Println(p,r,e)
}


