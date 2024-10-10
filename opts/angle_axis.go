package opts

import "github.com/iray-anriy/go-echarts/v2/types"

type AngleAxis struct {
	PolarAxisBase
	Clockwise types.Bool `json:"clockwise,omitempty"`
}
