package matchers_test

import (
	"sort"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/xchapter7x/clarity/pkg/matchers"
)

func TestMatchers(t *testing.T) {
	RegisterTestingT(t)
	for _, tt := range []struct {
		unmarshal matchers.Unmarshaller
		name      string
		skip      bool
	}{
		{name: "hcl1", unmarshal: matchers.GetUnmarshallerVersion(1)},
		{name: "hcl2", unmarshal: matchers.GetUnmarshallerVersion(2)},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skip("skipping", tt.name)
			}
			unmarshal := tt.unmarshal
			t.Run("m.ReadTerraform", func(t *testing.T) {
				m := &matchers.Match{}
				controlHCL := multiTFControl()
				t.Run("should read all the terraform files into the matcher", func(t *testing.T) {
					err := m.ReadTerraform("testdata/multitf", unmarshal)
					Expect(err).NotTo(HaveOccurred())
					sort.SliceStable(controlHCL, func(i, j int) bool { return controlHCL[i].InstanceName < controlHCL[j].InstanceName })
					sort.SliceStable(m.HCLEntries, func(i, j int) bool { return m.HCLEntries[i].InstanceName < m.HCLEntries[j].InstanceName })
					Expect(len(m.HCLEntries)).To(Equal(len(controlHCL)))
					attributesMatch(m.HCLEntries, controlHCL)
					Expect(m.HCLEntries[0].HCLType).To(BeEquivalentTo(controlHCL[0].HCLType))
					Expect(m.HCLEntries[0].ComponentName).To(BeEquivalentTo(controlHCL[0].ComponentName))
					Expect(m.HCLEntries[0].InstanceName).To(BeEquivalentTo(controlHCL[0].InstanceName))
				})
			})

			t.Run("m.AOfTypeNamed", func(t *testing.T) {
				t.Run("should return a filtered list of definitions", func(t *testing.T) {
					m := &matchers.Match{}
					m.ReadTerraform("testdata", unmarshal)
					err := m.AOfTypeNamed("google_compute_network", "resource", "my-custom-network")
					Expect(err).NotTo(HaveOccurred())
					control := controlInstanceMatch()
					Expect(len(m.MatchingEntries)).To(Equal(len(control)))
					attributesMatch(m.MatchingEntries, control)
					Expect(m.MatchingEntries[0].HCLType).To(BeEquivalentTo(control[0].HCLType))
					Expect(m.MatchingEntries[0].ComponentName).To(BeEquivalentTo(control[0].ComponentName))
					Expect(m.MatchingEntries[0].InstanceName).To(BeEquivalentTo(control[0].InstanceName))
				})
			})

			t.Run("m.AlwaysAttributeEqualsInt", func(t *testing.T) {
				t.Run("should fail but not panic if the attribute doesnt exist", func(t *testing.T) {
					m := &matchers.Match{MatchingEntries: controlMissingAttribute()}
					Expect(m.AlwaysAttributeEqualsInt("name", 5)).To(HaveOccurred())
				})

				t.Run("should fail if there are any elements that dont match", func(t *testing.T) {
					m := &matchers.Match{MatchingEntries: controlAlwaysMissInt()}
					Expect(m.AlwaysAttributeEqualsInt("name", 5)).To(HaveOccurred())
				})

				t.Run("should succeed all elements match", func(t *testing.T) {
					m := &matchers.Match{MatchingEntries: controlAlwaysMatchInt()}
					Expect(m.AlwaysAttributeEqualsInt("name", 5)).NotTo(HaveOccurred())
					Expect(m.MatchingEntries).To(BeEquivalentTo(controlAlwaysMatchInt()))
				})
			})

			t.Run("m.AlwaysAttributeEquals", func(t *testing.T) {
				t.Run("should fail but not panic if the attribute doesnt exist", func(t *testing.T) {
					m := &matchers.Match{MatchingEntries: controlMissingAttribute()}
					Expect(m.AlwaysAttributeEquals("name", "my-custom-network")).To(HaveOccurred())
				})

				t.Run("should fail if there are any elements that dont match", func(t *testing.T) {
					m := &matchers.Match{MatchingEntries: controlAlwaysMiss()}
					Expect(m.AlwaysAttributeEquals("name", "my-custom-network")).To(HaveOccurred())
				})

				t.Run("should succeed all elements match", func(t *testing.T) {
					m := &matchers.Match{MatchingEntries: controlAlwaysMatch()}
					Expect(m.AlwaysAttributeEquals("name", "my-custom-network")).NotTo(HaveOccurred())
					Expect(m.MatchingEntries).To(BeEquivalentTo(controlAlwaysMatch()))
				})
			})

			t.Run("m.AlwaysAttributeDoesNotEqual", func(t *testing.T) {
				t.Run("should fail but not panic if the attribute doesnt exist", func(t *testing.T) {
					m := &matchers.Match{MatchingEntries: controlMissingAttribute()}
					Expect(m.AlwaysAttributeDoesNotEqual("name", "def-doesnt-exist")).To(HaveOccurred())
				})

				t.Run("should fail if there are any elements that match", func(t *testing.T) {
					m := &matchers.Match{MatchingEntries: controlAlwaysMatch()}
					Expect(m.AlwaysAttributeDoesNotEqual("name", "my-custom-network")).To(HaveOccurred())
				})

				t.Run("should succeed NO elements match", func(t *testing.T) {
					m := &matchers.Match{MatchingEntries: controlAlwaysMiss()}
					Expect(m.AlwaysAttributeDoesNotEqual("name", "def-doesnt-exist")).NotTo(HaveOccurred())
					Expect(m.MatchingEntries).To(BeEquivalentTo(controlAlwaysMiss()))
				})
			})

			t.Run("m.AlwaysAttributeDoesNotEqualInt", func(t *testing.T) {
				t.Run("should fail but not panic if the attribute doesnt exist", func(t *testing.T) {
					m := &matchers.Match{MatchingEntries: controlMissingAttribute()}
					Expect(m.AlwaysAttributeDoesNotEqualInt("name", 5)).To(HaveOccurred())
				})

				t.Run("should fail if there are any elements that match", func(t *testing.T) {
					m := &matchers.Match{MatchingEntries: controlAlwaysMatchInt()}
					Expect(m.AlwaysAttributeDoesNotEqualInt("name", 5)).To(HaveOccurred())
				})

				t.Run("should succeed NO elements match", func(t *testing.T) {
					m := &matchers.Match{MatchingEntries: controlAlwaysMissInt()}
					Expect(m.AlwaysAttributeDoesNotEqualInt("name", 24)).NotTo(HaveOccurred())
					Expect(m.MatchingEntries).To(BeEquivalentTo(controlAlwaysMissInt()))
				})
			})

			t.Run("m.AlwaysAttributeRegex", func(t *testing.T) {
				t.Run("should fail but not panic if the attribute doesnt exist", func(t *testing.T) {
					m := &matchers.Match{MatchingEntries: controlMissingAttribute()}
					Expect(m.AlwaysAttributeRegex("name", "off")).To(HaveOccurred())
				})

				t.Run("should fail if there are any elements that dont match", func(t *testing.T) {
					m := &matchers.Match{MatchingEntries: controlAlwaysMiss()}
					Expect(m.AlwaysAttributeRegex("name", "my-custom-network")).To(HaveOccurred())
				})

				t.Run("should succeed when all elements match", func(t *testing.T) {
					m := &matchers.Match{MatchingEntries: controlAlwaysMatch()}
					Expect(m.AlwaysAttributeRegex("name", "my-custom-network")).NotTo(HaveOccurred())
					Expect(m.MatchingEntries).To(BeEquivalentTo(controlAlwaysMatch()))
				})
			})

			t.Run("m.AlwaysAttributeGreaterThan", func(t *testing.T) {
				t.Run("should fail but not panic if the attribute doesnt exist", func(t *testing.T) {
					m := &matchers.Match{MatchingEntries: controlMissingAttribute()}
					Expect(m.AlwaysAttributeGreaterThan("name", 5)).To(HaveOccurred())
				})

				t.Run("should fail if there are any elements that dont match", func(t *testing.T) {
					m := &matchers.Match{MatchingEntries: controlAlwaysMissInt()}
					Expect(m.AlwaysAttributeGreaterThan("name", 4)).To(HaveOccurred())
				})

				t.Run("should succeed when all elements match", func(t *testing.T) {
					m := &matchers.Match{MatchingEntries: controlAlwaysMatchInt()}
					Expect(m.AlwaysAttributeGreaterThan("name", 2)).NotTo(HaveOccurred())
					Expect(m.MatchingEntries).To(BeEquivalentTo(controlAlwaysMatchInt()))
				})
			})

			t.Run("m.AlwaysAttributeLessThan", func(t *testing.T) {
				t.Run("should fail but not panic if the attribute doesnt exist", func(t *testing.T) {
					m := &matchers.Match{MatchingEntries: controlMissingAttribute()}
					Expect(m.AlwaysAttributeLessThan("name", 5)).To(HaveOccurred())
				})

				t.Run("should fail if there are any elements that dont match", func(t *testing.T) {
					m := &matchers.Match{MatchingEntries: controlAlwaysMissInt()}
					Expect(m.AlwaysAttributeLessThan("name", 4)).To(HaveOccurred())
				})

				t.Run("should succeed when all elements match", func(t *testing.T) {
					m := &matchers.Match{MatchingEntries: controlAlwaysMatchInt()}
					Expect(m.AlwaysAttributeLessThan("name", 6)).NotTo(HaveOccurred())
					Expect(m.MatchingEntries).To(BeEquivalentTo(controlAlwaysMatchInt()))
				})
			})

			t.Run("m.AOfType", func(t *testing.T) {
				t.Run("should return a filtered list of definitions", func(t *testing.T) {
					m := resourceMatcherByType("google_compute_network", unmarshal)
					control := controlInstanceMatch()
					Expect(len(m.MatchingEntries)).To(Equal(len(control)))
					attributesMatch(m.MatchingEntries, control)
					Expect(m.MatchingEntries[0].HCLType).To(BeEquivalentTo(control[0].HCLType))
					Expect(m.MatchingEntries[0].ComponentName).To(BeEquivalentTo(control[0].ComponentName))
					Expect(m.MatchingEntries[0].InstanceName).To(BeEquivalentTo(control[0].InstanceName))
				})
			})

			t.Run("m.AttributeExists", func(t *testing.T) {
				t.Run("when there are no matches", func(t *testing.T) {
					m := resourceMatcherByType("google_compute_subnetwork", unmarshal)
					Expect(m.AttributeExists("bbbbbasdf")).To(HaveOccurred())
				})

				t.Run("when there are matches", func(t *testing.T) {
					m := resourceMatcherByType("google_compute_network", unmarshal)
					Expect(m.AttributeExists("name")).NotTo(HaveOccurred())
					control := controlInstanceMatch()
					Expect(len(m.MatchingEntries)).To(Equal(len(control)))
					attributesMatch(m.MatchingEntries, control)
					Expect(m.MatchingEntries[0].HCLType).To(BeEquivalentTo(control[0].HCLType))
					Expect(m.MatchingEntries[0].ComponentName).To(BeEquivalentTo(control[0].ComponentName))
					Expect(m.MatchingEntries[0].InstanceName).To(BeEquivalentTo(control[0].InstanceName))
				})

				t.Run("when attributes is not derived from a output block", func(t *testing.T) {
					m := &matchers.Match{}
					control := controlOutputMatch()
					err := m.ReadTerraform("testdata/outputs", unmarshal)
					Expect(err).NotTo(HaveOccurred())
					err = m.AOfType("rds_password", "output")
					Expect(err).NotTo(HaveOccurred())
					Expect(m.AttributeExists("sensitive")).NotTo(HaveOccurred())
					Expect(m.AttributeExists("value")).NotTo(HaveOccurred())
					Expect(len(m.MatchingEntries)).To(Equal(len(control)))
					attributesMatch(m.MatchingEntries, control)
					Expect(m.MatchingEntries[0].HCLType).To(BeEquivalentTo(control[0].HCLType))
					Expect(m.MatchingEntries[0].ComponentName).To(BeEquivalentTo(control[0].ComponentName))
					Expect(m.MatchingEntries[0].InstanceName).To(BeEquivalentTo(control[0].InstanceName))
				})
			})

			t.Run("m.AttributeGreaterThan", func(t *testing.T) {
				t.Run("when there are no matches", func(t *testing.T) {
					m := resourceMatcherByType("blah_blah", unmarshal)
					Expect(m.AttributeGreaterThan("nadafield", 81)).To(HaveOccurred())
				})

				t.Run("when a valid attribute has non gt value", func(t *testing.T) {
					m := resourceMatcherByType("blah_blah", unmarshal)
					Expect(m.AttributeGreaterThan("port", 81)).To(HaveOccurred())
				})

				t.Run("when there are matches with a correct value", func(t *testing.T) {
					m := resourceMatcherByType("blah_blah", unmarshal)
					Expect(m.AttributeGreaterThan("port", 79)).NotTo(HaveOccurred())
					control := controlInstanceMatchIntString()
					Expect(len(m.MatchingEntries)).To(Equal(len(control)))
					attributesMatch(m.MatchingEntries, control)
					Expect(m.MatchingEntries[0].HCLType).To(BeEquivalentTo(control[0].HCLType))
					Expect(m.MatchingEntries[0].ComponentName).To(BeEquivalentTo(control[0].ComponentName))
					Expect(m.MatchingEntries[0].InstanceName).To(BeEquivalentTo(control[0].InstanceName))
				})

				t.Run("when there are matches with a correct value and the attr value is a int type", func(t *testing.T) {
					m := resourceMatcherByType("foo", unmarshal)
					Expect(m.AttributeGreaterThan("port", 79)).NotTo(HaveOccurred())
					control := controlInstanceMatchInt()
					Expect(len(m.MatchingEntries)).To(Equal(len(control)))
					attributesMatch(m.MatchingEntries, control)
					Expect(m.MatchingEntries[0].HCLType).To(BeEquivalentTo(control[0].HCLType))
					Expect(m.MatchingEntries[0].ComponentName).To(BeEquivalentTo(control[0].ComponentName))
					Expect(m.MatchingEntries[0].InstanceName).To(BeEquivalentTo(control[0].InstanceName))
				})
			})

			t.Run("m.ItOccursAtLeastTimes", func(t *testing.T) {
				t.Run("when we dont have enough to satisfy at least", func(t *testing.T) {
					m := resourceMatcherByType("blah_blah", unmarshal)
					Expect(m.ItOccursAtLeastTimes(5)).To(HaveOccurred())
				})

				t.Run("when we have enough to satisfy at least", func(t *testing.T) {
					m := resourceMatcherByType("blah_blah", unmarshal)
					Expect(m.ItOccursAtLeastTimes(1)).NotTo(HaveOccurred())
				})
			})

			t.Run("m.ItOccursAtMostTimes", func(t *testing.T) {
				t.Run("when we have too many to satisfy at most", func(t *testing.T) {
					m := resourceMatcherByType("blah_blah", unmarshal)
					Expect(m.ItOccursAtMostTimes(0)).To(HaveOccurred())
				})

				t.Run("when we have enough to satisfy at most", func(t *testing.T) {
					m := resourceMatcherByType("blah_blah", unmarshal)
					Expect(m.ItOccursAtMostTimes(3)).NotTo(HaveOccurred())
				})
			})

			t.Run("m.ItOccursExactlyTimes", func(t *testing.T) {
				t.Run("when we dont have exactly the count", func(t *testing.T) {
					m := resourceMatcherByType("blah_blah", unmarshal)
					Expect(m.ItOccursExactlyTimes(2)).To(HaveOccurred())
				})

				t.Run("when we have exactly the count", func(t *testing.T) {
					m := resourceMatcherByType("blah_blah", unmarshal)
					Expect(m.ItOccursExactlyTimes(1)).NotTo(HaveOccurred())
				})
			})

			t.Run("m.AttributeRegex", func(t *testing.T) {
				t.Run("when we do not have a attr value which satisfies the regex", func(t *testing.T) {
					m := resourceMatcherByType("blah_blah", unmarshal)
					Expect(m.AttributeRegex("port", "abc")).To(HaveOccurred())
				})

				t.Run("when our attribute value satisfies the regex", func(t *testing.T) {
					m := resourceMatcherByType("blah_blah", unmarshal)
					Expect(m.AttributeRegex("port", "8*")).NotTo(HaveOccurred())
				})

				t.Run("when our attribute is a complex object and it satisfies the regex", func(t *testing.T) {
					m := &matchers.Match{}
					err := m.ReadTerraform("testdata", unmarshal)
					Expect(err).NotTo(HaveOccurred())
					err = m.AOfTypeNamed("google_compute_firewall", "resource", "allow-all-internal")
					Expect(err).NotTo(HaveOccurred())
					Expect(m.AttributeExists("allow")).NotTo(HaveOccurred())
					Expect(m.AttributeRegex("allow", "protocol.*tcp")).NotTo(HaveOccurred())
				})
			})

			t.Run("m.AttributeLessThan", func(t *testing.T) {
				t.Run("when there are no matches", func(t *testing.T) {
					m := resourceMatcherByType("blah_blah", unmarshal)
					Expect(m.AttributeLessThan("nadafield", 79)).To(HaveOccurred())
				})

				t.Run("when a valid attribute has non gt value", func(t *testing.T) {
					m := resourceMatcherByType("blah_blah", unmarshal)
					Expect(m.AttributeLessThan("port", 79)).To(HaveOccurred())
				})

				t.Run("when there are matches with a correct value", func(t *testing.T) {
					m := resourceMatcherByType("blah_blah", unmarshal)
					Expect(m.AttributeLessThan("port", 81)).NotTo(HaveOccurred())
					control := controlInstanceMatchIntString()
					Expect(len(m.MatchingEntries)).To(Equal(len(control)))
					attributesMatch(m.MatchingEntries, control)
					Expect(m.MatchingEntries[0].HCLType).To(BeEquivalentTo(control[0].HCLType))
					Expect(m.MatchingEntries[0].ComponentName).To(BeEquivalentTo(control[0].ComponentName))
					Expect(m.MatchingEntries[0].InstanceName).To(BeEquivalentTo(control[0].InstanceName))
				})

				t.Run("when there are matches with a correct value and attr value is type int", func(t *testing.T) {
					m := resourceMatcherByType("foo", unmarshal)
					control := controlInstanceMatchInt()
					Expect(m.AttributeLessThan("port", 81)).NotTo(HaveOccurred())
					Expect(len(m.MatchingEntries)).To(Equal(len(control)))
					attributesMatch(m.MatchingEntries, control)
					Expect(m.MatchingEntries[0].HCLType).To(BeEquivalentTo(control[0].HCLType))
					Expect(m.MatchingEntries[0].ComponentName).To(BeEquivalentTo(control[0].ComponentName))
					Expect(m.MatchingEntries[0].InstanceName).To(BeEquivalentTo(control[0].InstanceName))
				})
			})

			t.Run("m.AttributeEqualsInt", func(t *testing.T) {
				t.Run("when there are no matches", func(t *testing.T) {
					m := resourceMatcherByType("foo", unmarshal)
					Expect(m.AttributeEqualsInt("bbbbbasdf", 0)).To(HaveOccurred())
				})

				t.Run("when a valid attribute has non equal value", func(t *testing.T) {
					m := resourceMatcherByType("foo", unmarshal)
					Expect(m.AttributeEqualsInt("port", 43)).To(HaveOccurred())
				})

				t.Run("when there are matches with a correct value", func(t *testing.T) {
					m := resourceMatcherByType("foo", unmarshal)
					control := controlInstanceMatchInt()
					Expect(m.AttributeEqualsInt("port", 80)).NotTo(HaveOccurred())
					Expect(len(m.MatchingEntries)).To(Equal(len(control)))
					attributesMatch(m.MatchingEntries, control)
					Expect(m.MatchingEntries[0].HCLType).To(BeEquivalentTo(control[0].HCLType))
					Expect(m.MatchingEntries[0].ComponentName).To(BeEquivalentTo(control[0].ComponentName))
					Expect(m.MatchingEntries[0].InstanceName).To(BeEquivalentTo(control[0].InstanceName))
				})
			})

			t.Run("m.AttributeDoesNotEqualInt", func(t *testing.T) {
				t.Run("when there are no matches", func(t *testing.T) {
					m := resourceMatcherByType("foo", unmarshal)
					Expect(m.AttributeDoesNotEqualInt("bbbbbasdf", 0)).To(HaveOccurred())
				})

				t.Run("when a valid attribute has non equal value", func(t *testing.T) {
					m := resourceMatcherByType("foo", unmarshal)
					Expect(m.AttributeDoesNotEqualInt("port", 80)).To(HaveOccurred())
				})

				t.Run("when there are matches with a correct value", func(t *testing.T) {
					m := resourceMatcherByType("foo", unmarshal)
					control := controlInstanceMatchInt()
					Expect(m.AttributeDoesNotEqualInt("port", 81)).NotTo(HaveOccurred())
					Expect(len(m.MatchingEntries)).To(Equal(len(control)))
					attributesMatch(m.MatchingEntries, control)
					Expect(m.MatchingEntries[0].HCLType).To(BeEquivalentTo(control[0].HCLType))
					Expect(m.MatchingEntries[0].ComponentName).To(BeEquivalentTo(control[0].ComponentName))
					Expect(m.MatchingEntries[0].InstanceName).To(BeEquivalentTo(control[0].InstanceName))
				})
			})

			t.Run("m.AttributeEquals", func(t *testing.T) {
				t.Run("when there are no matches", func(t *testing.T) {
					m := resourceMatcherByType("google_compute_network", unmarshal)
					Expect(m.AttributeEquals("bbbbbasdf", "my-custom-network")).To(HaveOccurred())
				})

				t.Run("when a valid attribute has non equal value", func(t *testing.T) {
					m := resourceMatcherByType("google_compute_network", unmarshal)
					Expect(m.AttributeEquals("name", "'wrong-name'")).To(HaveOccurred())
				})

				t.Run("when there are matches with a correct value", func(t *testing.T) {
					m := resourceMatcherByType("google_compute_network", unmarshal)
					Expect(m.AttributeEquals("name", "my-custom-network")).NotTo(HaveOccurred())
					control := controlInstanceMatch()
					Expect(len(m.MatchingEntries)).To(Equal(len(control)))
					attributesMatch(m.MatchingEntries, control)
					Expect(m.MatchingEntries[0].HCLType).To(BeEquivalentTo(control[0].HCLType))
					Expect(m.MatchingEntries[0].ComponentName).To(BeEquivalentTo(control[0].ComponentName))
					Expect(m.MatchingEntries[0].InstanceName).To(BeEquivalentTo(control[0].InstanceName))
				})
			})

			t.Run("m.AttributeDoesNotEqual", func(t *testing.T) {
				t.Run("when there are no matches", func(t *testing.T) {
					m := resourceMatcherByType("google_compute_network", unmarshal)
					Expect(m.AttributeDoesNotEqual("bbbbbasdf", "my-custom-network")).To(HaveOccurred())
				})

				t.Run("when a valid attribute has non equal value", func(t *testing.T) {
					m := resourceMatcherByType("google_compute_network", unmarshal)
					Expect(m.AttributeDoesNotEqual("name", "wrong-name")).NotTo(HaveOccurred())
					control := controlInstanceMatch()
					Expect(len(m.MatchingEntries)).To(Equal(len(control)))
					attributesMatch(m.MatchingEntries, control)
					Expect(m.MatchingEntries[0].HCLType).To(BeEquivalentTo(control[0].HCLType))
					Expect(m.MatchingEntries[0].ComponentName).To(BeEquivalentTo(control[0].ComponentName))
					Expect(m.MatchingEntries[0].InstanceName).To(BeEquivalentTo(control[0].InstanceName))
				})

				t.Run("when there are matches with a correct value", func(t *testing.T) {
					m := resourceMatcherByType("google_compute_network", unmarshal)
					Expect(m.AttributeDoesNotEqual("name", "my-custom-network")).To(HaveOccurred())
				})
			})
		})
	}
}

