package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/caibinsong/wedding/config"
	"log"
	"net/http"
)

type HttpClient struct {
	client *http.Client
}

func NewHttpClient() *HttpClient {
	return &HttpClient{
		client: &http.Client{},
	}
}

//发送请求
func (this *HttpClient) Post(url string, header map[string]string, request interface{}) (*http.Response, error) {
	//解析主体body信息
	reqBytes, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("post request to json err: %v", err)
	}
	log.Println(string(reqBytes))
	//创建request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("http new request error: %v", err)
	}

	//设置header头
	req.Header.Set("Content-Type", "application/json")
	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}

	//发送请求
	return this.client.Do(req)
}

//微信用户信息
func (this *HttpClient) GetWXUserInfoResponse(userid int64) (*config.WXUserInfoResponse, error) {
	//header头部信息
	var header map[string]string = map[string]string{"ServerName": config.GetConfig().ServerName,
		"MethodName": config.WX_GetUserInfo,
		"userId":     fmt.Sprintf("%d", userid)}

	//body主体信息
	var data map[string]interface{} = map[string]interface{}{"user_id": userid, "app_id": 4}
	var request *config.GetWXUserInfo = &config.GetWXUserInfo{ActionName: "user_info", Data: data}
	var response *config.WXUserInfoResponse = &config.WXUserInfoResponse{}

	//发送post请求
	resp, err := this.Post(config.GetConfig().WXUserInfoUrl, header, request)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	//解析返回信息
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Println("json decode error: ", err)
		return nil, err
	}

	//判断是否成功
	if response.Code != 0 {
		log.Println("login robot result code error: ", response.Code, response.Msg)
		return nil, fmt.Errorf(response.Msg)
	}
	return response, nil
}

//微信用户信息列表
func (this *HttpClient) GetWXUserListResponse(userlist []int) (*config.WXUserListResponse, error) {
	//header头部信息
	var header map[string]string = map[string]string{"ServerName": config.GetConfig().ServerName,
		"MethodName": config.WX_GetUserList}

	//body主体信息
	var request *config.GetWXUserList = &config.GetWXUserList{ActionName: "get_user_list"}
	request.Data.UserList = userlist

	//发送post请求
	resp, err := this.Post(config.GetConfig().WXUserListUrl, header, request)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	//解析返回信息
	var response *config.WXUserListResponse = &config.WXUserListResponse{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Println("json decode error: ", err.Error())
		return nil, err
	}

	//判断是否成功
	if response.Code != 0 {
		log.Println("login robot result code error: ", response.Code, response.Msg)
		return nil, err
	}
	return response, nil
}

//广播服务
func (this *HttpClient) RoomSvr(serverName, methodname string, data map[string]interface{}) error {
	//头部信息
	var header map[string]string = map[string]string{"ServerName": serverName,
		"MethodName": methodname}

	//发送post请求
	resp, err := this.Post(config.GetConfig().RoomSvrUrl, header, data)
	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("广播失败")
	}
	defer resp.Body.Close()

	//解析返回信息
	var response *config.Response = &config.Response{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Println("json decode error: ", err)
		return fmt.Errorf("广播失败")
	}

	//判断是否成功
	if response.Code != 0 {
		log.Println("login robot result code error: ", response.Code, response.Msg)
		return errors.New(fmt.Sprint(response.Code))
	}
	return nil
}

//广播服务
func (this *HttpClient) AccessCtrlSvr(serverName, methodname string, data map[string]interface{}) error {
	//头部信息
	var header map[string]string = map[string]string{"ServerName": serverName,
		"MethodName": methodname}

	//发送post请求
	resp, err := this.Post("http://172.17.0.13:7777/rpc", header, data)
	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("广播失败")
	}
	defer resp.Body.Close()
	//解析返回信息
	var response *config.ResponseRoomSvr = &config.ResponseRoomSvr{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Println("json decode has error: ", err)
		return fmt.Errorf("广播失败")
	}

	//判断是否成功
	if response.Code != 0 {
		log.Println("login robot result code error: ", response.Code, response.Data)
		return errors.New(fmt.Sprint(response.Code))
	}
	return nil
}
