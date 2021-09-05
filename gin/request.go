package gin

import (
	"encoding/json"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"

	_ "github.com/go-playground/validator/v10"

	"github.com/jurycu/probix/utils"
	"github.com/jurycu/probix/utils/httpclient"
)

type dataJson struct {
	Values []metricsValues `json:"values" validate:"required"`
}

type metricsValues struct {
	Labels map[string]string `json:"labels"`
	Value  float64           `json:"value"`
}

func requestHTTP(target, method, body, metricsName, metricsHelp string, registry *prometheus.Registry) (bool, error) {
	var data dataJson
	var resp []byte
	var err error
	if method == "GET" {
		resp, err = httpclient.HttpGet(target, nil)
		if err != nil {
			return false, err
		}
	} else {
		if body == "" {
			resp, err = httpclient.HttpPost(target, nil)
			if err != nil {
				return false, err
			}
		} else {
			var jsonBytes map[string]string
			if err := json.Unmarshal([]byte(body), &jsonBytes); err != nil {
				return false, err
			}
			resp, err = httpclient.HttpPost(target, jsonBytes)
			if err != nil {
				return false, err
			}
		}
	}
	if err = json.Unmarshal(resp, &data); err != nil {
		return false, err
	}
	//检查数据是否合法
	if err = utils.Validator(data); err != nil {
		return false, err
	}
	//
	//注册指标
	var labelsInit []string
	if len(data.Values) == 0 {
		return false, nil
	}

	for labels, _ := range data.Values[0].Labels {
		labelsInit = append(labelsInit, labels)
	}

	var (
		probixMetrics = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: metricsName,
			Help: metricsHelp,
		}, labelsInit)
	)
	registry.MustRegister(probixMetrics)
	for _, metrics := range data.Values {
		for validatorLabels, _ := range metrics.Labels {
			isDifferent := false
			for _, label := range labelsInit {
				if validatorLabels == label {
					isDifferent = true
				}
			}
			//检测label是否合法
			if !isDifferent {
				return true, fmt.Errorf("dataJson's labels have different")
			}
		}
		probixMetrics.With(metrics.Labels).Set(metrics.Value)
	}

	return true, nil
}
