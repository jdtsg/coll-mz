package main

import(
	"github.com/fotomxq/ftmp-libs"
	"github.com/fotomxq/coll-mz/libs"
)

//日志处理
var log ftmplibs.Log
//错误
//某些声明下，直接复用而不是重新声明error变量
var err error

//启动脚本
func main(){
	//开始提示
	log.AddLog("* _ * * _ * 脚本开始运行 * _ * * _ *")
	log.AddLog("初始化参数中...")
	//设定错误前缀
	log.SetErrorPrefix("发生一个错误 : ")
	//获取配置数据
	config := new(ftmplibs.Config)
	err = config.LoadFile("content/config/config.json")
	if err != nil{
		log.AddErrorLog(err)
		return
	}
	//激活服务器
	collmzLibs.Router()
}