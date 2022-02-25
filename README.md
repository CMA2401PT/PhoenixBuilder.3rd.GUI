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

## 运行
你可以很简单的使用
```
go build main.go
```
编译出对应平台的程序

## 编译发行版
首先你需要安装必须的工具  
```
go get fyne.io/fyne/v2/cmd/fyne
go install fyne.io/fyne/v2/cmd/fyne
```

### 对于Windows/Linux/Mac

```
fyne package -os linux
fyne package -os windows
fyne package -os darwin
```
### 对于 android：
1. 准备环境，ndk，adb，并设置环境变量 ANDROID_NDK_HOME
2. 编译（windows上似乎无法正常工作）
```
fyne package -os android/arm64 -appID phoenixbuilder.third.gui -release true
```
3. 安装测试
```
fyne install -os android
```
### 对于ios：
你需要一个许可证文件：
```
fyne release -os ios -certificate "Apple Distribution" -profile "My App Distribution" -appID "phoenixbuilder.third.gui"
```

## 另一种编译方式(fyne-cross)
安装环境和工具
```
go get github.com/fyne-io/fyne-cross
go install github.com/fyne-io/fyne-cross
```
安装 docker，并想办法确保网络连接
编译 (输出在 fyne-cross/dist 目录下)
```
Linux:
fyne-cross linux -arch=amd64 -app-build 169 -app-id "fastbuilder.third.gui" -app-version 0.0.4 -icon unbundled_assets/Icon.png  -name "FastBuilder_3rd_Gui"

MacOS:
fyne-cross darwin -arch=amd64 -app-build 169 -app-id "fastbuilder.third.gui" -app-version 0.0.4 -icon unbundled_assets/Icon.png  -name "FastBuilder_3rd_Gui"

Windows:
fyne-cross windows -arch=amd64 -app-build 169 -app-id "fastbuilder.third.gui" -app-version 0.0.4 -icon unbundled_assets/Icon.png  -name "FastBuilder_3rd_Gui.exe"

Android:
fyne-cross android -arch=arm64 -app-build 169 -app-id "fastbuilder.third.gui" -app-version 0.0.4 -icon unbundled_assets/Icon.png  -name "FastBuilder_3rd_Gui"

IOS:
你需要创建一个开发者账号，并建立一个同名 Xcode项目"fastbuilder.third.gui"接着
fyne-cross ios -app-build 169 -app-id "fastbuilder.third.gui" -app-version 0.0.4 -icon unbundled_assets/Icon.png  -name "FastBuilder-3rd-Gui"
```

### 2022.2.25补充
现在配置好环境后，输入
```
bash fyne_cross_compile.sh
```
即可自动打包全平台的分发了

### 更多
参考  
https://developer.fyne.io/started/cross-compiling   
https://developer.fyne.io/started/packaging   
的编译说明

## 致谢
感谢 Ruphane 在该程序开发和测试中的帮助  
感谢 CodePwn 帮忙测试和反馈问题  
以及 fyne 库的开发者