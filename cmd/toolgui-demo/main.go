package main

import (
	_ "embed"
	"fmt"
	"log"

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

func SourceCodePage(s *framework.Session, c *framework.Container, _ *framework.Container) error {
	tccontent.Title(c, "Example for ToolGUI")
	tccontent.Code(c, code, "go")
	return nil
}

func MainPage(s *framework.Session, c *framework.Container, _ *framework.Container) error {
	tccontent.Markdown(c, readme)
	return nil
}

func SidebarPage(s *framework.Session, c *framework.Container, sidebar *framework.Container) error {
	if tcinput.Checkbox(s, c, "Show sidebar") {
		tccontent.Text(sidebar, "Sidebar is here")
	}

	tccontent.Code(c, code, "go")
	return nil
}

func ContentPage(s *framework.Session, c *framework.Container, _ *framework.Container) error {
	headerCompCol, headerCodeCol := tclayout.Column2(c, "header_of_rows")
	tccontent.Subtitle(headerCompCol, "Component")
	tccontent.Subtitle(headerCodeCol, "Code")

	titleCompCol, titleCodeCol := tclayout.Column2(c, "show_title")
	tcmisc.Echo(titleCodeCol, code, func() {
		tccontent.Title(titleCompCol, "Title")
	})

	tccontent.Divider(c)

	subtitleCompCol, subtitleCodeCol := tclayout.Column2(c, "show_subtitle")
	tcmisc.Echo(subtitleCodeCol, code, func() {
		tccontent.Subtitle(subtitleCompCol, "Subtitle")
	})

	tccontent.Divider(c)

	textCompCol, textCodeCol := tclayout.Column2(c, "show_text")
	tcmisc.Echo(textCodeCol, code, func() {
		tccontent.Text(textCompCol, "Text")
	})

	tccontent.Divider(c)

	imageCompCol, imageCodeCol := tclayout.Column2(c, "show_image")
	tcmisc.Echo(imageCodeCol, code, func() {
		tccontent.ImageByURI(imageCompCol, "https://http.cat/100")
	})

	tccontent.Divider(c)

	dividerCompCol, dividerCodeCol := tclayout.Column2(c, "show_divier")
	tcmisc.Echo(dividerCodeCol, code, func() {
		tccontent.Divider(dividerCompCol)
	})

	return nil
}

func DataPage(s *framework.Session, c *framework.Container, _ *framework.Container) error {
	headerCompCol, headerCodeCol := tclayout.Column2(c, "header_of_rows")
	tccontent.Subtitle(headerCompCol, "Component")
	tccontent.Subtitle(headerCodeCol, "Code")

	jsonCompCol, jsonCodeCol := tclayout.Column2(c, "show_json")

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

	tccontent.Divider(c)

	tableCompCol, tableCodeCol := tclayout.Column2(c, "show_table")
	tcmisc.Echo(tableCodeCol, code, func() {
		tcdata.Table(tableCompCol, []string{"a", "b"},
			[][]string{{"1", "2"}, {"3", "4"}})
	})

	return nil
}

func LayoutPage(s *framework.Session, c *framework.Container, _ *framework.Container) error {
	headerCompCol, headerCodeCol := tclayout.Column2(c, "header_of_rows")
	tccontent.Subtitle(headerCompCol, "Component")
	tccontent.Subtitle(headerCodeCol, "Code")

	colCompCol, colCodeCol := tclayout.Column2(c, "show_col")
	tcmisc.Echo(colCodeCol, code, func() {
		cols := tclayout.Column(colCompCol, "cols", 3)
		for i, col := range cols {
			tccontent.Text(col, fmt.Sprintf("col-%d", i))
		}
	})

	tccontent.Divider(c)

	boxCompCol, boxCodeCol := tclayout.Column2(c, "show_box")
	tcmisc.Echo(boxCodeCol, code, func() {
		box := tclayout.Box(boxCompCol, "box")
		tccontent.Text(box, "A box!")
	})

	return nil
}

func InputPage(s *framework.Session, c *framework.Container, _ *framework.Container) error {
	headerCompCol, headerCodeCol := tclayout.Column2(c, "header_of_rows")
	tccontent.Subtitle(headerCompCol, "Component")
	tccontent.Subtitle(headerCodeCol, "Code")

	textareaCompCol, textareaCodeCol := tclayout.Column2(c, "show_textarea")
	tcmisc.Echo(textareaCodeCol, code, func() {
		textareaValue := tcinput.Textarea(s, textareaCompCol, "Textarea")
		tccontent.Text(textareaCompCol, "Textarea Value: "+textareaValue)
	})

	tccontent.Divider(c)

	textboxCompCol, textboxCodeCol := tclayout.Column2(c, "show_textbox")
	tcmisc.Echo(textboxCodeCol, code, func() {
		textboxValue := tcinput.Textbox(s, textboxCompCol, "Textbox")
		tccontent.Text(textboxCompCol, "Textbox Value: "+textboxValue)
	})

	tccontent.Divider(c)

	tccontent.Text(c, "Currently we can only upload file that is smaller than 100k.")

	fileuploadCompCol, fileuploadCodeCol := tclayout.Column2(c, "show_fileupload")
	tcmisc.Echo(fileuploadCodeCol, code, func() {
		fileObj := tcinput.Fileupload(s, fileuploadCompCol, "Fileupload")
		tccontent.Text(fileuploadCompCol, "Fileupload filename: "+fileObj.Name)
		tccontent.ImageByURI(fileuploadCompCol, fileObj.Body)
	})

	tccontent.Divider(c)

	checkboxCompCol, checkboxCodeCol := tclayout.Column2(c, "show_checkbox")
	tcmisc.Echo(checkboxCodeCol, code, func() {
		checkboxValue := tcinput.Checkbox(s, checkboxCompCol, "Checkbox")
		if checkboxValue {
			tccontent.Text(checkboxCompCol, "Checkbox Value: true")
		} else {
			tccontent.Text(checkboxCompCol, "Checkbox Value: false")
		}
	})

	tccontent.Divider(c)

	buttonCompCol, buttonCodeCol := tclayout.Column2(c, "show_button")
	tcmisc.Echo(buttonCodeCol, code, func() {
		btnClicked := tcinput.Button(s, buttonCompCol, "button")
		if btnClicked {
			tccontent.Text(buttonCompCol, "Button Value: true")
		} else {
			tccontent.Text(buttonCompCol, "Button Value: false")
		}
	})

	tccontent.Divider(c)

	selectCompCol, selectCodeCol := tclayout.Column2(c, "show_select")
	tcmisc.Echo(selectCodeCol, code, func() {
		selValue := tcinput.Select(s, selectCompCol,
			"Select", []string{"Value1", "Value2"})
		tccontent.Text(selectCompCol, "Select Value: "+selValue)
	})

	tccontent.Divider(c)

	radioCompCol, radioCodeCol := tclayout.Column2(c, "show_radio")
	tcmisc.Echo(radioCodeCol, code, func() {
		radioValue := tcinput.Radio(s, radioCompCol,
			"Radio", []string{"Value3", "Value4"})
		tccontent.Text(radioCompCol, "Radio Value: "+radioValue)
	})

	tccontent.Divider(c)

	datepickerCompCol, datepickerCodeCol := tclayout.Column2(c, "show_datepicker")
	tcmisc.Echo(datepickerCodeCol, code, func() {
		dateValue := tcinput.Datepicker(s, datepickerCompCol, "Datepicker")
		tccontent.Text(datepickerCompCol, "Datepicker Value: "+dateValue)
	})

	tccontent.Divider(c)

	timepickerCompCol, timepickerCodeCol := tclayout.Column2(c, "show_timepicker")
	tcmisc.Echo(timepickerCodeCol, code, func() {
		dateValue := tcinput.Timepicker(s, timepickerCompCol, "Timepicker")
		tccontent.Text(timepickerCompCol, "Timepicker Value: "+dateValue)
	})

	tccontent.Divider(c)

	datetimepickerCompCol, datetimepickerCodeCol := tclayout.Column2(c, "show_datetimepicker")
	tcmisc.Echo(datetimepickerCodeCol, code, func() {
		dateValue := tcinput.Datetimepicker(s, datetimepickerCompCol, "Weekpicker")
		tccontent.Text(datetimepickerCompCol, "Weekpicker Value: "+dateValue)
	})

	return nil
}

func MiscPage(s *framework.Session, c *framework.Container, _ *framework.Container) error {
	headerCompCol, headerCodeCol := tclayout.Column2(c, "header_of_rows")
	tccontent.Subtitle(headerCompCol, "Component")
	tccontent.Subtitle(headerCodeCol, "Code")

	tccontent.Divider(c)

	echoCompCol, echoCodeCol := tclayout.Column2(c, "show_echo")
	tcmisc.Echo(echoCodeCol, code, func() {
		tcmisc.Echo(echoCompCol, code, func() {
			tccontent.Text(echoCompCol, "hello echo")
		})
	})

	tccontent.Divider(c)

	msgCompCol, msgCodeCol := tclayout.Column2(c, "show_msg")
	tcmisc.Echo(msgCodeCol, code, func() {
		tcmisc.Info(msgCompCol, "Title of msg", "body of msg")
	})

	tccontent.Divider(c)

	prgbarCompCol, prgbarCodeCol := tclayout.Column2(c, "show_progress_bar")
	tcmisc.Echo(prgbarCodeCol, code, func() {
		tcmisc.ProgressBar(prgbarCompCol, 30, "progress_bar")
	})

	return nil
}

func main() {
	e := executor.NewWebExecutor()
	e.AddPage("index", "Index", MainPage)
	e.AddPage("content", "Content", ContentPage)
	e.AddPage("data", "Data", DataPage)
	e.AddPage("input", "Input", InputPage)
	e.AddPage("layout", "Layout", LayoutPage)
	e.AddPage("misc", "Misc", MiscPage)
	e.AddPage("sidebar", "Sidebar", SidebarPage)
	e.AddPage("code", "Source Code", SourceCodePage)
	log.Println("Starting service...")
	e.StartService(":3000")
}
