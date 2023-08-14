package event

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/bidon-io/bidon-backend/internal/sdkapi/geocoder"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
)

type Event interface {
	Topic() Topic
	Payload() (map[string]any, error)
	Children() []Event
}

func NewConfig(request *schema.ConfigRequest, geoData geocoder.GeoData) Event {
	return &simpleEvent[*schema.ConfigRequest]{
		timestamp: generateTimestamp(),
		topic:     ConfigTopic,
		request:   request,
		geoData:   geoData,
	}
}

func NewShow(request *schema.ShowRequest, geoData geocoder.GeoData) Event {
	return &simpleEvent[*schema.ShowRequest]{
		timestamp: generateTimestamp(),
		topic:     ShowTopic,
		request:   request,
		geoData:   geoData,
	}
}

func NewClick(request *schema.ClickRequest, geoData geocoder.GeoData) Event {
	return &simpleEvent[*schema.ClickRequest]{
		timestamp: generateTimestamp(),
		topic:     ClickTopic,
		request:   request,
		geoData:   geoData,
	}
}

func NewReward(request *schema.RewardRequest, geoData geocoder.GeoData) Event {
	return &simpleEvent[*schema.RewardRequest]{
		timestamp: generateTimestamp(),
		topic:     RewardTopic,
		request:   request,
		geoData:   geoData,
	}
}

func NewStats(request *schema.StatsRequest, geoData geocoder.GeoData) Event {
	return &statsEvent{
		simpleEvent[*schema.StatsRequest]{
			timestamp: generateTimestamp(),
			topic:     StatsTopic,
			request:   request,
			geoData:   geoData,
		},
	}
}

func NewLoss(request *schema.LossRequest, geoData geocoder.GeoData) Event {
	return &simpleEvent[*schema.LossRequest]{
		timestamp: generateTimestamp(),
		topic:     LossTopic,
		request:   request,
		geoData:   geoData,
	}
}

func NewWin(request *schema.WinRequest, geoData geocoder.GeoData) Event {
	return &simpleEvent[*schema.WinRequest]{
		timestamp: generateTimestamp(),
		topic:     WinTopic,
		request:   request,
		geoData:   geoData,
	}
}

type Topic string

const (
	ConfigTopic Topic = "config"
	ShowTopic   Topic = "show"
	ClickTopic  Topic = "click"
	RewardTopic Topic = "reward"
	StatsTopic  Topic = "stats"
	LossTopic   Topic = "loss"
	WinTopic    Topic = "win"
)

type simpleEvent[T mapper] struct {
	timestamp float64
	topic     Topic
	request   T
	geoData   geocoder.GeoData
}

func (e *simpleEvent[T]) Topic() Topic {
	return e.topic
}

func (e *simpleEvent[T]) Payload() (map[string]any, error) {
	return prepareEventPayload(e.timestamp, e.request, e.geoData)
}

func (e *simpleEvent[T]) Children() []Event {
	return nil
}

type statsEvent struct {
	simpleEvent[*schema.StatsRequest]
}

func (s *statsEvent) Payload() (map[string]any, error) {
	payload, err := s.simpleEvent.Payload()

	payload["event_type"] = "stats"

	return payload, err
}

func (s *statsEvent) Children() []Event {
	children := make([]Event, 0)

	for roundIndex, round := range s.request.Stats.Rounds {
		for demandIndex := range round.Demands {
			children = append(children, &demandResultEvent{
				simpleEvent: s.simpleEvent,
				roundIndex:  roundIndex,
				demandIndex: demandIndex,
			})
		}

		children = append(children, &roundResultEvent{
			simpleEvent: s.simpleEvent,
			roundIndex:  roundIndex,
		})
	}

	return children
}

type roundResultEvent struct {
	simpleEvent[*schema.StatsRequest]
	roundIndex int
}

func (r roundResultEvent) Payload() (map[string]any, error) {
	payload, err := r.simpleEvent.Payload()

	round := r.request.Stats.Rounds[r.roundIndex]
	winnerDemand := roundWinnerDemand(round)

	payload["event_type"] = "round_result"
	payload["timestamp"] = roundTimestamp(round, r.timestamp)
	if round.WinnerID != "" {
		payload["stats__result__status"] = "SUCCESS"
	} else {
		payload["stats__result__status"] = "FAIL"
	}
	payload["stats__result__winner_id"] = round.WinnerID
	payload["stats__result__ad_unit_id"] = winnerDemand.AdUnitID
	payload["stats__result__ecpm"] = round.WinnerECPM
	payload["round_id"] = round.ID
	payload["pricefloor"] = round.PriceFloor

	return payload, err
}

