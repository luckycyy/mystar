package music_controller

import (
	"log"
	"net/http"
)

func PlayMusicToA() {
// 计时开始前往A点修复
	log.Println("playMusicToA")
	resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.46&group=action_st&st=%E8%AE%A1%E6%97%B6%E5%BC%80%E5%A7%8B%E5%89%8D%E5%BE%80A%E7%82%B9%E4%BF%AE%E5%A4%8D&user_action=true")
	if err != nil {
		print(err)
	}
	resp.Body.Close()
}

func PlayMusicToB() {
 //计时开始前往A点修复
	log.Println("playMusicToB")
	resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.46&group=action_st&st=%E8%AE%A1%E6%97%B6%E5%BC%80%E5%A7%8B%E5%89%8D%E5%BE%80B%E7%82%B9%E4%BF%AE%E5%A4%8D&user_action=true")
	if err != nil {
		print(err)
	}
	resp.Body.Close()
}

func PlayFailedToA() {
//修复失败前往A点修复
	log.Println("playFailedToA")
	resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.46&group=action_st&st=%E4%BF%AE%E5%A4%8D%E5%A4%B1%E8%B4%A5%E5%89%8D%E5%BE%80A%E7%82%B9%E4%BF%AE%E5%A4%8D&user_action=true")
	if err != nil {
		print(err)
	}
	resp.Body.Close()

}

func PlayFailedToB() {
//修复失败前往B点修复
	log.Println("playFailedToB")
	resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.46&group=action_st&st=%E4%BF%AE%E5%A4%8D%E5%A4%B1%E8%B4%A5%E5%89%8D%E5%BE%80B%E7%82%B9%E4%BF%AE%E5%A4%8D&user_action=true")
	if err != nil {
		print(err)
	}
	resp.Body.Close()
}
func PlayStep1ok() {
//折返跑一阶段成功
	log.Println("playStep1ok")
	resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.46&group=action_st&st=%E6%8A%98%E8%BF%94%E8%B7%91%E4%B8%80%E9%98%B6%E6%AE%B5%E6%88%90%E5%8A%9F&user_action=true")
	if err != nil {
		print(err)
	}
	resp.Body.Close()
}
func PlayStep2ok() {
//折返跑二阶段成功
	log.Println("playStep2ok")
	resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.46&group=action_st&st=%E6%8A%98%E8%BF%94%E8%B7%91%E4%BA%8C%E9%98%B6%E6%AE%B5%E6%88%90%E5%8A%9F&user_action=true")
	if err != nil {
		print(err)
	}
	resp.Body.Close()
}
func PlayFinish() {
//折返跑完成
	log.Println("playFinish")
	resp, err := http.Get("http://192.168.1.21:1235/jdq_status/report_st?ip=192.168.1.46&group=action_st&st=%E6%8A%98%E8%BF%94%E8%B7%91%E5%AE%8C%E6%88%90&user_action=true")
	if err != nil {
		print(err)
	}
	resp.Body.Close()
}



//func PlayMusicToA() { //播放声音  计时开始前往A点修复
//	log.Println("test playMusicToA")
//}
//
//func PlayMusicToB() {
//	log.Println("test playMusicToB")
//}
//
//func PlayFailedToA() {
//	log.Println("test playFailedToA")
//}
//
//func PlayFailedToB() {
//	log.Println("test playFailedToB")
//}
//func PlayStep1ok() {
//	log.Println("test playStep1ok")
//}
//func PlayStep2ok() {
//	log.Println("test playStep2ok")
//}
//func PlayFinish() {
//	log.Println("test playFinish")
//}
