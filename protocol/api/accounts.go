package api

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// AuthorityWeight is the weight assigned to a key or account in a Steem
// authority. Steem encodes it as a 16-bit integer.
type AuthorityWeight uint16

// KeyAuth is a single (public key, weight) entry in an account authority.
//
// On the wire (condenser_api.get_accounts), key_auths is serialized as a JSON
// array of two-element arrays: [["STMxxx", 1], ...]. KeyAuth implements a
// custom unmarshaler that flattens each [key, weight] pair into its fields,
// so consumers can work with a typed slice rather than nested raw arrays.
type KeyAuth struct {
	PubKey string
	Weight AuthorityWeight
}

// UnmarshalJSON parses the wire form ["STMxxx", 1] into KeyAuth.
func (k *KeyAuth) UnmarshalJSON(data []byte) error {
	// A JSON null leaves the value at its zero value (standard json behavior).
	if string(data) == "null" {
		return nil
	}
	var pair []json.RawMessage
	if err := json.Unmarshal(data, &pair); err != nil {
		return errors.Wrap(err, "key_auths entry must be a [key, weight] array")
	}
	if len(pair) != 2 {
		return errors.Errorf("key_auths entry must have exactly 2 elements, got %d", len(pair))
	}
	var key string
	if err := json.Unmarshal(pair[0], &key); err != nil {
		return errors.Wrap(err, "invalid public key in key_auths")
	}
	var weight AuthorityWeight
	if err := json.Unmarshal(pair[1], &weight); err != nil {
		return errors.Wrap(err, "invalid weight in key_auths")
	}
	k.PubKey = key
	k.Weight = weight
	return nil
}

// MarshalJSON emits the wire form ["STMxxx", 1], the inverse of UnmarshalJSON,
// so a round-trip through JSON preserves the condenser_api array-of-pairs shape.
func (k KeyAuth) MarshalJSON() ([]byte, error) {
	return json.Marshal([2]interface{}{k.PubKey, k.Weight})
}

// AccountAuthEntry is a weighted account name in an account authority. On the
// wire, account_auths is [["name", weight], ...], parsed the same way as
// key_auths.
type AccountAuthEntry struct {
	Name   string
	Weight AuthorityWeight
}

// UnmarshalJSON parses the wire form ["name", 1] into AccountAuthEntry.
func (a *AccountAuthEntry) UnmarshalJSON(data []byte) error {
	// A JSON null leaves the value at its zero value (standard json behavior).
	if string(data) == "null" {
		return nil
	}
	var pair []json.RawMessage
	if err := json.Unmarshal(data, &pair); err != nil {
		return errors.Wrap(err, "account_auths entry must be a [name, weight] array")
	}
	if len(pair) != 2 {
		return errors.Errorf("account_auths entry must have exactly 2 elements, got %d", len(pair))
	}
	var name string
	if err := json.Unmarshal(pair[0], &name); err != nil {
		return errors.Wrap(err, "invalid account name in account_auths")
	}
	var weight AuthorityWeight
	if err := json.Unmarshal(pair[1], &weight); err != nil {
		return errors.Wrap(err, "invalid weight in account_auths")
	}
	a.Name = name
	a.Weight = weight
	return nil
}

// MarshalJSON emits the wire form ["name", 1], the inverse of UnmarshalJSON.
func (a AccountAuthEntry) MarshalJSON() ([]byte, error) {
	return json.Marshal([2]interface{}{a.Name, a.Weight})
}

// Authority models a Steem account authority (owner / active / posting).
// weight_threshold is the total weight required to authorize an action under
// this authority.
type Authority struct {
	WeightThreshold uint32             `json:"weight_threshold"`
	AccountAuths    []AccountAuthEntry `json:"account_auths"`
	KeyAuths        []KeyAuth          `json:"key_auths"`
}

// ExtendedAccount models the subset of a condenser_api.get_accounts response
// entry that conveyor reads (see conveyor/src/user-search/user.ts UserAccount).
//
// Reputation is decoded as json.RawMessage because the chain returns it
// inconsistently as either a JSON string (legacy, large number) or a number,
// and consumers must handle both.
type ExtendedAccount struct {
	Name        string          `json:"name"`
	Created     string          `json:"created"`
	Reputation  json.RawMessage `json:"reputation"`
	VotingPower int16           `json:"voting_power"`
	Balance     string          `json:"balance"`
	Posting     Authority       `json:"posting"`
	Active      Authority       `json:"active"`
	Owner       Authority       `json:"owner"`
}
