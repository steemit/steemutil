package wif

type testData struct {
	WIF           string
	PrivateKeyHex string
	PublicKey     string
}

var data = []testData{
	{
		WIF:           "5JWHY5DxTF6qN5grTtChDCYBmWHfY9zaSsw4CxEKN5eZpH9iBma",
		PrivateKeyHex: "5ad2b8df2c255d4a2996ee7d065e013e1bbb35c075ee6e5208aca44adc9a9d4c",
		PublicKey:     "STM7jNh5ejQoqHqWcGWFJ1v4F5CzsG3EiBuz1VooCng1cH5QpJD27",
	},
	{
		WIF:           "5KPipdRzoxrp6dDqsBfMD6oFZG356trVHV5QBGx3rABs1zzWWs8",
		PrivateKeyHex: "cf9d6121ed458f24ea456ad7ff700da39e86688988cfe5c6ed6558642cf1e32f",
		PublicKey:     "STM7W7ACQDZJZ6rZGKeT9auipnSiSxFxJ4k71QXmrhY9HbvYsNnQ2",
	},
}
