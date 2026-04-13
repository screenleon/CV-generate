package generators

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/go-pdf/fpdf"
	"github.com/screenleon/cv-generate/models"
)

// GeneratePDF generates a PDF CV using the specified template style.
// Returns the PDF bytes or an error.
func GeneratePDF(data *models.CVData) ([]byte, error) {
	switch strings.ToLower(data.Template) {
	case "japan":
		return generateJapanPDF(data)
	default:
		return generateSimplePDF(data)
	}
}

// -------------------------------------------------------------------
// Simple template
// -------------------------------------------------------------------

func generateSimplePDF(data *models.CVData) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 20, 20)
	pdf.AddPage()

	// --- Header: Name ---
	pdf.SetFont("Helvetica", "B", 24)
	pdf.SetTextColor(30, 30, 30)
	pdf.CellFormat(0, 12, data.Name, "", 1, "C", false, 0, "")

	// --- Contact line ---
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
	if len(contactParts) > 0 {
		pdf.SetFont("Helvetica", "", 9)
		pdf.SetTextColor(80, 80, 80)
		pdf.CellFormat(0, 6, strings.Join(contactParts, "  |  "), "", 1, "C", false, 0, "")
	}

	// Divider line
	pdf.SetDrawColor(60, 120, 200)
	pdf.SetLineWidth(0.8)
	pdf.Line(20, pdf.GetY()+3, 190, pdf.GetY()+3)
	pdf.Ln(7)

	// --- Summary ---
	if data.Summary != "" {
		addSimpleSection(pdf, "SUMMARY")
		pdf.SetFont("Helvetica", "", 10)
		pdf.SetTextColor(50, 50, 50)
		pdf.MultiCell(0, 5, data.Summary, "", "L", false)
		pdf.Ln(4)
	}

	// --- Experience ---
	if len(data.Experience) > 0 {
		addSimpleSection(pdf, "EXPERIENCE")
		for _, exp := range data.Experience {
			endDate := exp.EndDate
			if endDate == "" {
				endDate = "Present"
			}
			pdf.SetFont("Helvetica", "B", 11)
			pdf.SetTextColor(30, 30, 30)
			pdf.CellFormat(130, 6, exp.Position, "", 0, "L", false, 0, "")
			pdf.SetFont("Helvetica", "", 9)
			pdf.SetTextColor(100, 100, 100)
			pdf.CellFormat(0, 6, exp.StartDate+" – "+endDate, "", 1, "R", false, 0, "")

			pdf.SetFont("Helvetica", "I", 10)
			pdf.SetTextColor(60, 60, 60)
			pdf.CellFormat(0, 5, exp.Company, "", 1, "L", false, 0, "")

			if exp.Description != "" {
				pdf.SetFont("Helvetica", "", 10)
				pdf.SetTextColor(50, 50, 50)
				pdf.MultiCell(0, 5, exp.Description, "", "L", false)
			}
			pdf.Ln(3)
		}
	}

	// --- Education ---
	if len(data.Education) > 0 {
		addSimpleSection(pdf, "EDUCATION")
		for _, edu := range data.Education {
			pdf.SetFont("Helvetica", "B", 11)
			pdf.SetTextColor(30, 30, 30)
			degree := edu.Degree
			if edu.Field != "" {
				degree += " – " + edu.Field
			}
			pdf.CellFormat(130, 6, degree, "", 0, "L", false, 0, "")
			pdf.SetFont("Helvetica", "", 9)
			pdf.SetTextColor(100, 100, 100)
			pdf.CellFormat(0, 6, edu.StartDate+" – "+edu.EndDate, "", 1, "R", false, 0, "")

			pdf.SetFont("Helvetica", "I", 10)
			pdf.SetTextColor(60, 60, 60)
			pdf.CellFormat(0, 5, edu.Institution, "", 1, "L", false, 0, "")
			pdf.Ln(3)
		}
	}

	// --- Skills ---
	if len(data.Skills) > 0 {
		addSimpleSection(pdf, "SKILLS")
		pdf.SetFont("Helvetica", "", 10)
		pdf.SetTextColor(50, 50, 50)
		pdf.MultiCell(0, 5, strings.Join(data.Skills, "  •  "), "", "L", false)
		pdf.Ln(4)
	}

	// --- Languages ---
	if len(data.Languages) > 0 {
		addSimpleSection(pdf, "LANGUAGES")
		for _, lang := range data.Languages {
			label := lang.Name
			if lang.Proficiency != "" {
				label += " – " + lang.Proficiency
			}
			pdf.SetFont("Helvetica", "", 10)
			pdf.SetTextColor(50, 50, 50)
			pdf.CellFormat(0, 5, label, "", 1, "L", false, 0, "")
		}
		pdf.Ln(2)
	}

	return pdfToBytes(pdf)
}

