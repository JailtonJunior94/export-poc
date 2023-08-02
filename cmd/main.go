package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jailtonjunior94/export-poc/configs"
	"github.com/jailtonjunior94/export-poc/internal/entities"
	"github.com/jailtonjunior94/export-poc/pkg/excel"
	migration "github.com/jailtonjunior94/export-poc/pkg/migrate"
	database "github.com/jailtonjunior94/export-poc/pkg/postgres"
)

func main() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	dbConn, err := database.NewPostgresDatabase(config)
	if err != nil {
		panic(err)
	}
	defer dbConn.Close()

	migrate, err := migration.NewMigrate(dbConn, config.MigratePath, config.DBName)
	if err != nil {
		panic(err)
	}

	if err = migrate.ExecuteMigration(); err != nil {
		panic(err)
	}

	courses, err := Courses(context.Background(), dbConn)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	provider := excel.NewProvider()
	xls := provider.NewFile(ctx)
	defaultSheet := xls.NewSheet(ctx, "Sheet1")

	export := make([]*CoursesExportRow, len(courses))
	for index, course := range courses {
		export[index] = NewCoursesExport(course.Name, course.Description, course.Category.Name, course.Category.Description)
		row := export[index]

		err := defaultSheet.Write(ctx, *row)
		if err != nil {
			panic(err)
		}
	}

	buffer, err := xls.WriteToBuffer(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(buffer)

	// category := entities.Category{
	// 	ID:          uuid.New().String(),
	// 	Name:        "Category",
	// 	Description: "Description",
	// }

	// _, err = dbConn.Exec("insert into categories values ($1, $2, $3)", category.ID, category.Name, category.Description)
	// if err != nil {
	// 	panic(err)
	// }

	// for i := 0; i < 1000000; i++ {
	// 	course := entities.Course{
	// 		ID:          uuid.New().String(),
	// 		Name:        fmt.Sprintf("Course %d", i),
	// 		Description: fmt.Sprintf("Description %d", i),
	// 	}

	// 	_, err = dbConn.Exec("insert into courses values ($1, $2, $3, $4)", course.ID, category.ID, course.Name, course.Description)
	// 	if err != nil {
	// 		log.Println(err)
	// 		continue
	// 	}
	// }
}

func Courses(ctx context.Context, conn *sql.DB) ([]entities.Course, error) {
	query := "SELECT c.id, c.name, c.description, c.category_id, ca.name as category_name, ca.description as category_description FROM courses c JOIN categories ca ON c.category_id = ca.id"
	rows, err := conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var items []entities.Course
	for rows.Next() {
		var i entities.Course
		if err := rows.Scan(&i.ID, &i.Name, &i.Description, &i.Category.ID, &i.Category.Name, &i.Category.Description); err != nil {
			return nil, err
		}
		items = append(items, i)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

type CoursesExportRow struct {
	Course              interface{} `column:"A" header:"Curso"`
	Description         interface{} `column:"B" header:"Descrição do Curso"`
	Category            interface{} `column:"C" header:"Categoria do Curso"`
	CategoryDescription interface{} `column:"D" header:"Descrição da Categoria"`
}

func NewCoursesExport(course, description, category, categoryDescription string) *CoursesExportRow {
	return &CoursesExportRow{
		Course:              course,
		Description:         description,
		Category:            category,
		CategoryDescription: categoryDescription,
	}
}
