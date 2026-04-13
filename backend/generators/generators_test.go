package generators_test

import (
	"bytes"
	"testing"

	"github.com/screenleon/cv-generate/generators"
	"github.com/screenleon/cv-generate/models"
)

// sampleCV returns a fully-populated CVData for use in tests.
func sampleCV(template, format string) *models.CVData {
	return &models.CVData{
		Name:        "Tanaka Yuki",
		Email:       "yuki@example.com",
		Phone:       "+81-90-1234-5678",
		Address:     "Tokyo, Japan",
		Website:     "https://example.com",
		BirthDate:   "1990-04-01",
		Gender:      "Female",
		Nationality: "Japanese",
		Summary:     "Experienced engineer with a passion for Go and distributed systems.",
		Experience: []models.Experience{
			{
				Company:     "Acme Corp",
				Position:    "Software Engineer",
				StartDate:   "2018-04",
				EndDate:     "2022-03",
				Description: "Built microservices in Go.",
			},
			{
				Company:   "Beta Inc",
				Position:  "Senior Engineer",
				StartDate: "2022-04",
				EndDate:   "",
			},
		},
		Education: []models.Education{
			{
				Institution: "University of Tokyo",
				Degree:      "Bachelor of Engineering",
				Field:       "Computer Science",
				StartDate:   "2008-04",
				EndDate:     "2012-03",
			},
		},
		Skills:   []string{"Go", "Docker", "Kubernetes", "SQL"},
		Languages: []models.Language{
			{Name: "Japanese", Proficiency: "Native"},
			{Name: "English", Proficiency: "Business"},
		},
		Template: template,
		Format:   format,
	}
}

// ── PDF tests ─────────────────────────────────────────────

func TestGeneratePDF_Simple(t *testing.T) {
	data := sampleCV("simple", "pdf")
	got, err := generators.GeneratePDF(data)
	if err != nil {
		t.Fatalf("GeneratePDF(simple) returned error: %v", err)
	}
	if len(got) == 0 {
		t.Fatal("GeneratePDF(simple) returned empty bytes")
	}
	// PDF files start with %PDF
	if !bytes.HasPrefix(got, []byte("%PDF")) {
		t.Errorf("GeneratePDF(simple) output does not start with %%PDF header")
	}
}

func TestGeneratePDF_Japan(t *testing.T) {
	data := sampleCV("japan", "pdf")
	got, err := generators.GeneratePDF(data)
	if err != nil {
		t.Fatalf("GeneratePDF(japan) returned error: %v", err)
	}
	if len(got) == 0 {
		t.Fatal("GeneratePDF(japan) returned empty bytes")
	}
	if !bytes.HasPrefix(got, []byte("%PDF")) {
		t.Errorf("GeneratePDF(japan) output does not start with %%PDF header")
	}
}

func TestGeneratePDF_MinimalData(t *testing.T) {
	data := &models.CVData{Name: "Jane Doe", Template: "simple", Format: "pdf"}
	got, err := generators.GeneratePDF(data)
	if err != nil {
		t.Fatalf("GeneratePDF with minimal data returned error: %v", err)
	}
	if !bytes.HasPrefix(got, []byte("%PDF")) {
		t.Errorf("GeneratePDF minimal output does not start with %%PDF header")
	}
}

// ── Word tests ────────────────────────────────────────────

func TestGenerateWord_Simple(t *testing.T) {
	data := sampleCV("simple", "word")
	got, err := generators.GenerateWord(data)
	if err != nil {
		t.Fatalf("GenerateWord(simple) returned error: %v", err)
	}
	if len(got) == 0 {
		t.Fatal("GenerateWord(simple) returned empty bytes")
	}
	// DOCX files are ZIP archives; ZIP magic bytes are PK\x03\x04
	if !bytes.HasPrefix(got, []byte("PK\x03\x04")) {
		t.Errorf("GenerateWord(simple) output does not start with ZIP/PK header")
	}
}

func TestGenerateWord_Japan(t *testing.T) {
	data := sampleCV("japan", "word")
	got, err := generators.GenerateWord(data)
	if err != nil {
		t.Fatalf("GenerateWord(japan) returned error: %v", err)
	}
	if len(got) == 0 {
		t.Fatal("GenerateWord(japan) returned empty bytes")
	}
	if !bytes.HasPrefix(got, []byte("PK\x03\x04")) {
		t.Errorf("GenerateWord(japan) output does not start with ZIP/PK header")
	}
}

func TestGenerateWord_MinimalData(t *testing.T) {
	data := &models.CVData{Name: "Jane Doe", Template: "simple", Format: "word"}
	got, err := generators.GenerateWord(data)
	if err != nil {
		t.Fatalf("GenerateWord with minimal data returned error: %v", err)
	}
	if !bytes.HasPrefix(got, []byte("PK\x03\x04")) {
		t.Errorf("GenerateWord minimal output does not start with ZIP/PK header")
	}
}

// ── XML escape test ───────────────────────────────────────

func TestGenerateWord_XMLSpecialChars(t *testing.T) {
	data := &models.CVData{
		Name:     "O'Brien & <Test>",
		Email:    "test@example.com",
		Summary:  `Summary with "quotes" and <tags> & ampersands`,
		Template: "simple",
		Format:   "word",
	}
	_, err := generators.GenerateWord(data)
	if err != nil {
		t.Fatalf("GenerateWord with special XML chars returned error: %v", err)
	}
}
