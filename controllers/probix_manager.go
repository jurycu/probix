package controllers

import (
	v1prometheus "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ferulaxv1alpha1 "github.com/jurycu/probix/api/v1alpha1"
)

const (
	FINALIZER = "delete-probix"
)

func getDeleteInstance(instance ferulaxv1alpha1.Probix) *v1prometheus.PodMonitor {
	return &v1prometheus.PodMonitor{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "monitoring.coreos.com/v1",
			Kind:       "PodMonitor",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
	}
}

func getUpdatePodMonitor(instance ferulaxv1alpha1.Probix, oldPM v1prometheus.PodMonitor) *v1prometheus.PodMonitor {
	var pmes []v1prometheus.PodMetricsEndpoint
	for _, target := range instance.Spec.Targets {
		var pme v1prometheus.PodMetricsEndpoint
		pme.Port = "probix"
		pme.Interval = instance.Spec.Interval
		pme.Path = "/probix"
		pme.ScrapeTimeout = instance.Spec.ScrapeTimeout
		pme.Params = map[string][]string{
			"target":      []string{target.Target},
			"method":      []string{target.Method},
			"body":        []string{target.Body},
			"metricsName": []string{target.MetricsName},
			"metricsHelp": []string{target.MetricsHelp},
		}
		pmes = append(pmes, pme)
	}

	return &v1prometheus.PodMonitor{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "monitoring.coreos.com/v1",
			Kind:       "PodMonitor",
		},

		ObjectMeta: metav1.ObjectMeta{
			Name:            instance.Name,
			Namespace:       instance.Namespace,
			Labels:          instance.Labels,
			ResourceVersion: oldPM.ObjectMeta.ResourceVersion,
		},
		Spec: v1prometheus.PodMonitorSpec{
			Selector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"k8s-app": "probix",
				},
			},
			PodMetricsEndpoints: pmes,
			JobLabel:            instance.Name,
		},
	}
}

func getCreatePodMonitor(instance ferulaxv1alpha1.Probix) *v1prometheus.PodMonitor {
	var pmes []v1prometheus.PodMetricsEndpoint
	for _, target := range instance.Spec.Targets {
		var pme v1prometheus.PodMetricsEndpoint
		pme.Port = "probix"
		pme.Interval = instance.Spec.Interval
		pme.Path = "/probix"
		pme.ScrapeTimeout = instance.Spec.ScrapeTimeout
		pme.Params = map[string][]string{
			"target":      []string{target.Target},
			"method":      []string{target.Method},
			"body":        []string{target.Body},
			"metricsName": []string{target.MetricsName},
			"metricsHelp": []string{target.MetricsHelp},
		}
		pmes = append(pmes, pme)
	}

	return &v1prometheus.PodMonitor{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "monitoring.coreos.com/v1",
			Kind:       "PodMonitor",
		},

		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
			Labels:    instance.Labels,
		},
		Spec: v1prometheus.PodMonitorSpec{
			Selector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"k8s-app": "probix",
				},
			},
			PodMetricsEndpoints: pmes,
		},
	}
}
