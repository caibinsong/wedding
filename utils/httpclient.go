package utils

import (
	"bytes"
	"encoding/json"
<<<<<<< HEAD
	"fmt"
	"github.com/caibinsong/wedding/config"
=======
	"errors"
	"fmt"
	"github.com/Amniversary/wedding-logic-redpacket/config"
	"io/ioutil"
>>>>>>> dd12374ac95f08e4145cdb3fa4b628e5d98bd4d3
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
<<<<<<< HEAD
func (this *HttpClient) Post(url string, header map[string]string, request interface{}) (*http.Response, error) {
	//解析主体body信息
=======
func (this *HttpClient) Post(url string, header map[string]string, request interface{}) ([]byte, error) {
>>>>>>> dd12374ac95f08e4145cdb3fa4b628e5d98bd4d3
	reqBytes, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("post request to json err: %v", err)
	}
<<<<<<< HEAD

	//创建request
=======
	if url == "" {
		return nil, errors.New("url is \"\"")
	}
>>>>>>> dd12374ac95f08e4145cdb3fa4b628e5d98bd4d3
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("http new request error: %v", err)
	}
<<<<<<< HEAD

	//设置header头
=======
>>>>>>> dd12374ac95f08e4145cdb3fa4b628e5d98bd4d3
	req.Header.Set("Content-Type", "application/json")
	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}
<<<<<<< HEAD

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
	var request *config.GetWXUserInfo = &config.GetWXUserInfo{ActionName: "get_user_info", Data: "get_user_info"}
	var response *config.WXUserInfoResponse = &config.WXUserInfoResponse{}

	//发送post请求
	resp, err := this.Post(config.GetConfig().WXUserInfoUrl, header, request)
=======
	resp, err := this.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http do request error: %v", err)
	}
	defer resp.Body.Close()
	rspBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ioutil ReadAll error: %v", err)
	}
	return rspBody, nil
}

func (this *HttpClient) GetWXUserInfoResponse(userid int64) (*config.WXUserInfoResponse, error) {
	var header map[string]string = map[string]string{"ServerName": config.ServerName,
		"MethodName": config.WX_GetUserInfo,
		"userId":     fmt.Sprintf("%d", userid)}

	var request *config.GetWXUserInfo = &config.GetWXUserInfo{ActionName: "get_user_info", Data: "get_user_info"}

	var response *config.WXUserInfoResponse = &config.WXUserInfoResponse{}

	body, err := this.Post(config.WX_USER_INFO_URL, header, request)
>>>>>>> dd12374ac95f08e4145cdb3fa4b628e5d98bd4d3
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
<<<<<<< HEAD
	defer resp.Body.Close()

	//解析返回信息
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Println("json decode error: ", err)
		return nil, err
	}

	//判断是否成功
=======
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println("json decode error: ", err, string(body))
		return nil, err
	}
>>>>>>> dd12374ac95f08e4145cdb3fa4b628e5d98bd4d3
	if response.Code != 0 {
		log.Println("login robot result code error: ", response.Code, response.Msg)
		return nil, fmt.Errorf(response.Msg)
	}
	return response, nil
}

<<<<<<< HEAD
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
=======
func (this *HttpClient) GetWXUserListResponse(userlist []int) (*config.WXUserListResponse, error) {
	var header map[string]string = map[string]string{"ServerName": config.ServerName,
		"MethodName": config.WX_GetUserList}

	var request *config.GetWXUserList = &config.GetWXUserList{ActionName: "get_user_list"}
	request.Data.UserList = userlist
	var response *config.WXUserListResponse = &config.WXUserListResponse{}

	body, err := this.Post(config.WX_USER_LIST_URL, header, request)
>>>>>>> dd12374ac95f08e4145cdb3fa4b628e5d98bd4d3
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
<<<<<<< HEAD
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
=======
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println("json decode error: %v [%s]", err, string(body))
		return nil, err
	}
	if response.Code != 0 {
		log.Println("login robot result code error: %d %s", response.Code, response.Msg)
>>>>>>> dd12374ac95f08e4145cdb3fa4b628e5d98bd4d3
		return nil, err
	}
	return response, nil
}

<<<<<<< HEAD
//广播服务
func (this *HttpClient) RoomSvr(data map[string]interface{}) error {
	//头部信息
	var header map[string]string = map[string]string{"ServerName": config.RoomSvr_ServerName,
		"MethodName": config.RoomSvr_MethodName}

	//发送post请求
	resp, err := this.Post(config.GetConfig().RoomSvrUrl, header, data)
=======
func (this *HttpClient) RoomSvr(data map[string]interface{}) error {
	var header map[string]string = map[string]string{"ServerName": config.RoomSvr_ServerName,
		"MethodName": config.RoomSvr_MethodName}

	var response *config.Response = &config.Response{}

	body, err := this.Post(config.RoomSvr_URL, header, data)
>>>>>>> dd12374ac95f08e4145cdb3fa4b628e5d98bd4d3
	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("广播失败")
	}
<<<<<<< HEAD
	defer resp.Body.Close()

	//解析返回信息
	var response *config.Response = &config.Response{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Println("json decode error: ", err)
		return fmt.Errorf("广播失败")
	}

	//判断是否成功
=======
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println("json decode error: ", err, string(body))
		return fmt.Errorf("广播失败")
	}
>>>>>>> dd12374ac95f08e4145cdb3fa4b628e5d98bd4d3
	if response.Code != 0 {
		log.Println("login robot result code error: ", response.Code, response.Msg)
		return fmt.Errorf(response.Msg)
	}
	return nil
}
