package event

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/bidon-io/bidon-backend/internal/sdkapi/geocoder"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"golang.org/x/exp/slices"
)

type Event interface {
	Topic() Topic
	Payload() (map[string]any, error)
	Children() []Event
}

func NewConfig(request *schema.ConfigRequest, geoData geocoder.GeoData) Event {
	return &configEvent{
		timestamp: generateTimestamp(),
		request:   request,
		geoData:   geoData,
	}
}

func NewShow(request *schema.ShowRequest, geoData geocoder.GeoData) Event {
	return &showEvent{
		timestamp: generateTimestamp(),
		request:   request,
		geoData:   geoData,
	}
}

func NewStats(request *schema.StatsRequest, geoData geocoder.GeoData) Event {
	return &statsEvent{
		timestamp: generateTimestamp(),
		request:   request,
		geoData:   geoData,
	}
}

func NewClick(request *schema.ClickRequest, geoData geocoder.GeoData) Event {
	return &clickEvent{
		timestamp: generateTimestamp(),
		request:   request,
		geoData:   geoData,
	}
}

func NewReward(request *schema.RewardRequest, geoData geocoder.GeoData) Event {
	return &rewardEvent{
		timestamp: generateTimestamp(),
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
)

type configEvent struct {
	timestamp float64
	request   *schema.ConfigRequest
	geoData   geocoder.GeoData
}

func (c *configEvent) Topic() Topic {
	return ConfigTopic
}

func (c *configEvent) Payload() (map[string]any, error) {
	return prepareEventPayload(c.timestamp, c.request, c.geoData)
}

func (c *configEvent) Children() []Event {
	return nil
}

type showEvent struct {
	timestamp float64
	request   *schema.ShowRequest
	geoData   geocoder.GeoData
}

func (e *showEvent) Topic() Topic {
	return ShowTopic
}

func (e *showEvent) Payload() (map[string]any, error) {
	payload, err := prepareEventPayload(e.timestamp, e.request, e.geoData)

	if _, found := payload["show"]; !found {
		payload["show"] = payload["bid"]
	}

	return payload, err
}

func (e *showEvent) Children() []Event {
	return nil
}

type clickEvent struct {
	timestamp float64
	request   *schema.ClickRequest
	geoData   geocoder.GeoData
}

func (e *clickEvent) Topic() Topic {
	return ClickTopic
}

func (e *clickEvent) Payload() (map[string]any, error) {
	payload, err := prepareEventPayload(e.timestamp, e.request, e.geoData)

	if _, found := payload["show"]; !found {
		payload["show"] = payload["bid"]
	}

	return payload, err
}

func (e *clickEvent) Children() []Event {
	return nil
}

type rewardEvent struct {
	timestamp float64
	request   *schema.RewardRequest
	geoData   geocoder.GeoData
}

func (e *rewardEvent) Topic() Topic {
	return RewardTopic
}

func (e *rewardEvent) Payload() (map[string]any, error) {
	payload, err := prepareEventPayload(e.timestamp, e.request, e.geoData)

	if _, found := payload["show"]; !found {
		payload["show"] = payload["bid"]
	}

	return payload, err
}

func (e *rewardEvent) Children() []Event {
	return nil
}

type statsEvent struct {
	timestamp float64
	request   *schema.StatsRequest
	geoData   geocoder.GeoData
}

func (s *statsEvent) Topic() Topic {
	return StatsTopic
}

func (s *statsEvent) Payload() (map[string]any, error) {
	payload, err := prepareEventPayload(s.timestamp, s.request, s.geoData)

	payload["event_type"] = "stats"

	return payload, err
}

func (s *statsEvent) Children() []Event {
	children := make([]Event, 0)

	for roundIndex, round := range s.request.Stats.Rounds {
		for demandIndex := range round.Demands {
			children = append(children, &demandResultEvent{
				timestamp:   s.timestamp,
				request:     s.request,
				geoData:     s.geoData,
				roundIndex:  roundIndex,
				demandIndex: demandIndex,
			})
		}

		children = append(children, &roundResultEvent{
			timestamp:  s.timestamp,
			request:    s.request,
			geoData:    s.geoData,
			roundIndex: roundIndex,
		})
	}

	return children
}

type roundResultEvent struct {
	timestamp  float64
	request    *schema.StatsRequest
	geoData    geocoder.GeoData
	roundIndex int
}

func (r roundResultEvent) Topic() Topic {
	return StatsTopic
}

func (r roundResultEvent) Payload() (map[string]any, error) {
	payload, err := prepareEventPayload(r.timestamp, r.request, r.geoData)

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

func (r roundResultEvent) Children() []Event {
	return nil
}

type demandResultEvent struct {
	timestamp   float64
	request     *schema.StatsRequest
	geoData     geocoder.GeoData
	roundIndex  int
	demandIndex int
}

func (r *demandResultEvent) Topic() Topic {
	return StatsTopic
}

func (r *demandResultEvent) Payload() (map[string]any, error) {
	payload, err := prepareEventPayload(r.timestamp, r.request, r.geoData)

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

func (r *demandResultEvent) Children() []Event {
	return nil
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
