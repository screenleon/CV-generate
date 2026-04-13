package generators

import (
	"archive/zip"
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/screenleon/cv-generate/models"
)

// GenerateWord generates a DOCX CV using the specified template style.
// Returns the DOCX bytes or an error.
func GenerateWord(data *models.CVData) ([]byte, error) {
	switch strings.ToLower(data.Template) {
	case "japan":
		return generateJapanDocx(data)
	default:
		return generateSimpleDocx(data)
	}
}

// -------------------------------------------------------------------
// DOCX helpers
// -------------------------------------------------------------------

// docxFile represents a single file inside the DOCX ZIP archive.
type docxFile struct {
	name    string
	content string
}

// buildDocx assembles the DOCX ZIP archive from the provided files.
func buildDocx(files []docxFile) ([]byte, error) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for _, f := range files {
		w, err := zw.Create(f.name)
		if err != nil {
			return nil, err
		}
		if _, err := w.Write([]byte(f.content)); err != nil {
			return nil, err
		}
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// xmlEscape replaces characters that are invalid inside XML text content.
func xmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}

// -------------------------------------------------------------------
// Shared DOCX boilerplate files
// -------------------------------------------------------------------

const contentTypesXML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/word/document.xml"
    ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
  <Override PartName="/word/styles.xml"
    ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.styles+xml"/>
</Types>`

const relsXML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1"
    Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument"
    Target="word/document.xml"/>
</Relationships>`

const wordRelsXML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1"
    Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles"
    Target="styles.xml"/>
</Relationships>`

const simpleStylesXML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:styles xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:style w:type="paragraph" w:styleId="Normal" w:default="1">
    <w:name w:val="Normal"/>
    <w:rPr><w:sz w:val="20"/></w:rPr>
  </w:style>
  <w:style w:type="paragraph" w:styleId="Heading1">
    <w:name w:val="heading 1"/>
    <w:rPr><w:b/><w:sz w:val="36"/><w:color w:val="1E78D4"/></w:rPr>
  </w:style>
  <w:style w:type="paragraph" w:styleId="Heading2">
    <w:name w:val="heading 2"/>
    <w:rPr><w:b/><w:sz w:val="24"/><w:color w:val="1E78D4"/></w:rPr>
  </w:style>
</w:styles>`

// -------------------------------------------------------------------
// Simple template
// -------------------------------------------------------------------

