package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type App struct {
	output   *widget.Label
	savePath *widget.Entry
}

var myApp App

func main() {
	a := app.New()
	w := a.NewWindow("Paper download")

	output, savePath, linksEntry, downloadBtn := myApp.makeUI(&w)

	savePathBtn := widget.NewButton("Select", func() {
		dialog.ShowFolderOpen(func(list fyne.ListableURI, err error) {
			if err == nil && list != nil {
				if strings.Contains(list.String(), "file:") {
					myApp.savePath.SetText(list.String()[7:])
				} else {
					myApp.savePath.SetText(list.String())
				}
			}
		}, w)
	})

	content := container.NewVBox(
		savePath,
		savePathBtn,
		output,
		linksEntry,
		downloadBtn,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(400, 300))
	w.ShowAndRun()
}

func (app *App) makeUI(w *fyne.Window) (*widget.Label, *widget.Entry, *widget.Entry, *widget.Button) {
	output := widget.NewLabel("Paste Link...")
	savePath := widget.NewEntry()
	savePath.SetPlaceHolder("Please select save path")

	linksEntry := widget.NewMultiLineEntry()
	downloadBtn := widget.NewButton("Download", func() {
		links := linksEntry.Text
		// Split links by space, newline, or comma
		re := regexp.MustCompile("[,|\\s]+")
		linkList := re.Split(links, -1)

		// Download and save each paper
		for _, link := range linkList {
			err := download(link, app.savePath.Text)
			if err != nil {
				fmt.Printf("Download fail %s", link)
			}
		}
		dialog.ShowInformation("Info", "Download Finish!", *w)
	})

	app.output = output
	app.savePath = savePath

	return output, savePath, linksEntry, downloadBtn
}

func transPdfLink(link string) string {
	// arxiv, if link contain arxiv.org
	if (link != "") && (strings.Contains(link, "arxiv.org")) {
		link = strings.Replace(link, "abs", "pdf", -1)
	}

	return strings.TrimSpace(link)
}

func getPdfName(link string) string {
	// arxiv, if link contain arxiv.org
	if link != "" {
		// split link by slash
		linkList := strings.Split(link, "/")
		// get the last element
		pdfName := linkList[len(linkList)-1]
		return pdfName
	}

	return ""
}

func download(link string, savePath string) error {

	pdfName := getPdfName(link)

	// trans to pdf download link
	link = transPdfLink(link)

	resp, err := http.Get(link)
	if err != nil {
		return err
	}

	fileSavePath := savePath + "/" + pdfName + ".pdf"
	out, err := os.Create(fileSavePath)
	defer out.Close()
	if err != nil {
		return err
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
