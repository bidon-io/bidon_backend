package openrtb

import (
	"encoding/json"
	"github.com/prebid/openrtb/v19/adcom1"
	"github.com/prebid/openrtb/v19/openrtb2"
)

// All these code was created to Patch the Segment struct to allow Signal attribute

type BidRequest struct {
	User *User `json:"user,omitempty"`

	ID      string                  `json:"id"`
	Imp     []openrtb2.Imp          `json:"imp"`
	Site    *openrtb2.Site          `json:"site,omitempty"`
	App     *openrtb2.App           `json:"app,omitempty"`
	DOOH    *openrtb2.DOOH          `json:"dooh,omitempty"`
	Device  *openrtb2.Device        `json:"device,omitempty"`
	Test    int8                    `json:"test,omitempty"`
	AT      int64                   `json:"at,omitempty"`
	TMax    int64                   `json:"tmax,omitempty"`
	WSeat   []string                `json:"wseat,omitempty"`
	BSeat   []string                `json:"bseat,omitempty"`
	AllImps int8                    `json:"allimps,omitempty"`
	Cur     []string                `json:"cur,omitempty"`
	WLang   []string                `json:"wlang,omitempty"`
	WLangB  []string                `json:"wlangb,omitempty"`
	BCat    []string                `json:"bcat,omitempty"`
	CatTax  adcom1.CategoryTaxonomy `json:"cattax,omitempty"`
	BAdv    []string                `json:"badv,omitempty"`
	BApp    []string                `json:"bapp,omitempty"`
	Source  *openrtb2.Source        `json:"source,omitempty"`
	Regs    *openrtb2.Regs          `json:"regs,omitempty"`
	Ext     json.RawMessage         `json:"ext,omitempty"`
}

type User struct {
	Data []Data `json:"data,omitempty"`

	ID         string          `json:"id,omitempty"`
	BuyerUID   string          `json:"buyeruid,omitempty"`
	Yob        int64           `json:"yob,omitempty"`
	Gender     string          `json:"gender,omitempty"`
	Keywords   string          `json:"keywords,omitempty"`
	KwArray    []string        `json:"kwarray,omitempty"`
	CustomData string          `json:"customdata,omitempty"`
	Geo        *openrtb2.Geo   `json:"geo,omitempty"`
	Consent    string          `json:"consent,omitempty"`
	EIDs       []openrtb2.EID  `json:"eids,omitempty"`
	Ext        json.RawMessage `json:"ext,omitempty"`
}

type Data struct {
	Segment []Segment `json:"segment,omitempty"`

	ID   string          `json:"id,omitempty"`
	Name string          `json:"name,omitempty"`
	Ext  json.RawMessage `json:"ext,omitempty"`
}

type Segment struct {
	Signal string `json:"signal,omitempty"`

	ID    string          `json:"id,omitempty"`
	Name  string          `json:"name,omitempty"`
	Value string          `json:"value,omitempty"`
	Ext   json.RawMessage `json:"ext,omitempty"`
}
