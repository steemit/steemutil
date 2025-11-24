package broadcast

import (
	"testing"
)

func TestOperationsDataNotEmpty(t *testing.T) {
	if len(OperationsData) == 0 {
		t.Error("OperationsData should not be empty")
	}
}

func TestOperationsDataStructure(t *testing.T) {
	for i, op := range OperationsData {
		if op.Operation == "" {
			t.Errorf("Operation %d: Operation field is empty", i)
		}
		if len(op.Roles) == 0 {
			t.Errorf("Operation %d: Roles field is empty", i)
		}
	}
}

func TestOperationsDataHasCommonOperations(t *testing.T) {
	operations := make(map[string]bool)
	for _, op := range OperationsData {
		operations[op.Operation] = true
	}

	expectedOperations := []string{
		"vote",
		"comment",
		"transfer",
		"account_create",
		"account_update",
		"custom_json",
		"comment_options",
	}

	for _, op := range expectedOperations {
		if !operations[op] {
			t.Errorf("OperationsData missing operation: %s", op)
		}
	}
}

func TestOperationsDataRoles(t *testing.T) {
	validRoles := map[string]bool{
		"posting": true,
		"active":  true,
		"owner":   true,
	}

	for i, op := range OperationsData {
		for _, role := range op.Roles {
			if !validRoles[role] {
				t.Errorf("Operation %d: Invalid role %s", i, role)
			}
		}
	}
}

func TestBroadcastOperationStructure(t *testing.T) {
	op := BroadcastOperation{
		Roles:     []string{"posting", "active", "owner"},
		Operation: "vote",
		Params:    []string{"voter", "author", "permlink", "weight"},
	}

	if len(op.Roles) != 3 {
		t.Error("BroadcastOperation.Roles not set correctly")
	}
	if op.Operation != "vote" {
		t.Error("BroadcastOperation.Operation not set correctly")
	}
	if len(op.Params) != 4 {
		t.Error("BroadcastOperation.Params not set correctly")
	}
}

func TestOperationsDataNoSMT(t *testing.T) {
	// Verify no SMT operations are included
	for _, op := range OperationsData {
		if len(op.Operation) >= 4 && op.Operation[:4] == "smt_" {
			t.Errorf("OperationsData contains SMT operation: %s", op.Operation)
		}
	}
}

