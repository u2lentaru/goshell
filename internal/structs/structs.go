package structs

//Command struct
type Command struct {
	Id          int    `json:"id"`
	CommandText string `json:"commandtext"`
	ScriptText  string `json:"scripttext"`
}

// //CommandAdd struct
// type CommandAdd struct {
// 	CommandText string `json:"commandtext"`
// 	ScriptText  string `json:"scriptttext"`
// }

//Command_count  struct
type Command_count struct {
	Values []Command `json:"values"`
	Count  int       `json:"count"`
}

//Result struct
type Result struct {
	Id     int    `json:"id"`
	IdC    int    `json:"idc"`
	Output string `json:"output"`
	TS     string `json:"ts"`
}

//Result_count  struct
type Result_count struct {
	Values []Result `json:"values"`
	Count  int      `json:"count"`
}
