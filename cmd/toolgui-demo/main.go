package main

import (
	_ "embed"
	"fmt"
	"log"

	"github.com/mudream4869/toolgui/toolgui/component"
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
	component.Title(c, "Example for ToolGUI")
	component.Code(c, code, "go")
	return nil
}

func MainPage(s *framework.Session, c *framework.Container, _ *framework.Container) error {
	component.Markdown(c, readme)
	return nil
}

func SidebarPage(s *framework.Session, c *framework.Container, sidebar *framework.Container) error {
	if component.Checkbox(s, c, "Show sidebar") {
		component.Text(sidebar, "Sidebar is here")
	}

	component.Code(c, code, "go")
	return nil
}

func ContentPage(s *framework.Session, c *framework.Container, _ *framework.Container) error {
	headerCompCol, headerCodeCol := component.Column2(c, "header_of_rows")
	component.Subtitle(headerCompCol, "Component")
	component.Subtitle(headerCodeCol, "Code")

	titleCompCol, titleCodeCol := component.Column2(c, "show_title")
	tcmisc.Echo(titleCodeCol, code, func() {
		component.Title(titleCompCol, "Title")
	})

	component.Divider(c)

	subtitleCompCol, subtitleCodeCol := component.Column2(c, "show_subtitle")
	tcmisc.Echo(subtitleCodeCol, code, func() {
		component.Subtitle(subtitleCompCol, "Subtitle")
	})

	component.Divider(c)

	textCompCol, textCodeCol := component.Column2(c, "show_text")
	tcmisc.Echo(textCodeCol, code, func() {
		component.Text(textCompCol, "Text")
	})

	component.Divider(c)

	imageCompCol, imageCodeCol := component.Column2(c, "show_image")
	tcmisc.Echo(imageCodeCol, code, func() {
		component.ImageByURL(imageCompCol, "https://http.cat/100")
	})

	component.Divider(c)

	dividerCompCol, dividerCodeCol := component.Column2(c, "show_divier")
	tcmisc.Echo(dividerCodeCol, code, func() {
		component.Divider(dividerCompCol)
	})

	return nil
}

func DataPage(s *framework.Session, c *framework.Container, _ *framework.Container) error {
	headerCompCol, headerCodeCol := component.Column2(c, "header_of_rows")
	component.Subtitle(headerCompCol, "Component")
	component.Subtitle(headerCodeCol, "Code")

	jsonCompCol, jsonCodeCol := component.Column2(c, "show_json")

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

		component.JSON(jsonCompCol, &DemoJSON{})
	})

	component.Divider(c)

	tableCompCol, tableCodeCol := component.Column2(c, "show_table")
	tcmisc.Echo(tableCodeCol, code, func() {
		component.Table(tableCompCol, []string{"a", "b"},
			[][]string{{"1", "2"}, {"3", "4"}})
	})

	return nil
}

func LayoutPage(s *framework.Session, c *framework.Container, _ *framework.Container) error {
	headerCompCol, headerCodeCol := component.Column2(c, "header_of_rows")
	component.Subtitle(headerCompCol, "Component")
	component.Subtitle(headerCodeCol, "Code")

	colCompCol, colCodeCol := component.Column2(c, "show_col")
	tcmisc.Echo(colCodeCol, code, func() {
		cols := component.Column(colCompCol, "cols", 3)
		for i, col := range cols {
			component.Text(col, fmt.Sprintf("col-%d", i))
		}
	})

	component.Divider(c)

	boxCompCol, boxCodeCol := component.Column2(c, "show_box")
	tcmisc.Echo(boxCodeCol, code, func() {
		box := component.Box(boxCompCol, "box")
		component.Text(box, "A box!")
	})

	return nil
}

func InputPage(s *framework.Session, c *framework.Container, _ *framework.Container) error {
	headerCompCol, headerCodeCol := component.Column2(c, "header_of_rows")
	component.Subtitle(headerCompCol, "Component")
	component.Subtitle(headerCodeCol, "Code")

	textareaCompCol, textareaCodeCol := component.Column2(c, "show_textarea")
	tcmisc.Echo(textareaCodeCol, code, func() {
		textareaValue := component.Textarea(s, textareaCompCol, "Textarea")
		component.Text(textareaCompCol, "Textarea Value: "+textareaValue)
	})

	component.Divider(c)

	textboxCompCol, textboxCodeCol := component.Column2(c, "show_textbox")
	tcmisc.Echo(textboxCodeCol, code, func() {
		textboxValue := component.Textbox(s, textboxCompCol, "Textbox")
		component.Text(textboxCompCol, "Textbox Value: "+textboxValue)
	})

	component.Divider(c)

	component.Text(c, "Currently we can only upload file that is smaller than 100k.")

	fileuploadCompCol, fileuploadCodeCol := component.Column2(c, "show_fileupload")
	tcmisc.Echo(fileuploadCodeCol, code, func() {
		fileObj := component.Fileupload(s, fileuploadCompCol, "Fileupload")
		component.Text(fileuploadCompCol, "Fileupload filename: "+fileObj.Name)
		component.ImageByURL(fileuploadCompCol, fileObj.Body)
	})

	component.Divider(c)

	checkboxCompCol, checkboxCodeCol := component.Column2(c, "show_checkbox")
	tcmisc.Echo(checkboxCodeCol, code, func() {
		checkboxValue := component.Checkbox(s, checkboxCompCol, "Checkbox")
		if checkboxValue {
			component.Text(checkboxCompCol, "Checkbox Value: true")
		} else {
			component.Text(checkboxCompCol, "Checkbox Value: false")
		}
	})

	component.Divider(c)

	buttonCompCol, buttonCodeCol := component.Column2(c, "show_button")
	tcmisc.Echo(buttonCodeCol, code, func() {
		btnClicked := component.Button(s, buttonCompCol, "button")
		if btnClicked {
			component.Text(buttonCompCol, "Button Value: true")
		} else {
			component.Text(buttonCompCol, "Button Value: false")
		}
	})

	component.Divider(c)

	selectCompCol, selectCodeCol := component.Column2(c, "show_select")
	tcmisc.Echo(selectCodeCol, code, func() {
		selValue := component.Select(s, selectCompCol,
			"Select", []string{"Value1", "Value2"})
		component.Text(selectCompCol, "Select Value: "+selValue)
	})

	component.Divider(c)

	radioCompCol, radioCodeCol := component.Column2(c, "show_radio")
	tcmisc.Echo(radioCodeCol, code, func() {
		radioValue := component.Radio(s, radioCompCol,
			"Radio", []string{"Value3", "Value4"})
		component.Text(radioCompCol, "Radio Value: "+radioValue)
	})

	return nil
}

func main() {
	e := executor.NewWebExecutor()
	e.AddPage("index", "Index", MainPage)
	e.AddPage("content", "Content", ContentPage)
	e.AddPage("data", "Data", DataPage)
	e.AddPage("layout", "Layout", LayoutPage)
	e.AddPage("sidebar", "Sidebar", SidebarPage)
	e.AddPage("input", "Input", InputPage)
	e.AddPage("code", "Source Code", SourceCodePage)
	log.Println("Starting service...")
	e.StartService(":3000")
}