const simpleDocxTmpl = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
<w:body>
  <!-- Name -->
  <w:p>
    <w:pPr><w:jc w:val="center"/><w:pStyle w:val="Heading1"/></w:pPr>
    <w:r><w:rPr><w:b/><w:sz w:val="48"/></w:rPr><w:t>{{.Name}}</w:t></w:r>
  </w:p>
  <!-- Contact -->
  <w:p>
    <w:pPr><w:jc w:val="center"/></w:pPr>
    <w:r><w:rPr><w:color w:val="505050"/><w:sz w:val="18"/></w:rPr>
      <w:t>{{.ContactLine}}</w:t>
    </w:r>
  </w:p>
  <w:p><w:r><w:t></w:t></w:r></w:p>
{{if .Summary}}
  <!-- Summary -->
  <w:p><w:pPr><w:pStyle w:val="Heading2"/></w:pPr>
    <w:r><w:t>SUMMARY</w:t></w:r>
  </w:p>
  <w:p><w:r><w:rPr><w:sz w:val="20"/></w:rPr><w:t xml:space="preserve">{{.Summary}}</w:t></w:r></w:p>
  <w:p><w:r><w:t></w:t></w:r></w:p>
{{end}}
{{if .Experience}}
  <!-- Experience -->
  <w:p><w:pPr><w:pStyle w:val="Heading2"/></w:pPr>
    <w:r><w:t>EXPERIENCE</w:t></w:r>
  </w:p>
  {{range .Experience}}
  <w:p>
    <w:r><w:rPr><w:b/><w:sz w:val="22"/></w:rPr><w:t xml:space="preserve">{{.Position}} — {{.Company}}</w:t></w:r>
    <w:r><w:rPr><w:color w:val="808080"/><w:sz w:val="18"/></w:rPr><w:t xml:space="preserve">  ({{.StartDate}} – {{.EndDateDisplay}})</w:t></w:r>
  </w:p>
  {{if .Description}}<w:p><w:r><w:rPr><w:sz w:val="20"/></w:rPr><w:t xml:space="preserve">{{.Description}}</w:t></w:r></w:p>{{end}}
  <w:p><w:r><w:t></w:t></w:r></w:p>
  {{end}}
{{end}}
{{if .Education}}
  <!-- Education -->
  <w:p><w:pPr><w:pStyle w:val="Heading2"/></w:pPr>
    <w:r><w:t>EDUCATION</w:t></w:r>
  </w:p>
  {{range .Education}}
  <w:p>
    <w:r><w:rPr><w:b/><w:sz w:val="22"/></w:rPr><w:t xml:space="preserve">{{.Degree}}{{if .Field}} – {{.Field}}{{end}}, {{.Institution}}</w:t></w:r>
    <w:r><w:rPr><w:color w:val="808080"/><w:sz w:val="18"/></w:rPr><w:t xml:space="preserve">  ({{.StartDate}} – {{.EndDate}})</w:t></w:r>
  </w:p>
  {{end}}
  <w:p><w:r><w:t></w:t></w:r></w:p>
{{end}}
{{if .Skills}}
  <!-- Skills -->
  <w:p><w:pPr><w:pStyle w:val="Heading2"/></w:pPr>
    <w:r><w:t>SKILLS</w:t></w:r>
  </w:p>
  <w:p><w:r><w:rPr><w:sz w:val="20"/></w:rPr><w:t>{{.SkillsLine}}</w:t></w:r></w:p>
  <w:p><w:r><w:t></w:t></w:r></w:p>
{{end}}
{{if .Languages}}
  <!-- Languages -->
  <w:p><w:pPr><w:pStyle w:val="Heading2"/></w:pPr>
    <w:r><w:t>LANGUAGES</w:t></w:r>
  </w:p>
  {{range .Languages}}<w:p><w:r><w:rPr><w:sz w:val="20"/></w:rPr><w:t>{{.Name}}{{if .Proficiency}} – {{.Proficiency}}{{end}}</w:t></w:r></w:p>{{end}}
{{end}}
  <w:sectPr/>
