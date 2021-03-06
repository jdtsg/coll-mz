package controller

import (
	"fmt"
	"time"
)

//log struct
//You need to set the related configuration.
//IP address if the output is set, otherwise you can leave empty.
type Log struct {
	logDirSrc         string
	isSendErrorToFmt  bool
	isSendMsgToFmt    bool
	isSendErrorToFile bool
	isSendMsgToFile   bool
	isAppendTime      bool
	isAppendIP        bool
	ip                string
	oneFileName string
	isForward bool
	lastSrc string
}

//Initialize the configuration
//This function must be executed before using Log.
func (this *Log) init(logDirSrc string, isSendErrorToFmt bool, isSendMsgToFmt bool, isSendErrorToFile bool, isSendMsgToFile bool, isAppendTime bool, isAppendIP bool) {
	this.logDirSrc = logDirSrc
	if IsFolder(logDirSrc) == false{
		err = CreateDir(logDirSrc)
		if err != nil{
			this.SendFmtPrintln("The log module was initialized and the storage directory could not be created.Error : " + err.Error())
		}
	}
	this.isSendErrorToFmt = isSendErrorToFmt
	this.isSendMsgToFmt = isSendMsgToFmt
	this.isSendErrorToFile = isSendErrorToFile
	this.isSendMsgToFile = isSendMsgToFile
	this.isAppendTime = isAppendTime
	this.isAppendIP = isAppendIP
	this.oneFileName = ""
	this.isForward = false
}

//New log
//The log is automatically sent according to the settings.
func (this *Log) NewLog(msg string, err error) {
	if this.isAppendIP == true {
		this.UpdateIP()
	}
	if msg != "" {
		if this.isSendMsgToFmt == true {
			this.SendFmtPrintln(msg)
		}
		if this.isSendMsgToFile == true {
			this.SendFile(msg)
		}
	}
	if err != nil {
		if this.isSendErrorToFmt == true {
			this.SendFmtPrintln("Error : " + err.Error())
		}
		if this.isSendErrorToFile == true {
			this.SendFile("Error : " + err.Error())
		}
	}
}

//Update the IP address
func (this *Log) UpdateIP() {
	this.ip = IPAddrsGetExternal()
	if this.ip == "" {
		this.ip = IPAddrsGetInternal()
	}
}

//Send logs to the console
func (this *Log) SendFmtPrintln(msg string) {
	if this.isAppendTime == true {
		msg = this.GetNowTime() + " " + this.ip + " " + msg
	}
	fmt.Println(msg)
}

//Send log to file
func (this *Log) SendFile(content string) {
	if this.logDirSrc == "" {
		this.SendFmtPrintln("The log directory path is not provided.")
		return
	}
	if this.isAppendTime == true {
		content = this.GetNowTime() + " " + this.ip + " " + content + "\n"
	}
	var src string
	if this.oneFileName != ""{
		err = CreateDir(this.logDirSrc)
		if err != nil{
			this.SendFmtPrintln(err.Error())
			return
		}
		src = this.logDirSrc + GetPathSep() + this.oneFileName + ".log"
	}else{
		src, err = GetTimeDirSrc(this.logDirSrc, ".log")
		if err != nil {
			this.SendFmtPrintln("Unable to create log save directory path.")
			return
		}
	}
	contentByte := []byte(content)
	err = WriteFileAppend(src, contentByte, this.isForward)
	if err != nil{
		this.SendFmtPrintln(err.Error())
	}
	this.lastSrc = src
}

//Set oneFileName
func (this *Log) SetOneFileName(name string) {
	this.oneFileName = name
}

//Set isForward
func (this *Log) SetIsForward(b bool) {
	this.isForward = b
}

//Gets the current time
func (this *Log) GetNowTime() string {
	return time.Now().String()
}
