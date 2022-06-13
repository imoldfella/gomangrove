package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	cp "github.com/otiai10/copy"
)

var builder = NewBuilder()

// NewSHA256 ...
func hash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])[0:16]
}

func writeBlob(dir string, data []byte, ext string) string {
	h := hash(data)
	os.WriteFile(dir+"/"+h+ext, data, 0666)
	return h
}
func writeHtml(dir string, data string) string {
	return writeBlob(dir, []byte(data), ".html")

}
func writeMd(dir string, data string) string {
	return writeBlob(dir, []byte(data), ".md")
}

func (sb *Subject) href(lesson, page int) string {
	if (lesson == 0 && page == 0) || page < 0 || page >= len(sb.Lesson[lesson].Slide) {
		return sb.Hash + ".html"
	}
	return fmt.Sprintf("%s.%d.%d.html", sb.Hash, lesson, page)
}

func (sb *Subject) Write(dir string) string {
	p := dir + "/blob"
	navbar := ""
	// this is hard to go backwards!
	// the link is in the page.
	// maybe we need to use one hash for the entire lesson.

	for z, o := range sb.Lesson {
		for i, c := range o.Slide {
			html := Md(c)
			page := builder.Slide(navbar, html, sb.href(z, i+1))
			os.WriteFile(p+"/"+sb.href(z, i), []byte(page), 0666)

		}
		// this should only link to the first slide.
		sb.Lesson[z].Link = sb.href(z, 0)
	}

	// we need to write a lesson index and return a link
	if false {
		b := builder.LessonSorter(sb)
		os.WriteFile(p+"/"+sb.Hash+".html", []byte(b), 0666)
	}

	return "blob/" + sb.Hash + ".html"
}

func Md(mdx string) string {

	//md := []byte("## markdown document")
	md := markdown.NormalizeNewlines([]byte(mdx))
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(extensions)
	unsafe := string(markdown.ToHTML(md, parser, nil))

	return unsafe
	//return string(bluemonday.UGCPolicy().SanitizeBytes(unsafe))
}

// I probably want to hash the textbook and then create the pages as
// hash.x.x? not necessarily a win because I'm not sharing with textbooks almost
// the same. but it does let be potentially extract dynamically.

func loadTextbook(path string) *Subject {
	input, _ := os.ReadFile(path)
	inputs := string(input)
	r := &Subject{
		Hash: hash(input),
	}

	rest, err := frontmatter.Parse(strings.NewReader(inputs), &r.FrontMatter)
	if err != nil {
		panic("bad frontmatter")
	}
	// we can split the rest on \n--\n to get slides.
	// the first slide is the table of contents, there might not be more.
	o := strings.Split(string(rest), "\n---\n")
	isTitle := func(c string) bool {
		return c[0:2] == "# "
	}
	for i, x := range o {
		// to find the start of the lesson find a slide with \n#{sp}
		if i == 0 || isTitle(x) {
			r.Lesson = append(r.Lesson, Lesson{})
		}
		bx := &r.Lesson[len(r.Lesson)-1]
		bx.Number = len(r.Lesson)
		bx.Slide = append(bx.Slide, x)
	}
	// r.Contents = r.Lesson[0]
	// r.Lesson = r.Lesson[1:]
	return r
}

// we need to turn each Folder into an html page.
// we need to add to site.xml (how should we do this incrementally?)
var next = 1

