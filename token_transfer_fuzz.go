package hedera

func FuzzTokenTransfetFromBytes(Data []byte) int {
	_, err := TokenTransferFromBytes(Data)
	if err == nil {
		return 1
	}
	return 0
}
