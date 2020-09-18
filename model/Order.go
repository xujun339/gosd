package model

import (
	"database/sql"
	"fmt"
	"time"
)

type NullTime struct {
	sql.NullTime
}

type Order struct {
	OrderId int64 `json:"order_id" db:"order_id"`
	CityId int32 `json:"city_id" db:"city_id"`
	OrderStatus int8 `json:"order_status" db:"order_status"`
	Latlong sql.NullString `json:"latlong"`
	RerouteTime NullTime `json:"reroute_time" db:"reroute_time"`
	//RerouteTime time.Time `json:"reroute_time" db:"reroute_time"`
}

func (this NullTime) MarshalJSON() ([]byte, error) {
	if this.Valid {
		stamp := fmt.Sprintf("\"%s\"", time.Time(this.Time).Format("2006-01-02 15:04:05"))
		return []byte(stamp), nil
	}
	return []byte("null"), nil


}
