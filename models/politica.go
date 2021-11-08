package models

type Politica struct {
	Version   string
	Id        string
	Statement []struct {
		Sid       string
		Effect    string
		Principal struct {
			AWS string
		}
		Action    []string
		Resource  string
		Condition map[string]map[string]string
	}
}
