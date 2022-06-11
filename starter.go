package main

// here we want to turn r's spreadsheet into a set of syllabus files
// then we will add these syllabus files to an empty school to create the
// first school.

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/kosa3/pexels-go"
	faker "github.com/manveru/faker"
)

type Slide struct {
	content string
}

// func buildSubjectTitle(topic string) string {
// 	return fmt.Sprintf(`# %s
// 	  **Education is the most powerful weapon which you can use to change the world.** \u2014  Nelson Mandela
// 		`, topic)
// }

//
func buildTitle(n int, topic string) string {
	return fmt.Sprintf(`
# %d. %s
**Education is the most powerful weapon which you can use to change the world.** —  Nelson Mandela
		`, n, topic)
}
func buildSlide(lessonNumber int, page int) string {
	fake, err := faker.New("en")
	if err != nil {
		panic(err)
	}

	slideTitle := fmt.Sprintf("\n##  %d.%d %s\n", lessonNumber, page, fake.Sentence(5, false))

	bullets := []string{slideTitle}
	for i := 0; i < 5; i++ {
		s := fake.Sentence(5, false)
		var o = "\n * " + s
		bullets = append(bullets, o)
	}
	return mergeStrings(bullets)
}

func mergeStrings(s []string) string {
	var b bytes.Buffer

	for _, o := range s {
		b.WriteString(o)
	}
	return b.String()
}

func filefy(s string) string {
	return strings.ReplaceAll(s, " ", "-")
}

type FrontMatter struct {
	Title    string `yaml:"title"`
	Subtitle string `yaml:"subtitle"`
	MinGrade int    `yaml:"minGrade"`
	MaxGrade int    `yaml:"maxGrade"`
}

func writeFrontMatter(w io.Writer, f *FrontMatter) {
	fmt.Fprintf(w,
		`---
title: %s
subtitle: %s
minGrade: %d
maxGrade: %d`, f.Title, f.Subtitle, f.MinGrade, f.MaxGrade)
}

// eventually this should zip, but for now no images?
var names = map[string]bool{}

func buildSyllabus(sc *SchoolJson, root string, grade string, s []string) {
	if len(s) == 0 {
		return
	}
	names[s[0]] = true
	log.Printf("%s,%s,%s", root, grade, s[0])
	var b bytes.Buffer

	// create a directory for each syllabus
	// it has name min-max-subject-title

	// we need to split the grade
	o := strings.Split(grade, "-")
	minGrade, _ := strconv.Atoi(o[0])
	maxGrade := minGrade
	if len(o) > 1 {
		maxGrade, _ = strconv.Atoi(o[1])
	}

	// the second line gives us the ability to create electives
	// but do we need that? it should probably be in pawpaw.
	// we should probably have departments.
	o = strings.Split(s[0], "\n")
	title := o[0]
	subtitle := ""
	if len(o) > 1 {
		subtitle = o[1]
	}

	// write the front matter
	writeFrontMatter(&b, &FrontMatter{
		Title:    title,
		Subtitle: subtitle,
		MinGrade: minGrade,
		MaxGrade: maxGrade,
	})

	// here the path doesn't matter, should we just increment?
	// we are going to import anyway, which will use the hash of the name

	// write a greek lesson in markdown
	// b.WriteString("\n---\n")
	// b.WriteString(buildSubjectTitle(title))

	for j := 1; j < len(s); j++ {
		b.WriteString("\n---")
		b.WriteString(buildTitle(j, s[j]))
		for i := 1; i < 4; i++ {
			b.WriteString("\n---")
			b.WriteString(buildSlide(j, i))
		}
	}

	// write the subject to CAS name
	content := b.Bytes()
	name := hash(content)
	path := fmt.Sprintf("%s/%s.md", root, name)
	os.WriteFile(path, b.Bytes(), 0666)

	paths := []string{}

	gradeName := func(x int) string {
		if x == 0 {
			return "Grade K"
		} else {
			return fmt.Sprintf("Grade %d", x)
		}
	}
	for g := minGrade; g <= maxGrade; g++ {
		// give the subject a path for each grade it's specfied for
		paths = append(paths)

		sc.Subject = append(sc.Subject, &SubjectLinkJson{
			Title: title,
			Sort:  title,
			Hash:  name,
			Image: "",
			Path:  gradeName(g),
			Pin:   false,
		})
	}
}

// the school builder needs the json, we need to build that here
// as well as the textbooks.
func schoolStarter(in string, out string) {

	os.RemoveAll(out)
	os.Mkdir(out, os.ModePerm)
	j := &SchoolJson{}
	// this is awkward, but add a folder for each grade so we
	// can control the sort key.
	for i := 0; i <= 12; i++ {
		name := fmt.Sprintf("Grade %d", i)
		sort := fmt.Sprintf("Grade %02d", i)
		if i == 0 {
			name = "Grade K"
		}
		j.Subject = append(j.Subject, &SubjectLinkJson{
			Title: name,
			Sort:  sort,
			Hash:  "",
			Image: "",
			Path:  "",
			Pin:   false,
		})
	}

	dat, err := os.Open(in)

	if err != nil {
		panic("no r.csv")
	}
	r := csv.NewReader(dat)

	grade := ""
	prev := ""
	topics := []string{}
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if record[1] == prev {
			topics = append(topics, record[2])
		} else {
			grade = record[0]

			buildSyllabus(j, out, grade, topics)
			prev = record[1]
			topics = []string{record[2]}
		}

		//fmt.Println(record)

	}
	buildSyllabus(j, out, grade, topics)
	j.Welcome = "Welcome to mangrove sample school"
	ph := Photos()
	v := []*pexels.Photo{}
	for _, o := range ph {
		v = append(v, o[0])
		log.Printf("lead image %d", o[0].ID)
	}

	for i := range j.Subject {
		ap, ok := ph[j.Subject[i].Title]
		if !ok {
			k := rand.Intn(len(v))
			log.Printf("random %d", v[k].ID)
			j.Subject[i].Image = fmt.Sprintf("%d.jpeg", v[k].ID)
			v[k] = v[len(v)-1]
			v = v[0 : len(v)-1]

		} else {
			j.Subject[i].Image = fmt.Sprintf("%d.jpeg", ap[0].ID)
		}
	}

	b, _ := json.Marshal(j)
	os.WriteFile(out+"/index.json", b, 0666)

	var b2 bytes.Buffer
	for k, _ := range names {
		b2.WriteString(k + "\n")
	}
	os.WriteFile(out+"subjects.txt", b2.Bytes(), 0666)
}
func Photos() map[string][]*pexels.Photo {
	var subject = map[string][]*pexels.Photo{}
	b, _ := os.ReadFile("photo.json")
	json.Unmarshal(b, &subject)
	return subject
}
