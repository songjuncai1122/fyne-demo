package main

import (
	"fmt"
	"fyne-demo/pkg/cryptoaes"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"time"
)

// 自定义字体主题
type myTheme struct {
	base fyne.Theme
}

func (m *myTheme) Font(style fyne.TextStyle) fyne.Resource {
	//fontPath := "HarmonyOS_Sans_Condensed_Medium.ttf" // 将此路径替换为您的中文字体路径
	//fontResource, err := fyne.LoadResourceFromPath(fontPath)
	//if err != nil {
	//	panic("无法加载字体: " + err.Error())
	//}
	//return fontResource

	return resourceHarmonyOSSansCondensedMediumTtf
}

func (m *myTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return m.base.Color(name, variant)
}

func (m *myTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return m.base.Icon(name)
}

func (m *myTheme) Size(name fyne.ThemeSizeName) float32 {
	return m.base.Size(name)
}

func main() {
	a := app.New()

	// 使用自定义主题，解决中文乱码并根据主题调整文字、按钮和生成码颜色
	currentTheme := &myTheme{base: theme.DefaultTheme()}
	a.Settings().SetTheme(currentTheme)

	w := a.NewWindow("授权码生成工具 V1.0.0")

	// 日期选择（年、月、日下拉选择框）
	currentYear := time.Now().Year()
	var years []string
	for i := currentYear; i <= currentYear+10; i++ {
		years = append(years, fmt.Sprintf("%d", i))
	}
	months := []string{"01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12"}
	var days []string
	for i := 1; i <= 31; i++ {
		days = append(days, fmt.Sprintf("%02d", i))
	}

	yearSelect := widget.NewSelect(years, func(value string) {})
	yearSelect.PlaceHolder = "选择年份"
	monthSelect := widget.NewSelect(months, func(value string) {})
	monthSelect.PlaceHolder = "选择月份"
	daySelect := widget.NewSelect(days, func(value string) {})
	daySelect.PlaceHolder = "选择日期"

	// 项目下拉选择
	projectOptions := []string{"项目A", "项目B", "项目C"}
	projectSelect := widget.NewSelect(projectOptions, func(value string) {})
	projectSelect.PlaceHolder = "请选择项目"

	// 授权码显示区域，使用 Label 控件展示生成后的授权码，并增加背景色和边框
	authCodeLabel := widget.NewLabel("")
	authCodeLabel.Wrapping = fyne.TextWrapWord // 设置自动换行
	authCodeLabel.Hide()                       // 初始隐藏

	// 创建带颜色的矩形作为背景
	authCodeBackground := canvas.NewRectangle(color.RGBA{R: 240, G: 240, B: 240, A: 255}) // 初始为浅灰色背景

	// 将 label 和背景叠加
	stack := container.NewStack(authCodeBackground, container.NewPadded(authCodeLabel))

	// 复制按钮
	copyButton := widget.NewButton("复制授权码", func() {
		clipboard := w.Clipboard() // 获取当前窗口的剪贴板对象
		clipboard.SetContent(authCodeLabel.Text)
		dialog.ShowInformation("已复制", "授权码已复制到剪贴板", w)
	})
	copyButton.Hide() // 初始隐藏

	// 退出授权按钮，点击关闭程序，并添加 "X" 图标
	exitButton := widget.NewButtonWithIcon("退出授权", theme.CancelIcon(), func() {
		a.Quit() // 退出应用程序
	})

	// 授权按钮
	authButton := widget.NewButtonWithIcon("生成授权", theme.ConfirmIcon(), func() {
		year := yearSelect.Selected
		month := monthSelect.Selected
		day := daySelect.Selected
		project := projectSelect.Selected

		if year == "" || month == "" || day == "" || project == "" {
			dialog.ShowError(fmt.Errorf("请完整填写日期和项目"), w)
			return
		}

		// 授权码生成逻辑（示例）
		date := fmt.Sprintf("%s-%s-%s", year, month, day)
		authCode := fmt.Sprintf("%s-%s-%d", project, date, time.Now().Unix())
		authAesCode, _ := cryptoaes.Encrypt("gHfSWrHzBkQ4Bk4a", authCode)
		authCodeLabel.SetText(authAesCode)

		// 显示授权码和复制按钮
		authCodeLabel.Show()
		copyButton.Show()

		// 刷新 Label 确保授权码显示
		authCodeLabel.Refresh()
	})

	// 皮肤选择
	themeOptions := []string{"默认", "浅色主题", "深色主题"}
	themeSelect := widget.NewSelect(themeOptions, func(value string) {
		switch value {
		case "默认":
			currentTheme.base = theme.DefaultTheme()
			authCodeBackground.FillColor = color.RGBA{R: 240, G: 240, B: 240, A: 255} // 浅色主题下，背景色为浅灰色
		case "浅色主题":
			currentTheme.base = theme.LightTheme()
			authCodeBackground.FillColor = color.RGBA{R: 240, G: 240, B: 240, A: 255} // 浅色主题下，背景色为浅灰色
		case "深色主题":
			currentTheme.base = theme.DarkTheme()
			authCodeBackground.FillColor = color.RGBA{R: 50, G: 50, B: 50, A: 255} // 深色主题下，背景色为深灰色
		}
		a.Settings().SetTheme(currentTheme) // 重新设置自定义主题
		authCodeBackground.Refresh()
	})
	themeSelect.PlaceHolder = "选择皮肤"

	// 布局设计
	buttons := container.NewGridWithColumns(2, authButton, exitButton) // 生成授权和退出授权按钮各占50%宽度

	form := container.NewVBox(
		projectSelect, // 将项目选择框放在时间选择框的上面
		container.NewAdaptiveGrid(3, yearSelect, monthSelect, daySelect),
		buttons,                         // 将生成授权和退出授权按钮放在同一行，各占50%
		stack,                           // 使用带颜色背景的堆叠布局
		container.NewCenter(copyButton), // 复制按钮居中显示
	)

	// 创建底部布局
	bottomBar := container.NewBorder(nil, nil, nil, themeSelect) // 添加皮肤选择器

	// 给内容增加边距和内边距
	paddedForm := container.NewVBox(
		container.NewPadded(form),
	)

	// 整体布局，确保底部栏和内容布局分开
	content := container.NewBorder(nil, bottomBar, nil, nil, paddedForm)

	w.SetContent(content)
	w.Resize(fyne.NewSize(350, 400)) // 调整窗口大小以适应高度设置
	w.CenterOnScreen()
	w.ShowAndRun()
}
