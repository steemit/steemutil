package encoder

import (
	"encoding/binary"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type TransactionMarshaller interface {
	MarshalTransaction(*Encoder) error
}

type Encoder struct {
	w io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w}
}

func (encoder *Encoder) EncodeVarint(i int64) error {
	if i >= 0 {
		return encoder.EncodeUVarint(uint64(i))
	}

	b := make([]byte, binary.MaxVarintLen64)
	n := binary.PutVarint(b, i)
	return encoder.writeBytes(b[:n])
}

func (encoder *Encoder) EncodeUVarint(i uint64) error {
	b := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(b, i)
	return encoder.writeBytes(b[:n])
}

func (encoder *Encoder) EncodeNumber(v interface{}) error {
	if err := binary.Write(encoder.w, binary.LittleEndian, v); err != nil {
		return errors.Wrapf(err, "encoder: failed to write number: %v", v)
	}
	return nil
}

func (encoder *Encoder) Encode(v interface{}) error {
	// if the Transaction v has MarshalTransaction method
	if marshaller, ok := v.(TransactionMarshaller); ok {
		return marshaller.MarshalTransaction(encoder)
	}

	// Check if it's an Operation interface using reflection to avoid import cycle
	if op := encoder.checkOperation(v); op != nil {
		return encoder.encodeOperation(op)
	}

	switch v := v.(type) {
	case int:
		return encoder.EncodeNumber(v)
	case int8:
		return encoder.EncodeNumber(v)
	case int16:
		return encoder.EncodeNumber(v)
	case int32:
		return encoder.EncodeNumber(v)
	case int64:
		return encoder.EncodeNumber(v)

	case uint:
		return encoder.EncodeNumber(v)
	case uint8:
		return encoder.EncodeNumber(v)
	case uint16:
		return encoder.EncodeNumber(v)
	case uint32:
		return encoder.EncodeNumber(v)
	case uint64:
		return encoder.EncodeNumber(v)

	case string:
		return encoder.encodeString(v)

	case bool:
		// Encode bool as uint8 (0 or 1)
		if v {
			return encoder.EncodeNumber(uint8(1))
		}
		return encoder.EncodeNumber(uint8(0))

	default:
		// Try reflection-based encoding for structs and other types
		return encoder.encodeByReflection(v)
	}
}

func (encoder *Encoder) encodeString(v string) error {
	if err := encoder.EncodeUVarint(uint64(len(v))); err != nil {
		return errors.Wrapf(err, "encoder: failed to write string: %v", v)
	}

	return encoder.writeString(v)
}

func (encoder *Encoder) writeBytes(bs []byte) error {
	if _, err := encoder.w.Write(bs); err != nil {
		return errors.Wrapf(err, "encoder: failed to write bytes: %v", bs)
	}
	return nil
}

// WriteBytes writes raw bytes to the encoder.
func (encoder *Encoder) WriteBytes(bs []byte) error {
	return encoder.writeBytes(bs)
}

func (encoder *Encoder) writeString(s string) error {
	if _, err := io.Copy(encoder.w, strings.NewReader(s)); err != nil {
		return errors.Wrapf(err, "encoder: failed to write string: %v", s)
	}
	return nil
}

// operationInterface represents the Operation interface methods without importing protocol package
type operationInterface struct {
	getTypeCode func() uint64
	getData     func() interface{}
}

// checkOperation checks if v implements Operation interface using reflection
func (encoder *Encoder) checkOperation(v interface{}) *operationInterface {
	rv := reflect.ValueOf(v)
	if !rv.IsValid() {
		return nil
	}

	// Check if v has Type() and Data() methods
	typeMethod := rv.MethodByName("Type")
	dataMethod := rv.MethodByName("Data")

	if !typeMethod.IsValid() || !dataMethod.IsValid() {
		return nil
	}

	// Check if Type() returns something with Code() method
	typeResult := typeMethod.Call(nil)
	if len(typeResult) == 0 {
		return nil
	}

	typeValue := typeResult[0]
	codeMethod := typeValue.MethodByName("Code")
	if !codeMethod.IsValid() {
		return nil
	}

	// Create closure to get type code
	getTypeCode := func() uint64 {
		codeResult := codeMethod.Call(nil)
		if len(codeResult) == 0 {
			return 0
		}
		codeValue := codeResult[0]
		switch codeValue.Kind() {
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return codeValue.Uint()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return uint64(codeValue.Int())
		default:
			return 0
		}
	}

	// Create closure to get data
	getData := func() interface{} {
		result := dataMethod.Call(nil)
		if len(result) > 0 {
			return result[0].Interface()
		}
		return nil
	}

	return &operationInterface{
		getTypeCode: getTypeCode,
		getData:     getData,
	}
}

