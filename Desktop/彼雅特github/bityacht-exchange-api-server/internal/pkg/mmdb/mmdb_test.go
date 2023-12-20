package mmdb

import (
	"testing"
)

func TestLookup(t *testing.T) {
	ips := []string{
		"127.0.0.1",
		"dsadasdsa",
		"192.168.127.0",
		"60.250.152.235",
	}

	for _, ip := range ips {
		if city, err := LookupCity(ip); err != nil {
			t.Error(err)
		} else {
			t.Log(city)
		}
	}
}

func TestListAllCity(t *testing.T) {
	allCity := make(map[string]struct{})

	networks := reader.Networks()
	for networks.Next() {
		var record CityResult

		_, err := networks.Network(&record)
		if err != nil {
			t.Error(err)
			break
		}

		allCity[record.String()] = struct{}{}
	}

	cityArr := make([]string, 0, len(allCity))
	for city := range allCity {
		cityArr = append(cityArr, city)
	}

	t.Log(cityArr)
}
