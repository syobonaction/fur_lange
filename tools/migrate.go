package tools

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/syobonaction/fur_lange/configs"
	"github.com/syobonaction/fur_lange/models"
)

func MigratePgsql(args ...string) error {
	partners := GetPartners()

	InsertPgsql(partners)

	return nil
}

func GetPartners() []models.Partner {
	url := fmt.Sprintf("http://%s/partners", configs.EnvMongoURL())
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

	var response models.Response

	error := json.Unmarshal(body, &response)

	if error != nil {
		fmt.Println(error)
	}

	return response.Data.Partners
}

func InsertPgsql(partners []models.Partner) {
	pgsqlinfo := configs.EnvPgURI()
	db, err := sql.Open("postgres", pgsqlinfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// ID COUNTERS
	addressID := 0
	expertiseID := 0
	pathID := 0
	certificationID := 0
	solutionID := 0
	programID := 0
	serviceID := 0
	competencyID := 0

	for i, partner := range partners {
		sqlStatement := `
			INSERT INTO partners ("ID", name, literalname, customertype, primarypath, programstatus, description, launches, status, active, validated, industry)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
			RETURNING "ID"
		`
		id := i
		var industry string
		if len(partner.Industry) > 0 {
			industry = partner.Industry[0]
		}
		err = db.QueryRow(sqlStatement, id, partner.Name, partner.LiteralName, partner.CustomerType, partner.PrimaryPath, partner.ProgramStatus, partner.Description, partner.Launches, partner.Status, partner.Active, partner.Validated, industry).Scan(&id)
		if err != nil {
			panic(err)
		}

		// ADDRESS ENTRY
		for _, address := range partner.Address {
			sqlStatement := `
				INSERT INTO partners_address ("ID", parent_fk, country, city, street, postalcode, state, locationtype)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
				RETURNING "ID"
			`
			var locationType string
			if len(address.LocationType) > 0 {
				locationType = address.LocationType[0]
			}
			err = db.QueryRow(sqlStatement, addressID, i, address.Country, address.City, address.Street, address.Postalcode, address.State, locationType).Scan(&id)
			if err != nil {
				panic(err)
			}
			addressID++
			fmt.Println("Partner address inserted.")
		}

		// TECHNOLOGY EXPERTISE ENTRY
		if len(partner.TechnologyExpertise) > 0 {
			for _, expertise := range partner.TechnologyExpertise {
				sqlStatement := `
					INSERT INTO partners_expertise ("ID", parent_fk, expertise)
					VALUES ($1, $2, $3)
					RETURNING "ID"
				`
				err = db.QueryRow(sqlStatement, expertiseID, i, expertise).Scan(&id)
				if err != nil {
					panic(err)
				}
				expertiseID++
				fmt.Println("Partner expertise inserted.")
			}
		}

		// PATHS ENTRY
		if len(partner.Paths) > 0 {
			for _, path := range partner.Paths {
				sqlStatement := `
					INSERT INTO partners_path ("ID", parent_fk, name, tier, stage)
					VALUES ($1, $2, $3, $4, $5)
					RETURNING "ID"
				`
				err = db.QueryRow(sqlStatement, pathID, i, path.Name, path.Tier, path.Stage).Scan(&id)
				if err != nil {
					panic(err)
				}
				pathID++
				fmt.Println("Partner path inserted.")
			}
		}

		// CERTIFICATIONS ENTRY
		if len(partner.Certifications) > 0 {
			for _, certification := range partner.Certifications {
				sqlStatement := `
					INSERT INTO partners_certification ("ID", parent_fk, name)
					VALUES ($1, $2, $3)
					RETURNING "ID"
				`
				err = db.QueryRow(sqlStatement, certificationID, i, certification).Scan(&id)
				if err != nil {
					panic(err)
				}
				certificationID++
				fmt.Println("Partner certification inserted.")
			}
		}

		// SOLUTIONS ENTRY
		if len(partner.Solutions) > 0 {
			for _, solution := range partner.Solutions {
				sqlStatement := `
					INSERT INTO partners_solution ("ID", parent_fk, name, status, type, title, description, level, date)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
					RETURNING "ID"
				`
				var solutionName string
				if len(solution.Name) > 0 {
					solutionName = solution.Name[0]
				}
				err = db.QueryRow(sqlStatement, solutionID, i, solutionName, solution.Status, solution.Type, solution.Title, solution.Description, solution.Level, solution.Date).Scan(&id)
				if err != nil {
					panic(err)
				}
				solutionID++
				fmt.Println("Partner solution inserted.")
			}
		}

		// PROGRAM MEMBERSHIP ENTRY
		if len(partner.ProgramMembership) > 0 {
			for _, program := range partner.ProgramMembership {
				sqlStatement := `
					INSERT INTO partners_program ("ID", parent_fk, name)
					VALUES ($1, $2, $3)
					RETURNING "ID"
				`
				err = db.QueryRow(sqlStatement, programID, i, program).Scan(&id)
				if err != nil {
					panic(err)
				}
				programID++
				fmt.Println("Partner program inserted.")
			}
		}

		// SERVICE MEMBERSHIP ENTRY
		if len(partner.ServiceMembership) > 0 {
			for _, service := range partner.ServiceMembership {
				sqlStatement := `
					INSERT INTO partners_service ("ID", parent_fk, name)
					VALUES ($1, $2, $3)
					RETURNING "ID"
				`
				err = db.QueryRow(sqlStatement, serviceID, i, service).Scan(&id)
				if err != nil {
					panic(err)
				}
				serviceID++
				fmt.Println("Partner service inserted.")
			}
		}

		// COMPETENCY MEMBERSHIP ENTRY
		if len(partner.CompetencyMembership) > 0 {
			for _, competency := range partner.CompetencyMembership {
				sqlStatement := `
					INSERT INTO partners_competency ("ID", parent_fk, name)
					VALUES ($1, $2, $3)
					RETURNING "ID"
				`
				err = db.QueryRow(sqlStatement, competencyID, i, competency).Scan(&id)
				if err != nil {
					panic(err)
				}
				competencyID++
				fmt.Println("Partner competency inserted.")
			}
		}

		fmt.Println("New partner record added!")
	}

	fmt.Println("Migrate completed!")
}