</w:body>
</w:document>`

type simpleDocxData struct {
	models.CVData
	ContactLine string
	SkillsLine  string
	Experience  []expView
	Education   []models.Education
	Languages   []models.Language
}

type expView struct {
	models.Experience
	EndDateDisplay string
}

func generateSimpleDocx(data *models.CVData) ([]byte, error) {
	// Build contact line
	contactParts := []string{}
	if data.Email != "" {
		contactParts = append(contactParts, data.Email)
	}
	if data.Phone != "" {
		contactParts = append(contactParts, data.Phone)
	}
	if data.Address != "" {
		contactParts = append(contactParts, data.Address)
	}
	if data.Website != "" {
		contactParts = append(contactParts, data.Website)
	}

	exps := make([]expView, len(data.Experience))
	for i, e := range data.Experience {
		end := e.EndDate
		if end == "" {
			end = "Present"
		}
		exps[i] = expView{Experience: e, EndDateDisplay: end}
		exps[i].Experience.Company = xmlEscape(e.Company)
		exps[i].Experience.Position = xmlEscape(e.Position)
		exps[i].Experience.Description = xmlEscape(e.Description)
	}

	td := simpleDocxData{
		CVData:      *data,
		ContactLine: xmlEscape(strings.Join(contactParts, " | ")),
		SkillsLine:  xmlEscape(strings.Join(data.Skills, " • ")),
		Experience:  exps,
		Education:   data.Education,
		Languages:   data.Languages,
	}
	td.CVData.Name = xmlEscape(data.Name)
	td.CVData.Summary = xmlEscape(data.Summary)

	tmpl, err := template.New("doc").Parse(simpleDocxTmpl)
	if err != nil {
		return nil, fmt.Errorf("template parse error: %w", err)
	}
	var docBuf bytes.Buffer
	if err := tmpl.Execute(&docBuf, td); err != nil {
		return nil, fmt.Errorf("template execute error: %w", err)
	}

	return buildDocx([]docxFile{
		{name: "[Content_Types].xml", content: contentTypesXML},
		{name: "_rels/.rels", content: relsXML},
		{name: "word/_rels/document.xml.rels", content: wordRelsXML},
		{name: "word/styles.xml", content: simpleStylesXML},
		{name: "word/document.xml", content: docBuf.String()},
	})
}

// -------------------------------------------------------------------
// Japan template
// -------------------------------------------------------------------

const japanStylesXML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:styles xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:style w:type="paragraph" w:styleId="Normal" w:default="1">
    <w:name w:val="Normal"/>
    <w:rPr><w:sz w:val="20"/></w:rPr>
  </w:style>
  <w:style w:type="paragraph" w:styleId="JapanTitle">
    <w:name w:val="Japan Title"/>
    <w:rPr><w:b/><w:sz w:val="40"/><w:color w:val="28509F"/></w:rPr>
  </w:style>
  <w:style w:type="paragraph" w:styleId="JapanSection">
    <w:name w:val="Japan Section"/>
    <w:rPr><w:b/><w:sz w:val="24"/><w:color w:val="FFFFFF"/></w:rPr>
  </w:style>
</w:styles>`

const japanDocxTmpl = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"
            xmlns:w14="http://schemas.microsoft.com/office/word/2010/wordml">
<w:body>
  <!-- Title -->
  <w:p>
    <w:pPr><w:jc w:val="center"/><w:pStyle w:val="JapanTitle"/></w:pPr>
    <w:r><w:t>履歴書 / Curriculum Vitae</w:t></w:r>
  </w:p>
  <w:p><w:r><w:t></w:t></w:r></w:p>

  <!-- Personal Information header -->
  <w:p>
    <w:pPr><w:pStyle w:val="JapanSection"/>
      <w:shd w:val="clear" w:color="auto" w:fill="28509F"/>
    </w:pPr>
    <w:r><w:rPr><w:b/><w:color w:val="FFFFFF"/><w:sz w:val="24"/></w:rPr>
      <w:t>Personal Information / 個人情報</w:t>
    </w:r>
  </w:p>

  <!-- Personal fields as a 2-column table -->
  <w:tbl>
    <w:tblPr>
      <w:tblW w:w="9000" w:type="dxa"/>
      <w:tblBorders>
        <w:top w:val="single" w:sz="4" w:space="0" w:color="A0B4D4"/>
        <w:left w:val="single" w:sz="4" w:space="0" w:color="A0B4D4"/>
        <w:bottom w:val="single" w:sz="4" w:space="0" w:color="A0B4D4"/>
        <w:right w:val="single" w:sz="4" w:space="0" w:color="A0B4D4"/>
        <w:insideH w:val="single" w:sz="4" w:space="0" w:color="A0B4D4"/>
        <w:insideV w:val="single" w:sz="4" w:space="0" w:color="A0B4D4"/>
      </w:tblBorders>
    </w:tblPr>
    {{range .PersonalFields}}
    <w:tr>
      <w:tc>
        <w:tcPr>
          <w:tcW w:w="2500" w:type="dxa"/>
          <w:shd w:val="clear" w:color="auto" w:fill="E0E8F5"/>
        </w:tcPr>
        <w:p><w:r><w:rPr><w:b/><w:color w:val="28509F"/><w:sz w:val="18"/></w:rPr>
          <w:t xml:space="preserve">{{.Label}}</w:t>
        </w:r></w:p>
      </w:tc>
      <w:tc>
        <w:tcPr><w:tcW w:w="6500" w:type="dxa"/></w:tcPr>
        <w:p><w:r><w:rPr><w:sz w:val="18"/></w:rPr>
          <w:t xml:space="preserve">{{.Value}}</w:t>
        </w:r></w:p>
      </w:tc>
    </w:tr>
    {{end}}
  </w:tbl>
  <w:p><w:r><w:t></w:t></w:r></w:p>

