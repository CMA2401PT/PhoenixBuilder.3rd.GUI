# PhoenixBuilder.3rd.GUI

## 说明
PhoenixBuilder.3rd.GUI 是第三方开发者开发的套壳 FastBuilder  
提供了一个带有界面的，完全跨平台的图形化Fastbuilder  
其核心来自 Fastbuilder https://github.com/LNSSPsd/PhoenixBuilder  
图形界面/跨平台编译技术来自 Fyne: https://github.com/fyne-io/fyne    

~~PhoenixBuilderHeadless 是原FB项目的无头版本，尽量减少对原项目的修改：  
https://github.com/LNSSPsd/PhoenixBuilder~~

没办法，为了能让项目在安卓上编译，不得不修改fastbuilder->现在的fb文件夹  
PhoenixBuilder GPLv3 协议，项目的核心:  
https://github.com/LNSSPsd/PhoenixBuilder  
根据协议要求，本项目同为 GPL v3 协议    
除了该项目外，本项目：
- 使用了 go-raknet 源代码MIT协议  
https://github.com/Sandertv/go-raknet 
- 使用了gophertunnel 源代码完成数据包解析等操作，MIT协议  
https://github.com/Sandertv/gophertunnel
- 使用了 dragonfly MC 服务器框架源代码完成对MC服务器的模拟  
https://github.com/df-mc/dragonfly  
- 字体来自 Consolas-with-Yahei （因为需要跨平台，所以内嵌了20+MB的字体文件）  
https://github.com/crvdgc/Consolas-with-Yahei 

## 编译
### 对于Windows/Linux/Mac
```
fyne package -os linux
fyne package -os windows
fyne package -os darwin
```
或者，也可以简单的go build 一下
```
go build main.go
```
### 对于 android：
1. 准备环境，ndk，adb，另外：
```
go get fyne.io/fyne/v2/cmd/fyne
go install fyne.io/fyne/v2/cmd/fyne
```
2. 根据你的系统，可能需要一些额外操作，总之，当在命令行中输入 fyne 有输出即为成功
3. 编译
```
fyne package -os android -appID my.domain.appname （官方说明，windows上似乎无法正常工作）
fyne package -os android/arm64 -appID phoenixbuilder.third.gui -release true （个人建议）
```
4. 安装测试，windows上似乎无法正常工作
```
fyne install -os android
```
### 对于ios：
你需要一个许可证文件：
```
fyne release -os ios -certificate "Apple Distribution" -profile "My App Distribution" -appID "com.example.myapp"
```
### 更多
参考  
https://developer.fyne.io/started/cross-compiling   
https://developer.fyne.io/started/packaging   
的编译说明

## 致谢
感谢 Ruphane 在该程序开发和测试中的帮助  
感谢 CodePwn 帮忙测试和反馈问题  
以及 fyne 库的开发者