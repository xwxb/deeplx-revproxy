package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/xwxb/deeplx-revproxy/config"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

func sortEndpointsByWeightDesc(eps []config.Endpoint) []config.Endpoint {
	sort.Slice(eps, func(i, j int) bool {
		return eps[i].Weight > eps[j].Weight
	})
	return eps
}

// 处理url，如果里面没找到 translate 字符串，那么则在最后面加上 /translate
func normalizeURL(url string) string {
	if !strings.Contains(url, "translate") {
		url = url + "/translate"
	}
	return url

}

type Request struct {
	TransText   string `json:"text"`
	SourceLang  string `json:"source_lang"`
	TargetLang  string `json:"target_lang"`
	TagHandling string `json:"tag_handling"`
}

type NormalResponse struct {
	Code         int      `json:"code"`
	ID           int      `json:"id"`
	Data         string   `json:"data"`
	Alternatives []string `json:"alternatives"`
	SourceLang   string   `json:"source_lang"`
	TargetLang   string   `json:"target_lang"`
	Method       string   `json:"method"`
}

func ProxyHandler(c *gin.Context) {
	req := &Request{}
	err := c.Bind(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	eps := sortEndpointsByWeightDesc(config.Global.Endpoints)
	// 顺序遍历，如果有一个200成功就直接返回
	for _, ep := range eps {
		url := normalizeURL(ep.URL)
		log.Println("Proxying request to", url)

		// todo 内部实现一个重试 + 超时
		resp, err := sendRequest(url, req, ep.Timeout)
		if err == nil && resp.Code == 200 {
			c.JSON(http.StatusOK, resp)
			return
		}
		log.Println("Failed to proxy request to", url, "with error:", err, "response:", resp)
	}

	c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to proxy request, all endpoints are down."})
}

func sendRequest(url string, req *Request, timeout int) (*NormalResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	reqJsonBytes, err := json.Marshal(req)
	if err != nil {
		return &NormalResponse{
			Code: http.StatusInternalServerError,
		}, err
	}

	// 转发请求
	client := http.Client{}
	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqJsonBytes))
	if err != nil {
		return &NormalResponse{
			Code: http.StatusInternalServerError,
		}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(httpReq.WithContext(ctx))
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) { // 好像不管用？
			return &NormalResponse{
				Code: http.StatusRequestTimeout,
			}, err
		} else {
			return &NormalResponse{
				Code: http.StatusBadGateway,
			}, err
		}
	}
	defer resp.Body.Close()

	// 解析 resp 到 NormalResponse 结构体
	normalResp := &NormalResponse{}
	err = json.NewDecoder(resp.Body).Decode(normalResp)
	if err != nil {
		return &NormalResponse{
			Code: http.StatusInternalServerError,
		}, err
	}

	return normalResp, nil
}
