package models

// request for check api
type CheckRequest struct {
	// current k8s resource version
	Version Version `json:"version"`
	// source configuration
	Source Source `json:"source"`
}

// request for in api
type InRequest struct {
	// current k8s resource version
	Version Version `json:"version"`
	// source configuration
	Source Source `json:"source"`
	// param configuration
	InParams `json:"params"`
}

// response for in api
type InResponse struct {
	// current k8s resource version
	Version Version `json:"version"`
	// resource metadata
	Metadata []Metadata `json:"metadata"`
}

// param for in api. nothing at the moment
type InParams struct {
}

// request for out api
type OutRequest struct {
	// source configuration
	Source Source `json:"source"`
	// param configuration
	Params OutParams `json:"params"`
}

// response for out api
type OutResponse struct {
	// latest k8s resource version
	Version Version `json:"version"`
	// resource metadata
	Metadata []Metadata `json:"metadata"`
}

// param for out api
type OutParams struct {
	// manifest paths
	Paths []string `json:"paths"`
	// if true, this deployment is executed by kustomize
	Kustomize bool `json:"kustomize"`
	// wait time seconds until ready
	StatusCheckTimeout int32 `json:"status_check_timeout"`
	// if true, delete resources
	Delete bool `json:"delete"`
	// if true, rollback to previous deployment
	Undo bool `json:"undo"`
	// kubectl timeout seconds
	CommandTimeout int32 `json:"command_timeout"`
	// if true, execute as dry-run=server
	ServerDryRun bool `json:"server_dry_run"`
	// if true, run diff command instead of apply
	Diff bool `json:"diff"`
}

// concourse metadata
type Metadata struct {
	// name of metadata
	Name string `json:"name"`
	// value of metadata
	Value string `json:"value"`
}

// represent k8s resource version
type Version struct {
	// target k8s revision
	Revision string `json:"ref"`
}

// source configuration
type Source struct {
	// k8s API server URL
	ApiServerUrl string `json:"api_server_url"`
	// k8s API server certificate
	ApiServerCA string `json:"api_server_cert"`
	// client certificate
	ClientCert string `json:"client_cert"`
	// client private key
	ClientKey string `json:"client_key"`
	// client token for k8s service account
	ClientToken string `json:"client_token"`
	// if true, skip unsecured TLS verification
	SkipTLSVerify bool `json:"skip_tls_verify"`
	// target namespace
	Namespace string `json:"namespace"`
	// if true, write verbose info
	Debug bool `json:"debug"`
	// monitoring resources
	WatchResources []WatchResource `json:"watch_resources"`
	// represent raw kubeconfig object
	Kubeconfig string `json:"kubeconfig"`
}

// monitoring resource identifier
type WatchResource struct {
	// resource kind
	Kind string `json:"kind"`
	// resource name
	Name string `json:"name"`
}
