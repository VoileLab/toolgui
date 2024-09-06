package main

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"image/jpeg"
	"log"
	"strings"

	"github.com/mudream4869/toolgui/toolgui/executor"
	"github.com/mudream4869/toolgui/toolgui/framework"
	"github.com/mudream4869/toolgui/toolgui/tgcomp"
)

//go:embed main.go
var code string

const readme = `
# [ToolGUI](https://github.com/mudream4869/toolgui)

This Go package provides a framework for rapidly building interactive data
dashboards and web applications. It aims to offer a similar development
experience to Streamlit for Python users.

> [!WARNING]
> ⚠️ Under Development:
> 
> The API for this package is still under development,
> and may be subject to changes in the future.`

func SourceCodePage(p *framework.Params) error {
	tgcomp.Title(p.Main, "Example for ToolGUI")
	tgcomp.Code(p.Main, code, "go")
	return nil
}

func MainPage(p *framework.Params) error {
	tgcomp.Markdown(p.Main, readme)
	return nil
}

func SidebarPage(p *framework.Params) error {
	if tgcomp.Checkbox(p.State, p.Main, "Show sidebar") {
		tgcomp.Text(p.Sidebar, "Sidebar is here")
	}

	tgcomp.Code(p.Main, code, "go")
	return nil
}

func ContentPage(p *framework.Params) error {
	headerCompCol, headerCodeCol := tgcomp.Column2(p.Main, "header_of_rows")
	tgcomp.Subtitle(headerCompCol, "Component")
	tgcomp.Subtitle(headerCodeCol, "Code")

	titleCompCol, titleCodeCol := tgcomp.Column2(p.Main, "show_title")
	tgcomp.Echo(titleCodeCol, code, func() {
		tgcomp.Title(titleCompCol, "Title")
	})

	tgcomp.Divider(p.Main)

	subtitleCompCol, subtitleCodeCol := tgcomp.Column2(p.Main, "show_subtitle")
	tgcomp.Echo(subtitleCodeCol, code, func() {
		tgcomp.Subtitle(subtitleCompCol, "Subtitle")
	})

	tgcomp.Divider(p.Main)

	textCompCol, textCodeCol := tgcomp.Column2(p.Main, "show_text")
	tgcomp.Echo(textCodeCol, code, func() {
		tgcomp.Text(textCompCol, "Text")
	})

	tgcomp.Divider(p.Main)

	imageCompCol, imageCodeCol := tgcomp.Column2(p.Main, "show_image")
	tgcomp.Echo(imageCodeCol, code, func() {
		tgcomp.ImageByURI(imageCompCol, "https://http.cat/100")
	})

	tgcomp.Divider(p.Main)

	dividerCompCol, dividerCodeCol := tgcomp.Column2(p.Main, "show_divier")
	tgcomp.Echo(dividerCodeCol, code, func() {
		tgcomp.Divider(dividerCompCol)
	})

	tgcomp.Divider(p.Main)

	linkCompCol, linkCodeCol := tgcomp.Column2(p.Main, "show_link")
	tgcomp.Echo(linkCodeCol, code, func() {
		tgcomp.Link(linkCompCol, "Link", "https://www.example.com/")
	})

	tgcomp.Divider(p.Main)

	downloadButtonCompCol, downloadButtonCodeCol := tgcomp.Column2(p.Main, "show_download_button")
	tgcomp.Echo(downloadButtonCodeCol, code, func() {
		tgcomp.DownloadButton(downloadButtonCompCol, "Download", []byte("123"), "123.txt")
	})

	return nil
}

func DataPage(p *framework.Params) error {
	headerCompCol, headerCodeCol := tgcomp.Column2(p.Main, "header_of_rows")
	tgcomp.Subtitle(headerCompCol, "Component")
	tgcomp.Subtitle(headerCodeCol, "Code")

	jsonCompCol, jsonCodeCol := tgcomp.Column2(p.Main, "show_json")

	tgcomp.Echo(jsonCodeCol, code, func() {
		type DemoJSONHeader struct {
			Type int
		}

		type DemoJSON struct {
			Header   DemoJSONHeader
			IntValue int
			URL      string
			IsOk     bool
		}

		tgcomp.JSON(jsonCompCol, &DemoJSON{})
	})

	tgcomp.Divider(p.Main)

	tableCompCol, tableCodeCol := tgcomp.Column2(p.Main, "show_table")
	tgcomp.Echo(tableCodeCol, code, func() {
		tgcomp.Table(tableCompCol, []string{"a", "b"},
			[][]string{{"1", "2"}, {"3", "4"}})
	})

	return nil
}