{{if .Summary}}
  <!-- Self PR -->
  <w:p>
    <w:pPr><w:shd w:val="clear" w:color="auto" w:fill="28509F"/></w:pPr>
    <w:r><w:rPr><w:b/><w:color w:val="FFFFFF"/><w:sz w:val="24"/></w:rPr>
      <w:t>Self PR / 自己PR</w:t>
    </w:r>
  </w:p>
  <w:p><w:r><w:rPr><w:sz w:val="20"/></w:rPr>
    <w:t xml:space="preserve">{{.Summary}}</w:t>
  </w:r></w:p>
  <w:p><w:r><w:t></w:t></w:r></w:p>
{{end}}

{{if .Education}}
  <!-- Education -->
  <w:p>
    <w:pPr><w:shd w:val="clear" w:color="auto" w:fill="28509F"/></w:pPr>
    <w:r><w:rPr><w:b/><w:color w:val="FFFFFF"/><w:sz w:val="24"/></w:rPr>
      <w:t>Education / 学歴</w:t>
    </w:r>
  </w:p>
  <w:tbl>
    <w:tblPr>
      <w:tblW w:w="9000" w:type="dxa"/>
      <w:tblBorders>
        <w:top w:val="single" w:sz="4" w:space="0" w:color="A0B4D4"/>
        <w:left w:val="single" w:sz="4" w:space="0" w:color="A0B4D4"/>
        <w:bottom w:val="single" w:sz="4" w:space="0" w:color="A0B4D4"/>
        <w:right w:val="single" w:sz="4" w:space="0" w:color="A0B4D4"/>
        <w:insideH w:val="single" w:sz="4" w:space="0" w:color="A0B4D4"/>
        <w:insideV w:val="single" w:sz="4" w:space="0" w:color="A0B4D4"/>
      </w:tblBorders>
    </w:tblPr>
    <w:tr>
      <w:tc><w:tcPr><w:tcW w:w="2500" w:type="dxa"/><w:shd w:val="clear" w:color="auto" w:fill="DCE6F5"/></w:tcPr>
        <w:p><w:r><w:rPr><w:b/><w:color w:val="28509F"/><w:sz w:val="18"/></w:rPr><w:t>Period / 期間</w:t></w:r></w:p>
      </w:tc>
      <w:tc><w:tcPr><w:tcW w:w="6500" w:type="dxa"/><w:shd w:val="clear" w:color="auto" w:fill="DCE6F5"/></w:tcPr>
        <w:p><w:r><w:rPr><w:b/><w:color w:val="28509F"/><w:sz w:val="18"/></w:rPr><w:t>Institution / 学校名</w:t></w:r></w:p>
      </w:tc>
    </w:tr>
    {{range .Education}}
    <w:tr>
      <w:tc><w:tcPr><w:tcW w:w="2500" w:type="dxa"/></w:tcPr>
        <w:p><w:r><w:rPr><w:sz w:val="18"/></w:rPr><w:t xml:space="preserve">{{.StartDate}} – {{.EndDate}}</w:t></w:r></w:p>
      </w:tc>
      <w:tc><w:tcPr><w:tcW w:w="6500" w:type="dxa"/></w:tcPr>
        <w:p><w:r><w:rPr><w:sz w:val="18"/></w:rPr><w:t xml:space="preserve">{{.Institution}}{{if .Degree}} ({{.Degree}}{{if .Field}} – {{.Field}}{{end}}){{end}}</w:t></w:r></w:p>
      </w:tc>
    </w:tr>
    {{end}}
  </w:tbl>
  <w:p><w:r><w:t></w:t></w:r></w:p>
{{end}}