// encodeOperation encodes an Operation using reflection.
// It first encodes the operation type code, then encodes all fields in order.
func (encoder *Encoder) encodeOperation(op *operationInterface) error {
	// Encode the operation type code
	code := op.getTypeCode()
	if err := encoder.EncodeUVarint(code); err != nil {
		return errors.Wrap(err, "failed to encode operation type")
	}

	// Get the operation data (the actual struct)
	opData := op.getData()
	if opData == nil {
		return errors.New("operation data is nil")
	}

	// Encode the operation data using reflection
	return encoder.encodeByReflection(opData)
}

// encodeByReflection encodes a value using reflection, following the struct field order.
func (encoder *Encoder) encodeByReflection(v interface{}) error {
	if v == nil {
		return errors.New("cannot encode nil value")
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			// For nil pointers, encode as false (optional not present)
			return encoder.EncodeNumber(uint8(0))
		}
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Struct:
		return encoder.encodeStruct(rv)
	case reflect.Slice:
		return encoder.encodeSlice(rv)
	case reflect.Map:
		return encoder.encodeMap(rv)
	case reflect.Interface:
		// If it's an interface, get the underlying value
		if rv.IsNil() {
			return errors.New("cannot encode nil interface")
		}
		return encoder.encodeByReflection(rv.Elem().Interface())
	default:
		// For basic types, try to encode directly
		return encoder.Encode(rv.Interface())
	}
}

// encodeStruct encodes a struct by iterating over its fields in order.
func (encoder *Encoder) encodeStruct(rv reflect.Value) error {
	typ := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		fieldType := typ.Field(i)

		// Skip unexported fields
		if !field.CanInterface() {
			continue
		}

		// Skip fields with json tag "-"
		if jsonTag := fieldType.Tag.Get("json"); jsonTag == "-" {
			continue
		}

		fieldValue := field.Interface()

		// Handle pointer fields
		if field.Kind() == reflect.Ptr {
			if field.IsNil() {
				// For nil pointers in structs, we need to check if it's optional
				// In Steem, optional fields are encoded as: bool (has_value) + value
				// If nil, encode false (0)
				if err := encoder.EncodeNumber(uint8(0)); err != nil {
					return errors.Wrapf(err, "failed to encode nil pointer field %s", fieldType.Name)
				}
				continue
			}
			// Encode true (1) to indicate value is present, then encode the value
			if err := encoder.EncodeNumber(uint8(1)); err != nil {
				return errors.Wrapf(err, "failed to encode pointer presence for field %s", fieldType.Name)
			}
			fieldValue = field.Elem().Interface()
		}

		// Check if this is an asset field (Amount, AmountToSell, MinToReceive, etc.)
		// and the value is a string that looks like an asset (e.g., "0.001 STEEM")
		if field.Kind() == reflect.String {
			fieldName := fieldType.Name
			if fieldName == "Amount" || fieldName == "AmountToSell" || fieldName == "MinToReceive" ||
				fieldName == "VestingShares" || fieldName == "SBDAmount" || fieldName == "SteemAmount" ||
				fieldName == "RewardSteem" || fieldName == "RewardSBD" || fieldName == "RewardVests" {
				assetStr := field.String()
				// Check if it looks like an asset string (contains space and has numeric part)
				if strings.Contains(assetStr, " ") && len(strings.Split(assetStr, " ")) == 2 {
					// Try to parse as asset
					amount, precision, symbol, err := encoder.parseAssetString(assetStr)
					if err == nil {
						// Successfully parsed as asset, encode it
						if err := encoder.encodeAsset(amount, precision, symbol); err != nil {
							return errors.Wrapf(err, "failed to encode asset field %s", fieldType.Name)
						}
						continue
					}
				}
			}
		}

		// Recursively encode the field value
		if err := encoder.encodeByReflection(fieldValue); err != nil {
			return errors.Wrapf(err, "failed to encode field %s", fieldType.Name)
		}
	}

	return nil
}

