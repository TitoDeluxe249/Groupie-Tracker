package main

import (
	"C"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)
import (
	"net/url"
	"strconv"
	"strings"
	"time"
)

var client *http.Client

type Relations struct {
	Relations []DatesLocation `json:"index"`
}
type DatesLocation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}
type Concert struct {
	Location string
	Dates    string
}

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	RelationsUrl string   `json:"relations"`
	PastConcert  []Concert
	FuturConcert []Concert
}

func downloadImage(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	file, err := ioutil.TempFile("", "image")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		panic(err)
	}

	return file.Name()
}

func getAllArtists() []Artist {
	var artists []Artist
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&artists)
	if err != nil {
		return nil
	}

	return artists
}

func getRelations(artists []Artist) []Artist {
	var relations Relations
	var newArtists []Artist = artists

	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&relations)
	if err != nil {
		panic(err)
	}

	for _, element := range relations.Relations {
		for city, dates := range element.DatesLocations {
			for _, dateString := range dates {
				date, err := time.Parse("02-01-2006", dateString)
				if err != nil {
					panic(err)
				}
				if date.Before(time.Now()) {
					newArtists[element.ID-1].PastConcert = append(newArtists[element.ID-1].PastConcert, Concert{city, dateString})
				} else {
					newArtists[element.ID-1].FuturConcert = append(newArtists[element.ID-1].PastConcert, Concert{city, dateString})
				}
			}
		}
	}

	return newArtists
}

