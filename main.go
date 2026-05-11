package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// functions=========================================
type LyricsResponse struct { // LyricsResponse describes the JSON structure from API
	Lyrics string `json:"lyrics"`
}

func fetchLyrics(artist, song string) (string, error) { // fetchLyrics делает запрос в интернет и получает текст песни
	apiURL := fmt.Sprintf("https://api.lyrics.ovh/v1/%s/%s", url.PathEscape(artist), url.PathEscape(song))

	resp, err := http.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("song was not found")
	}

	body, _ := io.ReadAll(resp.Body)
	var result LyricsResponse
	json.Unmarshal(body, &result)

	return result.Lyrics, nil
}

func searchInLyrics(lyrics, query string) []string { // searchInLyrics looks for words or phrases
	if query == "" {
		return []string{"Введите слово для поиска во второе поле."}
	}

	lines := strings.Split(lyrics, "\n")
	var matches []string
	lowerQuery := strings.ToLower(strings.TrimSpace(query))

	for i, line := range lines {
		if strings.Contains(strings.ToLower(line), lowerQuery) {
			matches = append(matches, fmt.Sprintf("Строка %d: %s", i+1, strings.TrimSpace(line))) // adding the number of a string
		}
	}

	if len(matches) == 0 {
		return []string{"Not found"}
	}

	return matches
}

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
	//=================================
	// creating widgets.....
	var albumEditor widget.Editor
	albumEditor.SingleLine = true       // one string field
	albumEditor.Alignment = text.Middle // making placeholder in the widget centered
	//=================================
	var queryEditor widget.Editor
	queryEditor.SingleLine = true       // one string field
	queryEditor.Alignment = text.Middle // making placeholder in the widget centered
	//=================================
	var searchBtn widget.Clickable // search button
	searchBtnColor := color.NRGBA{R: 95, G: 0, B: 87, A: 180}
	//=================================
	var resultsList widget.List // listbox for results
	resultsList.Axis = layout.Vertical
	//=================================
	// Тестовые данные, чтобы увидеть, как выглядит список
	listItems := []string{
		"Write artist's name and song's name like here (Queen - Don't Stop Me Now)",
		"Write a word or a phrase in the second field",
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
						btn := material.Button(th, &searchBtn, "Search")
						btn.Background = searchBtnColor // applying colour to the button
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
