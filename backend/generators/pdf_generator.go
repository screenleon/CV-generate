package generators

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/screenleon/cv-generate/models"
)

// GeneratePDF generates a PDF CV using the specified template style.
// It renders an HTML page via a headless Chromium browser so that
// CJK characters (Japanese, Chinese, Korean) display correctly.
// Returns the PDF bytes or an error.
func GeneratePDF(data *models.CVData) ([]byte, error) {
	var htmlContent string
	var err error
	switch strings.ToLower(data.Template) {
	case "japan":
		htmlContent, err = renderJapanHTML(data)
	case "shokumu":
		htmlContent, err = renderShokumuHTML(data)
	default:
		htmlContent, err = renderSimpleHTML(data)
	}
	if err != nil {
		return nil, fmt.Errorf("render HTML: %w", err)
	}
	return htmlToPDF(htmlContent)
}

// htmlToPDF uses headless Chromium to print an HTML string to a PDF byte slice.
func htmlToPDF(htmlContent string) ([]byte, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("headless", true),
	)
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancelAlloc()

	ctx, cancelCtx := chromedp.NewContext(allocCtx)
	defer cancelCtx()

	ctx, cancelTimeout := context.WithTimeout(ctx, 60*time.Second)
	defer cancelTimeout()

	// Use a data URI so no file system access is needed.
	dataURI := "data:text/html;charset=utf-8," + urlEncode(htmlContent)

	var pdfBuf []byte
	err := chromedp.Run(ctx,
		chromedp.Navigate(dataURI),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfBuf, _, err = page.PrintToPDF().
				WithPrintBackground(true).
				WithPaperWidth(8.27).   // A4 width in inches
				WithPaperHeight(11.69). // A4 height in inches
				WithMarginTop(0.6).
				WithMarginBottom(0.6).
				WithMarginLeft(0.6).
				WithMarginRight(0.6).
				Do(ctx)
			return err
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("chromedp PrintToPDF: %w", err)
	}
	return pdfBuf, nil
}

// urlEncode percent-encodes an HTML string for use in a data URI.
// Only characters unsafe in a URI are escaped.
func urlEncode(s string) string {
	var buf bytes.Buffer
	for _, b := range []byte(s) {
		if isURLSafe(b) {
			buf.WriteByte(b)
		} else {
			fmt.Fprintf(&buf, "%%%02X", b)
		}
	}
	return buf.String()
}

func isURLSafe(b byte) bool {
	return (b >= 'A' && b <= 'Z') ||
		(b >= 'a' && b <= 'z') ||
		(b >= '0' && b <= '9') ||
		b == '-' || b == '_' || b == '.' || b == '~' ||
		b == '!' || b == '\'' || b == '(' || b == ')' || b == '*' ||
		b == '/' || b == ':' || b == '@' || b == ',' || b == ';' ||
		b == '+' || b == '=' || b == '?' || b == '#'
}

// -------------------------------------------------------------------
// Shared CSS and font stack
// -------------------------------------------------------------------

const sharedCSS = `
  * { box-sizing: border-box; margin: 0; padding: 0; }
  body {
    font-family: 'Noto Sans CJK JP', 'Noto Sans JP', 'IPAGothic',
                 'Hiragino Sans', 'Meiryo', 'Yu Gothic', sans-serif;
    font-size: 10pt;
    color: #222;
    line-height: 1.5;
  }
  table { width: 100%; border-collapse: collapse; }
  td, th { padding: 4px 8px; border: 1px solid #ccc; vertical-align: top; }
  @media print {
    body { print-color-adjust: exact; -webkit-print-color-adjust: exact; }
  }
`

// -------------------------------------------------------------------
// Simple template
// -------------------------------------------------------------------

