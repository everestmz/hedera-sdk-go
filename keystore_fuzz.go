package hedera

func FuzzKeystoreData(Data []byte) int {
	privateKey, err := PrivateKeyFromString(string(Data))
	if err != nil {
		return -1
	}

	// Passphrase taken from `passphrase` in keystore_test.go
	keyStore, err := newKeystore(privateKey.Bytes(), "HelloHashgraph!")
	if err != nil {
		return -1
	}

	ksPrivateKey, err := parseKeystore(keyStore, "HelloHashgraph!")
	if err != nil {
		return -1
	}

	if string(privateKey.keyData) != string(ksPrivateKey.keyData) {
		panic("original != decoded")
	}
	return 0
}

func FuzzKeystorePassphrase(Data []byte) int {
	// string taken from testPrivateKeyStr
	privateKey, err := PrivateKeyFromString("302e020100300506032b657004220420db484b828e64b2d8f12ce3c0a0e93a0b8cce7af1bb8f39c97732394482538e10")
	if err != nil {
		return -1
	}

	keyStore, err := newKeystore(privateKey.Bytes(), string(Data))
	if err != nil {
		return -1
	}

	ksPrivateKey, err := parseKeystore(keyStore, string(Data))
	if err != nil {
		return -1
	}

	if string(privateKey.keyData) != string(ksPrivateKey.keyData) {
		panic("original != decoded")
	}
	return 0
}
