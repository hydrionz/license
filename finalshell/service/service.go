package service

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/ebfe/keccak"
)

func md5Hash(msg string) string {
	hash := md5.New()
	hash.Write([]byte(msg))
	str := hex.EncodeToString(hash.Sum(nil))
	return str[8:24]
}

func keccak384Hash(msg string) string {
	enc := keccak.New384()
	enc.Write([]byte(msg))
	str := hex.EncodeToString(enc.Sum(nil))
	return str[12:28]
}

type FinalShellLicense struct {
	AdvancedBelow396 string `json:"advancedBelow396"`
	ProBelow396      string `json:"proBelow396"`
	AdvancedAbove396 string `json:"advancedAbove396"`
	ProAbove396      string `json:"proAbove396"`
	AdvancedAbove45  string `json:"advancedAbove45"`
	ProAbove45       string `json:"proAbove45"`
	AdvancedAbove46  string `json:"advancedAbove46"`
	ProAbove46       string `json:"proAbove46"`
}

// GenerateLicense generates license codes for different FinalShell versions and editions
func GenerateLicense(machineCode string) FinalShellLicense {
	// License for Version < 3.9.6 Advanced Edition
	advancedBelow396 := md5Hash("61305" + machineCode + "8552")
	// License for Version < 3.9.6 Professional Edition
	proBelow396 := md5Hash("2356" + machineCode + "13593")
	// License for 3.9.6 <= Version < 4.5 Advanced Edition
	advancedAbove396 := keccak384Hash(machineCode + "hSf(78cvVlS5E")
	// License for3.9.6 <= Version < 4.5 Professional Edition
	proAbove396 := keccak384Hash(machineCode + "FF3Go(*Xvbb5s2")
	// License for 4.5 <= Version < 4.6 Advanced Edition
	AdvancedAbove45 := keccak384Hash(machineCode + "wcegS3gzA$")
	// License for 4.5 <= Version < 4.6 Professional Edition
	proAbove45 := keccak384Hash(machineCode + "b(xxkHn%z);x")
	// License for Version >= 4.6 Advanced Edition
	AdvancedAbove46 := keccak384Hash(machineCode + "csSf5*xlkgYSX,y")
	// License for Version >= 4.6 Professional Edition
	proAbove46 := keccak384Hash(machineCode + "Scfg*ZkvJZc,s,Y")

	return FinalShellLicense{
		AdvancedBelow396: advancedBelow396,
		ProBelow396:      proBelow396,
		AdvancedAbove396: advancedAbove396,
		ProAbove396:      proAbove396,
		AdvancedAbove45:  AdvancedAbove45,
		ProAbove45:       proAbove45,
		AdvancedAbove46:  AdvancedAbove46,
		ProAbove46:       proAbove46,
	}
}
