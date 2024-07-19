package syncmsg

// BaseMessage 接收消息
type BaseMessage struct {
	MsgID              string `json:"msgid"`           // 消息 ID
	OpenKFID           string `json:"open_kfid"`       // 客服帐号 ID（msgtype 为 event，该字段不返回）
	ExternalUserID     string `json:"external_userid"` // 客户 UserID（msgtype 为 event，该字段不返回）
	ReceptionistUserID string `json:"servicer_userid"` // 接待客服 userID
	SendTime           uint64 `json:"send_time"`       // 消息发送时间
	Origin             uint32 `json:"origin"`          // 消息来源。3-微信客户发送的消息 4-系统推送的事件消息 5-接待人员在企业微信客户端发送的消息
}

// Text 文本消息
type Text struct {
	BaseMessage
	MsgType string `json:"msgtype"` // 消息类型，此时固定为：text
	Text    struct {
		Content string `json:"content"` // 文本内容
		MenuID  string `json:"menu_id"` // 客户点击菜单消息，触发的回复消息中附带的菜单 ID
	} `json:"text"` // 文本消息
}

// Image 图片消息
type Image struct {
	BaseMessage
	MsgType string `json:"msgtype"` // 消息类型，此时固定为：image
	Image   struct {
		MediaID string `json:"media_id"` // 图片文件 ID
	} `json:"image"` // 图片消息
}

// Voice 语音消息
type Voice struct {
	BaseMessage
	MsgType string `json:"msgtype"` // 消息类型，此时固定为：voice
	Voice   struct {
		MediaID string `json:"media_id"` // 语音文件 ID
	} `json:"voice"` // 语音消息
}

// Video 视频消息
type Video struct {
	BaseMessage
	MsgType string `json:"msgtype"` // 消息类型，此时固定为：video
	Video   struct {
		MediaID string `json:"media_id"` // 文件 ID
	} `json:"video"` // 视频消息
}

// File 文件消息
type File struct {
	BaseMessage
	MsgType string `json:"msgtype"` // 消息类型，此时固定为：file
	File    struct {
		MediaID string `json:"media_id"` // 文件 ID
	} `json:"file"` // 文件消息
}

// Location 地理位置消息
type Location struct {
	BaseMessage
	MsgType  string `json:"msgtype"` // 消息类型，此时固定为：location
	Location struct {
		Latitude  float32 `json:"latitude"`  // 纬度
		Longitude float32 `json:"longitude"` // 经度
		Name      string  `json:"name"`      // 位置名
		Address   string  `json:"address"`   // 地址详情说明
	} `json:"location"` // 地理位置消息
}

// Link 链接消息
type Link struct {
	BaseMessage
	MsgType string `json:"msgtype"` // 消息类型，此时固定为：link
	Link    struct {
		Title  string `json:"title"`   // 标题
		Desc   string `json:"desc"`    // 描述
		URL    string `json:"url"`     // 点击后跳转的链接
		PicURL string `json:"pic_url"` // 缩略图链接
	} `json:"link"` // 链接消息
}

// BusinessCard 名片消息
type BusinessCard struct {
	BaseMessage
	MsgType      string `json:"msgtype"` // 消息类型，此时固定为：business_card
	BusinessCard struct {
		UserID string `json:"userid"` // 名片 userid
	} `json:"business_card"` // 名片消息
}

// MiniProgram 小程序消息
type MiniProgram struct {
	BaseMessage
	MsgType     string `json:"msgtype"` // 消息类型，此时固定为：miniprogram
	MiniProgram struct {
		AppID        string `json:"appid"`          // 小程序 appid，必须是关联到企业的小程序应用
		Title        string `json:"title"`          // 小程序消息标题，最多 64 个字节，超过会自动截断
		ThumbMediaID string `json:"thumb_media_id"` // 小程序消息封面的 mediaid，封面图建议尺寸为 520*416
		PagePath     string `json:"pagepath"`       // 点击消息卡片后进入的小程序页面路径
	} `json:"miniprogram"` // 小程序消息
}

