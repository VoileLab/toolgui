package toolguidemo

import (
	_ "embed"

	"fmt"

	"github.com/mudream4869/toolgui/toolgui/component/tccontent"
	"github.com/mudream4869/toolgui/toolgui/component/tcdata"
	"github.com/mudream4869/toolgui/toolgui/component/tcinput"
	"github.com/mudream4869/toolgui/toolgui/component/tclayout"
	"github.com/mudream4869/toolgui/toolgui/component/tcmisc"
	"github.com/mudream4869/toolgui/toolgui/framework"
)

//go:embed demo.go
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

func sourceCodePage(s *framework.Session, c *framework.Container, _ *framework.Container) error {
	tccontent.Title(c, "Example for ToolGUI")
	tccontent.Code(c, code, "go")
	return nil
}

func mainPage(s *framework.Session, c *framework.Container, _ *framework.Container) error {
	tccontent.Markdown(c, readme)
	return nil
}

func sidebarPage(s *framework.Session, c *framework.Container, sidebar *framework.Container) error {
	if tcinput.Checkbox(s, c, "Show sidebar") {
		tccontent.Text(sidebar, "Sidebar is here")
	}

	tccontent.Code(c, code, "go")
	return nil
}

func contentPage(s *framework.Session, c *framework.Container, _ *framework.Container) error {
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

func dataPage(s *framework.Session, c *framework.Container, _ *framework.Container) error {
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

func layoutPage(s *framework.Session, c *framework.Container, _ *framework.Container) error {
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

func inputPage(s *framework.Session, c *framework.Container, _ *framework.Container) error {
	headerCompCol, headerCodeCol := tclayout.Column2(c, "header_of_rows")
	tccontent.Subtitle(headerCompCol, "Component")
	tccontent.Subtitle(headerCodeCol, "Code")

	textareaCompCol, textareaCodeCol := tclayout.Column2(c, "show_textarea")
	tcmisc.Echo(textareaCodeCol, code, func() {
		textareaValue := tcinput.Textarea(s, textareaCompCol, "Textarea")
		tccontent.TextWithID(textareaCompCol, "Value: "+textareaValue, "textarea_result")
	})

	tccontent.DividerWithID(c, "1")

	textboxCompCol, textboxCodeCol := tclayout.Column2(c, "show_textbox")
	tcmisc.Echo(textboxCodeCol, code, func() {
		textboxValue := tcinput.Textbox(s, textboxCompCol, "Textbox")
		tccontent.TextWithID(textboxCompCol, "Value: "+textboxValue, "textbox_result")
	})

	tccontent.DividerWithID(c, "2")

	tccontent.Text(c, "Currently we can only upload file that is smaller than 100k.")

	fileuploadCompCol, fileuploadCodeCol := tclayout.Column2(c, "show_fileupload")
	tcmisc.Echo(fileuploadCodeCol, code, func() {
		fileObj := tcinput.Fileupload(s, fileuploadCompCol, "Fileupload")
		tccontent.Text(fileuploadCompCol, "Fileupload filename: "+fileObj.Name)
		tccontent.ImageByURI(fileuploadCompCol, fileObj.Body)
	})

	tccontent.DividerWithID(c, "3")

	checkboxCompCol, checkboxCodeCol := tclayout.Column2(c, "show_checkbox")
	tcmisc.Echo(checkboxCodeCol, code, func() {
		checkboxValue := tcinput.Checkbox(s, checkboxCompCol, "Checkbox")
		if checkboxValue {
			tccontent.TextWithID(checkboxCompCol, "Value: true", "checkbox_result")
		} else {
			tccontent.TextWithID(checkboxCompCol, "Value: false", "checkbox_result")
		}
	})

	tccontent.DividerWithID(c, "4")

	buttonCompCol, buttonCodeCol := tclayout.Column2(c, "show_button")
	tcmisc.Echo(buttonCodeCol, code, func() {
		btnClicked := tcinput.Button(s, buttonCompCol, "button")
		if btnClicked {
			tccontent.TextWithID(buttonCompCol, "Value: true", "button_result")
		} else {
			tccontent.TextWithID(buttonCompCol, "Value: false", "button_result")
		}
	})

	tccontent.DividerWithID(c, "5")

	selectCompCol, selectCodeCol := tclayout.Column2(c, "show_select")
	tcmisc.Echo(selectCodeCol, code, func() {
		selValue := tcinput.Select(s, selectCompCol,
			"Select", []string{"Value1", "Value2"})
		tccontent.TextWithID(selectCompCol, "Value: "+selValue, "select_result")
	})

	tccontent.DividerWithID(c, "6")

	radioCompCol, radioCodeCol := tclayout.Column2(c, "show_radio")
	tcmisc.Echo(radioCodeCol, code, func() {
		radioValue := tcinput.Radio(s, radioCompCol,
			"Radio", []string{"Value3", "Value4"})
		tccontent.TextWithID(radioCompCol, "Value: "+radioValue, "radio_result")
	})

	tccontent.DividerWithID(c, "7")

	datepickerCompCol, datepickerCodeCol := tclayout.Column2(c, "show_datepicker")
	tcmisc.Echo(datepickerCodeCol, code, func() {
		dateValue := tcinput.Datepicker(s, datepickerCompCol, "Datepicker")
		tccontent.TextWithID(datepickerCompCol, "Value: "+dateValue, "datepicker_result")
	})

	tccontent.DividerWithID(c, "8")

	timepickerCompCol, timepickerCodeCol := tclayout.Column2(c, "show_timepicker")
	tcmisc.Echo(timepickerCodeCol, code, func() {
		dateValue := tcinput.Timepicker(s, timepickerCompCol, "Timepicker")
		tccontent.TextWithID(timepickerCompCol, "Value: "+dateValue, "timepicker_result")
	})

	tccontent.DividerWithID(c, "9")

	datetimepickerCompCol, datetimepickerCodeCol := tclayout.Column2(c, "show_datetimepicker")
	tcmisc.Echo(datetimepickerCodeCol, code, func() {
		dateValue := tcinput.Datetimepicker(s, datetimepickerCompCol, "Weekpicker")
		tccontent.TextWithID(datetimepickerCompCol, "Value: "+dateValue, "weekpicker_result")
	})

	return nil
}

func miscPage(s *framework.Session, c *framework.Container, _ *framework.Container) error {
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

func NewApp() *framework.App {
	app := framework.NewApp()
	app.AddPage("index", "Index", mainPage)
	app.AddPage("content", "Content", contentPage)
	app.AddPage("data", "Data", dataPage)
	app.AddPage("input", "Input", inputPage)
	app.AddPage("layout", "Layout", layoutPage)
	app.AddPage("misc", "Misc", miscPage)
	app.AddPage("sidebar", "Sidebar", sidebarPage)
	app.AddPage("code", "Source Code", sourceCodePage)
	return app
}