const simpleTmplSrc = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<style>
` + sharedCSS + `
  .header { background: #28508A; color: #fff; padding: 18px 24px 14px; }
  .header h1 { font-size: 22pt; font-weight: bold; margin-bottom: 4px; }
  .header .contact { font-size: 9pt; opacity: 0.9; }
  .body { padding: 16px 24px; }
  .section-title {
    font-size: 10.5pt; font-weight: bold; color: #28508A;
    border-bottom: 2px solid #28508A; margin: 14px 0 6px;
    padding-bottom: 2px; text-transform: uppercase; letter-spacing: .5px;
  }
  .summary { font-size: 10pt; }
  .exp-row { display: flex; justify-content: space-between; align-items: baseline; }
  .exp-company { font-weight: bold; font-size: 10.5pt; }
  .exp-date { font-size: 9pt; color: #666; white-space: nowrap; }
  .exp-position { font-style: italic; font-size: 10pt; color: #28508A; }
  .exp-desc { font-size: 9.5pt; color: #444; margin-top: 2px; }
  .exp-entry { margin-bottom: 10px; }
  .edu-row { display: flex; justify-content: space-between; align-items: baseline; }
  .edu-inst { font-weight: bold; font-size: 10pt; }
  .edu-date { font-size: 9pt; color: #666; }
  .edu-degree { font-size: 9.5pt; color: #555; }
  .edu-entry { margin-bottom: 8px; }
  .skills { font-size: 10pt; }
  .lang-entry { display: inline-block; margin-right: 18px; margin-bottom: 4px; font-size: 10pt; }
</style>
</head>
<body>
<div class="header">
  <h1>{{.Name}}</h1>
  {{if .Contact}}<div class="contact">{{.Contact}}</div>{{end}}
</div>
<div class="body">
  {{if .Summary}}
  <div class="section-title">Summary</div>
  <div class="summary">{{.Summary}}</div>
  {{end}}
  {{if .Experience}}
  <div class="section-title">Experience</div>
  {{range .Experience}}
  <div class="exp-entry">
    <div class="exp-row">
      <span class="exp-company">{{.Company}}</span>
      <span class="exp-date">{{.StartDate}} – {{.EndDate}}</span>
    </div>
    <div class="exp-position">{{.Position}}</div>
    {{if .Description}}<div class="exp-desc">{{.Description}}</div>{{end}}
  </div>
  {{end}}
  {{end}}
  {{if .Education}}
  <div class="section-title">Education</div>
  {{range .Education}}
  <div class="edu-entry">
    <div class="edu-row">
      <span class="edu-inst">{{.Institution}}</span>
      <span class="edu-date">{{.StartDate}} – {{.EndDate}}</span>
    </div>
    <div class="edu-degree">{{.Degree}}{{if .Field}} – {{.Field}}{{end}}</div>
  </div>
  {{end}}
  {{end}}
  {{if .Skills}}
  <div class="section-title">Skills</div>
  <div class="skills">{{.SkillsLine}}</div>
  {{end}}
  {{if .Languages}}
  <div class="section-title">Languages</div>
  {{range .Languages}}
  <div class="lang-entry">{{.Name}}{{if .Proficiency}} – {{.Proficiency}}{{end}}</div>
  {{end}}
  {{end}}
</div>
</body>
</html>`

type simpleHTMLData struct {
	Name       string
	Contact    string
	Summary    string
	Experience []expHTMLView
	Education  []models.Education
	Skills     []string
	SkillsLine string
	Languages  []models.Language
}

type expHTMLView struct {
	Company     string
	Position    string
	StartDate   string
	EndDate     string
	Description string
}

func renderSimpleHTML(data *models.CVData) (string, error) {
	parts := []string{}
	if data.Email != "" {
		parts = append(parts, data.Email)
	}
	if data.Phone != "" {
		parts = append(parts, data.Phone)
	}
	if data.Address != "" {
		parts = append(parts, data.Address)
	}
	if data.Website != "" {
		parts = append(parts, data.Website)
	}

	exps := make([]expHTMLView, len(data.Experience))
	for i, e := range data.Experience {
		end := e.EndDate
		if end == "" {
			end = "Present"
		}
		exps[i] = expHTMLView{
			Company:     e.Company,
			Position:    e.Position,
			StartDate:   e.StartDate,
			EndDate:     end,
			Description: e.Description,
		}
	}

	td := simpleHTMLData{
		Name:       data.Name,
		Contact:    strings.Join(parts, "  |  "),
		Summary:    data.Summary,
		Experience: exps,
		Education:  data.Education,
		Skills:     data.Skills,
		SkillsLine: strings.Join(data.Skills, "  •  "),
		Languages:  data.Languages,
	}

	tmpl, err := template.New("simple").Parse(simpleTmplSrc)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, td); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// -------------------------------------------------------------------
// Japan (履歴書) template
// -------------------------------------------------------------------

const japanTmplSrc = `<!DOCTYPE html>
<html lang="ja">
<head>
<meta charset="UTF-8">
<style>
` + sharedCSS + `
  .title { text-align: center; font-size: 18pt; font-weight: bold; margin: 10px 0; color: #222; }
  .section-header {
    background: #28508A; color: #fff; font-weight: bold;
    font-size: 11pt; padding: 4px 10px; margin: 12px 0 4px;
  }
  .info-table td.label {
    background: #EEF2FB; color: #28508A; font-weight: bold;
    width: 160px; white-space: nowrap;
  }
  .info-table td.value { background: #fff; }
  .pr-box { border: 1px solid #ccc; padding: 8px; font-size: 10pt; min-height: 40px; }
  .col-date { width: 120px; text-align: center; background: #DCE6F5; font-weight: bold; }
  .col-inst { background: #DCE6F5; font-weight: bold; }
  .col-company { width: 160px; background: #DCE6F5; font-weight: bold; }
  .col-role { background: #DCE6F5; font-weight: bold; }
  .desc-row td { font-size: 9pt; color: #555; font-style: italic; }
  .skills-box { border: 1px solid #ccc; padding: 8px; font-size: 10pt; }
  .lang-item { display: inline-block; margin-right: 16px; font-size: 10pt; }
</style>
</head>
<body>
<div class="title">履歴書 / Curriculum Vitae</div>

<div class="section-header">個人情報 / Personal Information</div>
<table class="info-table">
  <tr><td class="label">氏名 / Name</td><td class="value">{{.Name}}</td></tr>
  {{if .BirthDate}}<tr><td class="label">生年月日 / Date of Birth</td><td class="value">{{.BirthDate}}</td></tr>{{end}}
  {{if .Gender}}<tr><td class="label">性別 / Gender</td><td class="value">{{.Gender}}</td></tr>{{end}}
  {{if .Nationality}}<tr><td class="label">国籍 / Nationality</td><td class="value">{{.Nationality}}</td></tr>{{end}}
  {{if .Address}}<tr><td class="label">住所 / Address</td><td class="value">{{.Address}}</td></tr>{{end}}
  {{if .Phone}}<tr><td class="label">電話 / Phone</td><td class="value">{{.Phone}}</td></tr>{{end}}
  {{if .Email}}<tr><td class="label">メール / Email</td><td class="value">{{.Email}}</td></tr>{{end}}
  {{if .Website}}<tr><td class="label">Website</td><td class="value">{{.Website}}</td></tr>{{end}}
</table>

{{if .Summary}}
<div class="section-header">自己PR / Self PR</div>
<div class="pr-box">{{.Summary}}</div>
{{end}}

{{if .Education}}
<div class="section-header">学歴 / Education</div>
<table>
  <tr><th class="col-date">期間 / Period</th><th class="col-inst">学校名・学位 / Institution &amp; Degree</th></tr>
  {{range .Education}}
  <tr>
    <td style="text-align:center">{{.StartDate}} – {{.EndDate}}</td>
    <td>{{.Institution}}{{if .Degree}} ({{.Degree}}{{if .Field}} – {{.Field}}{{end}}){{end}}</td>
  </tr>
  {{end}}
</table>
{{end}}

{{if .Experience}}
<div class="section-header">職歴 / Work Experience</div>
<table>
  <tr>
    <th class="col-date">期間 / Period</th>
    <th class="col-company">会社名 / Company</th>
    <th class="col-role">役職 / Position</th>
  </tr>
  {{range .Experience}}
  <tr>
    <td style="text-align:center">{{.StartDate}} – {{.EndDate}}</td>
    <td>{{.Company}}</td>
    <td>{{.Position}}</td>
  </tr>
  {{if .Description}}
  <tr class="desc-row"><td colspan="3" style="padding-left:24px">{{.Description}}</td></tr>
  {{end}}
  {{end}}
</table>
{{end}}

{{if .Skills}}
<div class="section-header">スキル / Skills</div>
<div class="skills-box">{{.SkillsLine}}</div>
{{end}}

{{if .Languages}}
<div class="section-header">語学 / Languages</div>
<div style="padding: 6px 0">
  {{range .Languages}}<span class="lang-item">• {{.Name}}{{if .Proficiency}} ({{.Proficiency}}){{end}}</span>{{end}}
</div>
{{end}}
</body>
</html>`

type japanHTMLData struct {
	Name        string
	BirthDate   string
	Gender      string
	Nationality string
	Address     string
	Phone       string
	Email       string
	Website     string
	Summary     string
	Experience  []expHTMLView
	Education   []models.Education
	Skills      []string
	SkillsLine  string
	Languages   []models.Language
}

func renderJapanHTML(data *models.CVData) (string, error) {
	exps := make([]expHTMLView, len(data.Experience))
	for i, e := range data.Experience {
		end := e.EndDate
		if end == "" {
			end = "現在"
		}
		exps[i] = expHTMLView{
			Company:     e.Company,
			Position:    e.Position,
			StartDate:   e.StartDate,
			EndDate:     end,
			Description: e.Description,
		}
	}

	td := japanHTMLData{
		Name:        data.Name,
		BirthDate:   data.BirthDate,
		Gender:      data.Gender,
		Nationality: data.Nationality,
		Address:     data.Address,
		Phone:       data.Phone,
		Email:       data.Email,
		Website:     data.Website,
		Summary:     data.Summary,
		Experience:  exps,
		Education:   data.Education,
		Skills:      data.Skills,
		SkillsLine:  strings.Join(data.Skills, "  /  "),
		Languages:   data.Languages,
	}

	tmpl, err := template.New("japan").Parse(japanTmplSrc)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, td); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// -------------------------------------------------------------------
// Shokumu (職務経歴書) template
// -------------------------------------------------------------------

const shokumuTmplSrc = `<!DOCTYPE html>
<html lang="ja">
<head>
<meta charset="UTF-8">
<style>
` + sharedCSS + `
  .title { text-align: center; font-size: 18pt; font-weight: bold; margin: 10px 0; color: #222; }
  .section-header {
    background: #28508A; color: #fff; font-weight: bold;
    font-size: 11pt; padding: 4px 10px; margin: 14px 0 4px;
  }
  .info-table td.label {
    background: #EEF2FB; color: #28508A; font-weight: bold;
    width: 160px; white-space: nowrap;
  }
  .info-table td.value { background: #fff; }
  .summary-box { border: 1px solid #ccc; padding: 8px; font-size: 10pt; min-height: 40px; }
  .exp-company-header {
    background: #EBF0FA; font-weight: bold; font-size: 10.5pt;
    padding: 5px 10px; margin: 10px 0 0; border: 1px solid #c0cce8;
  }
  .exp-detail-table td.detail-label {
    background: #EEF2FB; color: #28508A; font-weight: bold;
    width: 130px; white-space: nowrap;
  }
  .exp-detail-table td.detail-value { background: #fff; }
  .skills-table td.skill-label {
    background: #EBF0FA; color: #28508A; font-weight: bold;
    width: 150px; white-space: nowrap;
  }
  .skills-table td.skill-value { background: #fff; }
  .edu-item { font-size: 10pt; margin: 3px 0; }
  .lang-item { display: inline-block; margin-right: 16px; font-size: 10pt; }
</style>
</head>
<body>
<div class="title">職務経歴書 / Professional Resume</div>

<div class="section-header">個人情報 / Personal Information</div>
<table class="info-table">
  <tr><td class="label">氏名 / Name</td><td class="value">{{.Name}}</td></tr>
  {{if .Email}}<tr><td class="label">メール / Email</td><td class="value">{{.Email}}</td></tr>{{end}}
  {{if .Phone}}<tr><td class="label">電話 / Phone</td><td class="value">{{.Phone}}</td></tr>{{end}}
</table>

{{if .Summary}}
<div class="section-header">職務概要 / Professional Summary</div>
<div class="summary-box">{{.Summary}}</div>
{{end}}

{{if .Experience}}
<div class="section-header">職務経歴 / Work Experience</div>
{{range $i, $exp := .Experience}}
<div class="exp-company-header">【{{add $i 1}}】{{$exp.Company}}  （{{$exp.StartDate}} – {{$exp.EndDate}}）</div>
<table class="exp-detail-table">
  {{if $exp.Project}}<tr><td class="detail-label">案件 / Project</td><td class="detail-value">{{$exp.Project}}</td></tr>{{end}}
  <tr><td class="detail-label">役職 / Position</td><td class="detail-value">{{$exp.Position}}</td></tr>
  {{if $exp.Role}}<tr><td class="detail-label">担当業務 / Role</td><td class="detail-value">{{$exp.Role}}</td></tr>{{end}}
  {{if $exp.Description}}<tr><td class="detail-label">詳細 / Details</td><td class="detail-value">{{$exp.Description}}</td></tr>{{end}}
  {{if $exp.TechStackLine}}<tr><td class="detail-label">技術 / Tech Stack</td><td class="detail-value">{{$exp.TechStackLine}}</td></tr>{{end}}
</table>
{{end}}
{{end}}

{{if .HasTechStack}}
<div class="section-header">スキル一覧 / Technical Skills</div>
<table class="skills-table">
  {{range .TechStackRows}}
  <tr><td class="skill-label">{{.Label}}</td><td class="skill-value">{{.Value}}</td></tr>
  {{end}}
</table>
{{else if .Skills}}
<div class="section-header">スキル / Skills</div>
<table class="skills-table">
  <tr><td class="skill-value" colspan="2">{{.SkillsLine}}</td></tr>
</table>
{{end}}

{{if .Education}}
<div class="section-header">学歴 / Education</div>
{{range .Education}}
<div class="edu-item">• {{.StartDate}} – {{.EndDate}}：{{.Institution}}{{if .Degree}} ({{.Degree}}{{if .Field}} – {{.Field}}{{end}}){{end}}</div>
{{end}}
{{end}}

{{if .Languages}}
<div class="section-header">語学 / Languages</div>
<div style="padding: 6px 0">
  {{range .Languages}}<span class="lang-item">• {{.Name}}{{if .Proficiency}} ({{.Proficiency}}){{end}}</span>{{end}}
</div>
{{end}}
</body>
</html>`

type shokumuHTMLData struct {
	Name         string
	Email        string
	Phone        string
	Summary      string
	Experience   []shokumuExpHTMLView
	Education    []models.Education
	Skills       []string
	SkillsLine   string
	Languages    []models.Language
	TechStackRows []labelValueHTML
	HasTechStack bool
}

type shokumuExpHTMLView struct {
	Company       string
	Position      string
	StartDate     string
	EndDate       string
	Project       string
	Role          string
	Description   string
	TechStackLine string
}

type labelValueHTML struct {
	Label string
	Value string
}

func renderShokumuHTML(data *models.CVData) (string, error) {
	exps := make([]shokumuExpHTMLView, len(data.Experience))
	for i, e := range data.Experience {
		end := e.EndDate
		if end == "" {
			end = "現在"
		}
		techLine := ""
		if len(e.TechStack) > 0 {
			techLine = strings.Join(e.TechStack, ", ")
		}
		exps[i] = shokumuExpHTMLView{
			Company:       e.Company,
			Position:      e.Position,
			StartDate:     e.StartDate,
			EndDate:       end,
			Project:       e.Project,
			Role:          e.Role,
			Description:   e.Description,
			TechStackLine: techLine,
		}
	}

	techStackRows := []labelValueHTML{}
	hasTechStack := false
	if data.TechStack != nil {
		if len(data.TechStack.Languages) > 0 {
			techStackRows = append(techStackRows, labelValueHTML{"言語 / Languages", strings.Join(data.TechStack.Languages, ", ")})
			hasTechStack = true
		}
		if len(data.TechStack.Frameworks) > 0 {
			techStackRows = append(techStackRows, labelValueHTML{"FW / Frameworks", strings.Join(data.TechStack.Frameworks, ", ")})
			hasTechStack = true
		}
		if len(data.TechStack.Databases) > 0 {
			techStackRows = append(techStackRows, labelValueHTML{"DB / Databases", strings.Join(data.TechStack.Databases, ", ")})
			hasTechStack = true
		}
		if len(data.TechStack.Infrastructure) > 0 {
			techStackRows = append(techStackRows, labelValueHTML{"インフラ / Infrastructure", strings.Join(data.TechStack.Infrastructure, ", ")})
			hasTechStack = true
		}
		if len(data.TechStack.Tools) > 0 {
			techStackRows = append(techStackRows, labelValueHTML{"ツール / Tools", strings.Join(data.TechStack.Tools, ", ")})
			hasTechStack = true
		}
	}

	td := shokumuHTMLData{
		Name:          data.Name,
		Email:         data.Email,
		Phone:         data.Phone,
		Summary:       data.Summary,
		Experience:    exps,
		Education:     data.Education,
		Skills:        data.Skills,
		SkillsLine:    strings.Join(data.Skills, " / "),
		Languages:     data.Languages,
		TechStackRows: techStackRows,
		HasTechStack:  hasTechStack,
	}

	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
	}
	tmpl, err := template.New("shokumu").Funcs(funcMap).Parse(shokumuTmplSrc)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, td); err != nil {
		return "", err
	}
	return buf.String(), nil
}

