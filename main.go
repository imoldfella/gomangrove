package main

import cp "github.com/otiai10/copy"

// home
// index.html
// sitemap.xml
// robot.txt

// upload subjects as a zip file?
// schools are just

// we need a school theme

// themes stored as json?
// theme can have a dictionary of subjects
//

// index.html is the school mangrove home page;

// then schools underneath that
// /school/bouncy-red-robin/index.html

// when should we use datum ids? should all the id's map the owning school?
// "there is a new version of this page"; query hash->id
// school should get direct edit of page or brain fry.
// when they click edit we can decide if they need to fork or not and just do it.

// grades can be hashed too, that way forking a school is just one page initially
// (changing the welcome message).

// /blob/hash.html#
// /blob/partialHash/hash.html
// /

// to build from we

// the grade pages will link to subject pages:

// there maybe more than one syllabus for that subject, if there is
// then the subject will link to the subject page, otherwise if there
// is only one syllabus, it will just link to that.

// we can have a separately generated (fake for now) subject
// content/syllabusHash/index.{n}.md

// /subject/syllabusid.html

// /syllabus/id/index.html   #use id here because it can link from multiple grades
//   content will be a link for each lesson

// /syllabus/id/{lesson}-{page}.html

// maybe we should do it in two phases
// 1. write a fictitious school from the spreadsheet
// 2. build the school from the directories.

//sylabus/{anyid}/index.md

// we want to hash each syllabus so that it can stay fixed and be shared.

const (
	ex = `
<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>http://www.example.com/foo.html</loc>
    <lastmod>2018-06-04</lastmod>
  </url>
</urlset>
`
)

const (
	dist = "/Users/jimhurd/yakdb/mangrove-stage/dist"
)

const target = "/Users/jimhurd/yaktemp/dist"
const csvr = "./r.csv"

func main() {
	schoolStarter(csvr, target+"/starter")
	schoolBuilder(target+"/starter", target+"/dist")
	cp.Copy(target+"/theme", target+"/dist/theme")
}
