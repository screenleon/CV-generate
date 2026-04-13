package models

// CVData holds all information needed to generate a CV.
type CVData struct {
	// Personal information
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	Website string `json:"website,omitempty"`

	// Japan-specific fields (used in 履歴書 template)
	BirthDate   string `json:"birth_date,omitempty"`
	Nationality string `json:"nationality,omitempty"`
	Gender      string `json:"gender,omitempty"`

	// Professional summary
	Summary string `json:"summary,omitempty"`

	// Work experience
	Experience []Experience `json:"experience,omitempty"`

	// Education
	Education []Education `json:"education,omitempty"`

	// Skills
	Skills []string `json:"skills,omitempty"`

	// Categorized skills (for 職務経歴書 template)
	TechStack *TechStack `json:"tech_stack,omitempty"`

	// Languages spoken
	Languages []Language `json:"languages,omitempty"`

	// Generation options
	Template string `json:"template"` // "simple", "japan" (履歴書), or "shokumu" (職務経歴書)
	Format   string `json:"format"`   // "pdf" or "word"
}

// Experience represents a single work experience entry.
type Experience struct {
	Company     string   `json:"company"`
	Position    string   `json:"position"`
	StartDate   string   `json:"start_date"`
	EndDate     string   `json:"end_date"` // empty means "Present"
	Description string   `json:"description,omitempty"`
	Project     string   `json:"project,omitempty"`     // Project name (for 職務経歴書)
	Role        string   `json:"role,omitempty"`        // Detailed role (for 職務経歴書)
	TechStack   []string `json:"tech_stack,omitempty"`  // Technologies used (for 職務経歴書)
}

// Education represents a single education entry.
type Education struct {
	Institution string `json:"institution"`
	Degree      string `json:"degree"`
	Field       string `json:"field,omitempty"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}

// Language represents a spoken language with proficiency level.
type Language struct {
	Name        string `json:"name"`
	Proficiency string `json:"proficiency,omitempty"` // e.g., "Native", "Business", "Conversational"
}

// TechStack represents categorized technical skills (for 職務経歴書 template).
type TechStack struct {
	Languages      []string `json:"languages,omitempty"`      // Programming languages (e.g., Java, Go, TypeScript)
	Frameworks     []string `json:"frameworks,omitempty"`     // Frameworks (e.g., Spring Boot, React)
	Databases      []string `json:"databases,omitempty"`      // Databases (e.g., PostgreSQL, MongoDB)
	Infrastructure []string `json:"infrastructure,omitempty"` // Infrastructure (e.g., AWS, Docker, Kubernetes)
	Tools          []string `json:"tools,omitempty"`          // Tools (e.g., Git, Jenkins, JIRA)
}