// encodeSlice encodes a slice by first encoding its length, then each element.
func (encoder *Encoder) encodeSlice(rv reflect.Value) error {
	length := rv.Len()
	if err := encoder.EncodeUVarint(uint64(length)); err != nil {
		return errors.Wrap(err, "failed to encode slice length")
	}

	for i := 0; i < length; i++ {
		elem := rv.Index(i)
		if err := encoder.encodeByReflection(elem.Interface()); err != nil {
			return errors.Wrapf(err, "failed to encode slice element at index %d", i)
		}
	}

	return nil
}

// parseAssetString parses an asset string like "0.001 STEEM" into amount, precision, and symbol.
func (encoder *Encoder) parseAssetString(assetStr string) (amount int64, precision uint8, symbol string, err error) {
	parts := strings.Split(strings.TrimSpace(assetStr), " ")
	if len(parts) != 2 {
		return 0, 0, "", errors.Errorf("invalid asset format: %s", assetStr)
	}

	amountStr := parts[0]
	symbol = strings.ToUpper(parts[1])

	// Parse amount and calculate precision
	dotIndex := strings.Index(amountStr, ".")
	if dotIndex == -1 {
		// No decimal point
		amount, err = strconv.ParseInt(amountStr, 10, 64)
		if err != nil {
			return 0, 0, "", errors.Wrapf(err, "failed to parse amount: %s", amountStr)
		}
		precision = 0
	} else {
		// Has decimal point
		// Remove decimal point and parse as integer
		amountWithoutDot := strings.Replace(amountStr, ".", "", 1)
		amount, err = strconv.ParseInt(amountWithoutDot, 10, 64)
		if err != nil {
			return 0, 0, "", errors.Wrapf(err, "failed to parse amount: %s", amountStr)
		}
		precision = uint8(len(amountStr) - dotIndex - 1)
	}

	return amount, precision, symbol, nil
}

// encodeAsset encodes an asset in the Steem binary format.
// Format: int64 amount (little-endian) + uint8 precision + 7 bytes symbol (null-padded)
func (encoder *Encoder) encodeAsset(amount int64, precision uint8, symbol string) error {
	// Encode amount as int64 (little-endian)
	if err := encoder.EncodeNumber(amount); err != nil {
		return errors.Wrap(err, "failed to encode asset amount")
	}

	// Encode precision as uint8
	if err := encoder.EncodeNumber(precision); err != nil {
		return errors.Wrap(err, "failed to encode asset precision")
	}

	// Encode symbol as 7 bytes (null-padded)
	symbolBytes := make([]byte, 7)
	symbolUpper := strings.ToUpper(symbol)
	copy(symbolBytes, symbolUpper)
	// Remaining bytes are already zero (null-padded)

	if err := encoder.WriteBytes(symbolBytes); err != nil {
		return errors.Wrap(err, "failed to encode asset symbol")
	}

	return nil
}

// encodeMap encodes a map. For Steem, maps are typically encoded as key-value pairs.
func (encoder *Encoder) encodeMap(rv reflect.Value) error {
	length := rv.Len()
	if err := encoder.EncodeUVarint(uint64(length)); err != nil {
		return errors.Wrap(err, "failed to encode map length")
	}

	// Iterate over map entries
	for _, key := range rv.MapKeys() {
		value := rv.MapIndex(key)

		// Encode key
		if err := encoder.encodeByReflection(key.Interface()); err != nil {
			return errors.Wrap(err, "failed to encode map key")
		}

		// Encode value
		if err := encoder.encodeByReflection(value.Interface()); err != nil {
			return errors.Wrap(err, "failed to encode map value")
		}
	}

	return nil
}