// EventMessage 事件消息
type EventMessage struct {
	BaseMessage
	MsgType string `json:"msgtype"` // 消息类型，此时固定为：event
	Event   struct {
		EventType string `json:"event_type"` // 事件类型
	} `json:"event"` // 事件消息
}

// EnterSessionEvent 用户进入会话事件
type EnterSessionEvent struct {
	BaseMessage
	MsgType string `json:"msgtype"` // 消息类型，此时固定为：event
	Event   struct {
		EventType      string `json:"event_type"`      // 事件类型。此处固定为：enter_session
		OpenKFID       string `json:"open_kfid"`       // 客服账号 ID
		ExternalUserID string `json:"external_userid"` // 客户 UserID
		Scene          string `json:"scene"`           // 进入会话的场景值，获取客服帐号链接开发者自定义的场景值
		SceneParam     string `json:"scene_param"`     // 进入会话的自定义参数，获取客服帐号链接返回的 url，开发者按规范拼接的 scene_param 参数
		WelcomeCode    string `json:"welcome_code"`    // 如果满足发送欢迎语条件（条件为：1. 企业没有在管理端配置了原生欢迎语；2. 用户在过去 48 小时里未收过欢迎语，且未向该用户发过消息），会返回该字段。可用该 welcome_code 调用发送事件响应消息接口给客户发送欢迎语。
	} `json:"event"` // 事件消息
}

// MsgSendFailEvent 消息发送失败事件
type MsgSendFailEvent struct {
	BaseMessage
	MsgType string `json:"msgtype"` // 消息类型，此时固定为：event
	Event   struct {
		EventType      string `json:"event_type"`      // 事件类型。此处固定为：msg_send_fail
		OpenKFID       string `json:"open_kfid"`       // 客服账号 ID
		ExternalUserID string `json:"external_userid"` // 客户 UserID
		FailMsgID      string `json:"fail_msgid"`      // 发送失败的消息 msgid
		FailType       uint32 `json:"fail_type"`       // 失败类型。0-未知原因 1-客服账号已删除 2-应用已关闭 4-会话已过期，超过 48 小时 5-会话已关闭 6-超过 5 条限制 7-未绑定视频号 8-主体未验证 9-未绑定视频号且主体未验证 10-用户拒收
	} `json:"event"` // 事件消息
}

// ReceptionistStatusChangeEvent 客服人员接待状态变更事件
type ReceptionistStatusChangeEvent struct {
	BaseMessage
	MsgType string `json:"msgtype"` // 消息类型，此时固定为：event
	Event   struct {
		EventType          string `json:"event_type"`      // 事件类型。此处固定为：servicer_status_change
		ReceptionistUserID string `json:"servicer_userid"` // 客服人员 userid
		OpenKFID           string `json:"open_kfid"`       // 客服帐号 ID
		Status             uint32 `json:"status"`          // 状态类型。1-接待中 2-停止接待
	} `json:"event"`
}

// SessionStatusChangeEvent 会话状态变更事件
type SessionStatusChangeEvent struct {
	BaseMessage
	MsgType string `json:"msgtype"` // 消息类型，此时固定为：event
	Event   struct {
		EventType             string `json:"event_type"`          // 事件类型。此处固定为：session_status_change
		OpenKFID              string `json:"open_kfid"`           // 客服账号 ID
		ExternalUserID        string `json:"external_userid"`     // 客户 UserID
		ChangeType            uint32 `json:"change_type"`         // 变更类型。1-从接待池接入会话 2-转接会话 3-结束会话
		OldReceptionistUserID string `json:"old_servicer_userid"` // 老的客服人员 userid。仅 change_type 为 2 和 3 有值
		NewReceptionistUserID string `json:"new_servicer_userid"` // 新的客服人员 userid。仅 change_type 为 1 和 2 有值
		MsgCode               string `json:"msg_code"`            // 用于发送事件响应消息的 code，仅 change_type 为 1 和 3 时，会返回该字段。可用该 msg_code 调用发送事件响应消息接口给客户发送回复语或结束语。
	} `json:"event"` // 事件消息
}
