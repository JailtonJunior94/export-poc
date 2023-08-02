package entities

type Category struct {
	ID          string
	Name        string
	Description string
}

type Course struct {
	ID          string
	Name        string
	Description string
	Category    Category
}
