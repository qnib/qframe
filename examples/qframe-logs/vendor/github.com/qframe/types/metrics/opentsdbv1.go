package qtypes_metrics

import (
	"fmt"
	"strings"
)

type OpenTSDBMetric struct {
	Metric string 			`json:"metric"`
	Timestamp int 			`json:"timestamp"`
	Value float64 			`json:"value"`
	Tags map[string]string 	`json:"tags"`
}

func (otm *OpenTSDBMetric) String() string {
	res := []string{
		string(otm.Metric),
		fmt.Sprintf("%f", otm.Value),
		fmt.Sprintf("%d", otm.Timestamp),
	}
	tags := []string{}
	for k,v := range otm.Tags {
		tags = append(tags, fmt.Sprintf("%s=%s", k, v))
	}
	if len(tags) != 0 {
		res = append(res, strings.Join(tags, ","))
	}
	return strings.Join(res, " ")
}
