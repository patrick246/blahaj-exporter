package client

import "time"

type Response struct {
	Availabilities []Availability `json:"availabilities"`
	Timestamp      time.Time      `json:"timestamp"`
	TraceID        string         `json:"traceId"`
}

type Availability struct {
	AvailableForCashCarry    bool          `json:"availableForCashCarry"`
	AvailableForClickCollect bool          `json:"availableForClickCollect"`
	BuyingOption             BuyingOptions `json:"buyingOption"`
	ClassUnitKey             ClassUnitKey  `json:"classUnitKey"`
	ItemKey                  ItemKey       `json:"itemKey"`
}

type BuyingOptions struct {
	CashCarry    BuyingOption `json:"cashCarry"`
	ClickCollect BuyingOption `json:"clickCollect"`
	HomeDelivery BuyingOption `json:"homeDelivery"`
}

type BuyingOption struct {
	Availability BuyingOptionAvailability `json:"availability"`
	Range        BuyingOptionRange        `json:"range"`
}

type BuyingOptionAvailability struct {
	Quantity       int64     `json:"quantity"`
	UpdateDateTime time.Time `json:"updateDateTime"`
}

type BuyingOptionRange struct {
	InRange bool `json:"inRange"`
}

type ClassUnitKey struct {
	ClassUnitCode string `json:"classUnitCode"`
	ClassUnitType string `json:"classUnitType"`
}

type ItemKey struct {
	ItemNo   string `json:"itemNo"`
	ItemType string `json:"itemType"`
}
