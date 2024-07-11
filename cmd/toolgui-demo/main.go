package main

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"image/jpeg"
	"log"
	"strings"

	"github.com/mudream4869/toolgui/toolgui/component/tccontent"
	"github.com/mudream4869/toolgui/toolgui/component/tcdata"
	"github.com/mudream4869/toolgui/toolgui/component/tcinput"
	"github.com/mudream4869/toolgui/toolgui/component/tclayout"
	"github.com/mudream4869/toolgui/toolgui/component/tcmisc"
	"github.com/mudream4869/toolgui/toolgui/executor"
	"github.com/mudream4869/toolgui/toolgui/framework"
)

//go:embed main.go
var code string

const readme = `
# ToolGUI

This Go package provides a framework for rapidly building interactive data
dashboards and web applications. It aims to offer a similar development
experience to Streamlit for Python users.

> [!WARNING]
> ⚠️ Under Development:
> 
> The API for this package is still under development,
> and may be subject to changes in the future.`

func SourceCodePage(p *framework.Params) error {
	tccontent.Title(p.Main, "Example for ToolGUI")
	tccontent.Code(p.Main, code, "go")
	return nil
}

func MainPage(p *framework.Params) error {
	tccontent.Markdown(p.Main, readme)
	return nil
}

func SidebarPage(p *framework.Params) error {
	if tcinput.Checkbox(p.State, p.Main, "Show sidebar") {
		tccontent.Text(p.Sidebar, "Sidebar is here")
	}

	tccontent.Code(p.Main, code, "go")
	return nil
}

func ContentPage(p *framework.Params) error {
	headerCompCol, headerCodeCol := tclayout.Column2(p.Main, "header_of_rows")
	tccontent.Subtitle(headerCompCol, "Component")
	tccontent.Subtitle(headerCodeCol, "Code")

	titleCompCol, titleCodeCol := tclayout.Column2(p.Main, "show_title")
	tcmisc.Echo(titleCodeCol, code, func() {
		tccontent.Title(titleCompCol, "Title")
	})

	tccontent.Divider(p.Main)

	subtitleCompCol, subtitleCodeCol := tclayout.Column2(p.Main, "show_subtitle")
	tcmisc.Echo(subtitleCodeCol, code, func() {
		tccontent.Subtitle(subtitleCompCol, "Subtitle")
	})

	tccontent.Divider(p.Main)

	textCompCol, textCodeCol := tclayout.Column2(p.Main, "show_text")
	tcmisc.Echo(textCodeCol, code, func() {
		tccontent.Text(textCompCol, "Text")
	})

	tccontent.Divider(p.Main)

	imageCompCol, imageCodeCol := tclayout.Column2(p.Main, "show_image")
	tcmisc.Echo(imageCodeCol, code, func() {
		tccontent.ImageByURI(imageCompCol, "https://http.cat/100")
	})

	tccontent.Divider(p.Main)

	dividerCompCol, dividerCodeCol := tclayout.Column2(p.Main, "show_divier")
	tcmisc.Echo(dividerCodeCol, code, func() {
		tccontent.Divider(dividerCompCol)
	})

	tccontent.Divider(p.Main)

	linkCompCol, linkCodeCol := tclayout.Column2(p.Main, "show_link")
	tcmisc.Echo(linkCodeCol, code, func() {
		tccontent.Link(linkCompCol, "Link", "https://www.example.com/")
	})

	return nil
}

func DataPage(p *framework.Params) error {
	headerCompCol, headerCodeCol := tclayout.Column2(p.Main, "header_of_rows")
	tccontent.Subtitle(headerCompCol, "Component")
	tccontent.Subtitle(headerCodeCol, "Code")

	jsonCompCol, jsonCodeCol := tclayout.Column2(p.Main, "show_json")

	tcmisc.Echo(jsonCodeCol, code, func() {
		type DemoJSONHeader struct {
			Type int
		}

		type DemoJSON struct {
			Header   DemoJSONHeader
			IntValue int
			URL      string
			IsOk     bool
		}

		tcdata.JSON(jsonCompCol, &DemoJSON{})
	})

	tccontent.Divider(p.Main)

	tableCompCol, tableCodeCol := tclayout.Column2(p.Main, "show_table")
	tcmisc.Echo(tableCodeCol, code, func() {
		tcdata.Table(tableCompCol, []string{"a", "b"},
			[][]string{{"1", "2"}, {"3", "4"}})
	})

	return nil
}

