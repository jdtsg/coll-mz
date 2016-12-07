package controller

import (
	"strconv"
)

//Collect local data
func (this *Coll) CollLocal() {
	//Gets the object
	thisChildren := &this.collList.local
	var collOperate CollOperate
	if this.CollStart(thisChildren,&collOperate) == false{
		return
	}
	defer this.CollEnd(thisChildren,&collOperate)
	//该采集模块比较特殊，为方便理解，所以将在代码内写入中文注释
	//该模版不会使用到colloperate.auto这些自动化采集模块，而是在该文件内单独构建通用处理模块
	//逻辑上将按照流程，采集各目录下的文件，每种文件都有各自的特点
	//在采集前，建议先运行一次，这样可以自动建立相关需要的目录
	collOperate.NewLog("在采集前，建议先运行一次，这样可以自动建立相关需要的目录。",nil)
	collLocalDir := this.dataSrc + GetPathSep() + "coll-local"
	if IsFolder(collLocalDir) == false{
		err = CreateDir(collLocalDir)
		if err != nil{
			collOperate.NewLog("初始化失败。",err)
			return
		}
	}
	collOperate.NewLog("目录构建完成后，您可以在" + collLocalDir + "目录下查看到所有对应文件夹，将手动收集到的内容放入对应文件夹即可。",nil)
	collOperate.NewLog("文件夹结构主要分别为：",nil)
	collOperate.NewLog("  txt : 文本文件，eg : txt/name.txt",nil)
	collOperate.NewLog("    采集所有txt文件，不能嵌套多级目录。",nil)
	collOperate.NewLog("  save-imgs-html : 保存的图片类网页，eg : save-imgs-html/name/xxx.jpg ;save-imgs-html/name.html",nil)
	collOperate.NewLog("    采集子目录、子目录下jpg|gif|jpeg|png文件，其他文件自动删除。",nil)
	collOperate.NewLog("  download-movie : 网上下载的视频，仅支持mp4文件，eg : name/name.mp4",nil)
	collOperate.NewLog("    采集子目录、子目录下mp4视频、子目录下cover.jpg索引图片。",nil)
	collOperate.NewLog("  manhua : 网上保存下来的漫画合集，eg : manhua/name/xxx.jpg",nil)
	collOperate.NewLog("    采集子目录、子目录下jpg文件、子目录下cover.jpg封面图片。",nil)
	collOperate.NewLog("开始整理.......",nil)
	if thisChildren.status == false{
		return
	}
	//采集txt文件夹下的数据
	this.CollLocalTxt(thisChildren,&collOperate,collLocalDir)
	//采集download-movie文件夹下的数据
	this.CollLocalDownloadMovie(thisChildren,&collOperate,collLocalDir)
	//采集manhua文件夹下的数据
	this.CollLocalManhua(thisChildren,&collOperate,collLocalDir)
	//采集的保存网页数据
	this.CollLocalSaveImgsHtml(thisChildren,&collOperate,collLocalDir)
}

//local专用通用启动模块
//在内部采集启动前，务必运行该模块
// collOperate *CollOperate
// collLocalDir string local总采集器路径
// name string 子采集器名称
// fileFilter string 可选，文件列表过滤文件，不支持二级子目录下文件过滤，但会返回相关目录名称，eg : jpg|gif|txt
// return string - 该采集项目存储目录路径
// return []string - 该目录下文件列表
func (this *Coll) CollLocalStart(collOperate *CollOperate,collLocalDir string,name string,fileFilter string) (string,[]string) {
	//构建存储路径
	dir := collLocalDir + GetPathSep() + name
	collOperate.NewLog(" ## 开始整理" + name + "文件夹数据 ## ",nil)
	var fileList []string
	//如果不存在目录，则创建
	if IsFolder(dir) == false{
		err = CreateDir(dir)
		if err != nil{
			collOperate.NewLog("初始化文件夹失败。",err)
			return "",fileList
		}
		return "",fileList
	}
	//查询下面是否存在文件
	fileNum,err := GetFileListCount(dir)
	if err != nil{
		collOperate.NewLog("无法获取文件数量。",err)
		return "",fileList
	}
	if fileNum < 1{
		collOperate.NewLog("目录下没有任何文件，请先添加文件后再尝试采集该部分内容。",nil)
		return "",fileList
	}
	//直接获取该目录下所有文件
	fileList,err = GetFileList(dir,fileFilter,true)
	if err != nil{
		collOperate.NewLog("无法获取目录下的文件列表。",nil)
		return "",fileList
	}
	//返回新目录
	return dir,fileList
}

