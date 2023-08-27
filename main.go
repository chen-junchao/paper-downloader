package main

import (
	"io"
	"log"
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
	hintLabel     *widget.Label
	savePathEntry *widget.Entry
}

var myApp App

func main() {
	a := app.New()
	w := a.NewWindow("Paper download")

	hintLabel, savePathEntry, linksEntry, downloadBtn, savePathBtn := myApp.makeUI(&w)

	content := container.NewVBox(
		savePathEntry,
		savePathBtn,
		hintLabel,
		linksEntry,
		downloadBtn,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(400, 300))
	w.ShowAndRun()
}

func (app *App) makeUI(w *fyne.Window) (*widget.Label, *widget.Entry, *widget.Entry, *widget.Button, *widget.Button) {
	hintLabel := widget.NewLabel("Paste Link...")
	savePathEntry := widget.NewEntry()
	savePathEntry.SetPlaceHolder("Please select save path")
	linksEntry := widget.NewMultiLineEntry()

	downloadBtn := widget.NewButton("Download", func() {
		links := linksEntry.Text
		// Split links by space, newline, or comma
		re := regexp.MustCompile("[,|\\s]+")
		linkList := re.Split(links, -1)

		// Download and save each paper
		for _, link := range linkList {
			err := download(link, app.savePathEntry.Text)
			if err != nil {
				log.Fatal(err)
			}
		}
		dialog.ShowInformation("Info", "Download Finish!", *w)
	})

	savePathBtn := widget.NewButton("Select", func() {
		dialog.ShowFolderOpen(func(list fyne.ListableURI, err error) {
			if err == nil && list != nil {
				if strings.Contains(list.String(), "file:") {
					myApp.savePathEntry.SetText(list.String()[7:])
				} else {
					myApp.savePathEntry.SetText(list.String())
				}
			}
		}, *w)
	})

	app.hintLabel = hintLabel
	app.savePathEntry = savePathEntry

	return hintLabel, savePathEntry, linksEntry, downloadBtn, savePathBtn
}

func transPdfLink(link string) string {
	// arxiv
	if (link != "") && (strings.Contains(link, "arxiv.org")) {
		link = strings.Replace(link, "abs", "pdf", -1)
	}

	return strings.TrimSpace(link)
}

func getPdfName(link string) string {
	if link != "" {
		linkList := strings.Split(link, "/")
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
