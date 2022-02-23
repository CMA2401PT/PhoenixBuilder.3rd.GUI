# PhoenixBuilder.3rd.GUI

### 说明

~~PhoenixBuilderHeadless 是原FB项目的无头版本，尽量减少对原项目的修改：  
https://github.com/LNSSPsd/PhoenixBuilder~~

没办法，为了能让项目在安卓上编译，不得不修改fastbuilder->现在的fb文件夹  
PhoenixBuilder GPLv3 协议，项目的核心:  
https://github.com/LNSSPsd/PhoenixBuilder

还使用了https://github.com/crvdgc/Consolas-with-Yahei的字体文件解决中文问题

最后推荐一个库 Fyne: https://github.com/fyne-io/fyne   
这个库还不够成熟，但是基础相当优秀


### 编译
####对于 mac/win/linux 
go build main 即可
####对于 android：
1. 准备环境，ndk，adb，另外：
```
go get fyne.io/fyne/v2/cmd/fyne
go install fyne.io/fyne/v2/cmd/fyne
```
2. 根据你的系统，可能需要一些额外操作，总之，当在命令行中输入 fyne 有输出即为成功
3. 编译
```
fyne package -os android -appID my.domain.appname （官方说明，windows上似乎无法正常工作）
fyne package -os android/arm64 -appID my.domain.appname （个人建议）
```
4. 安装测试，windows上似乎无法正常工作
```
fyne install -os android
```
####对于ios：
没有试过，官网说：
```
fyne release -os ios -certificate "Apple Distribution" -profile "My App Distribution" -appID "com.example.myapp"
```