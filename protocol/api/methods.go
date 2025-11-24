package api

// APIMethod represents a Steem API method definition.
type APIMethod struct {
	// API is the API name, e.g., "database_api", "network_broadcast_api"
	API string `json:"api"`

	// Method is the method name, e.g., "get_block", "broadcast_transaction"
	Method string `json:"method"`

	// MethodName is the Go method name (camelCase version of Method).
	// If not set, it will be generated from Method.
	MethodName string `json:"method_name,omitempty"`

	// Params is the list of parameter names for this method.
	// If empty, the method takes no parameters.
	Params []string `json:"params,omitempty"`

	// IsObject indicates whether the params should be passed as an object
	// instead of individual parameters.
	IsObject bool `json:"is_object,omitempty"`
}
