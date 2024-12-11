package externalcontact

import (
	"fmt"
	"github.com/silenceper/wechat/v2/util"
)

const (
	// contactListURL 获取已服务的外部联系人
	contactListURL = "https://qyapi.weixin.qq.com/cgi-bin/externalcontact/contact_list?access_token=%s"
)

// ContactListRequest 获取已服务的外部联系人请求
type ContactListRequest struct {
	Cursor string `json:"cursor,omitempty"`
	Limit  int    `json:"limit,omitempty"`
}

// ContactListResponse 获取已服务的外部联系人响应
type ContactListResponse struct {
	util.CommonError
	InfoList   []ContactInfo `json:"info_list"`
	NextCursor string        `json:"next_cursor"`
}

// ContactInfo 外部联系人信息
type ContactInfo struct {
	IsCustomer     bool   `json:"is_customer"`
	TmpOpenID      string `json:"tmp_openid"`
	ExternalUserID string `json:"external_userid,omitempty"`
	Name           string `json:"name,omitempty"`
	FollowUserID   string `json:"follow_userid"`
	ChatID         string `json:"chat_id,omitempty"`
	ChatName       string `json:"chat_name,omitempty"`
	AddTime        int64  `json:"add_time"`
}

// GetContactList 获取已服务的外部联系人
// see https://developer.work.weixin.qq.com/document/path/92228
func (r *Client) GetContactList(req *ContactListRequest) (*ContactListResponse, error) {
	accessToken, err := r.GetAccessToken()
	if err != nil {
		return nil, err
	}
	var response []byte
	if response, err = util.PostJSON(fmt.Sprintf(contactListURL, accessToken), req); err != nil {
		return nil, err
	}
	result := &ContactListResponse{}
	err = util.DecodeWithError(response, result, "GetContactList")
	return result, err
}
