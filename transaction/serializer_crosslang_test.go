package transaction

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/steemit/steemutil/protocol"
)

// TestSerializerCrossLangFixtures replays the golden fixtures produced by
// steem-js's serializer (and byte-verified against the C++ FC_REFLECT field
// order) through steemutil's encoder. This is the regression net for the
// class of bug where a field is silently encoded with the wrong binary type
// (e.g. daily_pay as a string instead of an asset) — the digest is then wrong
// and the chain rejects the signature, but nothing fails locally.
//
// The fixtures live in the steem-js checkout; the path is passed via env var
// STEEM_JS_FIXTURES since steemutil must not vendor them. When the variable
// is unset (e.g. CI) the test skips — the golden hexes are also asserted from
// the steem-js side in that repository's own test suite.
func TestSerializerCrossLangFixtures(t *testing.T) {
	dir := os.Getenv("STEEM_JS_FIXTURES")
	if dir == "" {
		t.Skip("STEEM_JS_FIXTURES not set; skipping cross-language golden replay")
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Skipf("steem-js fixture dir not available: %v", err)
	}

	// Ops whose Go structs are known to diverge from steem-js's current
	// serializer. pow: the C++ chain reference requires worker as a 33-byte
	// binary public key and input/signature/work as raw byte arrays, while
	// steem-js encodes all inner fields as strings (and adds an optional
	// presence byte to the nonce) — the two encodings are mutually
	// incompatible. steemutil follows the C++ reference, so this replay
	// cannot byte-match steem-js's pow fixture; see the wire-format note on
	// protocol.POW for details. pow is a long-deactivated mining op and no
	// live path serializes it, so the divergence is accepted rather than
	// replicating steem-js's (wrong) encoding here.
	skipOps := map[string]bool{"pow": true}

	checked := 0
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			t.Fatalf("read fixture %s: %v", e.Name(), err)
		}
		var fx struct {
			Name string `json:"name"`
			Tx   struct {
				RefBlockNum    protocol.UInt16 `json:"ref_block_num"`
				RefBlockPrefix protocol.UInt32 `json:"ref_block_prefix"`
				Expiration     *protocol.Time  `json:"expiration"`
				Operations     json.RawMessage `json:"operations"`
				Extensions     []interface{}   `json:"extensions"`
			} `json:"tx"`
			ExpectedHex string `json:"expected_hex"`
		}
		if err := json.Unmarshal(raw, &fx); err != nil {
			t.Fatalf("parse fixture %s: %v", e.Name(), err)
		}

		// operations: [["op_name", {json fields}]]
		var rawOps [][]json.RawMessage
		if err := json.Unmarshal(fx.Tx.Operations, &rawOps); err != nil {
			t.Fatalf("parse operations of %s: %v", fx.Name, err)
		}
		if len(rawOps) != 1 {
			continue // only single-op fixtures are supported by this replay
		}
		var opName string
		if err := json.Unmarshal(rawOps[0][0], &opName); err != nil {
			t.Fatalf("parse op name of %s: %v", fx.Name, err)
		}
		if skipOps[opName] {
			continue
		}
		template, ok := protocol.OperationTemplateByName(opName)
		if !ok {
			continue // op not modeled in Go yet — not this test's concern
		}
		// Some committed fixtures were produced by an older steemutil whose
		// StringInt64Map.MarshalJSON emitted a nil-filled prefix inside
		// key_auths/account_auths arrays. steem-js's cross-lang test filters
		// those nulls on load; do the same here instead of failing.
		opJSON := normalizeAuthNulls(rawOps[0][1])
		if err := json.Unmarshal(opJSON, template); err != nil {
			t.Fatalf("unmarshal op %s of %s: %v", opName, fx.Name, err)
		}
		op := template

		tx := &Transaction{
			RefBlockNum:    fx.Tx.RefBlockNum,
			RefBlockPrefix: fx.Tx.RefBlockPrefix,
			Expiration:     fx.Tx.Expiration,
			Operations:     protocol.Operations{op},
			Extensions:     fx.Tx.Extensions,
		}
		if tx.Extensions == nil {
			tx.Extensions = []interface{}{}
		}
		stx := NewSignedTransaction(tx)
		serialized, err := stx.Serialize()
		if err != nil {
			t.Fatalf("serialize %s: %v", fx.Name, err)
		}
		got := hex.EncodeToString(serialized)
		if got != fx.ExpectedHex {
			t.Errorf("%s (%s) mismatch:\n  expected: %s\n  got:      %s", fx.Name, opName, fx.ExpectedHex, got)
		}
		checked++
	}
	if checked == 0 {
		t.Fatal("no fixtures were checked")
	}
	t.Logf("verified %d fixtures against golden hex", checked)
}

// normalizeAuthNulls decodes the op JSON and drops null entries from
// authority key_auths/account_auths arrays (legacy fixture corruption).
func normalizeAuthNulls(raw json.RawMessage) json.RawMessage {
	var opData map[string]json.RawMessage
	if err := json.Unmarshal(raw, &opData); err != nil {
		return raw
	}
	for _, authKey := range []string{"owner", "active", "posting"} {
		authRaw, ok := opData[authKey]
		if !ok {
			continue
		}
		var auth map[string]json.RawMessage
		if err := json.Unmarshal(authRaw, &auth); err != nil {
			continue
		}
		for _, mapKey := range []string{"key_auths", "account_auths"} {
			listRaw, ok := auth[mapKey]
			if !ok {
				continue
			}
			var list []json.RawMessage
			if err := json.Unmarshal(listRaw, &list); err != nil {
				continue
			}
			filtered := make([]json.RawMessage, 0, len(list))
			for _, item := range list {
				if string(item) == "null" {
					continue
				}
				filtered = append(filtered, item)
			}
			out, err := json.Marshal(filtered)
			if err != nil {
				continue
			}
			auth[mapKey] = out
		}
		out, err := json.Marshal(auth)
		if err != nil {
			continue
		}
		opData[authKey] = out
	}
	out, err := json.Marshal(opData)
	if err != nil {
		return raw
	}
	return out
}
