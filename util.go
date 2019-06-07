package kubemap

import (
	"fmt"
	"strings"

	apps_v1beta1 "k8s.io/api/apps/v1beta1"
	apps_v1beta2 "k8s.io/api/apps/v1beta2"
	autoscaling_v1 "k8s.io/api/autoscaling/v1"
	batch_v1 "k8s.io/api/batch/v1"
	core_v1 "k8s.io/api/core/v1"
	ext_v1beta1 "k8s.io/api/extensions/v1beta1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ObjectMetaData returns metadata of a given k8s object
func objectMetaData(obj interface{}) meta_v1.ObjectMeta {
	//object := obj.(type)
	// switch object {
	switch object := obj.(type) {
	// case *apps_v1beta1.Deployment:
	// 	return object.ObjectMeta
	case *apps_v1beta2.Deployment:
		return object.ObjectMeta
	case *core_v1.ReplicationController:
		return object.ObjectMeta
	case *ext_v1beta1.ReplicaSet:
		return object.ObjectMeta
	case *apps_v1beta1.StatefulSet:
		return object.ObjectMeta
	case *ext_v1beta1.DaemonSet:
		return object.ObjectMeta
	case *core_v1.Service:
		return object.ObjectMeta
	case *core_v1.Pod:
		return object.ObjectMeta
	case *batch_v1.Job:
		return object.ObjectMeta
	case *core_v1.PersistentVolume:
		return object.ObjectMeta
	case *core_v1.PersistentVolumeClaim:
		return object.ObjectMeta
	case *core_v1.Namespace:
		return object.ObjectMeta
	case *core_v1.Secret:
		return object.ObjectMeta
	case *ext_v1beta1.Ingress:
		return object.ObjectMeta
	case *core_v1.Event:
		return object.ObjectMeta
	case *core_v1.ConfigMap:
		return object.ObjectMeta
	case *autoscaling_v1.HorizontalPodAutoscaler:
		return object.ObjectMeta
	}
	var objectMeta meta_v1.ObjectMeta
	return objectMeta
}

//RemoveDuplicateStrings returns unique string slice
func removeDuplicateStrings(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

//CopyMappedResource dep copies an object to create new one to avoid pointer references.
//This helps to keep store Thread safe even for Get operations
func copyMappedResource(resource MappedResource) MappedResource {

	copiedMappedResource := MappedResource{}

	for _, item := range resource.Kube.Ingresses {
		copiedMappedResource.Kube.Ingresses = append(copiedMappedResource.Kube.Ingresses, *item.DeepCopy())
	}

	for _, item := range resource.Kube.Services {
		copiedMappedResource.Kube.Services = append(copiedMappedResource.Kube.Services, *item.DeepCopy())
	}

	for _, item := range resource.Kube.Deployments {
		copiedMappedResource.Kube.Deployments = append(copiedMappedResource.Kube.Deployments, *item.DeepCopy())
	}

	for _, item := range resource.Kube.ReplicaSets {
		copiedMappedResource.Kube.ReplicaSets = append(copiedMappedResource.Kube.ReplicaSets, *item.DeepCopy())
	}

	for _, item := range resource.Kube.Pods {
		copiedMappedResource.Kube.Pods = append(copiedMappedResource.Kube.Pods, *item.DeepCopy())
	}

	copiedMappedResource.CommonLabel = resource.CommonLabel
	copiedMappedResource.CurrentType = resource.CurrentType
	copiedMappedResource.Namespace = resource.Namespace

	return copiedMappedResource
}

// func metaResourceKeyFunc(obj interface{}) (string, error) {
// 	object := obj.(MappedResource)

// 	if object.Services != nil {
// 		return object.Services[0].Namespace + "/" + object.CurrentType + "/" + object.Services[0].Name, nil
// 	} else if object.Deployments != nil {
// 		return object.Deployments[0].Namespace + "/" + object.CurrentType + "/" + object.Deployments[0].Name, nil
// 	} else if object.ReplicaSets != nil {
// 		return object.ReplicaSets[0].Namespace + "/" + object.CurrentType + "/" + object.ReplicaSets[0].Name, nil
// 	} else if object.Pods != nil {
// 		return object.Pods[0].Namespace + "/" + object.CurrentType + "/" + object.Pods[0].Name, nil
// 	} else if object.Ingresses != nil {
// 		//If just ingress object is created then there is nothing to map to it.
// 		//So there will only be one entry for Ingress
// 		return object.Ingresses[0].Namespace + "/" + object.CurrentType + "/" + object.Ingresses[0].Name, nil
// 	}

// 	return "", fmt.Errorf("Can't determine key for given object")
// }

// MetaNamespaceWithoutHashKeyFunc is a convenient default KeyFunc which knows how to make keys for MappedResource
func metaResourceKeyFunc(obj interface{}) (string, error) {
	object := obj.(MappedResource)

	if object.Kube.Services != nil {
		return object.Kube.Services[0].Namespace + "/" + object.CurrentType + "/" + object.Kube.Services[0].Name, nil
	} else if object.Kube.Deployments != nil {
		return object.Kube.Deployments[0].Namespace + "/" + object.CurrentType + "/" + object.Kube.Deployments[0].Name, nil
	} else if object.Kube.ReplicaSets != nil {
		return object.Kube.ReplicaSets[0].Namespace + "/" + object.CurrentType + "/" + object.Kube.ReplicaSets[0].Name, nil
	} else if object.Kube.Pods != nil {
		return object.Kube.Pods[0].Namespace + "/" + object.CurrentType + "/" + object.Kube.Pods[0].Name, nil
	} else if object.Kube.Ingresses != nil {
		//If just ingress object is created then there is nothing to map to it.
		//So there will only be one entry for Ingress
		return object.Kube.Ingresses[0].Namespace + "/" + object.CurrentType + "/" + object.Kube.Ingresses[0].Name, nil
	} else if object.Kube.Events != nil {
		return object.Kube.Events[0].Namespace + "/" + object.EventType + "/" + strings.ToLower(object.Kube.Events[0].InvolvedObject.Kind) + "/" + object.Kube.Events[0].Name, nil
	}

	return "", fmt.Errorf("Can't determine key for given object")
}
