package matchers_test

import (
	"testing"

	. "github.com/onsi/gomega"
	"github.com/xchapter7x/clarity/pkg/matchers"
)

func TestMatchers(t *testing.T) {
	RegisterTestingT(t)
	t.Run("m.ReadTerraform", func(t *testing.T) {
		m := &matchers.Match{}
		controlHCL := multiTFControl()
		t.Run("should read all the terraform files into the matcher", func(t *testing.T) {
			err := m.ReadTerraform("testdata/multitf")
			Expect(err).NotTo(HaveOccurred())
			Expect(m.HCLEntries).To(Equal(controlHCL))
		})
	})

	t.Run("m.AOfTypeNamed", func(t *testing.T) {
		t.Run("should return a filtered list of definitions", func(t *testing.T) {
			m := &matchers.Match{}
			m.ReadTerraform("testdata")
			err := m.AOfTypeNamed("google_compute_network", "resource", "my-custom-network")
			Expect(err).NotTo(HaveOccurred())
			Expect(m.MatchingEntries).To(BeEquivalentTo(controlInstanceMatch()))
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
			m := resourceMatcherByType("google_compute_network")
			Expect(m.MatchingEntries).To(BeEquivalentTo(controlInstanceMatch()))
		})
	})

	t.Run("m.AttributeExists", func(t *testing.T) {
		t.Run("when there are no matches", func(t *testing.T) {
			m := resourceMatcherByType("google_compute_subnetwork")
			Expect(m.AttributeExists("bbbbbasdf")).To(HaveOccurred())
		})

		t.Run("when there are matches", func(t *testing.T) {
			m := resourceMatcherByType("google_compute_network")
			Expect(m.AttributeExists("name")).NotTo(HaveOccurred())
			Expect(m.MatchingEntries).To(BeEquivalentTo(controlInstanceMatch()))
		})

		t.Run("when attributes is not derived from a output block", func(t *testing.T) {
			m := &matchers.Match{}
			err := m.ReadTerraform("testdata/outputs")
			Expect(err).NotTo(HaveOccurred())
			err = m.AOfType("rds_password", "output")
			Expect(err).NotTo(HaveOccurred())
			Expect(m.AttributeExists("sensitive")).NotTo(HaveOccurred())
			Expect(m.AttributeExists("value")).NotTo(HaveOccurred())
			Expect(m.MatchingEntries).To(BeEquivalentTo(controlOutputMatch()))
		})
	})

	t.Run("m.AttributeGreaterThan", func(t *testing.T) {
		t.Run("when there are no matches", func(t *testing.T) {
			m := resourceMatcherByType("blah_blah")
			Expect(m.AttributeGreaterThan("nadafield", 81)).To(HaveOccurred())
		})

		t.Run("when a valid attribute has non gt value", func(t *testing.T) {
			m := resourceMatcherByType("blah_blah")
			Expect(m.AttributeGreaterThan("port", 81)).To(HaveOccurred())
		})

		t.Run("when there are matches with a correct value", func(t *testing.T) {
			m := resourceMatcherByType("blah_blah")
			Expect(m.AttributeGreaterThan("port", 79)).NotTo(HaveOccurred())
			Expect(m.MatchingEntries).To(BeEquivalentTo(controlInstanceMatchIntString()))
		})

		t.Run("when there are matches with a correct value and the attr value is a int type", func(t *testing.T) {
			m := resourceMatcherByType("foo")
			Expect(m.AttributeGreaterThan("port", 79)).NotTo(HaveOccurred())
			Expect(m.MatchingEntries).To(BeEquivalentTo(controlInstanceMatchInt()))
		})
	})

	t.Run("m.ItOccursAtLeastTimes", func(t *testing.T) {
		t.Run("when we dont have enough to satisfy at least", func(t *testing.T) {
			m := resourceMatcherByType("blah_blah")
			Expect(m.ItOccursAtLeastTimes(5)).To(HaveOccurred())
		})

		t.Run("when we have enough to satisfy at least", func(t *testing.T) {
			m := resourceMatcherByType("blah_blah")
			Expect(m.ItOccursAtLeastTimes(1)).NotTo(HaveOccurred())
		})
	})

	t.Run("m.ItOccursAtMostTimes", func(t *testing.T) {
		t.Run("when we have too many to satisfy at most", func(t *testing.T) {
			m := resourceMatcherByType("blah_blah")
			Expect(m.ItOccursAtMostTimes(0)).To(HaveOccurred())
		})

		t.Run("when we have enough to satisfy at most", func(t *testing.T) {
			m := resourceMatcherByType("blah_blah")
			Expect(m.ItOccursAtMostTimes(3)).NotTo(HaveOccurred())
		})
	})

	t.Run("m.ItOccursExactlyTimes", func(t *testing.T) {
		t.Run("when we dont have exactly the count", func(t *testing.T) {
			m := resourceMatcherByType("blah_blah")
			Expect(m.ItOccursExactlyTimes(2)).To(HaveOccurred())
		})

		t.Run("when we have exactly the count", func(t *testing.T) {
			m := resourceMatcherByType("blah_blah")
			Expect(m.ItOccursExactlyTimes(1)).NotTo(HaveOccurred())
		})
	})

	t.Run("m.AttributeRegex", func(t *testing.T) {
		t.Run("when we do not have a attr value which satisfies the regex", func(t *testing.T) {
			m := resourceMatcherByType("blah_blah")
			Expect(m.AttributeRegex("port", "abc")).To(HaveOccurred())
		})

		t.Run("when our attribute value satisfies the regex", func(t *testing.T) {
			m := resourceMatcherByType("blah_blah")
			Expect(m.AttributeRegex("port", "8*")).NotTo(HaveOccurred())
		})

		t.Run("when our attribute is a complex object and it satisfies the regex", func(t *testing.T) {
			m := &matchers.Match{}
			err := m.ReadTerraform("testdata")
			Expect(err).NotTo(HaveOccurred())
			err = m.AOfTypeNamed("google_compute_firewall", "resource", "allow-all-internal")
			Expect(err).NotTo(HaveOccurred())
			Expect(m.AttributeExists("allow")).NotTo(HaveOccurred())
			Expect(m.AttributeRegex("allow", "protocol.*tcp")).NotTo(HaveOccurred())
		})
	})

	t.Run("m.AttributeLessThan", func(t *testing.T) {
		t.Run("when there are no matches", func(t *testing.T) {
			m := resourceMatcherByType("blah_blah")
			Expect(m.AttributeLessThan("nadafield", 79)).To(HaveOccurred())
		})

		t.Run("when a valid attribute has non gt value", func(t *testing.T) {
			m := resourceMatcherByType("blah_blah")
			Expect(m.AttributeLessThan("port", 79)).To(HaveOccurred())
		})

		t.Run("when there are matches with a correct value", func(t *testing.T) {
			m := resourceMatcherByType("blah_blah")
			Expect(m.AttributeLessThan("port", 81)).NotTo(HaveOccurred())
			Expect(m.MatchingEntries).To(BeEquivalentTo(controlInstanceMatchIntString()))
		})

		t.Run("when there are matches with a correct value and attr value is type int", func(t *testing.T) {
			m := resourceMatcherByType("foo")
			Expect(m.AttributeLessThan("port", 81)).NotTo(HaveOccurred())
			Expect(m.MatchingEntries).To(BeEquivalentTo(controlInstanceMatchInt()))
		})
	})

	t.Run("m.AttributeEqualsInt", func(t *testing.T) {
		t.Run("when there are no matches", func(t *testing.T) {
			m := resourceMatcherByType("foo")
			Expect(m.AttributeEqualsInt("bbbbbasdf", 0)).To(HaveOccurred())
		})

		t.Run("when a valid attribute has non equal value", func(t *testing.T) {
			m := resourceMatcherByType("foo")
			Expect(m.AttributeEqualsInt("port", 43)).To(HaveOccurred())
		})

		t.Run("when there are matches with a correct value", func(t *testing.T) {
			m := resourceMatcherByType("foo")
			Expect(m.AttributeEqualsInt("port", 80)).NotTo(HaveOccurred())
			Expect(m.MatchingEntries).To(BeEquivalentTo(controlInstanceMatchInt()))
		})
	})

	t.Run("m.AttributeDoesNotEqualInt", func(t *testing.T) {
		t.Run("when there are no matches", func(t *testing.T) {
			m := resourceMatcherByType("foo")
			Expect(m.AttributeDoesNotEqualInt("bbbbbasdf", 0)).To(HaveOccurred())
		})

		t.Run("when a valid attribute has non equal value", func(t *testing.T) {
			m := resourceMatcherByType("foo")
			Expect(m.AttributeDoesNotEqualInt("port", 80)).To(HaveOccurred())
		})

		t.Run("when there are matches with a correct value", func(t *testing.T) {
			m := resourceMatcherByType("foo")
			Expect(m.AttributeDoesNotEqualInt("port", 81)).NotTo(HaveOccurred())
			Expect(m.MatchingEntries).To(BeEquivalentTo(controlInstanceMatchInt()))
		})
	})

	t.Run("m.AttributeEquals", func(t *testing.T) {
		t.Run("when there are no matches", func(t *testing.T) {
			m := resourceMatcherByType("google_compute_network")
			Expect(m.AttributeEquals("bbbbbasdf", "my-custom-network")).To(HaveOccurred())
		})

		t.Run("when a valid attribute has non equal value", func(t *testing.T) {
			m := resourceMatcherByType("google_compute_network")
			Expect(m.AttributeEquals("name", "'wrong-name'")).To(HaveOccurred())
		})

		t.Run("when there are matches with a correct value", func(t *testing.T) {
			m := resourceMatcherByType("google_compute_network")
			Expect(m.AttributeEquals("name", "my-custom-network")).NotTo(HaveOccurred())
			Expect(m.MatchingEntries).To(BeEquivalentTo(controlInstanceMatch()))
		})
	})

	t.Run("m.AttributeDoesNotEqual", func(t *testing.T) {
		t.Run("when there are no matches", func(t *testing.T) {
			m := resourceMatcherByType("google_compute_network")
			Expect(m.AttributeDoesNotEqual("bbbbbasdf", "my-custom-network")).To(HaveOccurred())
		})

		t.Run("when a valid attribute has non equal value", func(t *testing.T) {
			m := resourceMatcherByType("google_compute_network")
			Expect(m.AttributeDoesNotEqual("name", "wrong-name")).NotTo(HaveOccurred())
			Expect(m.MatchingEntries).To(BeEquivalentTo(controlInstanceMatch()))
		})

		t.Run("when there are matches with a correct value", func(t *testing.T) {
			m := resourceMatcherByType("google_compute_network")
			Expect(m.AttributeDoesNotEqual("name", "my-custom-network")).To(HaveOccurred())
		})
	})
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

func resourceMatcherByType(resourceType string) *matchers.Match {
	m := &matchers.Match{}
	err := m.ReadTerraform("testdata")
	Expect(err).NotTo(HaveOccurred())
	err = m.AOfType(resourceType, "resource")
	Expect(err).NotTo(HaveOccurred())
	return m
}