func controlOutputMatch() []matchers.HCLEntry {
	return []matchers.HCLEntry{
		matchers.HCLEntry{
			HCLType:       "output",
			ComponentName: "rds_password",
			InstanceName:  "",
			Attributes: []map[string]interface{}{
				map[string]interface{}{
					"sensitive": true,
					"value":     "${element(concat(aws_db_instance.rds.*.password, list(\"\")), 0)}",
				},
			},
		},
	}
}

func controlMissingAttribute() []matchers.HCLEntry {
	return []matchers.HCLEntry{
		matchers.HCLEntry{
			HCLType:       "resource",
			ComponentName: "google_compute_network",
			InstanceName:  "my-custom-network",
			Attributes: []map[string]interface{}{
				map[string]interface{}{},
			},
		},
	}
}

func controlAlwaysMissInt() []matchers.HCLEntry {
	return []matchers.HCLEntry{
		matchers.HCLEntry{
			HCLType:       "resource",
			ComponentName: "google_compute_network",
			InstanceName:  "my-custom-network",
			Attributes: []map[string]interface{}{
				map[string]interface{}{
					"name": 5,
				},
			},
		},
		matchers.HCLEntry{
			HCLType:       "resource",
			ComponentName: "google_compute_network",
			InstanceName:  "my-custom-network",
			Attributes: []map[string]interface{}{
				map[string]interface{}{
					"name": 3,
				},
			},
		},
	}
}

