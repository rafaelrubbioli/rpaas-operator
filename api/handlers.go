package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/tsuru/rpaas-operator/pkg/apis/extensions/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func serviceCreate(c echo.Context) error {
	// TODO: add validations (name and plan required, name doesn't exist, plan exists)
	name := c.FormValue("name")
	plan := c.FormValue("plan")
	annotations := map[string]string{
		"user":        c.FormValue("user"),
		"team":        c.FormValue("team"),
		"description": c.FormValue("description"),
		"eventid":     c.FormValue("eventid"),
	}

	instance := &v1alpha1.RpaasInstance{
		TypeMeta: metav1.TypeMeta{
			Kind:       "RpaasInstance",
			APIVersion: "extensions.tsuru.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   "default",
			Annotations: annotations,
		},
		Spec: v1alpha1.RpaasInstanceSpec{
			PlanName: plan,
		},
	}
	err := sdk.Create(instance)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusCreated, "")
}

func servicePlans(c echo.Context) error {
	list := &v1alpha1.RpaasPlanList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "RpaasPlan",
			APIVersion: "extensions.tsuru.io/v1alpha1",
		},
	}
	err := sdk.List("default", list)
	if err != nil {
		return c.String(http.StatusNotFound, "")
	}
	ret := make([]map[string]string, len(list.Items))
	for i, item := range list.Items {
		ret[i] = map[string]string{
			"name":        item.ObjectMeta.Name,
			"description": "no plan description",
		}
	}
	return c.JSON(http.StatusOK, ret)
}

func serviceInfo(c echo.Context) error {
	plan := &v1alpha1.RpaasInstance{
		TypeMeta: metav1.TypeMeta{
			Kind:       "RpaasInstance",
			APIVersion: "extensions.tsuru.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.Param("instance"),
			Namespace: "default",
		},
	}
	err := sdk.Get(plan)
	if err != nil {
		return c.String(http.StatusNotFound, "")
	}
	replicas := "0"
	if plan.Spec.Replicas != nil {
		replicas = fmt.Sprintf("%d", *plan.Spec.Replicas)
	}
	routes := make([]string, len(plan.Spec.Locations))
	for i, loc := range plan.Spec.Locations {
		routes[i] = loc.Config.Value
	}
	ret := []map[string]string{
		{
			"label": "Address",
			"value": "x.x.x.x",
		},
		{
			"label": "Instances",
			"value": replicas,
		},
		{
			"label": "Routes",
			"value": strings.Join(routes, "\n"),
		},
	}
	return c.JSON(http.StatusOK, ret)
}