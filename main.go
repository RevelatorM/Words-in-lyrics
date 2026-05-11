package main

import (
	"log"
	"os"

	"image/color"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// functions=========================================
func layoutBackground(gtx layout.Context) layout.Dimensions { // background color function
	// Define the color (RGBA)
	bgColor := color.NRGBA{R: 95, G: 0, B: 87, A: 80} // A can accept only int

	// Fill the constraints area with the color
	paint.Fill(gtx.Ops, bgColor)

	return layout.Dimensions{Size: gtx.Constraints.Max}
}

// end of functions/begining of main=================
func main() {

	// main cycle for goroutine(thread) to create and draw the elements
	go func() {
		w := new(app.Window)                            // creating window
		w.Option(app.Title("Words in lyrics"))          // window`s title
		w.Option(app.Size(unit.Dp(1080), unit.Dp(640))) // size

		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	// app.Main has to be called in the main thread of the program
	app.Main()
}

func loop(w *app.Window) error {
	th := material.NewTheme() // creating theme
	var ops op.Ops

	// creating widgets.....
	var albumEditor widget.Editor
	albumEditor.SingleLine = true // one string field

	var queryEditor widget.Editor
	queryEditor.SingleLine = true

	var searchBtn widget.Clickable // search button

	var resultsList widget.List // listbox for results
	resultsList.Axis = layout.Vertical

	// Тестовые данные, чтобы увидеть, как выглядит список
	listItems := []string{
		"Здесь будут отображаться",
		"результаты вашего поиска...",
	}

	// --- Главный цикл обработки событий ---
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err

		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			layoutBackground(gtx) // drawing the background in main cycle

			// Обработка нажатия на кнопку
			if searchBtn.Clicked(gtx) {
				// Временно: просто добавляем введенные данные в список при нажатии
				listItems = []string{
					"Вы ищете альбом/песню: " + albumEditor.Text(),
					"Слово или фраза: " + queryEditor.Text(),
				}
			}

			// drawing the interface.....
			layout.Flex{Axis: layout.Vertical}.Layout(gtx,

				// 1. Поле для названия альбома
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return material.Editor(th, &albumEditor, "Введите название альбома/песни").Layout(gtx)
					})
				}),

				// 2. Поле для искомого слова
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return material.Editor(th, &queryEditor, "Введите слово или фразу").Layout(gtx)
					})
				}),

				// 3. Кнопка "Поиск"
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(th, &searchBtn, "Поиск")
						return btn.Layout(gtx)
					})
				}),

				// 4. Listbox (занимает всё оставшееся место благодаря layout.Flexed(1))
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						// Отрисовываем элементы списка
						return material.List(th, &resultsList).Layout(gtx, len(listItems), func(gtx layout.Context, index int) layout.Dimensions {
							return layout.Inset{Bottom: unit.Dp(5)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return material.Body1(th, listItems[index]).Layout(gtx)
							})
						})
					})
				}),
			)

			e.Frame(gtx.Ops)
		}
	}
}
