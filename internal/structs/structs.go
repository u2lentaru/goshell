package structs

//Command struct
type Command struct {
	Id          int    `json:"id"`
	CommandText string `json:"commandtext"`
	ResultText  string `json:"resulttext"`
}

//CommandAdd struct
type CommandAdd struct {
	CommandText string `json:"commandtext"`
	ResultText  string `json:"resulttext"`
}

//Command_count  struct
type Command_count struct {
	Values []Command `json:"values"`
	Count  int       `json:"count"`
}
