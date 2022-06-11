package main

// https://pace.dev/blog/2020/03/02/dynamically-generate-social-images-in-golang-by-mat-ryer.html
import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"unicode"

	"github.com/noelyahan/impexp"
	"github.com/noelyahan/mergi"

	"github.com/joho/godotenv"
	"github.com/kosa3/pexels-go"
)

var subjects = `Math - Algebra II
Science - Chemistry
Geography
History
Patriotism and Citizenship
Practical Arts and Library Skills
Music and Visual Arts
United States History
Health and Safety
Social Studies - Modern World History
Math - Geometry
Math - Pre-Calculus
Math - Statistics & Probability
Science - Physics 
Reading
Visual Arts and Music
Social Studies - Global Studies
Arithmetic
Spelling
Mathematics
Health Education 
English Language
Physical Education
Writing
Social Studies - United States History
Math - Algebra 1
Health Education
Science 
Social Science - Participation in Government and Economics
Science - Biology
Science
History and Geography (World)
Science - Earth and Space
Arithmetic`

var subject = map[string][]*pexels.Photo{}

func pexel(prompt string) []*pexels.Photo {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cli := pexels.NewClient(os.Getenv("PEXEL"))
	ctx := context.Background()
	ps, err := cli.PhotoService.Search(ctx, &pexels.PhotoParams{
		Query: prompt,
		//Orientation: "",
		Size: "small",
		// Color:       "",
		// Locale:      "",
		Page:    1,
		PerPage: 0,
	})
	if err != nil {
		log.Fatal(err)
	}
	_ = &pexels.Photo{
		ID:              0,
		Width:           0,
		Height:          0,
		URL:             "",
		Photographer:    "",
		PhotographerURL: "",
		PhotographerID:  0,
		AvgColor:        "",
		Liked:           false,
		Src:             pexels.Source{},
	}
	return ps.Photos
}

func Test_pexel(t *testing.T) {
	o := strings.Split(subjects, "\n")
	for _, j := range o {
		subject[j] = pexel(j)
	}

	m, _ := json.Marshal(subject)
	os.WriteFile("photo.json", m, 0666)
}

func makeName(x string) string {
	s := ""
	for _, c := range x {
		if unicode.IsLetter(c) {
			s = s + string(c)
		}
	}
	return s
}

func Test_fetch(t *testing.T) {
	b, _ := os.ReadFile("photo.json")
	json.Unmarshal(b, &subject)

	download := func(url string, to int) error {
		path := fmt.Sprintf("./photo/%d.jpeg", to)
		r, e := http.Get(url)
		if e != nil {
			return e
		}
		defer r.Body.Close()
		//Create a empty file
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()

		//Write the bytes to the fiel
		_, err = io.Copy(file, r.Body)
		return err
	}

	for _, v := range subject {
		for _, o := range v {
			download(o.Src.Tiny, o.ID)
		}
	}
}

func Test_thumb(t *testing.T) {
	b, _ := os.ReadFile("photo.json")
	json.Unmarshal(b, &subject)

	resize := func(in, out string) error {
		// Get the image content by passing image path url or file path
		img, err := mergi.Import(impexp.NewFileImporter(in))
		if err != nil {
			log.Fatalf("failed to open: %s", err)
		}

		// crop to a square
		// Now let's use the mergi's crop API

		// Set where to start the crop point
		r := img.Bounds()
		if r.Dx() > r.Dy() {
			startX := (r.Dx() - r.Dy()) / 2
			cropStartPoint := image.Pt(startX, 0)
			cropSize := image.Pt(r.Dy(), r.Dy())
			img, err = mergi.Crop(img, cropStartPoint, cropSize)

		} else {
			startY := (r.Dy() - r.Dx()) / 2
			cropStartPoint := image.Pt(0, startY)
			cropSize := image.Pt(r.Dx(), r.Dx())
			img, err = mergi.Crop(img, cropStartPoint, cropSize)
		}
		if err != nil {
			log.Fatalf("Mergi crop fails due to [%v]", err)
		}
		// Lets resize double of the given image's size
		width := uint(180)
		height := uint(180)

		img, err = mergi.Resize(img, width, height)
		if err != nil {
			log.Fatalf("failed to resize: %s", err)
		}

		// Let's save the image
		err = mergi.Export(impexp.NewFileExporter(img, out))
		if err != nil {
			log.Fatalf("failed to save: %s", err)
		}

		return nil
	}
	for _, v := range subject {
		for _, o := range v {
			in := fmt.Sprintf("./photo/%d.jpeg", o.ID)
			out := fmt.Sprintf("./thumb/%d.jpeg", o.ID)
			resize(in, out)
		}
	}

}
