package tgcomp

import "github.com/VoileLab/toolgui/toolgui/tgcomp/tcmisc"

// Echo will execute lambda and show the code in the lambda.
// To use Echo, we need to store the code in advance (usually by embedded).
//
//	//go:embed main.go
//	var code string
//	// ...
//	// ok, echo will execute and show `tccontent.Text(c, "hello echo")`
//	tcmisc.Echo(c, code, func() {
//		tccontent.Text(c, "hello echo")
//	})
//
//	// panic, since Echo only parse code line by line
//	tcmisc.Echo(c, code, func() {tccontent.Text(c, "hello echo")})
//
//	// panic, since Echo only parse code that start from caller
//	myFunc := func() {
//		tccontent.Text(c, "hello echo")
//	}
//	tcmisc.Echo(c, code, myFunc)
var Echo = tcmisc.Echo

var Message = tcmisc.Message

type MessageConf = tcmisc.MessageConf

var MessageWithConf = tcmisc.MessageWithConf

var ProgressBar = tcmisc.ProgressBar
