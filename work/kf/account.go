package kf

import (
	"encoding/json"
	"fmt"

	"github.com/silenceper/wechat/v2/util"
)

const (
	// 添加客服账号
	accountAddAddr = "https://qyapi.weixin.qq.com/cgi-bin/kf/account/add?access_token=%s"
	// 删除客服账号
	accountDelAddr = "https://qyapi.weixin.qq.com/cgi-bin/kf/account/del?access_token=%s"
	// 修改客服账号
	accountUpdateAddr = "https://qyapi.weixin.qq.com/cgi-bin/kf/account/update?access_token=%s"
	// 获取客服账号列表
	accountListAddr = "https://qyapi.weixin.qq.com/cgi-bin/kf/account/list?access_token=%s"
	// 获取客服账号链接
	addContactWayAddr = "https://qyapi.weixin.qq.com/cgi-bin/kf/add_contact_way?access_token=%s"
)

// AccountAddOptions 添加客服账号请求参数
type AccountAddOptions struct {
	Name    string `json:"name"`     // 客服帐号名称，不多于 16 个字符
	MediaID string `json:"media_id"` // 客服头像临时素材。可以调用上传临时素材接口获取，不多于 128 个字节
}

// AccountAddSchema 添加客服账号响应内容
type AccountAddSchema struct {
	util.CommonError
	OpenKFID string `json:"open_kfid"` // 新创建的客服张号 ID
}

// AccountAdd 添加客服账号
func (r *Client) AccountAdd(options AccountAddOptions) (info AccountAddSchema, err error) {
	var (
		accessToken string
		data        []byte
	)
	if accessToken, err = r.ctx.GetAccessToken(); err != nil {
		return
	}
	if data, err = util.PostJSON(fmt.Sprintf(accountAddAddr, accessToken), options); err != nil {
		return
	}
	if err = json.Unmarshal(data, &info); err != nil {
		return
	}
	if info.ErrCode != 0 {
		return info, NewSDKErr(info.ErrCode, info.ErrMsg)
	}
	return info, nil
}

// AccountDelOptions 删除客服账号请求参数
type AccountDelOptions struct {
	OpenKFID string `json:"open_kfid"` // 客服帐号 ID, 不多于 64 字节
}

// AccountDel 删除客服账号
func (r *Client) AccountDel(options AccountDelOptions) (info util.CommonError, err error) {
	var (
		accessToken string
		data        []byte
	)
	if accessToken, err = r.ctx.GetAccessToken(); err != nil {
		return
	}
	if data, err = util.PostJSON(fmt.Sprintf(accountDelAddr, accessToken), options); err != nil {
		return
	}
	if err = json.Unmarshal(data, &info); err != nil {
		return
	}
	if info.ErrCode != 0 {
		return info, NewSDKErr(info.ErrCode, info.ErrMsg)
	}
	return info, nil
}

// AccountUpdateOptions 修改客服账号请求参数
type AccountUpdateOptions struct {
	OpenKFID string `json:"open_kfid"` // 客服帐号 ID, 不多于 64 字节
	Name     string `json:"name"`      // 客服帐号名称，不多于 16 个字符
	MediaID  string `json:"media_id"`  // 客服头像临时素材。可以调用上传临时素材接口获取，不多于 128 个字节
}

// AccountUpdate 修复客服账号
func (r *Client) AccountUpdate(options AccountUpdateOptions) (info util.CommonError, err error) {
	var (
		accessToken string
		data        []byte
	)
	if accessToken, err = r.ctx.GetAccessToken(); err != nil {
		return
	}
	if data, err = util.PostJSON(fmt.Sprintf(accountUpdateAddr, accessToken), options); err != nil {
		return
	}
	if err = json.Unmarshal(data, &info); err != nil {
		return
	}
	if info.ErrCode != 0 {
		return info, NewSDKErr(info.ErrCode, info.ErrMsg)
	}
	return info, nil
}

// AccountInfoSchema 客服详情
type AccountInfoSchema struct {
	OpenKFID string `json:"open_kfid"` // 客服帐号 ID
	Name     string `json:"name"`      // 客服帐号名称
	Avatar   string `json:"avatar"`    // 客服头像 URL
}

// AccountListSchema 获取客服账号列表响应内容
type AccountListSchema struct {
	util.CommonError
	AccountList []AccountInfoSchema `json:"account_list"` // 客服账号列表
}

// AccountList 获取客服账号列表
func (r *Client) AccountList() (info AccountListSchema, err error) {
	var (
		accessToken string
		data        []byte
	)
	if accessToken, err = r.ctx.GetAccessToken(); err != nil {
		return
	}
	if data, err = util.HTTPGet(fmt.Sprintf(accountListAddr, accessToken)); err != nil {
		return
	}
	if err = json.Unmarshal(data, &info); err != nil {
		return
	}
	if info.ErrCode != 0 {
		return info, NewSDKErr(info.ErrCode, info.ErrMsg)
	}
	return info, nil
}

// AddContactWayOptions 获取客服账号链接
// 1.若 scene 非空，返回的客服链接开发者可拼接 scene_param=SCENE_PARAM 参数使用，用户进入会话事件会将 SCENE_PARAM 原样返回。其中 SCENE_PARAM 需要 urlencode，且长度不能超过 128 字节。
// 如 https://work.weixin.qq.com/kf/kfcbf8f8d07ac7215f?enc_scene=ENCGFSDF567DF&scene_param=a%3D1%26b%3D2
// 2.历史调用接口返回的客服链接（包含 encScene=XXX 参数），不支持 scene_param 参数。
// 3.返回的客服链接，不能修改或复制参数到其他链接使用。否则进入会话事件参数校验不通过，导致无法回调。
type AddContactWayOptions struct {
	OpenKFID string `json:"open_kfid"` // 客服帐号 ID, 不多于 64 字节
	Scene    string `json:"scene"`     // 场景值，字符串类型，由开发者自定义，不多于 32 字节，字符串取值范围 (正则表达式)：[0-9a-zA-Z_-]*
}

// AddContactWaySchema 获取客服账号链接响应内容
type AddContactWaySchema struct {
	util.CommonError
	URL string `json:"url"` // 客服链接，开发者可将该链接嵌入到 H5 页面中，用户点击链接即可向对应的微信客服帐号发起咨询。开发者也可根据该 url 自行生成需要的二维码图片
}

// AddContactWay 获取客服账号链接
func (r *Client) AddContactWay(options AddContactWayOptions) (info AddContactWaySchema, err error) {
	var (
		accessToken string
		data        []byte
	)
	if accessToken, err = r.ctx.GetAccessToken(); err != nil {
		return
	}
	if data, err = util.PostJSON(fmt.Sprintf(addContactWayAddr, accessToken), options); err != nil {
		return
	}
	if err = json.Unmarshal(data, &info); err != nil {
		return
	}
	if info.ErrCode != 0 {
		return info, NewSDKErr(info.ErrCode, info.ErrMsg)
	}
	return info, nil
}
