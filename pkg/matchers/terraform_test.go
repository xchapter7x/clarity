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

	t.Run("m.AOfType", func(t *testing.T) {
		t.Run("should return a filtered list of definitions", func(t *testing.T) {
			m := &matchers.Match{}
			m.ReadTerraform("testdata")
			err := m.AOfType("google_compute_network", "resource")
			Expect(err).NotTo(HaveOccurred())
			Expect(m.MatchingEntries).To(BeEquivalentTo(controlInstanceMatch()))
		})
	})

	t.Run("m.AttributeExists", func(t *testing.T) {
		t.Run("when there are no matches", func(t *testing.T) {
			m := &matchers.Match{}
			err := m.ReadTerraform("testdata")
			Expect(err).NotTo(HaveOccurred())
			err = m.AOfType("google_compute_subnetwork", "resource")
			Expect(err).NotTo(HaveOccurred())
			Expect(m.AttributeExists("bbbbbasdf")).To(HaveOccurred())
		})

		t.Run("when there are matches", func(t *testing.T) {
			m := &matchers.Match{}
			err := m.ReadTerraform("testdata")
			Expect(err).NotTo(HaveOccurred())
			err = m.AOfType("google_compute_network", "resource")
			Expect(err).NotTo(HaveOccurred())
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
			m := &matchers.Match{}
			err := m.ReadTerraform("testdata")
			Expect(err).NotTo(HaveOccurred())
			err = m.AOfType("blah_blah", "resource")
			Expect(err).NotTo(HaveOccurred())
			Expect(m.AttributeGreaterThan("nadafield", 81)).To(HaveOccurred())
		})

		t.Run("when a valid attribute has non gt value", func(t *testing.T) {
			m := &matchers.Match{}
			err := m.ReadTerraform("testdata")
			Expect(err).NotTo(HaveOccurred())
			err = m.AOfType("blah_blah", "resource")
			Expect(err).NotTo(HaveOccurred())
			Expect(m.AttributeGreaterThan("port", 81)).To(HaveOccurred())
		})

		t.Run("when there are matches with a correct value", func(t *testing.T) {
			m := &matchers.Match{}
			err := m.ReadTerraform("testdata")
			Expect(err).NotTo(HaveOccurred())
			err = m.AOfType("blah_blah", "resource")
			Expect(err).NotTo(HaveOccurred())
			Expect(m.AttributeGreaterThan("port", 79)).NotTo(HaveOccurred())
			Expect(m.MatchingEntries).To(BeEquivalentTo(controlInstanceMatchIntString()))
		})

		t.Run("when there are matches with a correct value and the attr value is a int type", func(t *testing.T) {
			m := &matchers.Match{}
			err := m.ReadTerraform("testdata")
			Expect(err).NotTo(HaveOccurred())
			err = m.AOfType("foo", "resource")
			Expect(err).NotTo(HaveOccurred())
			Expect(m.AttributeGreaterThan("port", 79)).NotTo(HaveOccurred())
			Expect(m.MatchingEntries).To(BeEquivalentTo(controlInstanceMatchInt()))
		})
	})

	t.Run("m.ItOccursAtLeastTimes", func(t *testing.T) {
		t.Run("when we dont have enough to satisfy at least", func(t *testing.T) {
			m := &matchers.Match{}
			err := m.ReadTerraform("testdata")
			Expect(err).NotTo(HaveOccurred())
			err = m.AOfType("blah_blah", "resource")
			Expect(m.ItOccursAtLeastTimes(5)).To(HaveOccurred())
		})
		t.Run("when we have enough to satisfy at least", func(t *testing.T) {
			m := &matchers.Match{}
			err := m.ReadTerraform("testdata")
			Expect(err).NotTo(HaveOccurred())
			err = m.AOfType("blah_blah", "resource")
			Expect(err).NotTo(HaveOccurred())
			Expect(m.ItOccursAtLeastTimes(1)).NotTo(HaveOccurred())
		})
	})

	t.Run("m.ItOccursAtMostTimes", func(t *testing.T) {
		t.Run("when we have too many to satisfy at most", func(t *testing.T) {
			m := &matchers.Match{}
			err := m.ReadTerraform("testdata")
			Expect(err).NotTo(HaveOccurred())
			err = m.AOfType("blah_blah", "resource")
			Expect(m.ItOccursAtMostTimes(0)).To(HaveOccurred())
		})
		t.Run("when we have enough to satisfy at most", func(t *testing.T) {
			m := &matchers.Match{}
			err := m.ReadTerraform("testdata")
			Expect(err).NotTo(HaveOccurred())
			err = m.AOfType("blah_blah", "resource")
			Expect(err).NotTo(HaveOccurred())
			Expect(m.ItOccursAtMostTimes(3)).NotTo(HaveOccurred())
		})
	})

	t.Run("m.ItOccursExactlyTimes", func(t *testing.T) {
		t.Run("when we dont have exactly the count", func(t *testing.T) {
			m := &matchers.Match{}
			err := m.ReadTerraform("testdata")
			Expect(err).NotTo(HaveOccurred())
			err = m.AOfType("blah_blah", "resource")
			Expect(m.ItOccursExactlyTimes(2)).To(HaveOccurred())
		})
		t.Run("when we have exactly the count", func(t *testing.T) {
			m := &matchers.Match{}
			err := m.ReadTerraform("testdata")
			Expect(err).NotTo(HaveOccurred())
			err = m.AOfType("blah_blah", "resource")
			Expect(err).NotTo(HaveOccurred())
			Expect(m.ItOccursExactlyTimes(1)).NotTo(HaveOccurred())
		})
	})

	t.Run("m.AttributeRegex", func(t *testing.T) {
		t.Run("when we do not have a attr value which satisfies the regex", func(t *testing.T) {
			m := &matchers.Match{}
			err := m.ReadTerraform("testdata")
			Expect(err).NotTo(HaveOccurred())
			err = m.AOfType("blah_blah", "resource")
			Expect(m.AttributeRegex("port", "abc")).To(HaveOccurred())
		})
		t.Run("when our attribute value satisfies the regex", func(t *testing.T) {
			m := &matchers.Match{}
			err := m.ReadTerraform("testdata")
			Expect(err).NotTo(HaveOccurred())
			err = m.AOfType("blah_blah", "resource")
			Expect(err).NotTo(HaveOccurred())
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
			m := &matchers.Match{}
			err := m.ReadTerraform("testdata")
			Expect(err).NotTo(HaveOccurred())
			err = m.AOfType("blah_blah", "resource")
			Expect(err).NotTo(HaveOccurred())
			Expect(m.AttributeLessThan("nadafield", 79)).To(HaveOccurred())
		})

		t.Run("when a valid attribute has non gt value", func(t *testing.T) {
			m := &matchers.Match{}
			err := m.ReadTerraform("testdata")
			Expect(err).NotTo(HaveOccurred())
			err = m.AOfType("blah_blah", "resource")
			Expect(err).NotTo(HaveOccurred())
			Expect(m.AttributeLessThan("port", 79)).To(HaveOccurred())
		})

		t.Run("when there are matches with a correct value", func(t *testing.T) {
			m := &matchers.Match{}
			err := m.ReadTerraform("testdata")
			Expect(err).NotTo(HaveOccurred())
			err = m.AOfType("blah_blah", "resource")
			Expect(err).NotTo(HaveOccurred())
			Expect(m.AttributeLessThan("port", 81)).NotTo(HaveOccurred())
			Expect(m.MatchingEntries).To(BeEquivalentTo(controlInstanceMatchIntString()))
		})

		t.Run("when there are matches with a correct value and attr value is type int", func(t *testing.T) {
			m := &matchers.Match{}
			err := m.ReadTerraform("testdata")
			Expect(err).NotTo(HaveOccurred())
			err = m.AOfType("foo", "resource")
			Expect(err).NotTo(HaveOccurred())
			Expect(m.AttributeLessThan("port", 81)).NotTo(HaveOccurred())
			Expect(m.MatchingEntries).To(BeEquivalentTo(controlInstanceMatchInt()))
		})
	})

	t.Run("m.AttributeEqualsInt", func(t *testing.T) {
		t.Run("when there are no matches", func(t *testing.T) {
			m := &matchers.Match{}
			err := m.ReadTerraform("testdata")
			Expect(err).NotTo(HaveOccurred())
			err = m.AOfType("foo", "resource")
			Expect(err).NotTo(HaveOccurred())
			Expect(m.AttributeEqualsInt("bbbbbasdf", 0)).To(HaveOccurred())
		})

		t.Run("when a valid attribute has non equal value", func(t *testing.T) {
			m := &matchers.Match{}
			err := m.ReadTerraform("testdata")
			Expect(err).NotTo(HaveOccurred())
			err = m.AOfType("foo", "resource")
			Expect(err).NotTo(HaveOccurred())
			Expect(m.AttributeEqualsInt("port", 43)).To(HaveOccurred())
		})

		t.Run("when there are matches with a correct value", func(t *testing.T) {
			m := &matchers.Match{}
			err := m.ReadTerraform("testdata")
			Expect(err).NotTo(HaveOccurred())
			err = m.AOfType("foo", "resource")
			Expect(err).NotTo(HaveOccurred())
			Expect(m.AttributeEqualsInt("port", 80)).NotTo(HaveOccurred())
			Expect(m.MatchingEntries).To(BeEquivalentTo(controlInstanceMatchInt()))
		})
	})

	t.Run("m.AttributeDoesNotEqualInt", func(t *testing.T) {
		t.Run("when there are no matches", func(t *testing.T) {
			m := &matchers.Match{}
			err := m.ReadTerraform("testdata")
			Expect(err).NotTo(HaveOccurred())
			err = m.AOfType("foo", "resource")
			Expect(err).NotTo(HaveOccurred())
			Expect(m.AttributeDoesNotEqualInt("bbbbbasdf", 0)).To(HaveOccurred())
		})

		t.Run("when a valid attribute has non equal value", func(t *testing.T) {
			m := &matchers.Match{}
			err := m.ReadTerraform("testdata")
			Expect(err).NotTo(HaveOccurred())
			err = m.AOfType("foo", "resource")
			Expect(err).NotTo(HaveOccurred())
			Expect(m.AttributeDoesNotEqualInt("port", 80)).To(HaveOccurred())
		})

		t.Run("when there are matches with a correct value", func(t *testing.T) {
			m := &matchers.Match{}
			err := m.ReadTerraform("testdata")
			Expect(err).NotTo(HaveOccurred())
			err = m.AOfType("foo", "resource")
			Expect(err).NotTo(HaveOccurred())
			Expect(m.AttributeDoesNotEqualInt("port", 81)).NotTo(HaveOccurred())
			Expect(m.MatchingEntries).To(BeEquivalentTo(controlInstanceMatchInt()))
		})
	})

	t.Run("m.AttributeEquals", func(t *testing.T) {
		t.Run("when there are no matches", func(t *testing.T) {
			m := &matchers.Match{}
			err := m.ReadTerraform("testdata")
			Expect(err).NotTo(HaveOccurred())
			err = m.AOfType("google_compute_network", "resource")
			Expect(err).NotTo(HaveOccurred())
			Expect(m.AttributeEquals("bbbbbasdf", "my-custom-network")).To(HaveOccurred())
		})

		t.Run("when a valid attribute has non equal value", func(t *testing.T) {
			m := &matchers.Match{}
			err := m.ReadTerraform("testdata")
			Expect(err).NotTo(HaveOccurred())
			err = m.AOfType("google_compute_network", "resource")
			Expect(err).NotTo(HaveOccurred())
			Expect(m.AttributeEquals("name", "'wrong-name'")).To(HaveOccurred())
		})

		t.Run("when there are matches with a correct value", func(t *testing.T) {
			m := &matchers.Match{}
			err := m.ReadTerraform("testdata")
			Expect(err).NotTo(HaveOccurred())
			err = m.AOfType("google_compute_network", "resource")
			Expect(err).NotTo(HaveOccurred())
			Expect(m.AttributeEquals("name", "my-custom-network")).NotTo(HaveOccurred())
			Expect(m.MatchingEntries).To(BeEquivalentTo(controlInstanceMatch()))
		})
	})

	t.Run("m.AttributeDoesNotEqual", func(t *testing.T) {
		t.Run("when there are no matches", func(t *testing.T) {
			m := &matchers.Match{}
			err := m.ReadTerraform("testdata")
			Expect(err).NotTo(HaveOccurred())
			err = m.AOfType("google_compute_network", "resource")
			Expect(err).NotTo(HaveOccurred())
			Expect(m.AttributeDoesNotEqual("bbbbbasdf", "my-custom-network")).To(HaveOccurred())
		})

		t.Run("when a valid attribute has non equal value", func(t *testing.T) {
			m := &matchers.Match{}
			err := m.ReadTerraform("testdata")
			Expect(err).NotTo(HaveOccurred())
			err = m.AOfType("google_compute_network", "resource")
			Expect(err).NotTo(HaveOccurred())
			Expect(m.AttributeDoesNotEqual("name", "wrong-name")).NotTo(HaveOccurred())
			Expect(m.MatchingEntries).To(BeEquivalentTo(controlInstanceMatch()))
		})

		t.Run("when there are matches with a correct value", func(t *testing.T) {
			m := &matchers.Match{}
			err := m.ReadTerraform("testdata")
			Expect(err).NotTo(HaveOccurred())
			err = m.AOfType("google_compute_network", "resource")
			Expect(err).NotTo(HaveOccurred())
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
