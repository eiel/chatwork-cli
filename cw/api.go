package main

import (
    "net/http"
    "net/url"
    "io/ioutil"
    "os"
    "strings"
)

const (
    DEFAULT_HOST    = "api.chatwork.com"
    DEFAULT_VERSION = "v2"
    DEFAULT_TOKEN_ENV   = "CW_API_TOKEN"
)

// APIへのリクエストに必要な情報を集めた構造体
type CwApi struct {
    // HTTPメソッド
    Method string

    // APIのホスト
    Host string

    // APIバージョン
    Version string

    // エンドポイントのパスまでの配列
    Paths []string

    // リクエストパラメタ
    Param url.Values

    // リクエストに認証情報をつけるオブジェクト
    Auth CwApiAuthorizer
}

func NewCwApi() *CwApi {
    api := CwApi{}
    api.Host    = DEFAULT_HOST
    api.Version = DEFAULT_VERSION
    return &api
}

// http.Requestをつくる
func (ca *CwApi) toRequest() (*http.Request, error) {
    url := "https://" + ca.Host + "/" + ca.Version + "/" + strings.Join(ca.Paths, "/")
    req, err := http.NewRequest(ca.Method, url, nil)

    if ca.Param != nil {
        query := ca.Param.Encode()
        if strings.ToUpper(ca.Method) == "GET" {
            req.URL.RawQuery = query
        } else {
            req.Body = ioutil.NopCloser(strings.NewReader(query))
        }
    }

    if err != nil {
        return req, err
    }
    ca.Auth.Authorize(req)
    return req, nil
}


// 何かの方法でリクエストに認証情報をつけるオブジェクトを示すinterface
type CwApiAuthorizer interface {
    Authorize(r *http.Request)
}

// 環境変数からAPIトークンを読み取るAuthorizerの実装
type TokenFromEnvAuthorizer struct {
    EnvName string
}

func (ta *TokenFromEnvAuthorizer) Authorize(r *http.Request) {
    token := os.Getenv(ta.EnvName)
    r.Header.Add("X-ChatWorkToken", token)
}
