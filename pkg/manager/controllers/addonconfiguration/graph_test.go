package addonconfiguration

import (
	"reflect"
	"testing"

	addonv1alpha1 "open-cluster-management.io/api/addon/v1alpha1"

	"open-cluster-management.io/addon-framework/pkg/addonmanager/addontesting"
)

type placementStrategy struct {
	configs  []addonv1alpha1.AddOnConfig
	clusters []string
}

func TestConfigurationGraph(t *testing.T) {
	cases := []struct {
		name                string
		defaultConfigs      []addonv1alpha1.ConfigMeta
		addons              []*addonv1alpha1.ManagedClusterAddOn
		placementStrategies []placementStrategy
		expected            []*addonNode
	}{
		{
			name:     "no output",
			expected: nil,
		},
		{
			name: "default config only",
			defaultConfigs: []addonv1alpha1.ConfigMeta{
				{ConfigGroupResource: addonv1alpha1.ConfigGroupResource{Group: "core", Resource: "Bar"}},
				{ConfigGroupResource: addonv1alpha1.ConfigGroupResource{Group: "core", Resource: "Foo"},
					DefaultConfig: &addonv1alpha1.ConfigReferent{Name: "test"}},
			},
			addons: []*addonv1alpha1.ManagedClusterAddOn{
				addontesting.NewAddon("test", "cluster1"),
				addontesting.NewAddon("test", "cluster2"),
			},
			expected: []*addonNode{
				{
					desiredConfigs: map[addonv1alpha1.ConfigGroupResource]addonv1alpha1.ConfigReference{
						{Group: "core", Resource: "Foo"}: {
							ConfigGroupResource: addonv1alpha1.ConfigGroupResource{Group: "core", Resource: "Foo"},
							ConfigReferent:      addonv1alpha1.ConfigReferent{Name: "test"},
						},
					},
					mca: addontesting.NewAddon("test", "cluster1"),
				},
				{
					desiredConfigs: map[addonv1alpha1.ConfigGroupResource]addonv1alpha1.ConfigReference{
						{Group: "core", Resource: "Foo"}: {
							ConfigGroupResource: addonv1alpha1.ConfigGroupResource{Group: "core", Resource: "Foo"},
							ConfigReferent:      addonv1alpha1.ConfigReferent{Name: "test"},
						},
					},
					mca: addontesting.NewAddon("test", "cluster2"),
				},
			},
		},
		{
			name: "with placement strategy",
			defaultConfigs: []addonv1alpha1.ConfigMeta{
				{ConfigGroupResource: addonv1alpha1.ConfigGroupResource{Group: "core", Resource: "Bar"},
					DefaultConfig: &addonv1alpha1.ConfigReferent{Name: "test"}},
				{ConfigGroupResource: addonv1alpha1.ConfigGroupResource{Group: "core", Resource: "Foo"},
					DefaultConfig: &addonv1alpha1.ConfigReferent{Name: "test"}},
			},
			addons: []*addonv1alpha1.ManagedClusterAddOn{
				addontesting.NewAddon("test", "cluster1"),
				addontesting.NewAddon("test", "cluster2"),
				addontesting.NewAddon("test", "cluster3"),
			},
			placementStrategies: []placementStrategy{
				{configs: []addonv1alpha1.AddOnConfig{
					{ConfigGroupResource: addonv1alpha1.ConfigGroupResource{Group: "core", Resource: "Bar"},
						ConfigReferent: addonv1alpha1.ConfigReferent{Name: "test1"}},
				}, clusters: []string{"cluster1"}},
				{configs: []addonv1alpha1.AddOnConfig{
					{ConfigGroupResource: addonv1alpha1.ConfigGroupResource{Group: "core", Resource: "Bar"},
						ConfigReferent: addonv1alpha1.ConfigReferent{Name: "test2"}},
					{ConfigGroupResource: addonv1alpha1.ConfigGroupResource{Group: "core", Resource: "Foo"},
						ConfigReferent: addonv1alpha1.ConfigReferent{Name: "test2"}},
				}, clusters: []string{"cluster2"}},
			},
			expected: []*addonNode{
				{
					desiredConfigs: map[addonv1alpha1.ConfigGroupResource]addonv1alpha1.ConfigReference{
						{Group: "core", Resource: "Bar"}: {
							ConfigGroupResource: addonv1alpha1.ConfigGroupResource{Group: "core", Resource: "Bar"},
							ConfigReferent:      addonv1alpha1.ConfigReferent{Name: "test1"},
						},
						{Group: "core", Resource: "Foo"}: {
							ConfigGroupResource: addonv1alpha1.ConfigGroupResource{Group: "core", Resource: "Foo"},
							ConfigReferent:      addonv1alpha1.ConfigReferent{Name: "test"},
						},
					},
					mca: addontesting.NewAddon("test", "cluster1"),
				},
				{
					desiredConfigs: map[addonv1alpha1.ConfigGroupResource]addonv1alpha1.ConfigReference{
						{Group: "core", Resource: "Bar"}: {
							ConfigGroupResource: addonv1alpha1.ConfigGroupResource{Group: "core", Resource: "Bar"},
							ConfigReferent:      addonv1alpha1.ConfigReferent{Name: "test2"},
						},
						{Group: "core", Resource: "Foo"}: {
							ConfigGroupResource: addonv1alpha1.ConfigGroupResource{Group: "core", Resource: "Foo"},
							ConfigReferent:      addonv1alpha1.ConfigReferent{Name: "test2"},
						},
					},
					mca: addontesting.NewAddon("test", "cluster2"),
				},
				{
					desiredConfigs: map[addonv1alpha1.ConfigGroupResource]addonv1alpha1.ConfigReference{
						{Group: "core", Resource: "Bar"}: {
							ConfigGroupResource: addonv1alpha1.ConfigGroupResource{Group: "core", Resource: "Bar"},
							ConfigReferent:      addonv1alpha1.ConfigReferent{Name: "test"},
						},
						{Group: "core", Resource: "Foo"}: {
							ConfigGroupResource: addonv1alpha1.ConfigGroupResource{Group: "core", Resource: "Foo"},
							ConfigReferent:      addonv1alpha1.ConfigReferent{Name: "test"},
						},
					},
					mca: addontesting.NewAddon("test", "cluster3"),
				},
			},
		},
		{
			name: "mca override",
			defaultConfigs: []addonv1alpha1.ConfigMeta{
				{ConfigGroupResource: addonv1alpha1.ConfigGroupResource{Group: "core", Resource: "Bar"},
					DefaultConfig: &addonv1alpha1.ConfigReferent{Name: "test"}},
				{ConfigGroupResource: addonv1alpha1.ConfigGroupResource{Group: "core", Resource: "Foo"},
					DefaultConfig: &addonv1alpha1.ConfigReferent{Name: "test"}},
			},
			addons: []*addonv1alpha1.ManagedClusterAddOn{
				newManagedClusterAddon("test", "cluster1", []addonv1alpha1.AddOnConfig{
					{ConfigGroupResource: addonv1alpha1.ConfigGroupResource{Group: "core", Resource: "Bar"},
						ConfigReferent: addonv1alpha1.ConfigReferent{Name: "test1"}},
				}, nil),
				addontesting.NewAddon("test", "cluster2"),
				addontesting.NewAddon("test", "cluster3"),
			},
			placementStrategies: []placementStrategy{
				{configs: []addonv1alpha1.AddOnConfig{
					{ConfigGroupResource: addonv1alpha1.ConfigGroupResource{Group: "core", Resource: "Foo"},
						ConfigReferent: addonv1alpha1.ConfigReferent{Name: "test1"}},
				}, clusters: []string{"cluster1"}},
				{configs: []addonv1alpha1.AddOnConfig{
					{ConfigGroupResource: addonv1alpha1.ConfigGroupResource{Group: "core", Resource: "Bar"},
						ConfigReferent: addonv1alpha1.ConfigReferent{Name: "test2"}},
					{ConfigGroupResource: addonv1alpha1.ConfigGroupResource{Group: "core", Resource: "Foo"},
						ConfigReferent: addonv1alpha1.ConfigReferent{Name: "test2"}},
				}, clusters: []string{"cluster2"}},
			},
			expected: []*addonNode{
				{
					desiredConfigs: map[addonv1alpha1.ConfigGroupResource]addonv1alpha1.ConfigReference{
						{Group: "core", Resource: "Bar"}: {
							ConfigGroupResource: addonv1alpha1.ConfigGroupResource{Group: "core", Resource: "Bar"},
							ConfigReferent:      addonv1alpha1.ConfigReferent{Name: "test1"},
						},
						{Group: "core", Resource: "Foo"}: {
							ConfigGroupResource: addonv1alpha1.ConfigGroupResource{Group: "core", Resource: "Foo"},
							ConfigReferent:      addonv1alpha1.ConfigReferent{Name: "test1"},
						},
					},
					mca: newManagedClusterAddon("test", "cluster1", []addonv1alpha1.AddOnConfig{
						{ConfigGroupResource: addonv1alpha1.ConfigGroupResource{Group: "core", Resource: "Bar"},
							ConfigReferent: addonv1alpha1.ConfigReferent{Name: "test1"}},
					}, nil),
				},
				{
					desiredConfigs: map[addonv1alpha1.ConfigGroupResource]addonv1alpha1.ConfigReference{
						{Group: "core", Resource: "Bar"}: {
							ConfigGroupResource: addonv1alpha1.ConfigGroupResource{Group: "core", Resource: "Bar"},
							ConfigReferent:      addonv1alpha1.ConfigReferent{Name: "test2"},
						},
						{Group: "core", Resource: "Foo"}: {
							ConfigGroupResource: addonv1alpha1.ConfigGroupResource{Group: "core", Resource: "Foo"},
							ConfigReferent:      addonv1alpha1.ConfigReferent{Name: "test2"},
						},
					},
					mca: addontesting.NewAddon("test", "cluster2"),
				},
				{
					desiredConfigs: map[addonv1alpha1.ConfigGroupResource]addonv1alpha1.ConfigReference{
						{Group: "core", Resource: "Bar"}: {
							ConfigGroupResource: addonv1alpha1.ConfigGroupResource{Group: "core", Resource: "Bar"},
							ConfigReferent:      addonv1alpha1.ConfigReferent{Name: "test"},
						},
						{Group: "core", Resource: "Foo"}: {
							ConfigGroupResource: addonv1alpha1.ConfigGroupResource{Group: "core", Resource: "Foo"},
							ConfigReferent:      addonv1alpha1.ConfigReferent{Name: "test"},
						},
					},
					mca: addontesting.NewAddon("test", "cluster3"),
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			graph := newGraph(c.defaultConfigs)
			for _, addon := range c.addons {
				graph.addAddonNode(addon)
			}
			for _, strategy := range c.placementStrategies {
				graph.addPlacementNode(strategy.configs, strategy.clusters)
			}

			actual := graph.addonToUpdate()
			if len(actual) != len(c.expected) {
				t.Errorf("output length is not correct, expected %v, got %v", len(c.expected), len(actual))
			}

			for _, ev := range c.expected {
				compared := false
				for _, v := range actual {
					if v == nil || ev == nil {
						t.Errorf("addonNode should not be nil")
					}
					if ev.mca != nil && v.mca != nil && ev.mca.Namespace == v.mca.Namespace {
						if !reflect.DeepEqual(v, ev) {
							t.Errorf("output is not correct, cluster %s, expected %v, got %v", v.mca.Namespace, ev, v)
						}
						compared = true
					}
				}

				if !compared {
					t.Errorf("not found addonNode %v", ev.mca)
				}
			}
		})
	}
}
