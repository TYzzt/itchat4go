package webservice

import (
	"fmt"
	e "itchat4go/enum"
	s "itchat4go/service"
	m "itchat4go/model"
	"time"
)

var (
	err        error
	loginMap   m.LoginMap
	contactMap map[string]m.User
	groupMap   map[string][]m.User /* 关键字为key的，群组数组 */
)

/**
	微信逻辑处理
	30s过期
 */
func go_listener(uuid string)   {

	timeout:=make(chan bool,1)
	processout:=make(chan bool,1)
	go func() {

		for {
			select {
			case <-timeout:  //过期结束
				return
			default:
			}

			fmt.Println("正在验证登陆... ...")
			status, msg := s.CheckLogin(uuid)

			if status == 200 {
				fmt.Println("登陆成功,处理登陆信息...")
				loginMap, err = s.ProcessLoginInfo(msg)
				if err != nil {
					panicErr(err)
				}

				fmt.Println("登陆信息处理完毕,正在初始化微信...")
				err = s.InitWX(&loginMap)
				if err != nil {
					panicErr(err)
				}

				fmt.Println("初始化完毕,通知微信服务器登陆状态变更...")
				err = s.NotifyStatus(&loginMap)
				if err != nil {
					panicErr(err)
				}

				fmt.Println("通知完毕,本次登陆信息：")
				fmt.Println(e.SKey + "\t\t" + loginMap.BaseRequest.SKey)
				fmt.Println(e.PassTicket + "\t\t" + loginMap.PassTicket)
				break
			} else if status == 201 {
				fmt.Println("请在手机上确认")
			} else if status == 408 {
				fmt.Println("请扫描二维码")
			} else {
				fmt.Println("aaaaaaa"+msg)
			}
		}
		fmt.Println("开始获取联系人信息...")
		contactMap, err = s.GetAllContact(&loginMap)
		if err != nil {
			panicErr(err)
		}

		fmt.Println(contactMap)

		for _,contact := range contactMap {
			if contact.NickName=="赵涛" {
				wxSendMsg := m.WxSendMsg{}
				wxSendMsg.Type = 1
				wxSendMsg.Content = "测试..."
				wxSendMsg.FromUserName = loginMap.SelfUserName
				wxSendMsg.ToUserName = contact.UserName
				wxSendMsg.LocalID = fmt.Sprintf("%d", time.Now().Unix())
				wxSendMsg.ClientMsgId = wxSendMsg.LocalID

				//加点延时，避免消息次序混乱，同时避免微信侦察到机器人
				time.Sleep(time.Second)

				go s.SendMsg(&loginMap, wxSendMsg)
				break
			}
		}
		processout <- true //程序结束
	}()

	//设置过期时间
	select {
	case <-processout:
		fmt.Println("正常退出")
	case <-time.After(30 * time.Second):
		timeout<-true
		fmt.Println("超时退出")
	}
}