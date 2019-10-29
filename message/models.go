/*
 * @Author: kidd
 * @Date: 7/29/19 10:51 AM
 */

package message

import (
	"github.com/exwallet/go-common/database/mysql/data"
	"github.com/exwallet/go-common/goutil/gotime"
)

type MsgType string

const (
	MsgTypeSMS      MsgType = "sms"
	MsgTypeEmail    MsgType = "email"
	MsgTypeSMSVoice MsgType = "smsVoice"

	/*
		WAITING_FOR_SEND(0, "待发送"),
		SUCCESS(1, "成功"),
		FAIL(-1, "失败"),
	*/
	StatusWaitingSend int64 = 0
	StatusSucc        int64 = 1
	StatusFail        int64 = -1
)

// index key:
type Message struct {
	Id          int64  `json:"id" pk:"1"`
	MsgType     string `json:"msgType"`     // 消息类型
	Gateway     string `json:"gateway"`     // 发送网关
	Tos         string `json:"tos"`         // 业务类型
	UserId      int64  `json:"userId"`      //
	Username    string `json:"username"`    //
	CountryCode string `json:"countryCode"` //
	ReceiveAddr string `json:"receiveAddr"` //
	Title       string `json:"title"`       //
	Content     string `json:"content"`     //
	Status      int64  `json:"status"`      // 消息状态
	FailTimes   int64  `json:"failTimes"`   //
	AddTime     int64  `json:"addTime"`     //
	SendTime    int64  `json:"sendTime"`    //
	AddIP       string `json:"addIp"`       //
	IsAdmin     int64  `json:"isAdmin"`     // 0前台, 1后台
}

//
func SendSMS(dao *data.Dao, tos string, userId int64, username string, countryCode string, mobile string, content string, addIP string, isAdmin int64) error {
	m := &Message{
		Id:          0,
		MsgType:     string(MsgTypeSMS),
		Gateway:     "",
		Tos:         tos,
		UserId:      userId,
		Username:    username,
		CountryCode: countryCode,
		ReceiveAddr: mobile,
		Title:       "",
		Content:     content,
		Status:      StatusWaitingSend,
		FailTimes:   0,
		AddTime:     gotime.UnixNowMillSec(),
		SendTime:    0,
		AddIP:       addIP,
		IsAdmin:     isAdmin,
	}
	_, err := dao.Insert(m)
	return err
}

//
func SendEmail(dao *data.Dao, tos string, userId int64, username string, to string, title string, content string, addIP string, isAdmin int64) error {
	m := &Message{
		Id:          0,
		MsgType:     string(MsgTypeEmail),
		Gateway:     "",
		Tos:         tos,
		UserId:      userId,
		Username:    username,
		CountryCode: "",
		ReceiveAddr: to,
		Title:       title,
		Content:     content,
		Status:      StatusWaitingSend,
		FailTimes:   0,
		AddTime:     gotime.UnixNowMillSec(),
		SendTime:    0,
		AddIP:       addIP,
		IsAdmin:     isAdmin,
	}
	_, err := dao.Insert(m)
	return err
}

//
func SendSMSVoice(dao *data.Dao, tos string, userId int64, username string, countryCode string, mobile string, content string, addIP string, isAdmin int64) error {
	m := &Message{
		Id:          0,
		MsgType:     string(MsgTypeSMSVoice),
		Gateway:     "",
		Tos:         tos,
		UserId:      userId,
		Username:    username,
		CountryCode: countryCode,
		ReceiveAddr: mobile,
		Title:       "",
		Content:     content,
		Status:      StatusWaitingSend,
		FailTimes:   0,
		AddTime:     gotime.UnixNowMillSec(),
		SendTime:    0,
		AddIP:       addIP,
		IsAdmin:     isAdmin,
	}
	_, err := dao.Insert(m)
	return err
}
