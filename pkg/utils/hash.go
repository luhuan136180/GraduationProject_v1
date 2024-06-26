package utils

import (
	"crypto"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

func MD5(str string) []byte {
	hashed, _ := Hash(crypto.MD5, str)
	return hashed
}

func MD5Hex(str string) string {
	hashed, _ := HashHex(crypto.MD5, str)
	return hashed
}

func SHA1(str string) []byte {
	hashed, _ := Hash(crypto.SHA1, str)
	return hashed
}

func SHA1Hex(str string) string {
	hashed, _ := HashHex(crypto.SHA1, str)
	return hashed
}

func SHA256(str string) []byte {
	hashed, _ := Hash(crypto.SHA256, str)
	return hashed
}

func SHA256Hex(str string) string {
	hashed, _ := HashHex(crypto.SHA256, str)
	return hashed
}

func Hash(hash crypto.Hash, str string) ([]byte, error) {
	if !hash.Available() {
		return nil, fmt.Errorf("invilid hash")
	}

	h := hash.New()
	h.Write([]byte(str))
	return h.Sum(nil), nil
}

func HashHex(h crypto.Hash, str string) (string, error) {
	hashed, err := Hash(h, str)
	if err != nil {
		return "", nil
	}

	return hex.EncodeToString(hashed), nil
}

func HashCollegeID(CollegeName string) string {
	return MD5Hex(fmt.Sprintf("%s$", CollegeName))
}

func HashProfessionID(CollegeHashID, ProfessionName string) string {
	return MD5Hex(fmt.Sprintf("%s$%s", CollegeHashID, ProfessionName))
}

func HashClassID(professionHashID, className string, classID int) string {
	return MD5Hex(fmt.Sprintf("%s$%s$%d", professionHashID, className, classID))
}

func CreateContentHashBySHA256(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	// 计算序列化后JSON数据的SHA-256哈希值
	hash := sha256.Sum256(data)

	// 将哈希值转换为十六进制字符串
	hashID := hex.EncodeToString(hash[:])

	return hashID, nil
}