func addSimpleSection(pdf *fpdf.Fpdf, title string) {
	pdf.SetFont("Helvetica", "B", 12)
	pdf.SetTextColor(60, 120, 200)
	pdf.CellFormat(0, 7, title, "", 1, "L", false, 0, "")
	pdf.SetDrawColor(60, 120, 200)
	pdf.SetLineWidth(0.3)
	pdf.Line(20, pdf.GetY(), 190, pdf.GetY())
	pdf.Ln(3)
}

// -------------------------------------------------------------------
// Japan (履歴書) template
// -------------------------------------------------------------------

func generateJapanPDF(data *models.CVData) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.AddPage()

	pageW := 180.0 // usable width (210 - 30 margins)

	// ---- Title ----
	pdf.SetFont("Helvetica", "B", 16)
	pdf.SetTextColor(20, 20, 20)
	pdf.CellFormat(0, 10, "履歴書 / Curriculum Vitae", "", 1, "C", false, 0, "")
	pdf.Ln(2)

	// ---- Personal information box ----
	pdf.SetFont("Helvetica", "B", 11)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFillColor(40, 80, 160)
	pdf.CellFormat(0, 7, "  Personal Information / 個人情報", "", 1, "L", true, 0, "")
	pdf.Ln(1)

	addJapanField(pdf, pageW, "Name / 氏名", data.Name)
	if data.BirthDate != "" {
		addJapanField(pdf, pageW, "Date of Birth / 生年月日", data.BirthDate)
	}
	if data.Gender != "" {
		addJapanField(pdf, pageW, "Gender / 性別", data.Gender)
	}
	if data.Nationality != "" {
		addJapanField(pdf, pageW, "Nationality / 国籍", data.Nationality)
	}
	addJapanField(pdf, pageW, "Address / 住所", data.Address)
	addJapanField(pdf, pageW, "Phone / 電話", data.Phone)
	addJapanField(pdf, pageW, "Email / メール", data.Email)
	if data.Website != "" {
		addJapanField(pdf, pageW, "Website", data.Website)
	}
	pdf.Ln(4)

	// ---- Summary / Motivation ----
	if data.Summary != "" {
		addJapanSectionHeader(pdf, "Self PR / 自己PR")
		pdf.SetFont("Helvetica", "", 10)
		pdf.SetTextColor(40, 40, 40)
		pdf.MultiCell(0, 5, data.Summary, "1", "L", false)
		pdf.Ln(4)
	}

	// ---- Education ----
	if len(data.Education) > 0 {
		addJapanSectionHeader(pdf, "Education / 学歴")
		pdf.SetFont("Helvetica", "B", 9)
		pdf.SetFillColor(220, 230, 245)
		pdf.SetTextColor(30, 30, 30)
		colDate := 35.0
		colRest := pageW - colDate
		pdf.CellFormat(colDate, 6, "Period / 期間", "1", 0, "C", true, 0, "")
		pdf.CellFormat(colRest, 6, "Institution / 学校名", "1", 1, "C", true, 0, "")
		pdf.SetFont("Helvetica", "", 9)
		pdf.SetFillColor(255, 255, 255)
		for _, edu := range data.Education {
			period := edu.StartDate + " – " + edu.EndDate
			institution := edu.Institution
			if edu.Degree != "" {
				institution += " (" + edu.Degree
				if edu.Field != "" {
					institution += " – " + edu.Field
				}
				institution += ")"
			}
			pdf.CellFormat(colDate, 6, period, "1", 0, "C", false, 0, "")
			pdf.CellFormat(colRest, 6, institution, "1", 1, "L", false, 0, "")
		}
		pdf.Ln(4)
	}

	// ---- Experience ----
	if len(data.Experience) > 0 {
		addJapanSectionHeader(pdf, "Work Experience / 職歴")
		pdf.SetFont("Helvetica", "B", 9)
		pdf.SetFillColor(220, 230, 245)
		pdf.SetTextColor(30, 30, 30)
		colDate := 35.0
		colCompany := 50.0
		colRole := pageW - colDate - colCompany
		pdf.CellFormat(colDate, 6, "Period / 期間", "1", 0, "C", true, 0, "")
		pdf.CellFormat(colCompany, 6, "Company / 会社名", "1", 0, "C", true, 0, "")
		pdf.CellFormat(colRole, 6, "Position / 役職", "1", 1, "C", true, 0, "")
		pdf.SetFont("Helvetica", "", 9)
		pdf.SetFillColor(255, 255, 255)
		for _, exp := range data.Experience {
			endDate := exp.EndDate
			if endDate == "" {
				endDate = "Present"
			}
			period := exp.StartDate + " – " + endDate
			pdf.CellFormat(colDate, 6, period, "1", 0, "C", false, 0, "")
			pdf.CellFormat(colCompany, 6, exp.Company, "1", 0, "L", false, 0, "")
			pdf.CellFormat(colRole, 6, exp.Position, "1", 1, "L", false, 0, "")
			if exp.Description != "" {
				pdf.SetFont("Helvetica", "I", 8)
				pdf.SetTextColor(80, 80, 80)
				pdf.MultiCell(0, 4, "  "+exp.Description, "LR", "L", false)
				pdf.SetFont("Helvetica", "", 9)
				pdf.SetTextColor(30, 30, 30)
				// close bottom border of description row
				pdf.CellFormat(0, 0, "", "B", 1, "", false, 0, "")
			}
		}
		pdf.Ln(4)
	}

	// ---- Skills ----
	if len(data.Skills) > 0 {
		addJapanSectionHeader(pdf, "Skills / スキル")
		pdf.SetFont("Helvetica", "", 10)
		pdf.SetTextColor(40, 40, 40)
		pdf.MultiCell(0, 5, strings.Join(data.Skills, "  /  "), "1", "L", false)
		pdf.Ln(4)
	}

	// ---- Languages ----
	if len(data.Languages) > 0 {
		addJapanSectionHeader(pdf, "Languages / 語学")
		for _, lang := range data.Languages {
			label := lang.Name
			if lang.Proficiency != "" {
				label += fmt.Sprintf("  (%s)", lang.Proficiency)
			}
			pdf.SetFont("Helvetica", "", 10)
			pdf.SetTextColor(40, 40, 40)
			pdf.CellFormat(0, 5, "  • "+label, "", 1, "L", false, 0, "")
		}
	}

	return pdfToBytes(pdf)
}

func addJapanSectionHeader(pdf *fpdf.Fpdf, title string) {
	pdf.SetFont("Helvetica", "B", 11)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFillColor(40, 80, 160)
	pdf.CellFormat(0, 7, "  "+title, "", 1, "L", true, 0, "")
	pdf.Ln(1)
}

func addJapanField(pdf *fpdf.Fpdf, pageW float64, label, value string) {
	if value == "" {
		return
	}
	labelW := 50.0
	valueW := pageW - labelW
	pdf.SetFont("Helvetica", "B", 9)
	pdf.SetFillColor(240, 244, 252)
	pdf.SetTextColor(40, 80, 160)
	pdf.CellFormat(labelW, 6, "  "+label, "1", 0, "L", true, 0, "")
	pdf.SetFont("Helvetica", "", 9)
	pdf.SetFillColor(255, 255, 255)
	pdf.SetTextColor(30, 30, 30)
	pdf.CellFormat(valueW, 6, "  "+value, "1", 1, "L", false, 0, "")
}

// pdfToBytes converts the FPDF object to a byte slice.
func pdfToBytes(pdf *fpdf.Fpdf) ([]byte, error) {
	if err := pdf.Error(); err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