func LayoutPage(p *framework.Params) error {
	headerCompCol, headerCodeCol := tclayout.Column2(p.Main, "header_of_rows")
	tccontent.Subtitle(headerCompCol, "Component")
	tccontent.Subtitle(headerCodeCol, "Code")

	colCompCol, colCodeCol := tclayout.Column2(p.Main, "show_col")
	tcmisc.Echo(colCodeCol, code, func() {
		cols := tclayout.Column(colCompCol, "cols", 3)
		for i, col := range cols {
			tccontent.Text(col, fmt.Sprintf("col-%d", i))
		}
	})

	tccontent.Divider(p.Main)

	boxCompCol, boxCodeCol := tclayout.Column2(p.Main, "show_box")
	tcmisc.Echo(boxCodeCol, code, func() {
		box := tclayout.Box(boxCompCol, "box")
		tccontent.Text(box, "A box!")
	})

	return nil
}

func InputPage(p *framework.Params) error {
	headerCompCol, headerCodeCol := tclayout.Column2(p.Main, "header_of_rows")
	tccontent.Subtitle(headerCompCol, "Component")
	tccontent.Subtitle(headerCodeCol, "Code")

	textareaCompCol, textareaCodeCol := tclayout.Column2(p.Main, "show_textarea")
	tcmisc.Echo(textareaCodeCol, code, func() {
		textareaValue := tcinput.Textarea(p.State, textareaCompCol, "Textarea")
		tccontent.TextWithID(textareaCompCol, "Value: "+textareaValue, "textarea_result")
	})

	tccontent.DividerWithID(p.Main, "1")

	textboxCompCol, textboxCodeCol := tclayout.Column2(p.Main, "show_textbox")
	tcmisc.Echo(textboxCodeCol, code, func() {
		textboxValue := tcinput.Textbox(p.State, textboxCompCol, "Textbox")
		tccontent.TextWithID(textboxCompCol, "Value: "+textboxValue, "textbox_result")
	})

	tccontent.DividerWithID(p.Main, "2")

	tccontent.Text(p.Main, "Currently we can only upload file that is smaller than 100k.")

	fileuploadCompCol, fileuploadCodeCol := tclayout.Column2(p.Main, "show_fileupload")
	tcmisc.Echo(fileuploadCodeCol, code, func() {
		fileObj := tcinput.Fileupload(p.State, fileuploadCompCol, "Fileupload")
		tccontent.Text(fileuploadCompCol, "Fileupload filename: "+fileObj.Name)
		tccontent.Text(fileuploadCompCol,
			fmt.Sprintf("Fileupload bytes length: %d", len(fileObj.Bytes)))
		if strings.HasSuffix(fileObj.Name, ".jpg") {
			img, err := jpeg.Decode(bytes.NewReader(fileObj.Bytes))
			if err == nil {
				tccontent.Image(fileuploadCompCol, img)
			}
		}
	})

	tccontent.DividerWithID(p.Main, "3")

	checkboxCompCol, checkboxCodeCol := tclayout.Column2(p.Main, "show_checkbox")
	tcmisc.Echo(checkboxCodeCol, code, func() {
		checkboxValue := tcinput.Checkbox(p.State, checkboxCompCol, "Checkbox")
		if checkboxValue {
			tccontent.TextWithID(checkboxCompCol, "Value: true", "checkbox_result")
		} else {
			tccontent.TextWithID(checkboxCompCol, "Value: false", "checkbox_result")
		}
	})

	tccontent.DividerWithID(p.Main, "4")

	buttonCompCol, buttonCodeCol := tclayout.Column2(p.Main, "show_button")
	tcmisc.Echo(buttonCodeCol, code, func() {
		btnClicked := tcinput.Button(p.State, buttonCompCol, "button")
		if btnClicked {
			tccontent.TextWithID(buttonCompCol, "Value: true", "button_result")
		} else {
			tccontent.TextWithID(buttonCompCol, "Value: false", "button_result")
		}
	})

	tccontent.DividerWithID(p.Main, "5")

	selectCompCol, selectCodeCol := tclayout.Column2(p.Main, "show_select")
	tcmisc.Echo(selectCodeCol, code, func() {
		selValue := tcinput.Select(p.State, selectCompCol,
			"Select", []string{"Value1", "Value2"})
		tccontent.TextWithID(selectCompCol, "Value: "+selValue, "select_result")
	})

	tccontent.DividerWithID(p.Main, "6")

	radioCompCol, radioCodeCol := tclayout.Column2(p.Main, "show_radio")
	tcmisc.Echo(radioCodeCol, code, func() {
		radioValue := tcinput.Radio(p.State, radioCompCol,
			"Radio", []string{"Value3", "Value4"})
		tccontent.TextWithID(radioCompCol, "Value: "+radioValue, "radio_result")
	})

	tccontent.DividerWithID(p.Main, "7")

	datepickerCompCol, datepickerCodeCol := tclayout.Column2(p.Main, "show_datepicker")
	tcmisc.Echo(datepickerCodeCol, code, func() {
		dateValue := tcinput.Datepicker(p.State, datepickerCompCol, "Datepicker")
		tccontent.TextWithID(datepickerCompCol, "Value: "+dateValue, "datepicker_result")
	})

	tccontent.DividerWithID(p.Main, "8")

	timepickerCompCol, timepickerCodeCol := tclayout.Column2(p.Main, "show_timepicker")
	tcmisc.Echo(timepickerCodeCol, code, func() {
		dateValue := tcinput.Timepicker(p.State, timepickerCompCol, "Timepicker")
		tccontent.TextWithID(timepickerCompCol, "Value: "+dateValue, "timepicker_result")
	})

	tccontent.DividerWithID(p.Main, "9")

	datetimepickerCompCol, datetimepickerCodeCol := tclayout.Column2(p.Main, "show_datetimepicker")
	tcmisc.Echo(datetimepickerCodeCol, code, func() {
		dateValue := tcinput.Datetimepicker(p.State, datetimepickerCompCol, "Datetimepicker")
		tccontent.TextWithID(datetimepickerCompCol, "Value: "+dateValue, "datetimepicker_result")
	})

	return nil
}

