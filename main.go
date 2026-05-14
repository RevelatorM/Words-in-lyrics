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

	"os/exec"
	"runtime"

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
	Lyrics string `json:"lyrics"` // `json:"lyrics"` json: is a tag and "lyrics" is a key so when the answear from API is read GO will look for this keyword
}

// openBrowser opens the provided link in the default browser
func openBrowser(link string) {
	var err error

	// Output the link to the terminal for verification (helps to debug)
	log.Printf("Trying to open link: %s\n", link)

	// Depending on the OS, use different methods to open the link
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", link).Start()
	case "windows":
		// Method 1: Absolute path to cmd.exe (bypasses terminal %PATH% issues)
		err = exec.Command("C:\\Windows\\System32\\cmd.exe", "/c", "start", "", link).Start()

		if err != nil {
			log.Printf("cmd.exe failed, trying explorer.exe...")
			// Method 2: Absolute path to explorer.exe
			err = exec.Command("C:\\Windows\\explorer.exe", link).Start()
		}

		if err != nil {
			log.Printf("explorer.exe failed, trying rundll32.exe...")
			// Method 3: Absolute path to rundll32.exe
			err = exec.Command("C:\\Windows\\System32\\rundll32.exe", "url.dll,FileProtocolHandler", link).Start()
		}

	case "darwin": // macOS
		err = exec.Command("open", link).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		log.Printf("Complete failure. Couldn't open the browser: %v", err)
	} else {
		log.Printf("The command to open the browser was successfully sent to the system!")
	}
}

func fetchLyrics(artist, song string) (string, error) { // fetchLyrics gets artist and song arguments, error
	apiURL := fmt.Sprintf("https://api.lyrics.ovh/v1/%s/%s", url.PathEscape(artist), url.PathEscape(song)) //fmt.Sprintf takes a template and puts it changin %s to it
	//url.PathEscape converts space signs into symbols like Linkin Park becomes Linkin%20Park for a link
	resp, err := http.Get(apiURL) // http.Get makes a request to the internet, result is being saved into the resp variable
	if err != nil {               // if err is not empty then program stops
		return "", err // returning empty string
	}
	defer resp.Body.Close() // defer is a destructor that closes connection

	if resp.StatusCode != http.StatusOK { // http.StatusOK (code 200) means everything is fine, if http.StatusOK is not 200 we write song was not found
		return "", fmt.Errorf("song was not found") // fmt.Errorf returning the message
	}

	body, _ := io.ReadAll(resp.Body) // _ means we are ignoring possible mistakes reading the whole body by exporting ReadAll from io
	var result LyricsResponse        // creating variable result from our structure LyricsResponse (LyricsResponse is our data type)
	json.Unmarshal(body, &result)

	return result.Lyrics, nil
}

func searchInLyrics(lyrics, query string) []string { // searchInLyrics looks for words or phrases, query is a word
	if query == "" { // if query is empty we ask for a word
		return []string{"Write a word in the second field to search"}
	}

	lines := strings.Split(lyrics, "\n")
	var matches []string                                    // empty list with matches
	lowerQuery := strings.ToLower(strings.TrimSpace(query)) // strings.TrimSpace cuts possible space before and after word
	// strings.ToLower makes all the letters small
	for i, line := range lines { // i is a string's number, line is a text
		if strings.Contains(strings.ToLower(line), lowerQuery) { // strings.Contains checks for a match, lowerQuery converts into a lower case
			matches = append(matches, fmt.Sprintf("Line %d: %s", i+1, strings.TrimSpace(line))) // append is a new note that will be added to matches if there is one
		}
	}

	if len(matches) == 0 { // if len = 0, means there is no matches
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
	var youtubeBtn widget.Clickable
	var youtubeURL string // Keeps the link
	var showYoutube bool  // A flag that shows if button's drawing is needed
	youtubeBtnColor := color.NRGBA{R: 205, G: 32, B: 31, A: 255}
	//=================================
	// Placeholders
	listItems := []string{
		"Write artist's name and song's name like here (Queen - Don't Stop Me Now)",
		"Write a word or a phrase in the second field",
	}

	// Main cycle
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err

		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			layoutBackground(gtx) // drawing the background in main cycle
			if youtubeBtn.Clicked(gtx) {
				go openBrowser(youtubeURL)
			}
			// Button click
			if searchBtn.Clicked(gtx) {
				// Showing loading message
				listItems = []string{"Loading please wait..."}

				// Reading the text
				targetText := albumEditor.Text()
				queryText := queryEditor.Text()

				// Starting search in the other thread
				go func(target, query string) {
					parts := strings.SplitN(target, "-", 2)
					if len(parts) != 2 {
						listItems = []string{"Error! please write: Artist - Song"}
						w.Invalidate()
						return
					}

					artist := strings.TrimSpace(parts[0])
					song := strings.TrimSpace(parts[1])

					lyrics, err := fetchLyrics(artist, song)
					if err != nil {
						listItems = []string{"Error! the song was not found"}
						showYoutube = false // No button if there is an error
					} else {
						listItems = searchInLyrics(lyrics, query)

						// Generating YouTube link
						// url.QueryEscape codes the spaces "+" and special symbols into the safe format
						searchQuery := url.QueryEscape(artist + " " + song)
						youtubeURL = "https://www.youtube.com/results?search_query=" + searchQuery
						showYoutube = true // drawing the button
					}

					w.Invalidate()
				}(targetText, queryText)
			}

			// drawing the interface.....
			layout.Flex{Axis: layout.Vertical}.Layout(gtx,

				// 1. Поле для названия альбома
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return material.Editor(th, &albumEditor, "Name: Artist - Song").Layout(gtx)
					})
				}),

				// 2. Поле для искомого слова
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return material.Editor(th, &queryEditor, "Write the word or a phrase").Layout(gtx)
					})
				}),

				// 3. Search button
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(th, &searchBtn, "Search")
						btn.Background = searchBtnColor // applying colour to the button
						return btn.Layout(gtx)
					})
				}),
				// YouTube button
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					// If the song was not found
					if !showYoutube {
						return layout.Dimensions{}
					}

					return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(th, &youtubeBtn, "YouTube")
						btn.Background = youtubeBtnColor // red colour
						return btn.Layout(gtx)
					})
				}),

				// 4. Listbox (occupies all the space because of layout.Flexed(1))
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						// Drawing list's elements
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