func LayoutPage(p *framework.Params) error {
	headerCompCol, headerCodeCol := tgcomp.Column2(p.Main, "header_of_rows")
	tgcomp.Subtitle(headerCompCol, "Component")
	tgcomp.Subtitle(headerCodeCol, "Code")

	colCompCol, colCodeCol := tgcomp.Column2(p.Main, "show_col")
	tgcomp.Echo(colCodeCol, code, func() {
		cols := tgcomp.Column(colCompCol, "cols", 3)
		for i, col := range cols {
			tgcomp.Text(col, fmt.Sprintf("col-%d", i))
		}
	})

	tgcomp.Divider(p.Main)

	boxCompCol, boxCodeCol := tgcomp.Column2(p.Main, "show_box")
	tgcomp.Echo(boxCodeCol, code, func() {
		box := tgcomp.Box(boxCompCol, "box")
		tgcomp.Text(box, "A box!")
	})

	return nil
}

func InputPage(p *framework.Params) error {
	headerCompCol, headerCodeCol := tgcomp.Column2(p.Main, "header_of_rows")
	tgcomp.Subtitle(headerCompCol, "Component")
	tgcomp.Subtitle(headerCodeCol, "Code")

	textareaCompCol, textareaCodeCol := tgcomp.Column2(p.Main, "show_textarea")
	tgcomp.Echo(textareaCodeCol, code, func() {
		textareaValue := tgcomp.Textarea(p.State, textareaCompCol, "Textarea", 5)
		tgcomp.TextWithID(textareaCompCol, "Value: "+textareaValue, "textarea_result")
	})

	tgcomp.DividerWithID(p.Main, "1")

	textboxCompCol, textboxCodeCol := tgcomp.Column2(p.Main, "show_textbox")
	tgcomp.Echo(textboxCodeCol, code, func() {
		textboxValue := tgcomp.Textbox(p.State, textboxCompCol, "Textbox")
		tgcomp.TextWithID(textboxCompCol, "Value: "+textboxValue, "textbox_result")
	})

	tgcomp.DividerWithID(p.Main, "2")

	fileuploadCompCol, fileuploadCodeCol := tgcomp.Column2(p.Main, "show_fileupload")
	tgcomp.Echo(fileuploadCodeCol, code, func() {
		fileObj := tgcomp.Fileupload(p.State, fileuploadCompCol,
			"Fileupload", ".jpg,.png")
		if fileObj == nil {
			return
		}

		tgcomp.Text(fileuploadCompCol, "Fileupload filename: "+fileObj.Name)
		tgcomp.Text(fileuploadCompCol,
			fmt.Sprintf("Fileupload bytes length: %d", len(fileObj.Bytes)))
		if strings.HasSuffix(fileObj.Name, ".jpg") {
			img, err := jpeg.Decode(bytes.NewReader(fileObj.Bytes))
			if err == nil {
				tgcomp.Image(fileuploadCompCol, img)
			}
		}
	})

	tgcomp.DividerWithID(p.Main, "3")

	checkboxCompCol, checkboxCodeCol := tgcomp.Column2(p.Main, "show_checkbox")
	tgcomp.Echo(checkboxCodeCol, code, func() {
		checkboxValue := tgcomp.Checkbox(p.State, checkboxCompCol, "Checkbox")
		tgcomp.TextWithID(checkboxCompCol,
			fmt.Sprint("Value: ", checkboxValue), "checkbox_result")
	})

	tgcomp.DividerWithID(p.Main, "4")

	buttonCompCol, buttonCodeCol := tgcomp.Column2(p.Main, "show_button")
	tgcomp.Echo(buttonCodeCol, code, func() {
		btnClicked := tgcomp.Button(p.State, buttonCompCol, "button")
		tgcomp.TextWithID(buttonCompCol,
			fmt.Sprint("Value: ", btnClicked), "button_result")
	})

	tgcomp.DividerWithID(p.Main, "5")

	selectCompCol, selectCodeCol := tgcomp.Column2(p.Main, "show_select")
	tgcomp.Echo(selectCodeCol, code, func() {
		selValue := tgcomp.Select(p.State, selectCompCol,
			"Select", []string{"Value1", "Value2"})
		tgcomp.TextWithID(selectCompCol, "Value: "+selValue, "select_result")
	})

	tgcomp.DividerWithID(p.Main, "6")

	radioCompCol, radioCodeCol := tgcomp.Column2(p.Main, "show_radio")
	tgcomp.Echo(radioCodeCol, code, func() {
		radioValue := tgcomp.Radio(p.State, radioCompCol,
			"Radio", []string{"Value3", "Value4"})
		tgcomp.TextWithID(radioCompCol, "Value: "+radioValue, "radio_result")
	})

	tgcomp.DividerWithID(p.Main, "7")

	datepickerCompCol, datepickerCodeCol := tgcomp.Column2(p.Main, "show_datepicker")
	tgcomp.Echo(datepickerCodeCol, code, func() {
		dateValue := tgcomp.Datepicker(p.State, datepickerCompCol, "Datepicker")
		tgcomp.TextWithID(datepickerCompCol, "Value: "+dateValue, "datepicker_result")
	})

	tgcomp.DividerWithID(p.Main, "8")

	timepickerCompCol, timepickerCodeCol := tgcomp.Column2(p.Main, "show_timepicker")
	tgcomp.Echo(timepickerCodeCol, code, func() {
		dateValue := tgcomp.Timepicker(p.State, timepickerCompCol, "Timepicker")
		tgcomp.TextWithID(timepickerCompCol, "Value: "+dateValue, "timepicker_result")
	})

	tgcomp.DividerWithID(p.Main, "9")

	datetimepickerCompCol, datetimepickerCodeCol := tgcomp.Column2(p.Main, "show_datetimepicker")
	tgcomp.Echo(datetimepickerCodeCol, code, func() {
		dateValue := tgcomp.Datetimepicker(p.State, datetimepickerCompCol, "Datetimepicker")
		tgcomp.TextWithID(datetimepickerCompCol, "Value: "+dateValue, "datetimepicker_result")
	})

	return nil
}

