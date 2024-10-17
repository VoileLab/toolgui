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

// Iframe and IframeWithID is experimental component, their feature is not stable.
// Use them with caution.

// Iframe show a iframe.
// script is used to allow the iframe to run javascript. (notice that this is not secure)
var Iframe = tcmisc.Iframe

// IframeWithID create a iframe component with a user specific id.
// script is used to allow the iframe to run javascript. (notice that this is not secure)
var IframeWithID = tcmisc.IframeWithID
