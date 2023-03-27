package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/resty.v1"
)

func main() {
	type tokenResp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			AccessToken           string `json:"access_token"`
			ExpiresIn             int    `json:"expires_in"`
			RefreshToken          string `json:"refresh_token"`
			RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
		} `json:"data"`
	}
	var (
		openApiUrlPrefix string = "https://ad.oceanengine.com/open_api/"
		uri              string = "oauth2/refresh_token/"
		// 请求Header
		contentType string = "application/json"
		// 请求参数
		appId        int64  = 1760961623693339                           // 开发者申请的应用APP_ID
		secret       string = "5438474df9c94fb8f4ba7ed6a2552c2ef795022d" // 开发者应用的私钥Secret
		grantType    string = "refresh_token"                            // 授权类型
		refreshToken string = "8b05d240318c3a78078668496e74e6b8ce39ba5e" // 刷新token
		// tokenResp
		ttTokenResp tokenResp
	)
	url := fmt.Sprintf("%s%s", openApiUrlPrefix, uri)

	resp, err := resty.New().SetRetryCount(3).R().
		SetHeaders(map[string]string{
			"Content-Type": contentType,
		}).
		SetFormData(map[string]string{
			"app_id":        fmt.Sprintf("%d", appId),
			"secret":        secret,
			"grant_type":    grantType,
			"refresh_token": refreshToken,
		}).
		Post(url)
	if err != nil {
		fmt.Println("Post err", err)
		return
	}
	fmt.Println("code:", resp.StatusCode())
	err = json.Unmarshal(resp.Body(), &ttTokenResp)
	if err != nil {
		fmt.Println("Unmarshal err:", err)
		return
	}

	fmt.Println("ttTokenResp:", ttTokenResp) // ttTokenResp: {0 OK {e88f206ab28a97ef494b853982d81739b81a1e37 86399 0d3e3b3a2c2b4bd342c536ed4b8a4ec0470db13c 2591999}}
	//"data": {
	//	"access_token": "e88f206ab28a97ef494b853982d81739b81a1e37",
	//		"expires_in": 86400,	// 24h
	//		"refresh_token": "0d3e3b3a2c2b4bd342c536ed4b8a4ec0470db13c",
	//		"refresh_token_expires_in": 2600000,	// 30天
	//}
}
