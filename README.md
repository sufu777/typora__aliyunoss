## 一个简单的Typora上传图片到阿里云oss的工具

使用方法：

- 将exe文件下载到本地。
- 在exe同级目录下创建配置文件config.json，文件内容如下文所述。
- 在typora的`文件`->`偏好设置`->`图像`->`上传服务设定`处选择上传服务为`Custom Command`
然后在命令处指定exe文件文件地址，如果路径中有空格，需要加上 "" ，例如`"C:\Program Files\uploader\tiu.go"`。
- 点击验证图片上传选项即可测试配置是否成功

配置文件详解：

> accessKeyId和accessKeySecret是在阿里云的RAM访问控制中为用户创建的AccessKey，该AccessKey
需要有Oss的文件上传权限。
> 
> bucket即是要上传图片的bucket的名字。
> 
>area是该bucket的地域简写，例如：cn-chengdu表示中国成都，可以在oss概览界面查看。
>path是存储的服务器的文件夹，不要以 / 开头否则会上传失败，例如：typora_imgs
> 
>customUrl是配置的自定义域名的地址，若您使用域名绑定功能，则需要将用户域名（CNAME 域名解析）指向到您的Bucket 域名上。
然后在此指定您的用户域名，没有使用则设置为空字符串“”。例如：`https://oos.imgs.icu/` ，注意最后要带上 `/ `

最后附上配置文件示例：
```json
{
  "accessKeyId": "your aliyun accessKeyId",
  "accessKeySecret": "your aliyun accessKeySecret",
  "bucket": "Bucket111",
  "area": "cn-chengdu",
  "path": "typora_imgs",
  "customUrl": ""
}
```
代码很简单不过多解释，如果要自己编译，无窗口运行程序需指定编译参数`-ldflags "-H windowsgui"`，最后的-o 指定输出可执行文件的名字。
```shell
go build -ldflags "-H windowsgui" -o tiu.exe
```