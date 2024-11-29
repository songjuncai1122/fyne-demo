## 打包命令

```shell
# 1. 安装 fyne 官方工具
go get fyne.io/fyne/v2/cmd/fyne

# 2. 打包字体
fyne bundle HarmonyOS_Sans_Condensed_Medium.ttf > bundle.go

# 3. 生成 Mac 包
fyne package -os darwin -icon license.png
```
