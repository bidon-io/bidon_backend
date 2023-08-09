package admin

type Country struct {
	ID int64 `json:"id"`
	CountryAttrs
}

type CountryAttrs struct {
	HumanName  string `json:"human_name"`
	Alpha2Code string `json:"alpha2_code"`
	Alpha3Code string `json:"alpha3_code"`
}

type CountryService = resourceService[Country, CountryAttrs]
