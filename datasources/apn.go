package datasources

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	model "github.com/syobonaction/fur_lange/models"
)

type Response struct {
	Message Message `json:"message"`
}

type Message struct {
	Results       []*Result `json:"results"`
	Refiners      []string  `json:"refiners"`
	LocationTypes []string  `json:"locationTypes"`
	OfferingTypes []string  `json:"offeringTypes"`
	Total         int       `json:"total"`
	Fuzziness     bool      `json:"fuzziness"`
}

type Result struct {
	ID     string  `json:"_id"`
	Source *Source `json:"_source"`
}

type Source struct {
	CustomerType                  string           `json:"customer_type"`
	AwsCertficationsCount         int              `json:"aws_certifications_count"`
	PartnerActive                 bool             `json:"partner_active"`
	TechnologyExpertise           []string         `json:"technology_expertise"`
	Industry                      []string         `json:"industry"`
	Refiners                      []string         `json:"refiners"`
	TargetClientBase              []string         `json:"target_client_base"`
	NumberOfEmployees             int              `json:"numberofemployees"`
	CustomerLaunchesCount         int              `json:"customer_launches_count"`
	Segment                       []string         `json:"segment"`
	DownloadUrl                   string           `json:"download_url"`
	PartnerPath                   *Path            `json:"partner_path"`
	ServicesCount                 int              `json:"services_count"`
	LiteralName                   string           `json:"literal_name"`
	ProfessionalServiceTypes      []string         `json:"professional_service_types"`
	SolutionCount                 int              `json:"solution_count"`
	CurrentProgramStatus          string           `json:"current_program_status"`
	SocioEconomicCategoriesCount  int              `json:"socio_economic_categories_count"`
	IsSaasVendor                  string           `json:"is_saas_vendor"`
	ProgramsCount                 int              `json:"programs_count"`
	Website                       string           `json:"website"`
	CompetenciesCount             int              `json:"competencies_count"`
	PublicSectorCategoriesCount   int              `json:"public_sector_categories_count"`
	CompetencyMembership          []string         `json:"competency_membership"`
	PublicSectorProgramCategories []string         `json:"public_sector_program_categories"`
	ProgramMembership             []string         `json:"program_membership"`
	PublicSectorContractCount     int              `json:"public_sector_contract_count"`
	Domain                        []string         `json:"domain"`
	ServiceMembership             []string         `json:"service_membership"`
	ReferenceCount                int              `json:"reference_count"`
	AwsCertfications              []string         `json:"aws_certifications"`
	UseCaseExpertise              []string         `json:"use_case_expertise"`
	BriefDescription              string           `json:"brief_description"`
	Description                   string           `json:"description"`
	References                    []*Reference     `json:"references"`
	Solutions                     []*Solution      `json:"solutions"`
	Name                          string           `json:"name"`
	NameAka                       []string         `json:"name_aka"`
	OfficeAddress                 []*OfficeAddress `json:"office_address"`
	OfficeAddressAka              []*OfficeAddress `json:"office_address_aka"`
	Timestamp                     string           `json:"timestamp"`
	Language                      string           `json:"language"`
	PartnerValidated              bool             `json:"partner_validated"`
	SolutionsNested               []*Solution      `json:"solution_nested"`
	ReferencesNested              []*Reference     `json:"references_nested"`
	SolutionsSolutionCount        int              `json:"solutions_solution_count"`
	SolutionsPracticeCount        int              `json:"solutions_practice_count"`
	ReferenceCasestudyCount       int              `json:"references_casestudy_count"`
	ReferenceReferenceCount       int              `json:"references_reference_count"`
}

type Path struct {
	PathDetail  []*PathDetail `json:"path_detail"`
	PrimaryPath string        `json:"primary_path"`
}

type PathDetail struct {
	PathTier  string `json:"path_detail"`
	PathStage string `json:"path_stage"`
	PathName  string `json:"path_name"`
}

type Reference struct {
	ReferenceID    string   `json:"reference_id"`
	Refiners       []string `json:"refiners"`
	Description    string   `json:"description"`
	ApprovalDate   string   `json:"approval_date"`
	ReferenceUrl   string   `json:"reference_url"`
	ExpirationDate string   `json:"expiration_date"`
	RecordType     string   `json:"record_type"`
	CustomerName   string   `json:"customer_types"`
	Title          string   `json:"title"`
}

