package models

type CheckRequest struct {
	Version Version `json:"version"`
	Source  Source  `json:"source"`
}

type InRequest struct {
	Version Version `json:"version"`
	Source  Source  `json:"source"`

	InParams `json:"params"`
}
type InResponse struct {
	Version  Version    `json:"version"`
	Metadata []Metadata `json:"metadata"`
}
type InParams struct {
}

type OutRequest struct {
	Source Source    `json:"source"`
	Params OutParams `json:"params"`
}
type OutResponse struct {
	Version  Version    `json:"version"`
	Metadata []Metadata `json:"metadata"`
}
type OutParams struct {
	Paths              []string `json:"paths"`
	Kustomize          bool     `json:"kustomize"`
	StatusCheckTimeout int32    `json:"status_check_timeout"`
	Delete             bool     `json:"delete"`
	Undo               bool     `json:"undo"`
	CommandTimeout     int32    `json:"command_timeout"`
}

type Metadata struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
type Version struct {
	Revision string `json:"ref"`
}

type Source struct {
	ApiServerUrl  string `json:"api_server_url"`
	ApiServerCA   string `json:"api_server_cert"`
	ClientCert    string `json:"client_cert"`
	ClientKey     string `json:"client_key"`
	ClientToken   string `json:"client_token"`
	SkipTLSVerify bool   `json:"skip_tls_verify"`
	Namespace     string `json:"namespace"`
	Debug         bool   `json:"debug"`

	WatchResources []WatchResource `json:"watch_resources"`
	Kubeconfig     string          `json:"kubeconfig"`
}

type WatchResource struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}
