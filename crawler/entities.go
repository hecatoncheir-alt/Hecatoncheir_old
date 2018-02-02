package crawler

// Company type for parse
type Company struct {
	ID   string
	IRI  string
	Name string
}

// Category of company structure for parse
type Category struct {
	ID   string
	Name string
}

// City of products for parse
type City struct {
	ID   string
	Name string
}