{{if .Experience}}
  <!-- Work Experience -->
  <w:p>
    <w:pPr><w:shd w:val="clear" w:color="auto" w:fill="28509F"/></w:pPr>
    <w:r><w:rPr><w:b/><w:color w:val="FFFFFF"/><w:sz w:val="24"/></w:rPr>
      <w:t>Work Experience / 職歴</w:t>
    </w:r>
  </w:p>
  <w:tbl>
    <w:tblPr>
      <w:tblW w:w="9000" w:type="dxa"/>
      <w:tblBorders>
        <w:top w:val="single" w:sz="4" w:space="0" w:color="A0B4D4"/>
        <w:left w:val="single" w:sz="4" w:space="0" w:color="A0B4D4"/>
        <w:bottom w:val="single" w:sz="4" w:space="0" w:color="A0B4D4"/>
        <w:right w:val="single" w:sz="4" w:space="0" w:color="A0B4D4"/>
        <w:insideH w:val="single" w:sz="4" w:space="0" w:color="A0B4D4"/>
        <w:insideV w:val="single" w:sz="4" w:space="0" w:color="A0B4D4"/>
      </w:tblBorders>
    </w:tblPr>
    <w:tr>
      <w:tc><w:tcPr><w:tcW w:w="2000" w:type="dxa"/><w:shd w:val="clear" w:color="auto" w:fill="DCE6F5"/></w:tcPr>
        <w:p><w:r><w:rPr><w:b/><w:color w:val="28509F"/><w:sz w:val="18"/></w:rPr><w:t>Period</w:t></w:r></w:p>
      </w:tc>
      <w:tc><w:tcPr><w:tcW w:w="3000" w:type="dxa"/><w:shd w:val="clear" w:color="auto" w:fill="DCE6F5"/></w:tcPr>
        <w:p><w:r><w:rPr><w:b/><w:color w:val="28509F"/><w:sz w:val="18"/></w:rPr><w:t>Company</w:t></w:r></w:p>
      </w:tc>
      <w:tc><w:tcPr><w:tcW w:w="4000" w:type="dxa"/><w:shd w:val="clear" w:color="auto" w:fill="DCE6F5"/></w:tcPr>
        <w:p><w:r><w:rPr><w:b/><w:color w:val="28509F"/><w:sz w:val="18"/></w:rPr><w:t>Position / Role</w:t></w:r></w:p>
      </w:tc>
    </w:tr>
    {{range .Experience}}
    <w:tr>
      <w:tc><w:tcPr><w:tcW w:w="2000" w:type="dxa"/></w:tcPr>
        <w:p><w:r><w:rPr><w:sz w:val="18"/></w:rPr><w:t xml:space="preserve">{{.StartDate}} – {{.EndDateDisplay}}</w:t></w:r></w:p>
      </w:tc>
      <w:tc><w:tcPr><w:tcW w:w="3000" w:type="dxa"/></w:tcPr>
        <w:p><w:r><w:rPr><w:sz w:val="18"/></w:rPr><w:t xml:space="preserve">{{.Company}}</w:t></w:r></w:p>
      </w:tc>
      <w:tc><w:tcPr><w:tcW w:w="4000" w:type="dxa"/></w:tcPr>
        <w:p><w:r><w:rPr><w:sz w:val="18"/></w:rPr><w:t xml:space="preserve">{{.Position}}</w:t></w:r></w:p>
      </w:tc>
    </w:tr>
    {{end}}
  </w:tbl>
  <w:p><w:r><w:t></w:t></w:r></w:p>
{{end}}

