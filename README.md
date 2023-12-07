# README

## 配置
修改config.yaml文件
```yaml
server:
  MaxFileSize: 3 # 限制上传文件大小
  UploadPath: "/wwwroot/picgo/images" # 上传文件保存路径
  HttpUlr: "http://0.0.0.0:10000/images/" # 图片调用URL
  Token: "e10adc3949ba59abbe56e057f20f883e" # 上传Token验证
```
服务端默认监听16001端口,可自行修改
```bash
err := r.Run(":16001")
if err != nil {
	return
}
```

## 使用
将imageUpdateServer注册为服务
### 创建一个 Systemd 服务单元文件
  使用文本编辑器（比如 nano 或 vim）创建一个 .service 文件，比如 imageUpdateServer.service：
```bash
sudo nano /etc/systemd/system/imageUpdateServer.service
```
###  编辑服务单元文件
在新创建的服务文件中，指定以下内容（以你的实际设置为准）：
```bash
[Unit]
Description=imageUpdateServer
After=network.target

[Service]
Type=simple
WorkingDirectory=/root/picgo/imageUpdateServer
ExecStart=/root/picgo

[Install]
WantedBy=multi-user.target
```
- Description：描述你的服务的信息。
- After：指定在哪些服务之后启动，这里使用了 network.target。
- Type：指定服务的类型。simple 适用于一般的启动脚本。
- WorkingDirectory：指定工作目录的路径。
- ExecStart：指定要运行的可执行文件的路径。
- WantedBy：指定服务安装的目标。multi-user.target 用于多用户模式下。
### 重载 systemd 并启用你的服务
保存并关闭文件后，重新加载 systemd 并启用你的服务：
```bash
sudo systemctl daemon-reload
sudo systemctl enable imageUpdateServer.service
```
这将使你的服务在系统启动时自动启动。
### 启动或停止服务
你可以使用 systemctl 命令启动、停止、重启或查看服务状态：
```bash
sudo systemctl start imageUpdateServer.service
sudo systemctl stop imageUpdateServer.service
sudo systemctl restart imageUpdateServer.service
sudo systemctl status imageUpdateServer.service
```
通过这些步骤，你的可执行文件将被配置为一个系统服务，并在系统启动时自动启动，同时会在指定的工作目录中运行。