package webservice

import (
	"log"
	"net/http"
	"html/template"
	e "itchat4go/enum"
	_ "itchat4go/model"
	s "itchat4go/service"
)

/**
页面参数
 */
type WebParm struct{
	URLSrc string  //二维码url
}

func BeginListene()  {
	routes() //监听
	log.Println("listener : Started : Listening on :4000")
	http.ListenAndServe(":4000", nil)
}


func routes() {
	http.HandleFunc("/", WxQrCode)
}

//二维码页面
func WxQrCode(rw http.ResponseWriter, r *http.Request) {

	/* 从微信服务器获取UUID */
	uuid, err := s.GetUUIDFromWX()
	if err != nil {
		panicErr(err)
	}

	var parm = WebParm{URLSrc:e.QRCODE_URL+uuid}

	go go_listener(uuid)  //微信api监听

	renderHTML(rw,"index.html",parm)
}
func renderHTML(w http.ResponseWriter, file string, data interface{}) {
	// 获取页面内容
	t, err := template.New(file).ParseFiles("../view/" + file)
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