package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestPodFullyParsing(t *testing.T) {
	assert := assert.New(t)
	t1, _ := time.Parse(time.RFC822Z, "08 Mar 18 17:44 +0300")
	k8sPod := v1.Pod{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:              "details-v1-3618568057-dnkjp",
			CreationTimestamp: meta_v1.NewTime(t1),
			Labels:            map[string]string{"apps": "details", "version": "v1"},
			Annotations: map[string]string{"kubernetes.io/created-by": "{\"kind\":\"SerializedReference\",\"apiVersion\":\"v1\",\"reference\":{\"kind\":\"ReplicaSet\",\"namespace\":\"bookinfo\",\"name\":\"details-v1-3618568057\",\"uid\":\"\",\"apiVersion\":\"extensions\",\"resourceVersion\":\"9068\"}}",
				"sidecar.istio.io/status": "{\"version\":\"\",\"initContainers\":[\"istio-init\",\"enable-core-dump\"],\"containers\":[\"istio-proxy\"],\"volumes\":[\"istio-envoy\",\"istio-certs\"]}"}},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				v1.Container{Name: "details", Image: "whatever"},
				v1.Container{Name: "istio-proxy", Image: "docker.io/istio/proxy:0.7.1"},
			},
			InitContainers: []v1.Container{
				v1.Container{Name: "istio-init", Image: "docker.io/istio/proxy_init:0.7.1"},
				v1.Container{Name: "enable-core-dump", Image: "alpine"},
			},
		}}

	pod := Pod{}
	pod.Parse(&k8sPod)
	assert.Equal("details-v1-3618568057-dnkjp", pod.Name)
	assert.Equal("2018-03-08T17:44:00+03:00", pod.CreatedAt)
	assert.Equal(map[string]string{"apps": "details", "version": "v1"}, pod.Labels)
	assert.Equal(Reference{Name: "details-v1-3618568057", Kind: "ReplicaSet"}, pod.CreatedBy)
	assert.Len(pod.IstioContainers, 1)
	assert.Equal("istio-proxy", pod.IstioContainers[0].Name)
	assert.Equal("docker.io/istio/proxy:0.7.1", pod.IstioContainers[0].Image)
	assert.Len(pod.IstioInitContainers, 2)
	assert.Equal("istio-init", pod.IstioInitContainers[0].Name)
	assert.Equal("docker.io/istio/proxy_init:0.7.1", pod.IstioInitContainers[0].Image)
	assert.Equal("enable-core-dump", pod.IstioInitContainers[1].Name)
	assert.Equal("alpine", pod.IstioInitContainers[1].Image)
}

func TestPodParsingMissingImage(t *testing.T) {
	assert := assert.New(t)
	t1, _ := time.Parse(time.RFC822Z, "08 Mar 18 17:44 +0300")
	k8sPod := v1.Pod{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:              "details-v1-3618568057-dnkjp",
			CreationTimestamp: meta_v1.NewTime(t1),
			Labels:            map[string]string{"apps": "details", "version": "v1"},
			Annotations: map[string]string{"kubernetes.io/created-by": "{\"kind\":\"SerializedReference\",\"apiVersion\":\"v1\",\"reference\":{\"kind\":\"ReplicaSet\",\"namespace\":\"bookinfo\",\"name\":\"details-v1-3618568057\",\"uid\":\"\",\"apiVersion\":\"extensions\",\"resourceVersion\":\"9068\"}}",
				"sidecar.istio.io/status": "{\"version\":\"\",\"initContainers\":[\"istio-init\",\"enable-core-dump\"],\"containers\":[\"istio-proxy\"],\"volumes\":[\"istio-envoy\",\"istio-certs\"]}"}},
	}

	pod := Pod{}
	pod.Parse(&k8sPod)
	assert.Equal("details-v1-3618568057-dnkjp", pod.Name)
	assert.Equal("2018-03-08T17:44:00+03:00", pod.CreatedAt)
	assert.Equal(map[string]string{"apps": "details", "version": "v1"}, pod.Labels)
	assert.Equal(Reference{Name: "details-v1-3618568057", Kind: "ReplicaSet"}, pod.CreatedBy)
	assert.Len(pod.IstioContainers, 1)
	assert.Equal("istio-proxy", pod.IstioContainers[0].Name)
	assert.Equal("", pod.IstioContainers[0].Image)
	assert.Len(pod.IstioInitContainers, 2)
	assert.Equal("istio-init", pod.IstioInitContainers[0].Name)
	assert.Equal("", pod.IstioInitContainers[0].Image)
	assert.Equal("enable-core-dump", pod.IstioInitContainers[1].Name)
	assert.Equal("", pod.IstioInitContainers[1].Image)
}

func TestPodParsingMissingAnnotations(t *testing.T) {
	assert := assert.New(t)
	t1, _ := time.Parse(time.RFC822Z, "08 Mar 18 17:44 +0300")
	k8sPod := v1.Pod{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:              "details-v1-3618568057-dnkjp",
			CreationTimestamp: meta_v1.NewTime(t1),
			Labels:            map[string]string{"apps": "details", "version": "v1"},
		}}

	pod := Pod{}
	pod.Parse(&k8sPod)
	assert.Equal("details-v1-3618568057-dnkjp", pod.Name)
	assert.Equal("2018-03-08T17:44:00+03:00", pod.CreatedAt)
	assert.Equal(map[string]string{"apps": "details", "version": "v1"}, pod.Labels)
	assert.Equal(Reference{Name: "", Kind: ""}, pod.CreatedBy)
	assert.Len(pod.IstioContainers, 0)
	assert.Len(pod.IstioInitContainers, 0)
}

func TestPodParsingInvalidAnnotations(t *testing.T) {
	assert := assert.New(t)
	t1, _ := time.Parse(time.RFC822Z, "08 Mar 18 17:44 +0300")
	k8sPod := v1.Pod{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:              "details-v1-3618568057-dnkjp",
			CreationTimestamp: meta_v1.NewTime(t1),
			Labels:            map[string]string{"apps": "details", "version": "v1"},
			Annotations: map[string]string{"kubernetes.io/created-by": "{whoops! bad json!",
				"sidecar.istio.io/status": "{\"version\":\"\",\"initContainers\":[{\"badkey\": \"Ooops! Not expected!\"}]}"}},
	}

	pod := Pod{}
	pod.Parse(&k8sPod)
	assert.Equal("details-v1-3618568057-dnkjp", pod.Name)
	assert.Equal("2018-03-08T17:44:00+03:00", pod.CreatedAt)
	assert.Equal(map[string]string{"apps": "details", "version": "v1"}, pod.Labels)
	assert.Equal(Reference{Name: "", Kind: ""}, pod.CreatedBy)
	assert.Len(pod.IstioContainers, 0)
	assert.Len(pod.IstioInitContainers, 0)
}