func main() {
	var artists []Artist = getAllArtists()
	var listArtistName []string
	var listConcertPast []string
	var listConcertFutur []string

	getRelations(artists)

	a := app.New()

	a.Settings().SetTheme(theme.DarkTheme())

	w := a.NewWindow("GROUPIE TRACKER")

	w.Resize(fyne.NewSize(1000, 700))

	mainMenu := fyne.NewMainMenu(
		fyne.NewMenu("Fichier",
			fyne.NewMenuItem("Exit", func() {
				a.Quit()
			}),
		),
		fyne.NewMenu("Support",

			fyne.NewMenuItem("A propos de", func() {
				u, _ := url.Parse("https://meilleurs-albums.com/principaux-concerts-en-2023/")
				_ = a.OpenURL(u)
			}),
			fyne.NewMenuItem("Documentation", func() {
				u, _ := url.Parse("https://www.nrj.be/article/23085/quels-artistes-ont-ete-les-plus-ecoutes-durant-l-annee-2022")
				_ = a.OpenURL(u)
			}),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Sponsor", func() {
				u, _ := url.Parse("https://nike.com/")
				_ = a.OpenURL(u)
			}),
		))

	w.SetMainMenu(mainMenu)

	artist := artists[0]
	aTitle := widget.NewLabelWithStyle("Info Artist :", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	aTitle.Move(fyne.NewPos(520, 0))

	aName := widget.NewLabel(" ")

	aMembers := widget.NewLabelWithStyle(" ", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})
	aMembers.Move(fyne.NewPos(520, 250))

	aImage := canvas.NewImageFromFile("")
	aImage.Resize(fyne.NewSize(150, 190))
	aImage.Move(fyne.NewPos(520, 50))

	aCreationDate := widget.NewLabel(" ")
	aFirstAlbum := widget.NewLabel(" ")
	aLabelPastConcert := widget.NewLabelWithStyle("Concerts Passés :", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	aLabelPastConcert.Move(fyne.NewPos(520, 450))
	aLabelFuturConcert := widget.NewLabelWithStyle("Concerts Futurs :", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	aLabelFuturConcert.Move(fyne.NewPos(750, 450))

	aPastConcerts := widget.NewList(
		func() int { return 1 },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(lii widget.ListItemID, co fyne.CanvasObject) {
		},
	)
	aPastConcerts.Resize(fyne.NewSize(250, 200))
	aPastConcerts.Move(fyne.NewPos(520, 500))

	aFuturConcerts := widget.NewList(
		func() int { return 1 },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(lii widget.ListItemID, co fyne.CanvasObject) {
		},
	)
	aFuturConcerts.Resize(fyne.NewSize(250, 200))
	aFuturConcerts.Move(fyne.NewPos(750, 500))

	for _, artist := range artists {
		listArtistName = append(listArtistName, artist.Name)
	}

	list := widget.NewList(
		func() int { return len(listArtistName) },
		func() fyne.CanvasObject { return widget.NewLabel("Liste des artistes") },
		func(lii widget.ListItemID, co fyne.CanvasObject) {
			co.(*widget.Label).SetText(listArtistName[lii])
		},
	)
	list.Resize(fyne.NewSize(280, 500))
	list.Move(fyne.NewPos(220, 0))

	searchEntry := widget.NewEntry()
	searchButton := widget.NewButton("Rechercher", func() {
		// Nouvelle liste pour les résultats de la recherche
		filteredList := []string{}

		// Parcourir la liste d'origine et ajouter les éléments correspondants à la nouvelle liste
		for _, item := range listArtistName {
			if strings.Contains(strings.ToLower(item), strings.ToLower(searchEntry.Text)) {
				filteredList = append(filteredList, item)
			}
		}

		// Mettre à jour la liste avec les résultats de la recherche
		list.Length = func() int {
			return len(filteredList)
		}
		list.CreateItem = func() fyne.CanvasObject {
			return widget.NewLabel("")
		}
		list.UpdateItem = func(index int, item fyne.CanvasObject) {
			item.(*widget.Label).SetText(filteredList[index])
		}
		list.Refresh()
	})
	clearButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		searchEntry.SetText("")

		// Mettre à jour la liste avec les résultats de la recherche
		list.Length = func() int {
			return len(listArtistName)
		}
		list.CreateItem = func() fyne.CanvasObject {
			return widget.NewLabel("")
		}
		list.UpdateItem = func(index int, item fyne.CanvasObject) {
			item.(*widget.Label).SetText(listArtistName[index])
		}
		list.Refresh()
	})

	searchBar := container.NewVBox(
		searchEntry,
		searchButton,
		clearButton,
	)
	searchBar.Resize(fyne.NewSize(200, 100))

	content := container.NewWithoutLayout(
		searchBar,
		list,
	)

	separator := widget.NewSeparator()
	separator.Move(fyne.NewPos(500, 0))

	nameContainer := container.NewVBox(
		aName,
		aFirstAlbum,
	)
	nameContainer.Move(fyne.NewPos(750, 50))

	infoArtist := container.NewWithoutLayout(
		aTitle,
		aImage,
		nameContainer,
		aMembers,
		aPastConcerts,
		aFuturConcerts,
		aLabelPastConcert,
		aLabelFuturConcert,
	)
	infoArtist.Resize(fyne.NewSize(480, 800))
	// infoArtist.Move(fyne.NewPos(520, 0))
	infoArtist.Hide()

	list.OnSelected = func(id widget.ListItemID) {

		artist = artists[id]
		aName.Text = artist.Name
		aName.Refresh()
		aMembersList := "Membre : \n"
		for _, member := range artist.Members {
			aMembersList += "- " + member + "\n"
		}
		aMembers.Text = aMembersList
		aMembers.Refresh()
		aCreationDate.Text = strconv.Itoa(artist.CreationDate)
		aImagePath := downloadImage(artist.Image)
		aImage.File = aImagePath
		aImage.Refresh()

		aFirstAlbum.Text = "Date première album : " + artist.FirstAlbum
		aFirstAlbum.Refresh()

		listConcertPast = nil

		for _, concert := range artists[id].PastConcert {
			listConcertPast = append(listConcertPast, concert.Location+" : "+concert.Dates)

		}
		aPastConcerts.Length = func() int {
			return len(listConcertPast)
		}
		aPastConcerts.CreateItem = func() fyne.CanvasObject {
			return widget.NewLabel("")
		}
		aPastConcerts.UpdateItem = func(index int, item fyne.CanvasObject) {
			item.(*widget.Label).SetText(listConcertPast[index])
		}
		aPastConcerts.Refresh()

		listConcertFutur = nil

		for _, concert := range artists[id].PastConcert {
			listConcertFutur = append(listConcertFutur, concert.Location+" : "+concert.Dates)

		}
		aPastConcerts.Length = func() int {
			return len(listConcertFutur)
		}
		aPastConcerts.CreateItem = func() fyne.CanvasObject {
			return widget.NewLabel("")
		}
		aPastConcerts.UpdateItem = func(index int, item fyne.CanvasObject) {
			item.(*widget.Label).SetText(listConcertFutur[index])
		}
		aPastConcerts.Refresh()

		if infoArtist.Hidden {
			infoArtist.Show()
		}

	}

	w.SetContent(
		container.NewWithoutLayout(
			content,
			separator,
			infoArtist),
	)

	w.ShowAndRun()
}
