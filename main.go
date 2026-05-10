package main

import (
	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
)

//functions=========================================

// end of functions/begining of main=================
func main() {
	go func(){
		w := new(app.Window)                           //creating a window
		w.Option(app.Title("App"))                     //giving a title
		w.Option(app.Size(unit.Dp(500), unit.Dp(500))) //setting up the start size of the window
		th := material.NewTheme()                      //creating a theme th is a variable
		w := new(app.Window)                           //creating a window
		w.Option(app.Title("App"))                     //giving a title
		w.Option(app.Size(unit.Dp(500), unit.Dp(500))) //setting up the start size of the window
		th := material.NewTheme()                      //creating a theme th is a variable

		//==================================
		var ops op.Ops
		var ed widget.Editor
		var searchBtn widget.Clickable // button
		//==================================
		var list widget.List //creating list analog of Listbox in C# winforms
		list.Axis = layout.Vertical

		//==================================
		for {
			switch e := w.Event().(type) {
			case app.DestroyEvent:
				os.Exit(0)

			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)
	}
}