// we need to walk the children so we get back the link
func walkFolders(p string, f *Folder, depth int) {
	log.Printf("visit %s, %s", f.Link.Path, f.Link.Title)
	for _, v := range f.Folder {
		walkFolders(p, v, depth+1)
	}
	for _, o := range f.More {
		log.Printf("link %s", o.Title)
	}
	// a this point all the children a walked, we can set our link to
	// whatever file we generate.
	if len(f.Title) == 0 {
		f.Title = defaultTitle
	}
	path := fmt.Sprintf("%s/%d.html", p, next)
	f.Link.Link = fmt.Sprintf("%d.html", next)

	loader := ""
	next++
	if len(f.Welcome) > 0 {
		path = fmt.Sprintf("%s/%s.html", p, "index")
		loader = `<script>
		if ('serviceWorker' in navigator) {
		   window.addEventListener('load', () => {
			  navigator.serviceWorker.register('/sw.js');
		   });
		 }
		</script>`
	}

	// we need bread crumbs here, if this is not the root.
	crumbs := ""
	if len(f.Link.Title) > 0 {
		crumbs = "<div class='crumb'>" + f.Link.Title + "</div>"
	}

	// create a page from the pin and more
	b := builder.Folder(f, crumbs, loader)
	os.WriteFile(path, []byte(b), 0666)

}

// the main thing is that the root needs to go in /index.html
// everything else can be a blob?
// should everything else go in a subdirectory to support more schools?
func schoolBuilder(in string, out string) {
	cp.Copy(in, out)
	// walk through the syllabus directory and create
	// a parallel directory with markdown converted to html
	// also create the index html files for linking to syllabus
	// and sitemap.xml file.
	os.Mkdir(out+"/blob", os.ModePerm)
	b, _ := os.ReadFile(in + "/index.json")
	var sc SchoolJson
	json.Unmarshal(b, &sc)

	root := &Folder{
		Title:   defaultTitle,
		Welcome: sc.Welcome,
		Pin:     []*SubjectLinkJson{},
		More:    []*SubjectLinkJson{},
		Folder:  map[string]*Folder{},
		Link:    &SubjectLinkJson{},
	}

	// we can build all the subjects. Skip this step if hash=0
	for _, o := range sc.Subject {
		if !o.Folder {
			s := loadTextbook(in + "/" + o.Content + ".md")
			o.Link = s.Write(out)
		}
	}

	// build nested folders
	for _, o := range sc.Subject {
		if len(o.Title) == 0 {
			continue
		}

		ok := false
		at := root
		if len(o.Path) > 0 {
			a := strings.Split(o.Path, "|")
			for _, pc := range a {
				at, ok = at.Folder[pc]
				if !ok {
					panic(o.Path)
				}
			}
		}

		if o.Pin {
			at.Pin = append(at.Pin, o)
		} else {
			at.More = append(at.More, o)
		}
		// a folder will not have a hash
		if o.Folder {
			at.Folder[o.Title] = &Folder{
				Title:   o.Title,
				Welcome: "",
				Pin:     []*SubjectLinkJson{},
				More:    []*SubjectLinkJson{},
				Folder:  map[string]*Folder{},
				Link:    o,
			}
		}

	}

	walkFolders(out, root, 0)

	// as we create each page of the index we need to to sort the list
	// by pin, name.

}

/*
// we need to create all the pages to the root if necessary
			at := root
			at.More = append(at.More, o)
			at.Folder[o.Title] = &Folder{
				Title:   o.Title,
				Welcome: "",
				Pin:     []SubjectLinkJson{},
				More:    []SubjectLinkJson{},
				Folder:  map[string]*Folder{},
			}
		//
		for _, k := range a {
			if len(k) == 0 {
				continue
			}
			find, ok := at.Folder[k]
			if !ok {
				find = &Folder{
					Title:   "",
					Welcome: "",
					Pin:     []SubjectLinkJson{},
					More:    []SubjectLinkJson{},
					Folder:  map[string]*Folder{},
				}
				at.Folder[k] = find
				at.More = append(at.More,
					SubjectLinkJson{
						Title: k,
						Sort:  k,
						Hash:  "",
						Image: "",
						Path:  "",
						Pin:   false,
					})
			}
			at = find
		}


	files, err := os.ReadDir(p)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		ext := path.Ext(file.Name())
		if ext == ".md" {
			s := loadTextbook(p + "/" + file.Name())
			s.Write(out)
		}

	}
}*/
