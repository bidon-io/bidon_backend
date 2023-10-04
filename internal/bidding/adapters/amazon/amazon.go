package amazon

import (
	"encoding/json"
	"fmt"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/bidding/adapters"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"github.com/gofrs/uuid/v5"
)

type Slot struct {
	SlotUUID   string `json:"slotUuid"`
	PricePoint string `json:"pricePoint"`
}

type PricePointsMap map[string]PricePoint

type PricePoint struct {
	Price      float64 `json:"price"`
	PricePoint string  `json:"price_point"`
}

type Adapter struct {
	PricePointsMap PricePointsMap
}

func (a *Adapter) FetchBids(br *schema.BiddingRequest) ([]*adapters.DemandResponse, error) {
	slotsJSON := br.Imp.Demands[adapter.AmazonKey]["token"].(string)
	var slots []Slot
	err := json.Unmarshal([]byte(slotsJSON), &slots)
	if err != nil {
		return nil, err
	}

	demandResponses := make([]*adapters.DemandResponse, 0, len(slots))
	for _, slot := range slots {
		if pricePoint, ok := a.PricePointsMap[slot.PricePoint]; ok {
			impId, _ := uuid.NewV4()
			demandResponse := adapters.DemandResponse{
				DemandID: adapter.AmazonKey,
				SlotUUID: slot.SlotUUID,
				Bid: &adapters.BidDemandResponse{
					ImpID: impId.String(),
					Price: pricePoint.Price,
				},
			}
			demandResponses = append(demandResponses, &demandResponse)
		} else {
			demandResponse := adapters.DemandResponse{
				DemandID: adapter.AmazonKey,
				Error:    fmt.Errorf("cannot find price point"),
				SlotUUID: slot.SlotUUID,
			}
			demandResponses = append(demandResponses, &demandResponse)
		}
	}

	return demandResponses, nil
}

func Builder(cfg adapter.ProcessedConfigsMap) (*Adapter, error) {
	amazonCfg := cfg[adapter.AmazonKey]

	pricePointsMapRaw, ok := amazonCfg["price_points_map"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("price_points_map is not set")
	}

	pricePointsMap := make(PricePointsMap)
	for key, value := range pricePointsMapRaw {
		var point PricePoint
		if v, ok := value.(map[string]any); ok {
			if price, ok := v["price"].(float64); ok {
				point.Price = price
			}
			if pricePoint, ok := v["price_point"].(string); ok {
				point.PricePoint = pricePoint
			}
			pricePointsMap[key] = point
		} else {
			return nil, fmt.Errorf("invalid structure for price_point: %s", key)
		}
	}

	adpt := &Adapter{
		PricePointsMap: pricePointsMap,
	}

	return adpt, nil
}
