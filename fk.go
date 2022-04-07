package ishumei

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"github.com/mozillazg/go-httpheader"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"text/template"
	"time"
)

const (
	// Version current go sdk version
	Version               = "0.0.1"
	userAgent             = "shumei-sdk-v4/" + Version
	contentTypeJson        = "application/json"
)

// BaseURL 访问各 API 所需的基础 URL
type BaseURL struct {
	//访问 数美风控 的基础 URL
	SMURL *url.URL
}

// Response API 响应
type Response struct {
	*http.Response
}

func newResponse(resp *http.Response) *Response {
	return &Response{
		Response: resp,
	}
}

/**
文本审核：	http://api-text-bj.fengkongcloud.com/text/v4
视频审核：	http://api-video-bj.fengkongcloud.com/video/v4
视频流审核：	http://api-videostream-bj.fengkongcloud.com/videostream/v4
音频审核：	http://api-audio-bj.fengkongcloud.com/audio/v4
音频流审核： 	http://api-audiostream-xjp.fengkongcloud.com/audiostream/v4
图片审核：	http://api-img-bj.fengkongcloud.com/image/v4
 */
var baseURLTemplate = template.Must(
	template.New("baseURLTemplate").Parse(
		"{{.Schema}}://api-{{.BusinessName}}-{{.Region}}.fengkongcloud.com",
	),
)

func NewBaseURL(businessName, region string, secure bool) (*url.URL, error) {
	schema := "https"
	if !secure {
		schema = "http"
	}

	if region == "" {
		return nil, fmt.Errorf("region[%v] is invalid", region)
	}
	if businessName == "" {
		return nil, fmt.Errorf("businessName[%v] is invalid", businessName)
	}
	w := bytes.NewBuffer(nil)
	baseURLTemplate.Execute(w, struct {
		Schema     string
		BusinessName string
		Region     string
	}{
		schema, businessName, region,
	})
	u, _ := url.Parse(w.String())
	return u, nil
}

// NewClient returns a new SM API client.
func NewClient(uri *BaseURL, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	baseURL := &BaseURL{}
	if uri != nil {
		baseURL.SMURL = uri.SMURL
	}
	c := &Client{
		client:    httpClient,
		UserAgent: userAgent,
		BaseURL:   baseURL,
		Conf: &Config{
			EnableCRC:        true,
			RequestBodyClose: false,
			RetryOpt: RetryOptions{
				Count:    3,
				Interval: time.Duration(0),
			},
		},
	}
	c.common.client = c
	c.SM = (*SmService)(&c.common)
	return c
}

func (c *Client) newRequest(ctx context.Context, baseURL *url.URL, uri, method string, body interface{}, optQuery interface{}, optHeader interface{}) (req *http.Request, err error) {
	uri, err = addURLOptions(uri, optQuery)
	if err != nil {
		return
	}
	u, _ := url.Parse(uri)
	urlStr := baseURL.ResolveReference(u).String()

	var reader io.Reader
	contentType := "application/json"
	if body != nil {
		// 上传文件
		if r, ok := body.(io.Reader); ok {
			reader = r
		} else {
			b, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			reader = bytes.NewReader(b)
		}
	}

	req, err = http.NewRequest(method, urlStr, reader)
	if err != nil {
		return
	}

	req.Header, err = addHeaderOptions(req.Header, optHeader)
	if err != nil {
		return
	}
	if v := req.Header.Get("Content-Length"); req.ContentLength == 0 && v != "" && v != "0" {
		req.ContentLength, _ = strconv.ParseInt(v, 10, 64)
	}

	if v := req.Header.Get("User-Agent"); v == "" || !strings.HasPrefix(v, userAgent) {
		if c.UserAgent != "" {
			req.Header.Set("User-Agent", c.UserAgent)
		}
	}

	if req.Header.Get("Content-Type") == "" && contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	if c.Host != "" {
		req.Host = c.Host
	}
	if c.Conf.RequestBodyClose {
		req.Close = true
	}
	return
}

