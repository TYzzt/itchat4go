package webservice

import (
	"log"
	"net/http"
	"html/template"
	e "itchat4go/enum"
	_ "itchat4go/model"
	s "itchat4go/service"
	"os"
	"strings"
	"path/filepath"
	"fmt"
)

/**
页面参数
 */
type WebParm struct{
	UrlSrc string  //二维码url
	Uuid string
	Qf_msg string
	Qf_state bool //群发还是测试
}


/**
开启web监听
 */
func BeginListene()  {
	routes() //路由
	log.Println("listener : Started : Listening on :4000")
	http.ListenAndServe(":4000", nil)
}


func routes() {
	http.HandleFunc("/wxSendProcess", WxQrCode)
	http.HandleFunc("/msgSetView", msgSetView)
}

//二维码页面
func WxQrCode(rw http.ResponseWriter, r *http.Request) {

	/* 从微信服务器获取UUID */
	uuid, err := s.GetUUIDFromWX()
	if err != nil {
		panicErr(err)
	}
	var parm = WebParm{UrlSrc:e.QRCODE_URL+uuid,Uuid:uuid}
	renderHTML(rw,"index.html",parm) //返回页面

	//微信api监听
	go go_listener(uuid)
}

/**
设置群发消息页面
 */
func msgSetView(rw http.ResponseWriter, req *http.Request)  {

	//设置消息
	post_msg := req.PostFormValue("msg")
	post_Qf_state := req.PostFormValue("Qf_state")
	if len(post_msg)>0 {
		Qf_msg = post_msg
	}
	if len(post_Qf_state)> 0 {
		Qf_state = post_Qf_state=="true"
	}
	var parm = WebParm{Qf_msg:Qf_msg,Qf_state:Qf_state}
	renderHTML(rw,"msgSetView.html",parm) //返回页面
}


func renderHTML(w http.ResponseWriter, file string, data interface{}) {
	var path = getCurrentDirectory()
	fmt.Println(path)
	// 获取页面内容
	t, err := template.ParseFiles(path+"/view/"+file)
	fmt.Println(t)
	checkErr(err)
	// 将页面渲染后反馈给客户端
	t.Execute(w, data)
}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

/**
获取运行路径
 */
func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panicErr(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}
