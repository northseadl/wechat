// Package device 设备相关接口
package device

import (
	"encoding/json"
	"fmt"

	"github.com/silenceper/wechat/v2/util"
)

const (
	// DeviceAdd 添加设备标识
	DeviceAdd = iota
	// DeviceUpgrade 更新设备标识
	DeviceUpgrade
)

type reqDeviceAuthorize struct {
	// 设备 id 的个数
	DeviceNum string `json:"device_num"`
	// 设备 id 的列表，json 的 array 格式，其 size 必须等于 device_num
	DeviceList []ReqDevice `json:"device_list"`
	// 请求操作的类型，限定取值为：0：设备授权（缺省值为 0）1：设备更新（更新已授权设备的各属性值）
	OpType string `json:"op_type,omitempty"`
	// 设备的产品编号（由微信硬件平台分配）。可在公众号设备功能管理页面查询。
	// 当 op_type 为‘0’，product_id 为‘1’时，不要填写 product_id 字段（会引起不必要错误）；
	// 当 op_typy 为‘0’，product_id 不为‘1’时，必须填写 product_id 字段；
	// 当 op_type 为 1 时，不要填写 product_id 字段。
	ProductID string `json:"product_id,omitempty"`
}

// ReqDevice 设备授权实体
type ReqDevice struct {
	// 设备的 device id
	ID string `json:"id"`
	// 设备的 mac 地址 格式采用 16 进制串的方式（长度为 12 字节），
	// 不需要 0X 前缀，如：1234567890AB
	Mac string `json:"mac"`
	//  支持以下四种连接协议：
	//	android classic bluetooth – 1
	//	ios classic bluetooth – 2
	//	ble – 3
	//	wifi -- 4
	//	一个设备可以支持多种连接类型，用符号"|"做分割，客户端优先选择靠前的连接方式（优先级按 | 关系的排序依次降低），举例：
	//	1：表示设备仅支持 andiod classic bluetooth 1|2：表示设备支持 android 和 ios 两种 classic bluetooth，但是客户端优先选择 android classic bluetooth 协议，如果 android classic bluetooth 协议连接失败，再选择 ios classic bluetooth 协议进行连接
	//	（注：安卓平台不同时支持 BLE 和 classic 类型）
	ConnectProtocol string `json:"connect_protocol"`
	// auth 及通信的加密 key，第三方需要将 key 烧制在设备上（128bit），格式采用 16 进制串的方式（长度为 32 字节），不需要 0X 前缀，如：1234567890ABCDEF1234567890ABCDEF
	AuthKey string `json:"auth_key"`
	// 断开策略，目前支持：1：退出公众号页面时即断开连接 2：退出公众号之后保持连接不断开
	CloseStrategy string `json:"close_strategy"`
	// 连接策略，32 位整型，按 bit 位置位，目前仅第 1bit 和第 3bit 位有效（bit 置 0 为无效，1 为有效；第 2bit 已被废弃），且 bit 位可以按或置位（如 1|4=5），各 bit 置位含义说明如下：
	// 1：（第 1bit 置位）在公众号对话页面，不停的尝试连接设备
	// 4：（第 3bit 置位）处于非公众号页面（如主界面等），微信自动连接。当用户切换微信到前台时，可能尝试去连接设备，连上后一定时间会断开
	ConnStrategy string `json:"conn_strategy"`
	// auth version，设备和微信进行 auth 时，会根据该版本号来确认 auth buf 和 auth key 的格式（各 version 对应的 auth buf 及 key 的具体格式可以参看“客户端蓝牙外设协议”），该字段目前支持取值：
	// 0：不加密的 version
	// 1：version 1
	AuthVer string `json:"auth_ver"`
	// 表示 mac 地址在厂商广播 manufacture data 里含有 mac 地址的偏移，取值如下：
	// -1：在尾部、
	// -2：表示不包含 mac 地址 其他：非法偏移
	ManuMacPos string `json:"manu_mac_pos"`
	// 表示 mac 地址在厂商 serial number 里含有 mac 地址的偏移，取值如下：
	// -1：表示在尾部
	// -2：表示不包含 mac 地址 其他：非法偏移
	SerMacPost string `json:"ser_mac_post"`
	// 精简协议类型，取值如下：计步设备精简协议：1（若该字段填 1，connect_protocol 必须包括 3。非精简协议设备切勿填写该字段）
	BleSimpleProtocol string `json:"ble_simple_protocol,omitempty"`
}

// ResBaseInfo 授权回调实体
type ResBaseInfo struct {
	BaseInfo struct {
		DeviceType string `json:"device_type"`
		DeviceID   string `json:"device_id"`
	} `json:"base_info"`
}

// 授权回调根信息
type resDeviceAuthorize struct {
	util.CommonError
	Resp []ResBaseInfo `json:"resp"`
}

// DeviceAuthorize 设备授权
func (d *Device) DeviceAuthorize(devices []ReqDevice, opType int, product string) (res []ResBaseInfo, err error) {
	var accessToken string
	accessToken, err = d.GetAccessToken()
	if err != nil {
		return
	}

	uri := fmt.Sprintf("%s?access_token=%s", uriAuthorize, accessToken)
	req := reqDeviceAuthorize{
		DeviceNum:  fmt.Sprintf("%d", len(devices)),
		DeviceList: devices,
		OpType:     fmt.Sprintf("%d", opType),
		ProductID:  product,
	}
	var response []byte
	response, err = util.PostJSON(uri, req)
	if err != nil {
		return nil, err
	}
	var result resDeviceAuthorize
	err = json.Unmarshal(response, &result)
	if err != nil {
		return
	}
	if result.ErrCode != 0 {
		err = fmt.Errorf("DeviceAuthorize Error , errcode=%d , errmsg=%s", result.ErrCode, result.ErrMsg)
		return
	}
	res = result.Resp
	return
}