{{if .Skills}}
  <!-- Skills -->
  <w:p>
    <w:pPr><w:shd w:val="clear" w:color="auto" w:fill="28509F"/></w:pPr>
    <w:r><w:rPr><w:b/><w:color w:val="FFFFFF"/><w:sz w:val="24"/></w:rPr>
      <w:t>Skills / スキル</w:t>
    </w:r>
  </w:p>
  <w:p><w:r><w:rPr><w:sz w:val="20"/></w:rPr><w:t>{{.SkillsLine}}</w:t></w:r></w:p>
  <w:p><w:r><w:t></w:t></w:r></w:p>
{{end}}

{{if .Languages}}
  <!-- Languages -->
  <w:p>
    <w:pPr><w:shd w:val="clear" w:color="auto" w:fill="28509F"/></w:pPr>
    <w:r><w:rPr><w:b/><w:color w:val="FFFFFF"/><w:sz w:val="24"/></w:rPr>
      <w:t>Languages / 語学</w:t>
    </w:r>
  </w:p>
  {{range .Languages}}
  <w:p><w:r><w:rPr><w:sz w:val="20"/></w:rPr>
    <w:t xml:space="preserve">• {{.Name}}{{if .Proficiency}} ({{.Proficiency}}){{end}}</w:t>
  </w:r></w:p>
  {{end}}
{{end}}

  <w:sectPr/>
</w:body>
</w:document>`

type japanDocxData struct {
	models.CVData
	PersonalFields []labelValue
	SkillsLine     string
	Experience     []expView
	Education      []models.Education
	Languages      []models.Language
}

type labelValue struct {
	Label string
	Value string
}

func generateJapanDocx(data *models.CVData) ([]byte, error) {
	// Build personal fields list
	fields := []labelValue{
		{Label: "Name / 氏名", Value: xmlEscape(data.Name)},
	}
	if data.BirthDate != "" {
		fields = append(fields, labelValue{"Date of Birth / 生年月日", xmlEscape(data.BirthDate)})
	}
	if data.Gender != "" {
		fields = append(fields, labelValue{"Gender / 性別", xmlEscape(data.Gender)})
	}
	if data.Nationality != "" {
		fields = append(fields, labelValue{"Nationality / 国籍", xmlEscape(data.Nationality)})
	}
	if data.Address != "" {
		fields = append(fields, labelValue{"Address / 住所", xmlEscape(data.Address)})
	}
	if data.Phone != "" {
		fields = append(fields, labelValue{"Phone / 電話", xmlEscape(data.Phone)})
	}
	if data.Email != "" {
		fields = append(fields, labelValue{"Email / メール", xmlEscape(data.Email)})
	}
	if data.Website != "" {
		fields = append(fields, labelValue{"Website", xmlEscape(data.Website)})
	}

	exps := make([]expView, len(data.Experience))
	for i, e := range data.Experience {
		end := e.EndDate
		if end == "" {
			end = "Present"
		}
		exps[i] = expView{Experience: e, EndDateDisplay: end}
		exps[i].Experience.Company = xmlEscape(e.Company)
		exps[i].Experience.Position = xmlEscape(e.Position)
	}

	td := japanDocxData{
		CVData:         *data,
		PersonalFields: fields,
		SkillsLine:     xmlEscape(strings.Join(data.Skills, " / ")),
		Experience:     exps,
		Education:      data.Education,
		Languages:      data.Languages,
	}
	td.CVData.Summary = xmlEscape(data.Summary)

	tmpl, err := template.New("doc").Parse(japanDocxTmpl)
	if err != nil {
		return nil, fmt.Errorf("template parse error: %w", err)
	}
	var docBuf bytes.Buffer
	if err := tmpl.Execute(&docBuf, td); err != nil {
		return nil, fmt.Errorf("template execute error: %w", err)
	}

	return buildDocx([]docxFile{
		{name: "[Content_Types].xml", content: contentTypesXML},
		{name: "_rels/.rels", content: relsXML},
		{name: "word/_rels/document.xml.rels", content: wordRelsXML},
		{name: "word/styles.xml", content: japanStylesXML},
		{name: "word/document.xml", content: docBuf.String()},
	})
}