type Solution struct {
	SolutionID              string   `json:"solution_id"`
	ValidationLevelDetailed string   `json:"validation_level_detailed"`
	OfferingTypeRaw         string   `json:"offering_type_raw"`
	OfferingStatus          string   `json:"offering_status"`
	Refiners                []string `json:"refiners"`
	Availability            string   `json:"availability"`
	RecordType              string   `json:"record_type"`
	SolutionUrlDup          string   `json:"solution_url_dup"`
	SolutionName            []string `json:"solution_name"`
	OfferingType            string   `json:"offering_type"`
	CreatedDate             string   `json:"created_date"`
	SolutionUrl             string   `json:"solution_url"`
	ValidationLevel         string   `json:"validation_level"`
	Proposition             string   `json:"proposition"`
	Description             string   `json:"description"`
	Title                   string   `json:"title"`
	BdUseCase               []string `json:"bd_use_case"`
}

type OfficeAddress struct {
	Country      string   `json:"country"`
	City         string   `json:"city"`
	Street       string   `json:"street"`
	Postalcode   string   `json:"postalcode"`
	State        string   `json:"state"`
	LocationType []string `json:"location_type"`
	Latlon       *Latlon  `json:"latlon"`
}

type Latlon struct {
	Lat float32 `json:"lat"`
	Lon float32 `json:"lon"`
}

func GetAWSPartners(start int, size int) []*model.Partner {
	starting_record := strconv.Itoa(start)
	record_size := strconv.Itoa(size)
	url := fmt.Sprintf("https://api.finder.partners.aws.a2z.com/search?locale=en&from=%s&size=%s", starting_record, record_size)
	fmt.Println(url)
	resp, err := http.Get(url)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Print(err.Error())
	}

	var response Response

	error := json.Unmarshal(body, &response)

	if error != nil {
		fmt.Println(error)
	}

	return formatResults(response)
}

func formatResults(results Response) []*model.Partner {
	partners := []*model.Partner{}

	for _, result := range results.Message.Results {
		s := result.Source

		addresses := []*model.Address{}
		for _, a := range s.OfficeAddress {
			address := &model.Address{
				Country:      a.Country,
				City:         a.City,
				Street:       a.Street,
				Postalcode:   a.Postalcode,
				State:        a.State,
				LocationType: a.LocationType,
			}

			addresses = append(addresses, address)
		}

		pathDetail := []*model.Path{}
		for _, p := range s.PartnerPath.PathDetail {
			path := &model.Path{
				Name:  p.PathName,
				Tier:  p.PathTier,
				Stage: p.PathStage,
			}

			pathDetail = append(pathDetail, path)
		}

		solutions := []*model.Solution{}
		for _, l := range s.Solutions {
			solution := &model.Solution{
				Name:        l.SolutionName,
				Status:      l.OfferingStatus,
				Type:        l.OfferingType,
				Title:       l.Title,
				Description: l.Description,
				Level:       l.ValidationLevel,
				Date:        l.CreatedDate,
			}

			solutions = append(solutions, solution)
		}

		partner := &model.Partner{
			Name:                 s.Name,
			LiteralName:          s.LiteralName,
			Address:              addresses,
			CustomerType:         s.CustomerType,
			Industry:             s.Industry,
			PrimaryPath:          s.PartnerPath.PrimaryPath,
			TechnologyExpertise:  s.TechnologyExpertise,
			Paths:                pathDetail,
			Certifications:       s.AwsCertfications,
			EmployeeCount:        s.NumberOfEmployees,
			Solutions:            solutions,
			ProgramStatus:        s.CurrentProgramStatus,
			ProgramMembership:    s.ProgramMembership,
			CompetencyMembership: s.CompetencyMembership,
			ServiceMembership:    s.ServiceMembership,
			Launches:             s.CustomerLaunchesCount,
			Status:               s.CurrentProgramStatus,
			Description:          s.Description,
			Active:               s.PartnerActive,
			Validated:            s.PartnerValidated,
		}

		partners = append(partners, partner)
	}

	return partners
}
