package http

import (
	"context"
	"encoding/json"
	"net/http"
	httpURL "net/url"
	"strings"

	"github.com/akatsukisun2020/wx_common/logger"
)

// 封装http的基本方法，参考： https://juejin.cn/post/7014046220976914468
type optionFunc func(*HttpClient)

// WithSetParamOption 设置get的参数
func WithSetParamOption(key, value string) optionFunc {
	return func(cli *HttpClient) {
		if cli.param == nil {
			cli.param = httpURL.Values{}
		}
		cli.param.Set(key, value)
	}
}

// WithSetBodyOption 设置set方法的参数
func WithSetBodyOption(body interface{}) optionFunc {
	return func(cli *HttpClient) {
		data, _ := json.Marshal(body)
		cli.body = string(data)
	}
}

type HttpClient struct {
	url   string         // 需要访问的目的地址
	param httpURL.Values // get方法的参数
	body  string         // post方法的负载
}

func NewHttpClient(url string, opts ...optionFunc) *HttpClient {
	httpCli := &HttpClient{
		url:   url,
		param: httpURL.Values{},
		body:  "",
	}

	for _, opt := range opts {
		opt(httpCli)
	}

	return httpCli
}

// Get 带参数的Get方法
func (cli *HttpClient) Get(ctx context.Context) error {
	geturl := cli.url + "/get"
	u, err := httpURL.ParseRequestURI(geturl)
	if err != nil {
		logger.Errorf("parse url requestUrl failed,err:%v", err)
		return err
	}
	u.RawQuery = cli.param.Encode()
	resp, err := http.Get(u.String())
	if err != nil {
		logger.Errorf("Get failed, err:%v", err)
		return err
	}
	logger.Debugf("Get.resp:%v", resp) // TODO 需要完善，将数据发送出去
	return nil
}

// Post ...
func (cli *HttpClient) Post(ctx context.Context) error {
	contentType := "application/json"
	resp, err := http.Post(cli.url, contentType, strings.NewReader(cli.body)) // todo：搜索 "http.Post" ==> 看etcd中的httpcli是怎么写的.
	if err != nil {
		logger.Errorf("post failed, err:%v\n", err)
		return err
	}
	logger.Debugf("Post.resp:%v", resp) // TODO 需要完善，将数据发送出去
	return nil
}
