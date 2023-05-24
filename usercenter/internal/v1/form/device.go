package form

import (
	"encoding/json"
	"github.com/pengcainiao/pengcainiao/usercenter/internal/v1/constant"
)

type DeviceInfoV2 struct {
	//DeviceName    string         `json:"device_name"`                                                          //设备名称
	ClientVersion string                  `json:"client_version" binding:"required"`                                                 //飞项应用的版本号
	ClientIP      string                  `json:"client_ip" `                                                                        //客户端IP
	DeviceID      string                  `json:"device_id,omitempty" db:"device_id"`                                                //设备ID
	OS            string                  `json:"os,omitempty"`                                                                      //Android、iOS、Windows、macOS、WeChat、web
	OSVersion     string                  `json:"os_version,omitempty"`                                                              //操作系统版本号
	Model         string                  `json:"model,omitempty"`                                                                   //设备名称，类似 HUAWEI||LIO-AN00
	Network       string                  `json:"network,omitempty"`                                                                 //wifi、wlan、
	Operator      string                  `json:"operator,omitempty"`                                                                //网络运营商，中国联通、中国移动
	SSID          string                  `json:"ssid,omitempty"`                                                                    //WiFi名称
	MacAddress    string                  `json:"mac_addr,omitempty"`                                                                //设备Mac地址
	Platform      constant.PlatformDefine `json:"platform" db:"platform" binding:"required,oneof=wechat pc mobile web pc_wechat h5 corp_wechat"` //所属平台 // 多端同步改动
	IsNewUser     bool                    `json:"-"`                                                                                 // 是否新用户
}

func (d DeviceInfoV2) String() string {
	b, _ := json.Marshal(d)
	return string(b)
}