func controlAlwaysMatchInt() []matchers.HCLEntry {
	return []matchers.HCLEntry{
		matchers.HCLEntry{
			HCLType:       "resource",
			ComponentName: "google_compute_network",
			InstanceName:  "my-custom-network",
			Attributes: []map[string]interface{}{
				map[string]interface{}{
					"name": 5,
				},
			},
		},
		matchers.HCLEntry{
			HCLType:       "resource",
			ComponentName: "google_compute_network",
			InstanceName:  "my-custom-network",
			Attributes: []map[string]interface{}{
				map[string]interface{}{
					"name": 5,
				},
			},
		},
	}
}

func controlAlwaysMiss() []matchers.HCLEntry {
	return []matchers.HCLEntry{
		matchers.HCLEntry{
			HCLType:       "resource",
			ComponentName: "google_compute_network",
			InstanceName:  "my-custom-network",
			Attributes: []map[string]interface{}{
				map[string]interface{}{
					"name": "my-custom-network",
				},
			},
		},
		matchers.HCLEntry{
			HCLType:       "resource",
			ComponentName: "google_compute_network",
			InstanceName:  "my-custom-network",
			Attributes: []map[string]interface{}{
				map[string]interface{}{
					"name": "epic fail",
				},
			},
		},
	}
}