func (c *Client) doAPI(ctx context.Context, req *http.Request, result interface{}, closeBody bool) (*Response, error) {
	var cancel context.CancelFunc
	if closeBody {
		ctx, cancel = context.WithCancel(ctx)
		defer cancel()
	}
	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		return nil, err
	}

	defer func() {
		if closeBody {
			// Close the body to let the Transport reuse the connection
			io.Copy(ioutil.Discard, resp.Body)
			resp.Body.Close()
		}
	}()

	response := newResponse(resp)

	err = checkResponse(resp)
	if err != nil {
		if !closeBody {
			resp.Body.Close()
		}
		return response, err
	}
	if result != nil {
		if w, ok := result.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(result)
			if err == io.EOF {
				err = nil
			}
		}
	}
	return response, err
}

type sendOptions struct {
	// 基础 URL
	baseURL *url.URL
	// URL 中除基础 URL 外的剩余部分
	uri string
	// 请求方法
	method string

	body interface{}
	// url 查询参数
	optQuery interface{}
	// http header 参数
	optHeader interface{}
	// 用 result 反序列化 resp.Body
	result interface{}
	// 是否禁用自动调用 resp.Body.Close()
	// 自动调用 Close() 是为了能够重用连接
	disableCloseBody bool
}

func (c *Client) doRetry(ctx context.Context, opt *sendOptions) (resp *Response, err error) {
	if opt.body != nil {
		if _, ok := opt.body.(io.Reader); ok {
			resp, err = c.send(ctx, opt)
			return
		}
	}
	count := 1
	if count < c.Conf.RetryOpt.Count {
		count = c.Conf.RetryOpt.Count
	}
	nr := 0
	interval := c.Conf.RetryOpt.Interval
	for nr < count {
		resp, err = c.send(ctx, opt)
		if err != nil {
			if resp != nil && resp.StatusCode <= 499 {
				dobreak := true
				for _, v := range c.Conf.RetryOpt.StatusCode {
					if resp.StatusCode == v {
						dobreak = false
						break
					}
				}
				if dobreak {
					break
				}
			}
			nr++
			if interval > 0 && nr < count {
				time.Sleep(interval)
			}
			continue
		}
		break
	}
	return

}

func (c *Client) send(ctx context.Context, opt *sendOptions) (resp *Response, err error) {
	req, err := c.newRequest(ctx, opt.baseURL, opt.uri, opt.method, opt.body, opt.optQuery, opt.optHeader)
	if err != nil {
		return
	}
	resp, err = c.doAPI(ctx, req, opt.result, !opt.disableCloseBody)
	return
}


// addURLOptions adds the parameters in opt as URL query parameters to s. opt
// must be a struct whose fields may contain "url" tags.
func addURLOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	// 保留原有的参数，并且放在前面。因为 cos 的 url 路由是以第一个参数作为路由的
	// e.g. /?uploads
	q := u.RawQuery
	rq := qs.Encode()
	if q != "" {
		if rq != "" {
			u.RawQuery = fmt.Sprintf("%s&%s", q, qs.Encode())
		}
	} else {
		u.RawQuery = rq
	}
	return u.String(), nil
}

// addHeaderOptions adds the parameters in opt as Header fields to req. opt
// must be a struct whose fields may contain "header" tags.
func addHeaderOptions(header http.Header, opt interface{}) (http.Header, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return header, nil
	}

	h, err := httpheader.Header(opt)
	if err != nil {
		return nil, err
	}

	for key, values := range h {
		for _, value := range values {
			header.Add(key, value)
		}
	}
	return header, nil
}

type RetryOptions struct {
	Count      int
	Interval   time.Duration
	StatusCode []int
}

type Config struct {
	EnableCRC        bool
	RequestBodyClose bool
	RetryOpt         RetryOptions
}

type Client struct {
	client *http.Client

	BaseURL   *BaseURL

	Host      string
	UserAgent string

	common service

	SM *SmService

	Conf *Config
}

type service struct {
	client *Client
}
