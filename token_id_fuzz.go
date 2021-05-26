package hedera

func FuzzTokenIDFromBytes(Data []byte) int {
	_, err := TokenIDFromBytes(Data)
	if err == nil {
		return 1
	}
	return 0
}
