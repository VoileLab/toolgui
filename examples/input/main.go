package main

import (
	"log"

	"github.com/mudream4869/toolgui/toolgui/component"
	"github.com/mudream4869/toolgui/toolgui/executor"
	"github.com/mudream4869/toolgui/toolgui/framework"
)

func Main(s *framework.Session, c *framework.Container) error {
	headerCompCol, headerCodeCol := component.Column2(c, "header")
	component.Subtitle(headerCompCol, "Component Column")
	component.Subtitle(headerCodeCol, "Code Column")

	textareaCompCol, textareaCodeCol := component.Column2(c, "show_textarea")
	textareaValue := component.Textarea(s, textareaCompCol, "Textarea")
	component.Text(textareaCompCol, "Textarea Value: "+textareaValue)
	component.Code(textareaCodeCol, `textareaValue := component.Textbox(s, textareaCompCol, "Textarea")
component.Text(c, "Textarea Value: "+textareaValue)`, "go")

	component.Divider(c)

	textboxCompCol, textboxCodeCol := component.Column2(c, "show_textbox")
	textboxValue := component.Textbox(s, textboxCompCol, "Textbox")
	component.Text(textboxCompCol, "Textbox Value: "+textboxValue)
	component.Code(textboxCodeCol, `textboxValue := component.Textbox(s, textboxCompCol, "Textbox")
component.Text(c, "Textbox Value: "+textboxValue)`, "go")

	component.Divider(c)

	fileuploadCompCol, fileuploadCodeCol := component.Column2(c, "show_fileupload")
	fileObj := component.Fileupload(s, fileuploadCompCol, "Fileupload")
	component.Text(fileuploadCompCol, "Fileupload filename: "+fileObj.Name)
	component.ImageByURL(fileuploadCompCol, fileObj.Body)
	component.Code(fileuploadCodeCol, `component.Textarea(s, fileuploadCompCol, "Fileupload")
component.Text(fileuploadCompCol, "Fileupload filename: "+fileObj.Name)
component.ImageByURL(fileuploadCompCol, fileObj.Body)`, "go")

	component.Divider(c)

	checkboxCompCol, checkboxCodeCol := component.Column2(c, "show_checkbox")
	checkboxValue := component.Checkbox(s, checkboxCompCol, "Checkbox")
	if checkboxValue {
		component.Text(checkboxCompCol, "Checkbox Value: true")
	} else {
		component.Text(checkboxCompCol, "Checkbox Value: false")
	}
	component.Code(checkboxCodeCol, `checkboxValue := component.Checkbox(s, checkboxCompCol, "Checkbox")
if checkboxValue {
	component.Text(c, "Checkbox Value: true")
} else {
	component.Text(c, "Checkbox Value: false")
}`, "go")

	component.Divider(c)

	buttonCompCol, buttonCodeCol := component.Column2(c, "show_button")
	btnClicked := component.Button(s, buttonCompCol, "button")
	if btnClicked {
		component.Text(buttonCompCol, "Button Value: true")
	} else {
		component.Text(buttonCompCol, "Button Value: false")
	}
	component.Code(buttonCodeCol, `btnClicked := component.Button(s, buttonCompCol, "button")
if btnClicked {
	component.Text(buttonCompCol, "Button Value: true")
} else {
	component.Text(buttonCompCol, "Button Value: false")
}`, "go")

	component.Divider(c)

	selectCompCol, selectCodeCol := component.Column2(c, "show_select")
	selValue := component.Select(s, selectCompCol,
		"Select", []string{"Value1", "Value2"})
	component.Text(selectCompCol, "Select Value: "+selValue)
	component.Code(selectCodeCol, `selValue := component.Select(s, selectCompCol,
	"Select", []string{"Value1", "Value2"})
component.Text(selectCodeCol, "Select Value: "+selValue)`, "go")

	return nil
}

func main() {
	e := executor.NewWebExecutor()
	e.AddPage("index", "Index", Main)
	log.Println("Starting service...")
	e.StartService(":3000")
}
