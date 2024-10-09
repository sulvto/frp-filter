package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
)

// 定义响应数据的结构
type IPInfoResponse struct {
    Code  string `json:"code"`
    Msg   string `json:"msg"`
    Data  struct {
        ContinentCN  string `json:"continentCN"`
        CountryCN    string `json:"countryCN"`
        ZoneCN       string `json:"zoneCN"`
        ProvinceCN   string `json:"provinceCN"`
        CityCN       string `json:"cityCN"`
        CountyCN     string `json:"countyCN"`
        TownCN       string `json:"townCN"`
        IspCN        string `json:"ispCN"`
        ContinentID  int    `json:"continentID"`
        CountryID    int    `json:"countryID"`
        ZoneID       int    `json:"zoneID"`
        ProvinceID   int    `json:"provinceID"`
        CityID       int    `json:"cityID"`
        CountyID     int    `json:"countyID"`
        IspID        int    `json:"ispID"`
        TownID        int    `json:"townID"`
        Latitude     string `json:"latitude"`
        Longitude    string `json:"longitude"`
        OverseasRegion bool   `json:"overseasRegion"`
    } `json:"data"`
}

// 获取 IP 信息的工具函数
func GetIPInfo(ip string) (*IPInfoResponse, error) {
    // 构建请求的 URL
    url := fmt.Sprintf("https://mesh.if.iqiyi.com/aid/ip/info?version=1.1.1&ip=%s", ip)

    // 创建 HTTP 客户端
    client := &http.Client{}

    // 发起 GET 请求
    req, err := http.NewRequest(http.MethodGet, url, nil)
    if err != nil {
        return nil, err
    }

    // 发送请求
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close() // 确保关闭响应体

    // 读取响应体
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    // 将 JSON 转换为 Go 的数据结构
    var ipInfoResponse IPInfoResponse
    err = json.Unmarshal(body, &ipInfoResponse)
    if err != nil {
        return nil, err
    }

    return &ipInfoResponse, nil
}