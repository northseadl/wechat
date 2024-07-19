// Package user migrate 用于微信公众号账号迁移，获取 openID 变化
// 参考文档：https://kf.qq.com/faq/1901177NrqMr190117nqYJze.html
package user

import (
	"errors"
	"fmt"

	"github.com/silenceper/wechat/v2/util"
)

const (
	changeOpenIDURL = "https://api.weixin.qq.com/cgi-bin/changeopenid"
)

// ChangeOpenIDResult OpenID 迁移变化
type ChangeOpenIDResult struct {
	OriOpenID string `json:"ori_openid"`
	NewOpenID string `json:"new_openid"`
	ErrMsg    string `json:"err_msg,omitempty"`
}

// ChangeOpenIDResultList OpenID 迁移变化列表
type ChangeOpenIDResultList struct {
	util.CommonError
	List []ChangeOpenIDResult `json:"result_list"`
}

// ListChangeOpenIDs 返回指定 OpenID 变化列表
// fromAppID 为老账号 AppID
// openIDs 为老账号的 openID，openIDs 限 100 个以内
// AccessToken 为新账号的 AccessToken
func (user *User) ListChangeOpenIDs(fromAppID string, openIDs ...string) (list *ChangeOpenIDResultList, err error) {
	list = &ChangeOpenIDResultList{}
	// list.List = make([]ChangeOpenIDResult, 0)
	if len(openIDs) > 100 {
		err = errors.New("openIDs length must be lt 100")
		return
	}

	if fromAppID == "" {
		err = errors.New("fromAppID is required")
		return
	}

	accessToken, err := user.GetAccessToken()
	if err != nil {
		return
	}

	uri := fmt.Sprintf("%s?access_token=%s", changeOpenIDURL, accessToken)
	var resp []byte
	var req struct {
		FromAppID  string   `json:"from_appid"`
		OpenidList []string `json:"openid_list"`
	}
	req.FromAppID = fromAppID
	req.OpenidList = append(req.OpenidList, openIDs...)
	resp, err = util.PostJSON(uri, req)
	if err != nil {
		return
	}

	err = util.DecodeWithError(resp, list, "ListChangeOpenIDs")
	return
}

// ListAllChangeOpenIDs  返回所有用户 OpenID 列表
// fromAppID 为老账号 AppID
// openIDs 为老账号的 openID
// AccessToken 为新账号的 AccessToken
func (user *User) ListAllChangeOpenIDs(fromAppID string, openIDs ...string) (list []ChangeOpenIDResult, err error) {
	list = make([]ChangeOpenIDResult, 0)
	chunks := util.SliceChunk(openIDs, 100)
	for _, chunk := range chunks {
		result, err := user.ListChangeOpenIDs(fromAppID, chunk...)
		if err != nil {
			return list, err
		}
		list = append(list, result.List...)
	}
	return
}