func controlAlwaysMatch() []matchers.HCLEntry {
	return []matchers.HCLEntry{
		matchers.HCLEntry{
			HCLType:       "resource",
			ComponentName: "google_compute_network",
			InstanceName:  "my-custom-network",
			Attributes: []map[string]interface{}{
				map[string]interface{}{
					"name": "my-custom-network",
				},
			},
		},
		matchers.HCLEntry{
			HCLType:       "resource",
			ComponentName: "google_compute_network",
			InstanceName:  "my-custom-network",
			Attributes: []map[string]interface{}{
				map[string]interface{}{
					"name": "my-custom-network",
				},
			},
		},
	}
}

func controlInstanceMatch() []matchers.HCLEntry {
	return []matchers.HCLEntry{
		matchers.HCLEntry{
			HCLType:       "resource",
			ComponentName: "google_compute_network",
			InstanceName:  "my-custom-network",
			Attributes: []map[string]interface{}{
				map[string]interface{}{
					"name": "my-custom-network",
				},
			},
		},
	}
}

func controlInstanceMatchInt() []matchers.HCLEntry {
	return []matchers.HCLEntry{
		matchers.HCLEntry{
			HCLType:       "resource",
			ComponentName: "foo",
			InstanceName:  "bar",
			Attributes: []map[string]interface{}{
				map[string]interface{}{
					"port": 80,
				},
			},
		},
	}
}