//local专用通用采集模块
func (this *Coll) CollLocalParentFiles(thisChildren *CollChildren,collOperate *CollOperate,parentTitle string,parentSrc string,fileSrcList []string) bool {
	//构建parent数据
	//Create parent directory data
	parentSha1 := collOperate.matchString.GetSha1(parentTitle + parentSrc)
	if parentSha1 == ""{
		collOperate.NewLog(this.lang.Get("coll-error-sha1") + " parent : " + parentTitle + " , src : " + parentSrc,nil)
		return false
	}
	//Check parent sha1 if the data already exists
	if collOperate.CheckDataSha1(parentSha1) == true{
		collOperate.NewLog(this.lang.Get("coll-error-repeat-sha1") + parentSrc + " , sha1 : " + parentSha1,nil)
		return true
	}
	//Create parent database data
	parentID := collOperate.CreateNewData(0,parentSha1,"",parentSrc,parentTitle,"folder","0")
	if parentID > 0{
		collOperate.NewLog(this.lang.Get("coll-new-id") + strconv.FormatInt(parentID,10) + " , src : " + parentSrc,nil)
	}else{
		collOperate.NewLog(this.lang.Get("coll-error-move-file") + parentSrc,nil)
		return false
	}
	//根据文件列表，将其保存并建立数据
	var errNum = 0
	for _,value := range fileSrcList{
		//获取文件SHA1值
		fileSha1,err := GetFileSha1(value)
		if err != nil{
			collOperate.NewLog("无法获取文件SHA1值。",nil)
			errNum += 1
			continue
		}
		//查询SHA1是否存在于数据库
		if collOperate.CheckDataSha1(fileSha1) == true{
			collOperate.NewLog("该文件已经存在，跳过。",nil)
			errNum += 1
			continue
		}
		//获取文件基本信息
		cacheFileInfo := make(map[string]string)
		cacheFileInfo["cache-src"] = value
		fileNames,err := GetFileNames(value)
		if err != nil{
			collOperate.NewLog("无法获取文件名称和类型信息。",err)
			errNum += 1
			continue
		}
		if fileNames["name"] == "" || fileNames["type"] == "" || fileNames["onlyName"] == ""{
			collOperate.NewLog("无法获取文件名称和类型内容。",err)
			errNum += 1
			continue
		}
		cacheFileInfo["full-name"] = fileSha1 + "." + fileNames["type"]
		fileSize := GetFileSize(value)
		if fileSize < 1{
			collOperate.NewLog("无法获取文件大小。",nil)
			errNum += 1
			continue
		}
		//构建文件存储路径
		newFileSrc := collOperate.SaveCacheToFile("txt",cacheFileInfo)
		//创建文件
		newID := collOperate.CreateNewData(parentID,fileSha1,newFileSrc,value,fileNames["onlyName"],fileNames["type"],strconv.FormatInt(fileSize,10))
		if newID > 0{
			collOperate.NewLog(this.lang.Get("coll-new-id") + strconv.FormatInt(newID,10) + " , URL : " + value,nil)
		}else{
			collOperate.NewLog(this.lang.Get("coll-error-move-file") + value,nil)
		}
	}
	//错误过多，则返回
	if errNum > 10{
		collOperate.NewLog("出现太多次错误了。",nil)
		return false
	}
	//返回
	return true
}

//local专用通用清理模块
//删除所有文件
func (this *Coll) CollLocalEnd(thisChildren *CollChildren,collOperate *CollOperate,name string,parentSrc string) {
	//删除子采集文件夹目录，之后再创建
	err = DeleteFile(parentSrc)
	if err != nil{
		collOperate.NewLog("删除文件夹失败。",err)
		return
	}
	err = CreateDir(parentSrc)
	if err != nil{
		collOperate.NewLog("删除后创建文件夹失败。",err)
		return
	}
	collOperate.NewLog(name + "采集结束。",nil)
}

//local文本数据采集器
func (this *Coll) CollLocalTxt(thisChildren *CollChildren,collOperate *CollOperate,collLocalDir string){
	//初始化获取
	name := "txt"
	dir,fileList := this.CollLocalStart(collOperate,collLocalDir,name,name)
	if dir == ""{
		return
	}
	//获取dir名称
	parentNames,err := GetFileNames(dir)
	if err != nil{
		collOperate.NewLog("无法获取parent名称和类型数据。",err)
		return
	}
	parentName := parentNames["only-name"]
	//重新建立文件数据，剔除所有目录、非txt文件
	var newFileList []string
	for _,v := range fileList {
		if IsFolder(v) == true{
			continue
		}
		names,err := GetFileNames(v)
		if err != nil{
			collOperate.NewLog("",err)
			continue
		}
		if names["type"] == "txt" {
			newFileList = append(newFileList,v)
		}
	}
	//开始构建数据
	b := this.CollLocalParentFiles(thisChildren,collOperate,parentName,dir,newFileList)
	if b == false{
		return
	}
	//收尾工作
	this.CollLocalEnd(thisChildren,collOperate,name,dir)
}

//local下载视频数据采集器
func (this *Coll) CollLocalDownloadMovie(thisChildren *CollChildren,collOperate *CollOperate,collLocalDir string){
	//初始化获取
	//dir,fileList := this.CollLocalStart(collOperate,collLocalDir,"download-movie")
}

//local漫画数据采集器
func (this *Coll) CollLocalManhua(thisChildren *CollChildren,collOperate *CollOperate,collLocalDir string){
	//初始化获取
	//dir,fileList := this.CollLocalStart(collOperate,collLocalDir,"manhua")
}

//local保存网页数据采集器
func (this *Coll) CollLocalSaveImgsHtml(thisChildren *CollChildren,collOperate *CollOperate,collLocalDir string){
	//初始化获取
	//dir,fileList := this.CollLocalStart(collOperate,collLocalDir,"save-imgs-html")
}