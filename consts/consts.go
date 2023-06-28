package consts

const ADDRESS_PREFIX = "STM"

const STEEM_ASSET_SYMBOL_PRECISION_BITS = 4
const STEEM_ASSET_CONTROL_BITS = 1
const STEEM_NAI_SHIFT = (STEEM_ASSET_SYMBOL_PRECISION_BITS + STEEM_ASSET_CONTROL_BITS)
const SMT_MAX_NAI = 99999999
const SMT_MIN_NAI = 1
const SMT_MIN_NON_RESERVED_NAI = 10000000
const STEEM_ASSET_SYMBOL_NAI_LENGTH = 10
const STEEM_ASSET_SYMBOL_NAI_STRING_LENGTH = (STEEM_ASSET_SYMBOL_NAI_LENGTH + 2)
const SMT_MAX_NAI_POOL_COUNT = 10
const SMT_MAX_NAI_GENERATION_TRIES = 100

// One's place is used for check digit, which means NAI 0-9 all have NAI data of 0 which is invalid
// This space is safe to use because it would alwasys result in failure to convert from NAI
const STEEM_NAI_SBD = 1
const STEEM_NAI_STEEM = 2
const STEEM_NAI_VESTS = 3

const STEEM_PRECISION_SBD = 3
const STEEM_PRECISION_STEEM = 3
const STEEM_PRECISION_VESTS = 6

func STEEM_ASSET_NUM_SBD() uint32 {
	return uint32(((SMT_MAX_NAI + STEEM_NAI_SBD) << STEEM_NAI_SHIFT) | STEEM_PRECISION_SBD)
}
func STEEM_ASSET_NUM_STEEM() uint32 {
	return uint32(((SMT_MAX_NAI + STEEM_NAI_STEEM) << STEEM_NAI_SHIFT) | STEEM_PRECISION_STEEM)
}
func STEEM_ASSET_NUM_VESTS() uint32 {
	return uint32(((SMT_MAX_NAI + STEEM_NAI_VESTS) << STEEM_NAI_SHIFT) | STEEM_PRECISION_VESTS)
}
func calculateSymbolU64(symbol string) uint64 {
	var result uint64
	for _, char := range symbol {
		result = (result << 8) | uint64(char)
	}
	return result
}
func VESTS_SYMBOL_U64() uint64 {
	return calculateSymbolU64("VESTS")
}
func STEEM_SYMBOL_U64() uint64 {
	return calculateSymbolU64("STEEM")
}
func SBD_SYMBOL_U64() uint64 {
	return calculateSymbolU64("SBD")
}
