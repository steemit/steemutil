package broadcast

// BroadcastOperation represents a Steem broadcast operation metadata.
type BroadcastOperation struct {
	// Roles is the list of key roles required to sign this operation.
	// Common roles: "posting", "active", "owner"
	Roles []string `json:"roles"`

	// Operation is the operation name, e.g., "vote", "comment", "transfer"
	Operation string `json:"operation"`

	// Params is the list of parameter names for this operation.
	Params []string `json:"params,omitempty"`
}
