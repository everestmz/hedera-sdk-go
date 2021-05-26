package hedera

func FuzzTokenInfoFromBytes(Data []byte) int {
	_, err := TokenInfoFromBytes(Data)
	if err == nil {
		return 1
	}
	return 0
}
