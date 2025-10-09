package geo

import "strings"

const DefaultRegion = "us"

var alpha3ToRegionMapping = map[string]string{
	// US region countries
	"ABW": "us", "AIA": "us", "ARG": "us", "ATG": "us", "BES": "us", "BHS": "us", "BLM": "us", "BLZ": "us",
	"BOL": "us", "BRA": "us", "BRB": "us", "CAN": "us", "CHL": "us", "COL": "us", "CRI": "us", "CUB": "us",
	"CUW": "us", "CYM": "us", "DMA": "us", "DOM": "us", "ECU": "us", "GLP": "us", "GRD": "us", "GRL": "us",
	"GTM": "us", "GUF": "us", "GUY": "us", "HND": "us", "HTI": "us", "JAM": "us", "KNA": "us", "LCA": "us",
	"MAF": "us", "MEX": "us", "MSR": "us", "MTQ": "us", "NIC": "us", "PAN": "us", "PER": "us", "PRI": "us",
	"PRY": "us", "SLV": "us", "SUR": "us", "SXM": "us", "TCA": "us", "TST": "us", "TTO": "us", "UMI": "us",
	"URY": "us", "USA": "us", "VCT": "us", "VEN": "us", "VGB": "us", "VIR": "us",

	// Asia region countries
	"AFG": "asia", "ARE": "asia", "ARM": "asia", "ASM": "asia", "ATA": "asia", "ATF": "asia", "AUS": "asia",
	"BGD": "asia", "BHR": "asia", "BRN": "asia", "BTN": "asia", "CCK": "asia", "CHN": "asia", "COK": "asia",
	"COM": "asia", "CXR": "asia", "FJI": "asia", "FSM": "asia", "GUM": "asia", "HKG": "asia", "HMD": "asia",
	"IDN": "asia", "IND": "asia", "IOT": "asia", "IRN": "asia", "IRQ": "asia", "ISR": "asia", "JPN": "asia",
	"KAZ": "asia", "KHM": "asia", "KIR": "asia", "KOR": "asia", "KWT": "asia", "LAO": "asia", "LBN": "asia",
	"LKA": "asia", "MAC": "asia", "MDV": "asia", "MHL": "asia", "MMR": "asia", "MNG": "asia", "MNP": "asia",
	"MYS": "asia", "MYT": "asia", "NCL": "asia", "NFK": "asia", "NIU": "asia", "NPL": "asia", "NRU": "asia",
	"NZL": "asia", "OMN": "asia", "PAK": "asia", "PCN": "asia", "PHL": "asia", "PLW": "asia", "PNG": "asia",
	"PRK": "asia", "PYF": "asia", "QAT": "asia", "SAU": "asia", "SGP": "asia", "SLB": "asia", "SSG": "asia",
	"SYC": "asia", "THA": "asia", "TJK": "asia", "TKL": "asia", "TKM": "asia", "TLS": "asia", "TON": "asia",
	"TUV": "asia", "TWN": "asia", "UZB": "asia", "VNM": "asia", "VUT": "asia", "WLF": "asia", "WSM": "asia",
	"YEM": "asia",

	// EU region countries
	"AGO": "eu", "ALA": "eu", "ALB": "eu", "AND": "eu", "AUT": "eu", "AZE": "eu", "BDI": "eu", "BEL": "eu",
	"BEN": "eu", "BFA": "eu", "BGR": "eu", "BIH": "eu", "BLR": "eu", "BMU": "eu", "BVT": "eu", "BWA": "eu",
	"CAF": "eu", "CHE": "eu", "CIV": "eu", "CMR": "eu", "COD": "eu", "COG": "eu", "CPV": "eu", "CYP": "eu",
	"CZE": "eu", "DEU": "eu", "DJI": "eu", "DNK": "eu", "DZA": "eu", "EGY": "eu", "ERI": "eu", "ESH": "eu",
	"ESP": "eu", "EST": "eu", "ETH": "eu", "FIN": "eu", "FRA": "eu", "FRO": "eu", "GAB": "eu", "GBR": "eu",
	"GEO": "eu", "GGY": "eu", "GHA": "eu", "GIB": "eu", "GIN": "eu", "GMB": "eu", "GNB": "eu", "GNQ": "eu",
	"GRC": "eu", "HRV": "eu", "HUN": "eu", "IMN": "eu", "IRL": "eu", "ISL": "eu", "ITA": "eu", "JEY": "eu",
	"JOR": "eu", "KEN": "eu", "KGZ": "eu", "LBR": "eu", "LBY": "eu", "LIE": "eu", "LSO": "eu", "LTU": "eu",
	"LUX": "eu", "LVA": "eu", "MAR": "eu", "MCO": "eu", "MDA": "eu", "MDG": "eu", "MKD": "eu", "MLI": "eu",
	"MLT": "eu", "MNE": "eu", "MOZ": "eu", "MRT": "eu", "MUS": "eu", "MWI": "eu", "NAM": "eu", "NER": "eu",
	"NGA": "eu", "NLD": "eu", "NOR": "eu", "POL": "eu", "PRT": "eu", "PSE": "eu", "REU": "eu", "ROU": "eu",
	"RUS": "eu", "RWA": "eu", "SDN": "eu", "SEN": "eu", "SGS": "eu", "SHN": "eu", "SJM": "eu", "SLE": "eu",
	"SMR": "eu", "SOM": "eu", "SRB": "eu", "SSD": "eu", "STP": "eu", "SVK": "eu", "SVN": "eu", "SWE": "eu",
	"SWZ": "eu", "SYR": "eu", "TCD": "eu", "TGO": "eu", "TUN": "eu", "TUR": "eu", "TZA": "eu", "UGA": "eu",
	"UKR": "eu", "VAT": "eu", "ZAF": "eu", "ZMB": "eu", "ZWE": "eu",
}

// Region returns the Moloco/Start.io region identifier for the provided ISO-3166-1 alpha-3 country code.
func Region(alpha3 string) string {
	code := strings.ToUpper(alpha3)
	if region, ok := alpha3ToRegionMapping[code]; ok {
		return region
	}
	return DefaultRegion
}
