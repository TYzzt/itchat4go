package webservice

import (
	"fmt"
	_ "itchat4go/enum"
	s "itchat4go/service"
	m "itchat4go/model"
	"time"
	"log"
)

var (
	err        error
	loginMap   m.LoginMap
	contactMap map[string]m.User
	groupMap   map[string][]m.User /* 关键字为key的，群组数组 */
	Qf_msg ="人是要有精神的~" //群发默认消息
	Qf_state = false // false测试 true群发
	Qf_bgr = "金钊百货批发13932872008" //群发报告人
)


/**
	微信逻辑处理
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
			/*	fmt.Println(e.SKey + "\t\t" + loginMap.BaseRequest.SKey)
				fmt.Println(e.PassTicket + "\t\t" + loginMap.PassTicket)*/
				log.Print(loginMap.SelfNickName+"登陆")
				break
			} else if status == 201 {
				fmt.Println("请在手机上确认")
			} else if status == 408 {
				fmt.Println("请扫描二维码")
			} else {
				fmt.Println("aaaaaaa"+msg)
				return
			}
		}
		fmt.Println("开始获取联系人信息...")
		contactMap, err = s.GetAllContact(&loginMap)
		if err != nil {
			panicErr(err)
		}

		/*fmt.Println(contactMap)*/

		if !Qf_state { //测试
		fmt.Println(len(contactMap))
			for _,contact := range contactMap {
				if contact.NickName=="赵涛" {
					wxSendMsg := m.WxSendMsg{}
					wxSendMsg.Type = 1
					wxSendMsg.Content = Qf_msg
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
		}else{  //群发
			var zInt = 0
			for _,contact := range contactMap {
				zInt++
				if  zInt  % 100== 0 {
					time.Sleep(time.Second*4)
				}

				wxSendMsg := m.WxSendMsg{}
				wxSendMsg.Type = 1
				wxSendMsg.Content = Qf_msg
				wxSendMsg.FromUserName = loginMap.SelfUserName
				wxSendMsg.ToUserName = contact.UserName
				wxSendMsg.LocalID = fmt.Sprintf("%d", time.Now().Unix())
				wxSendMsg.ClientMsgId = loginMap.SelfNickName+"你好/n"+wxSendMsg.LocalID

				//加点延时，避免消息次序混乱，同时避免微信侦察到机器人
				time.Sleep(time.Second)
				go s.SendMsg(&loginMap, wxSendMsg)
			}
		}
		processout <- true //程序结束
	}()

	//设置过期时间
	select {
	case <-processout:
		fmt.Println("正常退出")
	case <-time.After(3000 * time.Second):
		timeout<-true
		fmt.Println("超时退出")
	}
}