func controlInstanceMatchIntString() []matchers.HCLEntry {
	return []matchers.HCLEntry{
		matchers.HCLEntry{
			HCLType:       "resource",
			ComponentName: "blah_blah",
			InstanceName:  "new_thing",
			Attributes: []map[string]interface{}{
				map[string]interface{}{
					"port": "80",
				},
			},
		},
	}
}

func multiTFControl() []matchers.HCLEntry {
	return []matchers.HCLEntry{
		matchers.HCLEntry{
			HCLType:       "resource",
			ComponentName: "aws",
			InstanceName:  "blah",
			Attributes: []map[string]interface{}{
				map[string]interface{}{},
			},
		},
		matchers.HCLEntry{
			HCLType:       "resource",
			ComponentName: "aws",
			InstanceName:  "bar",
			Attributes: []map[string]interface{}{
				map[string]interface{}{},
			},
		},
	}
}

func attributesMatch(object1, object2 interface{}) {
	var attr1, attr2 []map[string]interface{}
	switch attr := object1.(type) {
	case map[string]interface{}:
		attr1 = []map[string]interface{}{attr}
	case []map[string]interface{}:
		attr1 = attr
	}

	switch attr := object2.(type) {
	case map[string]interface{}:
		attr2 = []map[string]interface{}{attr}
	case []map[string]interface{}:
		attr2 = attr
	}
	Expect(attr1).To(BeEquivalentTo(attr2))
}

func resourceMatcherByType(resourceType string, unmarshal matchers.Unmarshaller) *matchers.Match {
	m := &matchers.Match{}
	err := m.ReadTerraform("testdata", unmarshal)
	Expect(err).NotTo(HaveOccurred())
	err = m.AOfType(resourceType, "resource")
	Expect(err).NotTo(HaveOccurred())
	return m
}
