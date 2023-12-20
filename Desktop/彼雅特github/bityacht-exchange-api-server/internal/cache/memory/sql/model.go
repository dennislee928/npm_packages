package sqlcache

import (
	"bityacht-exchange-api-server/internal/database/sql/bankbranchs"
	"bityacht-exchange-api-server/internal/database/sql/banks"
	"bityacht-exchange-api-server/internal/database/sql/countries"
	"bityacht-exchange-api-server/internal/database/sql/industrialclassifications"
	"bityacht-exchange-api-server/internal/database/sql/levellimits"
	"bityacht-exchange-api-server/internal/database/sql/mainnets"
	"bityacht-exchange-api-server/internal/database/sql/users"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
)

var cache cacheStruct

type cacheStruct struct {
	IDVOptionsResponse json.RawMessage // Marshal from IDVOptionsResponse
	CountryMap         map[string]countries.Country
	ICMap              map[int64]industrialclassifications.IC
	LevelLimitMap      map[int64]levellimits.LevelLimit
	BankMap            map[string]banks.Bank
	BankBranchMap      map[string]map[string]bankbranchs.Branch
	MainnetMap         map[string]map[string]mainnets.Model

	SpotOptionsResponse json.RawMessage // Marshal from SpotOptionsResponse
	BankOptionsResponse json.RawMessage // Marshal from BankOptionsResponse

	mux sync.RWMutex
}

func Update() *errpkg.Error {
	if newCountries, newCountryMap, err := countries.GetCountryListAndMap(); err != nil {
		return err
	} else if newICs, newICMap, err := industrialclassifications.GetICListAndMap(); err != nil {
		return err
	} else if newMainnetMap, err := mainnets.GetMap(); err != nil {
		return err
	} else if newLevelLimits, err := levellimits.GetAllLevelLimit(); err != nil {
		return err
	} else if allBanks, err := banks.GetAllBanks(); err != nil {
		return err
	} else if branchListMapByBank, err := bankbranchs.GetBranchListMapByBank(); err != nil {
		return err
	} else {
		idvOptions := IDVOptionsResponse{
			Countries: newCountries,
			ICs:       newICs,
		}

		newIDVOptionsResponse, err := json.Marshal(idvOptions)
		if err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeJSONMarshal, Err: err}
		}

		spotOptions := SpotOptionsResponse{
			Mainnets: make(map[string]map[string]string, len(newMainnetMap)),
		}
		for currency, mainnetMapByCurrency := range newMainnetMap {
			spotOptions.Mainnets[currency] = make(map[string]string, len(mainnetMapByCurrency))

			for _, record := range mainnetMapByCurrency {
				spotOptions.Mainnets[currency][record.Mainnet] = record.Name
			}
		}

		newSpotOptionsResponse, err := json.Marshal(spotOptions)
		if err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeJSONMarshal, Err: err}
		}

		newLevelLimitMap := make(map[int64]levellimits.LevelLimit, len(newLevelLimits))
		for _, record := range newLevelLimits {
			newLevelLimitMap[getLevelLimitKey(record.Type, record.Level)] = record
		}

		newBankMap := make(map[string]banks.Bank, len(allBanks))
		newBankBranchMap := make(map[string]map[string]bankbranchs.Branch, len(allBanks))
		bankOptions := BankOptionsResponse{Banks: make([]BankInfo, len(allBanks))}
		for i, bank := range allBanks {
			newBankInfo := BankInfo{
				Bank:    bank,
				Branchs: branchListMapByBank[bank.Code],
			}
			bankOptions.Banks[i] = newBankInfo

			newBankMap[bank.Code] = bank
			newBankBranchMap[bank.Code] = make(map[string]bankbranchs.Branch, len(newBankInfo.Branchs))
			for _, branch := range newBankInfo.Branchs {
				newBankBranchMap[bank.Code][branch.Code] = branch
			}
		}

		newBankOptionsResponse, err := json.Marshal(bankOptions)
		if err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeJSONMarshal, Err: err}
		}

		cache.mux.Lock()
		defer cache.mux.Unlock()

		cache.IDVOptionsResponse = newIDVOptionsResponse
		cache.CountryMap = newCountryMap
		cache.ICMap = newICMap
		cache.SpotOptionsResponse = newSpotOptionsResponse
		cache.LevelLimitMap = newLevelLimitMap
		cache.BankMap = newBankMap
		cache.BankBranchMap = newBankBranchMap
		cache.BankOptionsResponse = newBankOptionsResponse
		cache.MainnetMap = newMainnetMap
	}

	return nil
}

func GetUserIDVOptionsResponse() json.RawMessage {
	cache.mux.RLock()
	defer cache.mux.RUnlock()

	return cache.IDVOptionsResponse
}

func GetSpotOptionsResponse() json.RawMessage {
	cache.mux.RLock()
	defer cache.mux.RUnlock()

	return cache.SpotOptionsResponse
}

func GetCountry(code string) (countries.Country, *errpkg.Error) {
	cache.mux.RLock()
	defer cache.mux.RUnlock()

	if country, ok := cache.CountryMap[code]; ok {
		return country, nil
	}

	return countries.Country{}, &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound, Err: errors.New("country code not found in sql cache")}
}

func GetIC(id int64) (industrialclassifications.IC, *errpkg.Error) {
	cache.mux.RLock()
	defer cache.mux.RUnlock()

	if ic, ok := cache.ICMap[id]; ok {
		return ic, nil
	}

	return industrialclassifications.IC{}, &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound, Err: errors.New("industrial classifications id not found in sql cache")}
}

func GetLevelLimit(usersType users.Type, level int32) (levellimits.LevelLimit, *errpkg.Error) {
	cache.mux.RLock()
	defer cache.mux.RUnlock()

	if levelLimit, ok := cache.LevelLimitMap[getLevelLimitKey(int32(usersType), level)]; ok {
		return levelLimit, nil
	}

	return levellimits.LevelLimit{}, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadRecord, Err: errors.New("user's type + level not found")}
}

func getLevelLimitKey(usersType int32, level int32) int64 {
	return int64(usersType)<<32 | int64(level)
}

func GetBankOptionsResponse() json.RawMessage {
	cache.mux.RLock()
	defer cache.mux.RUnlock()

	return cache.BankOptionsResponse
}

func GetBankAndBranch(bankCode string, branchCode string) (banks.Bank, bankbranchs.Branch, *errpkg.Error) {
	cache.mux.RLock()
	defer cache.mux.RUnlock()

	bank, ok := cache.BankMap[bankCode]
	if !ok {
		return banks.Bank{}, bankbranchs.Branch{}, &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound, Err: errors.New("bank not found")}
	}

	branch, ok := cache.BankBranchMap[bankCode][branchCode]
	if !ok {
		return bank, bankbranchs.Branch{}, &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound, Err: errors.New("branch not found")}
	}

	return bank, branch, nil
}

func GetBankBranch(bankCode string, branchCode string) (bankbranchs.Branch, *errpkg.Error) {
	cache.mux.RLock()
	defer cache.mux.RUnlock()

	if branch, ok := cache.BankBranchMap[bankCode][branchCode]; ok {
		return branch, nil
	}

	return bankbranchs.Branch{}, &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound, Err: errors.New("bank or branch not found")}
}

func GetMainnetRecord(currency string, mainnet string) (mainnets.Model, *errpkg.Error) {
	cache.mux.RLock()
	defer cache.mux.RUnlock()

	if mainnetRecord, ok := cache.MainnetMap[currency][mainnet]; ok {
		return mainnetRecord, nil
	}

	return mainnets.Model{}, &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound, Err: errors.New("currency or mainnet not found")}
}
