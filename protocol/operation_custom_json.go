package protocol

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/steemit/steemutil/encoder"
)

const (
	TypeFollow = "follow"
)

var customJSONDataObjects = map[string]interface{}{
	TypeFollow: &FollowOperation{},
}

// FC_REFLECT( steemit::chain::custom_json_operation,
//             (required_auths)
//             (required_posting_auths)
//             (id)
//             (json) )

// CustomJSONOperation represents custom_json operation data.
type CustomJSONOperation struct {
	RequiredAuths        []string `json:"required_auths"`
	RequiredPostingAuths []string `json:"required_posting_auths"`
	ID                   string   `json:"id"`
	JSON                 string   `json:"json"`
}

func (op *CustomJSONOperation) Type() OpType {
	return TypeCustomJSON
}

func (op *CustomJSONOperation) Data() interface{} {
	return op
}

func (op *CustomJSONOperation) UnmarshalData() (interface{}, error) {
	// Get the corresponding data object template.
	template, ok := customJSONDataObjects[op.ID]
	if !ok {
		// In case there is no corresponding template, return nil.
		return nil, nil
	}

	// Clone the template.
	opData := reflect.New(reflect.Indirect(reflect.ValueOf(template)).Type()).Interface()

	// Prepare the whole operation tuple.
	var bodyReader io.Reader
	if op.JSON[0] == '[' {
		rawTuple := make([]json.RawMessage, 2)
		if err := json.NewDecoder(strings.NewReader(op.JSON)).Decode(&rawTuple); err != nil {
			return nil, errors.Wrapf(err,
				"failed to unmarshal CustomJSONOperation.JSON: \n%v", op.JSON)
		}
		if rawTuple[1] == nil {
			return nil, errors.Errorf("invalid CustomJSONOperation.JSON: \n%v", op.JSON)
		}
		bodyReader = bytes.NewReader([]byte(rawTuple[1]))
	} else {
		bodyReader = strings.NewReader(op.JSON)
	}

	// Unmarshal into the new object instance.
	if err := json.NewDecoder(bodyReader).Decode(opData); err != nil {
		return nil, errors.Wrapf(err,
			"failed to unmarshal CustomJSONOperation.JSON: \n%v", op.JSON)
	}

	return opData, nil
}

// MarshalTransaction implements custom binary serialization for custom_json operation.
// This ensures required_auths and required_posting_auths are sorted before serialization,
// matching steem-js behavior for flat_set serialization.
func (op *CustomJSONOperation) MarshalTransaction(encoderObj *encoder.Encoder) error {
	if op == nil {
		return errors.New("cannot marshal nil CustomJSONOperation")
	}

	// Encode operation type code first (required for all operations)
	if err := encoderObj.EncodeUVarint(uint64(op.Type().Code())); err != nil {
		return errors.Wrap(err, "failed to encode operation type code")
	}

	// Sort required_auths (flat_set serialization requires sorted order)
	requiredAuths := make([]string, len(op.RequiredAuths))
	copy(requiredAuths, op.RequiredAuths)
	sort.Strings(requiredAuths)

	// Encode required_auths: varint32 length + each account name (string)
	if err := encoderObj.EncodeUVarint(uint64(len(requiredAuths))); err != nil {
		return errors.Wrap(err, "failed to encode required_auths length")
	}
	for _, account := range requiredAuths {
		if err := encoderObj.Encode(account); err != nil {
			return errors.Wrapf(err, "failed to encode required_auths account: %s", account)
		}
	}

	// Sort required_posting_auths (flat_set serialization requires sorted order)
	requiredPostingAuths := make([]string, len(op.RequiredPostingAuths))
	copy(requiredPostingAuths, op.RequiredPostingAuths)
	sort.Strings(requiredPostingAuths)

	// Encode required_posting_auths: varint32 length + each account name (string)
	if err := encoderObj.EncodeUVarint(uint64(len(requiredPostingAuths))); err != nil {
		return errors.Wrap(err, "failed to encode required_posting_auths length")
	}
	for _, account := range requiredPostingAuths {
		if err := encoderObj.Encode(account); err != nil {
			return errors.Wrapf(err, "failed to encode required_posting_auths account: %s", account)
		}
	}

	// Encode id (string)
	if err := encoderObj.Encode(op.ID); err != nil {
		return errors.Wrap(err, "failed to encode id")
	}

	// Encode json (string)
	if err := encoderObj.Encode(op.JSON); err != nil {
		return errors.Wrap(err, "failed to encode json")
	}

	return nil
}
