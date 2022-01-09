package main

import (
	"encoding/json"
	"fmt"
	ossClient "github.com/alibabacloud-go/tea-oss-sdk/client"
	ossUtils "github.com/alibabacloud-go/tea-oss-utils/service"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"time"
)

type config struct {
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	Bucket          string `json:"bucket"`
	Area            string `json:"area"`
	// Path 不要以/ 开头否则会上传失败
	Path            string `json:"path"`
	// CustomUrl 最后要带上 /
	CustomUrl       string `json:"customUrl"`
}

func init(){
	// 初始化logger
	logFile, err := os.OpenFile("uploader.log",os.O_WRONLY|os.O_CREATE|os.O_APPEND,0755)
	if err != nil {
		fmt.Println("无法创建日志文件")
	}
	// 开发时同时输出到文件和控制台
	//log.SetOutput(io.MultiWriter(logFile,os.Stderr))
	// 正式编译到使用的时候只输出到日志文件
	log.SetOutput(logFile)
	log.SetReportCaller(true)
}

func main(){
	// 程序的启动参数，第一个是程序本身，后面的才是参数
	inputs := getInputFilesPath()
	config := readConfig()
	cfg := new(ossClient.Config).
		SetAccessKeyId(config.AccessKeyId).
		SetAccessKeySecret(config.AccessKeySecret).
		SetRegionId(config.Area).
		SetType("access_key")
	cli, err := ossClient.NewClient(cfg)
	if err != nil {
		log.Fatalf("配置文件有误: %v",err)
	}
	putObjectRequestHeader := new(ossClient.PutObjectRequestHeader).SetStorageClass("Standard")
	runtimeOptions := new(ossUtils.RuntimeOptions).SetAutoretry(false).SetMaxIdleConns(3)
	failedFiles := make([]string,0)
	successFilesUrl := make([]string,0)
	// 可能有多个输入文件
	for _, inputPath := range inputs {
		var reader io.Reader
		reader, err = os.Open(inputPath)
		if err != nil {
			parsedUrl, err := url.Parse(inputPath)
			if err != nil {
				log.Errorf("读取文件 %s 失败: %v", inputPath, err)
				failedFiles = append(failedFiles, inputPath)
				continue
			}
			res, err := http.Get(parsedUrl.String())
			if err != nil {
				log.Errorf("下载来自网络的图片时出错：%v",err)
				failedFiles = append(failedFiles, inputPath)
				continue
			}
			reader = res.Body
		}
		fileAddress := getFileName(inputPath,config.Path)
		putObjectRequest := new(ossClient.PutObjectRequest).
			SetBucketName(config.Bucket).
			SetObjectName(fileAddress).
			SetHeader(putObjectRequestHeader).
			SetBody(reader)
		_, err = cli.PutObject(putObjectRequest, runtimeOptions)
		if err != nil {
			log.Errorf("上传文件 %s 失败：%v", inputPath,err)
			failedFiles = append(failedFiles, inputPath)
		}
		if len(config.CustomUrl) >0 {
			successFilesUrl = append(successFilesUrl,config.CustomUrl+fileAddress)
		}else {
			successFilesUrl = append(successFilesUrl,fmt.Sprintf("https://%s.oss-%s.aliyuncs.com/%s",config.Bucket,config.Area,fileAddress))
		}
	}
	// 所有都成功上传
	if len(failedFiles) == 0{
		fmt.Println("Upload Success:")
		for _,e := range successFilesUrl {
			fmt.Println(e)
		}
	}else {
		log.Errorf("文件上传失败：\n")
		log.Error(failedFiles)
	}
}
func readConfig() *config {
	// 在可执行文件的同级目录下寻找配置文件
	dir := filepath.Dir(os.Args[0])
	abs := filepath.Join(dir,"./config.json")
	cfgFile, err := os.Open(abs)
	if err != nil {
		log.Fatalf("找不到配置文件: %v",err)
	}
	var config = &config{}
	decoder := json.NewDecoder(cfgFile)
	err = decoder.Decode(config)
	if err != nil {
		log.Fatalf("解析配置文件出错: %v",err)
	}
	return config
}
func getInputFilesPath() []string{
	return os.Args[1:]
}
// dir 是存储到云端的文件夹，不能以/ 开头
func getFileName(inputPath string,dir string) string {
	fileSubFix := path.Ext(path.Base(inputPath))
	filename := time.Now().Format("20060102150405")
	return path.Join(dir, filename+fileSubFix)
}