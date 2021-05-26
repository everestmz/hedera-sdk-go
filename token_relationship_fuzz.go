package hedera

func FuzzTokenRelationshipFromBytes(Data []byte) int {
	_, err := TokenRelationshipFromBytes(Data)
	if err == nil {
		return 1
	}
	return 0
}
