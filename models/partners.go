package models

type Partners struct {
	Partners []Partner
	Total    int
}

type Partner struct {
	Name                 string
	LiteralName          string
	Address              []*Address
	CustomerType         string
	Industry             []string
	PrimaryPath          string
	TechnologyExpertise  []string
	Paths                []*Path
	Certifications       []string
	EmployeeCount        int
	Solutions            []*Solution
	ProgramStatus        string
	ProgramMembership    []string
	ServiceMembership    []string
	CompetencyMembership []string
	Description          string
	Launches             int
	Status               string
	Active               bool
	Validated            bool
}

type Address struct {
	Country      string
	City         string
	Street       string
	Postalcode   string
	State        string
	LocationType []string
}

type Path struct {
	Name  string
	Tier  string
	Stage string
}

type Solution struct {
	Name        []string
	Status      string
	Type        string
	Title       string
	Description string
	Level       string
	Date        string
}
