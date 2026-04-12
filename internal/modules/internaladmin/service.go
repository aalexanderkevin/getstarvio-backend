package internaladmin

import (
	"fmt"
	"strconv"
)

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetPlanConfig() (PlanConfigResponse, error) {
	p, err := s.repo.GetPrimaryPlanConfig()
	if err != nil {
		return PlanConfigResponse{}, err
	}

	return PlanConfigResponse{
		BusinessID: p.BusinessID,
		FreeBonus:  p.FreeBonus,
		SubCredits: p.SubCredits,
		SubPrice:   p.SubPrice,
		TopupPrice: p.TopupPrice,
		Tier1Price: p.Tier1Price,
		Tier1Creds: p.Tier1Credits,
		Tier2Price: p.Tier2Price,
		Tier2Creds: p.Tier2Credits,
		Tier3Price: p.Tier3Price,
		Tier3Creds: p.Tier3Credits,
	}, nil
}

func (s *Service) UpdatePlanConfig(data map[string]interface{}) error {
	p, err := s.repo.GetPrimaryPlanConfig()
	if err != nil {
		return err
	}

	payload := map[string]interface{}{}
	assignIfInt(payload, "free_bonus", data, "free_bonus", "freeBonus")
	assignIfInt(payload, "sub_credits", data, "sub_credits", "subCredits")
	assignIfInt(payload, "sub_price", data, "sub_price", "subPrice")
	assignIfInt(payload, "topup_price", data, "topup_price", "topupPrice")
	assignIfInt(payload, "tier1_price", data, "tier1_price", "tier1Price")
	assignIfInt(payload, "tier1_credits", data, "tier1_credits", "tier1Credits")
	assignIfInt(payload, "tier2_price", data, "tier2_price", "tier2Price")
	assignIfInt(payload, "tier2_credits", data, "tier2_credits", "tier2Credits")
	assignIfInt(payload, "tier3_price", data, "tier3_price", "tier3Price")
	assignIfInt(payload, "tier3_credits", data, "tier3_credits", "tier3Credits")

	if len(payload) == 0 {
		return fmt.Errorf("empty payload")
	}

	return s.repo.UpdatePlanConfig(p.BusinessID, payload)
}

func assignIfInt(dst map[string]interface{}, dstKey string, src map[string]interface{}, keys ...string) {
	for _, key := range keys {
		v, ok := src[key]
		if !ok {
			continue
		}
		n, err := toInt(v)
		if err != nil {
			continue
		}
		dst[dstKey] = n
		return
	}
}

func toInt(v interface{}) (int, error) {
	switch t := v.(type) {
	case int:
		return t, nil
	case int32:
		return int(t), nil
	case int64:
		return int(t), nil
	case float32:
		return int(t), nil
	case float64:
		return int(t), nil
	case string:
		n, err := strconv.Atoi(t)
		if err != nil {
			return 0, err
		}
		return n, nil
	default:
		return 0, fmt.Errorf("not numeric")
	}
}
