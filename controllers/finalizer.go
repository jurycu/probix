package controllers

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Finalizer interface {
	AddFinalizer(c client.Client, obj client.Object) error
	RemoveFinalizer(c client.Client, obj client.Object) error
}

type ProbixFinalizer struct{}

func (f ProbixFinalizer) AddFinalizer(c client.Client, obj client.Object) error {
	if obj.GetFinalizers() == nil || len(obj.GetFinalizers()) == 0 {
		metaMap := make(map[string]interface{})
		itemMap := make(map[string]interface{})
		itemMap["finalizers"] = []string{FINALIZER}
		metaMap["metadata"] = itemMap

		b, err := json.Marshal(metaMap)
		if err != nil {
			return err
		}
		_ = c.Patch(context.TODO(), obj, client.RawPatch(types.MergePatchType, b))
		return nil
	}
	return nil
}

func (f ProbixFinalizer) RemoveFinalizer(c client.Client, obj client.Object) error {
	if umap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj); err == nil {
		un := &unstructured.Unstructured{Object: umap}
		un.SetFinalizers([]string{})
		return c.Update(context.TODO(), un)
	}
	return fmt.Errorf("cannot convert this obj to unstructured")
}