func MiscPage(p *framework.Params) error {
	headerCompCol, headerCodeCol := tgcomp.Column2(p.Main, "header_of_rows")
	tgcomp.Subtitle(headerCompCol, "Component")
	tgcomp.Subtitle(headerCodeCol, "Code")

	tgcomp.Divider(p.Main)

	echoCompCol, echoCodeCol := tgcomp.Column2(p.Main, "show_echo")
	tgcomp.Echo(echoCodeCol, code, func() {
		tgcomp.Echo(echoCompCol, code, func() {
			tgcomp.Text(echoCompCol, "hello echo")
		})
	})

	tgcomp.Divider(p.Main)

	msgCompCol, msgCodeCol := tgcomp.Column2(p.Main, "show_msg")
	tgcomp.Echo(msgCodeCol, code, func() {
		tgcomp.Info(msgCompCol, "Title of msg", "body of msg")
	})

	tgcomp.Divider(p.Main)

	prgbarCompCol, prgbarCodeCol := tgcomp.Column2(p.Main, "show_progress_bar")
	tgcomp.Echo(prgbarCodeCol, code, func() {
		tgcomp.ProgressBar(prgbarCompCol, 30, "progress_bar")
	})

	tgcomp.Divider(p.Main)

	errorCompCol, errorCodeCol := tgcomp.Column2(p.Main, "show_error")
	if tgcomp.Button(p.State, errorCompCol, "Show error") {
		return errors.New("new error")
	}
	tgcomp.Code(errorCodeCol, `if tgcomp.Button(p.State, errorCompCol, "Show error") {
	return errors.New("New error")
}`, "go")

	tgcomp.Divider(p.Main)

	panicCompCol, panicCodeCol := tgcomp.Column2(p.Main, "show_panic")
	tgcomp.Echo(panicCodeCol, code, func() {
		if tgcomp.Button(p.State, panicCompCol, "Show panic") {
			panic("show panic")
		}
	})

	return nil
}

func main() {
	app := framework.NewApp()
	app.AddPage("index", "Index", MainPage)
	app.AddPage("content", "Content", ContentPage)
	app.AddPage("data", "Data", DataPage)
	app.AddPage("input", "Input", InputPage)
	app.AddPage("layout", "Layout", LayoutPage)
	app.AddPage("misc", "Misc", MiscPage)
	app.AddPage("sidebar", "Sidebar", SidebarPage)
	app.AddPage("code", "Source Code", SourceCodePage)

	e := executor.NewWebExecutor(app)
	log.Println("Starting service...")
	e.StartService(":3000")
}