type demandResultEvent struct {
	simpleEvent[*schema.StatsRequest]
	roundIndex  int
	demandIndex int
}

func (r *demandResultEvent) Payload() (map[string]any, error) {
	payload, err := r.simpleEvent.Payload()

	round := r.request.Stats.Rounds[r.roundIndex]
	demand := round.Demands[r.demandIndex]

	payload["event_type"] = "demand_result"
	payload["timestamp"] = demandTimestamp(demand, r.timestamp)
	payload["stats__result__status"] = demand.Status
	payload["stats__result__winner_id"] = demand.ID
	payload["stats__result__ad_unit_id"] = demand.AdUnitID
	payload["stats__result__ecpm"] = demand.ECPM
	payload["round_id"] = round.ID
	payload["pricefloor"] = round.PriceFloor

	return payload, err
}

func generateTimestamp() float64 {
	return float64(time.Now().UnixNano()) / 1e9
}

type mapper interface {
	Map() map[string]any
}

func prepareEventPayload(timestamp float64, requestMapper mapper, geoData geocoder.GeoData) (map[string]any, error) {
	requestMap := requestMapper.Map()

	requestMap["timestamp"] = timestamp

	geo, _ := requestMap["geo"].(map[string]any)
	requestMap["geo"] = enhanceEventGeo(geo, geoData)

	ext, _ := requestMap["ext"].(string)
	eventExt, err := unmarshalEventExt(ext)
	requestMap["ext"] = eventExt

	if _, showPresent := requestMap["show"]; !showPresent {
		if bid, bidPresent := requestMap["bid"]; bidPresent {
			requestMap["show"] = bid
		}
	}

	return smashMap(requestMap, nil), err
}

func enhanceEventGeo(geo map[string]any, geoData geocoder.GeoData) map[string]any {
	if geo == nil {
		geo = make(map[string]any)
	}

	if geoData != (geocoder.GeoData{}) {
		geo["ip"] = geoData.IPString
		geo["country"] = geoData.CountryCode
		geo["country_id"] = geoData.CountryID
	}

	return geo
}

func unmarshalEventExt(ext string) (map[string]any, error) {
	result := make(map[string]any)

	if ext == "" {
		return result, nil
	}

	err := json.Unmarshal([]byte(ext), &result)
	if err != nil {
		return result, fmt.Errorf("unmarshal ext: %v", err)
	}

	return result, nil
}

func smashMap(src, dst map[string]any, nesting ...string) map[string]any {
	if dst == nil {
		dst = make(map[string]any)
	}
	prefix := strings.Join(nesting, "__")

	for key, value := range src {
		switch mapValue := value.(type) {
		case map[string]any:
			n := slices.Clone(nesting)
			n = append(n, key)
			smashMap(mapValue, dst, n...)
		case []map[string]any:
			for i, v := range mapValue {
				n := slices.Clone(nesting)
				n = append(n, fmt.Sprintf("%s__%d", key, i))
				smashMap(v, dst, n...)
			}
		default:
			if prefix != "" {
				dst[fmt.Sprintf("%s__%s", prefix, key)] = value
			} else {
				dst[key] = value
			}
		}
	}

	return dst
}

func roundWinnerDemand(round schema.StatsRound) (winnerDemand schema.StatsDemand) {
	for _, demand := range round.Demands {
		if demand.Status == "WIN" {
			winnerDemand = demand
			break
		}
	}

	return
}

func roundTimestamp(round schema.StatsRound, statsTS float64) (roundTS float64) {
	for _, demand := range round.Demands {
		demandTS := demandTimestamp(demand, statsTS)
		if demandTS > roundTS {
			roundTS = demandTS
		}
	}

	return
}

func demandTimestamp(demand schema.StatsDemand, statsTS float64) (demandTS float64) {
	if demand.FillFinishTS != 0 {
		demandTS = float64(demand.FillFinishTS) / 1000
	} else if demand.BidFinishTS != 0 {
		demandTS = float64(demand.BidFinishTS) / 1000
	}

	// We don't really care what the timestamp is,
	// as long as it's less than the timestamp of the stats event and is not 0
	if demandTS == 0 || demandTS > statsTS {
		demandTS = statsTS
	}

	return
}
