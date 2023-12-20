package idv

import (
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"net/http"
)

var nationalIDPrefixMap = map[string]string{"A": "10", "B": "11", "C": "12", "D": "13", "E": "14", "F": "15", "G": "16", "H": "17", "I": "34", "J": "18", "K": "19", "M": "21", "N": "22", "O": "35", "P": "23", "Q": "24", "T": "27", "U": "28", "V": "29", "W": "32", "X": "30", "Z": "33", "L": "20", "R": "25", "S": "26", "Y": "31"}

var oldARCGenderMap = map[string]rune{"A": '0', "B": '1', "C": '2', "D": '3'} // = nationalIDPrefix[1]

var nationalIDWeights = [11]int32{1, 9, 8, 7, 6, 5, 4, 3, 2, 1, 1}

// Ref: https://wisdom-life.in/generator/taiwain-id-generator
// Ref: https://zh.wikipedia.org/wiki/%E4%B8%AD%E8%8F%AF%E6%B0%91%E5%9C%8B%E5%9C%8B%E6%B0%91%E8%BA%AB%E5%88%86%E8%AD%89
func validateNationalID(nationalID string, isARC bool) *errpkg.Error {
	cityVal, ok := nationalIDPrefixMap[nationalID[0:1]]
	if !ok || len(nationalID) != 10 {
		return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadNationalID}
	}

	if isARC {
		if val, ok := oldARCGenderMap[nationalID[1:2]]; ok {
			newID := []rune(nationalID)
			newID[1] = val
			nationalID = string(newID)
		}
	}

	nationalID = cityVal + nationalID[1:]
	var sum int32
	for index, val := range nationalID {
		iVal := val - '0'
		if iVal < 0 || iVal > 9 {
			return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadNationalID}
		}

		sum += iVal * nationalIDWeights[index]
	}

	if sum%10 != 0 {
		return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadNationalID, Data: sum}
	}

	return nil
}