func MiscPage(p *framework.Params) error {
	headerCompCol, headerCodeCol := tclayout.Column2(p.Main, "header_of_rows")
	tccontent.Subtitle(headerCompCol, "Component")
	tccontent.Subtitle(headerCodeCol, "Code")

	tccontent.Divider(p.Main)

	echoCompCol, echoCodeCol := tclayout.Column2(p.Main, "show_echo")
	tcmisc.Echo(echoCodeCol, code, func() {
		tcmisc.Echo(echoCompCol, code, func() {
			tccontent.Text(echoCompCol, "hello echo")
		})
	})

	tccontent.Divider(p.Main)

	msgCompCol, msgCodeCol := tclayout.Column2(p.Main, "show_msg")
	tcmisc.Echo(msgCodeCol, code, func() {
		tcmisc.Info(msgCompCol, "Title of msg", "body of msg")
	})

	tccontent.Divider(p.Main)

	prgbarCompCol, prgbarCodeCol := tclayout.Column2(p.Main, "show_progress_bar")
	tcmisc.Echo(prgbarCodeCol, code, func() {
		tcmisc.ProgressBar(prgbarCompCol, 30, "progress_bar")
	})

	tccontent.Divider(p.Main)

	errorCompCol, errorCodeCol := tclayout.Column2(p.Main, "show_error")
	if tcinput.Button(p.State, errorCompCol, "Show error") {
		return errors.New("new error")
	}
	tccontent.Code(errorCodeCol, `if tcinput.Button(p.State, errorCompCol, "Show error") {
	return errors.New("New error")
}`, "go")

	tccontent.Divider(p.Main)

	panicCompCol, panicCodeCol := tclayout.Column2(p.Main, "show_panic")
	tcmisc.Echo(panicCodeCol, code, func() {
		if tcinput.Button(p.State, panicCompCol, "Show panic") {
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
