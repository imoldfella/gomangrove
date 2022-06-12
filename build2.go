package main

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"text/template"

	"gopkg.in/yaml.v2"
)

var defaultTitle = "📚 mangrove"

type SchoolTheme struct {
	Image map[string]string
}

// sorting the grade strings is going to go wrong :(
// should we sort before creating the json then?
// folders are like subject links too.
// can we  capture that subject link here to use it as a link (instead of the
// subject hash)
type Folder struct {
	Title   string
	Welcome string
	// how do I turn a subject link into a href?
	Pin    []*SubjectLinkJson
	More   []*SubjectLinkJson
	Folder map[string]*Folder
	Link   *SubjectLinkJson
}

type Mangrove struct {
	Theme map[string]SchoolTheme
}

type SchoolJson struct {
	Welcome string
	Subject []*SubjectLinkJson
}

type SubjectLinkJson struct {
	Title string
	Sort  string
	//Hash  string
	Image   string
	Path    string
	Pin     bool
	Folder  bool
	Content string
	// we don't load this, we set it when build the page.
	Link string
}

// type SubjectLinkJson struct {
// 	Title string
// 	Sort  string
// 	Hash  string
// 	Image string
// 	Path  string
// 	Pin   bool
// 	Link  string
// }

type FrontMatter struct {
	Title    string `yaml:"title"`
	Subtitle string `yaml:"subtitle"`
	MinGrade int    `yaml:"minGrade"`
	MaxGrade int    `yaml:"maxGrade"`
}
type Subject struct {
	FrontMatter
	// it would be nice to be able to hash slides or lessons
	// but then we couldn't statically link within the subject
	// tradeoffs
	Hash     string
	Contents Lesson
	Lesson   []Lesson
}

// if we need it the lesson title is on the first slide #
type Lesson struct {
	Number int
	Slide  []string
	Link   string
}

type Builder struct {
	page,
	lessonList,
	pinList *template.Template

	theme SchoolTheme
}

func NewBuilder() *Builder {
	z, e := os.ReadFile("./css.yaml")
	if e != nil {
		panic(e)
	}
	var m = map[string]string{}
	e = yaml.Unmarshal(z, &m)
	if e != nil {
		panic(e)
	}

	make := func(name string, code string) *template.Template {
		t, e := template.New("navbar").Parse(code)
		if e != nil {
			panic(e)
		}
		return t
	}

	return &Builder{

		page:       make("page", m["page"]),
		pinList:    make("pinList", m["pinList"]),
		lessonList: make("lessonList", m["lessonList"]),
	}
}
func (d *Builder) Page(navbar, content string) string {
	var b bytes.Buffer
	d.page.Execute(&b, &PageInfo{
		Content: content,
	})
	return b.String()
}

func (d *Builder) Slide(navbar, content string, next string) string {
	link := next
	// <div class='button'>Next</div>
	return d.Page("📚 mangrove",
		fmt.Sprintf(`<a class="content" href="%s">%s</a>`, link, content))
}

func (d *Builder) Folder(f *Folder) string {
	var b bytes.Buffer

	if len(f.Pin) > 0 {
		sort.Slice(f.Pin, func(i, j int) bool {
			return f.Pin[i].Sort < f.Pin[j].Sort
		})
		d.pinList.Execute(&b, &f.Pin)
	}
	sort.Slice(f.More, func(i, j int) bool {
		return f.More[i].Sort < f.More[j].Sort
	})
	d.pinList.Execute(&b, &f.More)
	return d.Page(f.Title,
		b.String(),
	)
}

func (d *Builder) LessonSorter(pg *Subject) string {
	var b bytes.Buffer
	d.lessonList.Execute(&b, &pg.Lesson)
	return d.Page(pg.Title,
		b.String(),
	)
}

type PageInfo struct {
	Title,
	Content string
}
