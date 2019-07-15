package tools

type ModulePublic struct {
	Path     string `json:",omitempty"`
	Version  string `json:",omitempty"`
	Main     bool   `json:",omitempty"`
	Indirect bool   `json:",omitempty"`
}
