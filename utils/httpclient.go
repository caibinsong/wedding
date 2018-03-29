package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Amniversary/wedding-logic-redpacket/config"
	"io/ioutil"
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
func (this *HttpClient) Post(url string, header map[string]string, request interface{}) ([]byte, error) {
	reqBytes, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("post request to json err: %v", err)
	}
	if url == "" {
		return nil, errors.New("url is \"\"")
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("http new request error: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}
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
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println("json decode error: ", err, string(body))
		return nil, err
	}
	if response.Code != 0 {
		log.Println("login robot result code error: ", response.Code, response.Msg)
		return nil, fmt.Errorf(response.Msg)
	}
	return response, nil
}

func (this *HttpClient) GetWXUserListResponse(userlist []int) (*config.WXUserListResponse, error) {
	var header map[string]string = map[string]string{"ServerName": config.ServerName,
		"MethodName": config.WX_GetUserList}

	var request *config.GetWXUserList = &config.GetWXUserList{ActionName: "get_user_list"}
	request.Data.UserList = userlist
	var response *config.WXUserListResponse = &config.WXUserListResponse{}

	body, err := this.Post(config.WX_USER_LIST_URL, header, request)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println("json decode error: %v [%s]", err, string(body))
		return nil, err
	}
	if response.Code != 0 {
		log.Println("login robot result code error: %d %s", response.Code, response.Msg)
		return nil, err
	}
	return response, nil
}

func (this *HttpClient) RoomSvr(data map[string]interface{}) error {
	var header map[string]string = map[string]string{"ServerName": config.RoomSvr_ServerName,
		"MethodName": config.RoomSvr_MethodName}

	var response *config.Response = &config.Response{}

	body, err := this.Post(config.RoomSvr_URL, header, data)
	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("广播失败")
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println("json decode error: ", err, string(body))
		return fmt.Errorf("广播失败")
	}
	if response.Code != 0 {
		log.Println("login robot result code error: ", response.Code, response.Msg)
		return fmt.Errorf(response.Msg)
	}
	return nil
}
