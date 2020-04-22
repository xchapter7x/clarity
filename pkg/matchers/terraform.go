package matchers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl"
	hcl2 "github.com/instrumenta/conftest/parser/hcl2"
)

func NewMatch() *Match {
	return &Match{}
}

type Match struct {
	HCLEntries      []HCLEntry
	MatchingEntries []HCLEntry
}

type HCLEntry struct {
	HCLType       string
	ComponentName string
	InstanceName  string
	Attributes    interface{}
}

func parseHCL2(rawHCL map[string]interface{}) []HCLEntry {
	hclEntries := make([]HCLEntry, 0)
	for hclType, componentsArrayMap := range rawHCL {
		switch components := componentsArrayMap.(type) {
		case map[string]interface{}:
			hclEntries = addComponentsToEntries(hclEntries, hclType, components)

		case []map[string]interface{}:
			for _, componentsMap := range components {
				hclEntries = addComponentsToEntries(hclEntries, hclType, componentsMap)
			}
		}
	}
	return hclEntries
}

func addComponentsToEntries(hclEntries []HCLEntry, hclType string, componentsMap map[string]interface{}) []HCLEntry {
	for componentName, instancesFormatted := range componentsMap {
		if hclType == "output" || hclType == "module" || hclType == "provider" {
			hclEntries = append(hclEntries, HCLEntry{
				HCLType:       hclType,
				ComponentName: componentName,
				Attributes:    instancesFormatted,
			})
		} else {
			switch instancesCast := instancesFormatted.(type) {
			case []map[string]interface{}:
				instances := flattenArrayMap(instancesCast)
				for instanceName, attributes := range instances {
					hclEntries = append(hclEntries, HCLEntry{
						HCLType:       hclType,
						ComponentName: componentName,
						InstanceName:  instanceName,
						Attributes:    attributes,
					})
				}
			case map[string]interface{}:
				for instanceName, instance := range instancesCast {
					hclEntries = append(hclEntries, HCLEntry{
						HCLType:       hclType,
						ComponentName: componentName,
						InstanceName:  instanceName,
						Attributes:    instance,
					})
				}
			default:
				hclEntries = append(hclEntries, HCLEntry{
					HCLType:      hclType,
					InstanceName: componentName,
					Attributes:   instancesCast,
				})
			}
		}
	}
	return hclEntries
}

// Terrraform a simple matcher to show intent to init from terraform in
// current directory
func (m *Match) Terraform() error {
	pwd, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}
	return m.ReadTerraform(pwd, GetUnmarshallerVersion(1))
}

// HCL2 a simple matcher to show intent to init from a hcl2 terraform in
// current directory
func (m *Match) HCL2() error {
	pwd, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}
	return m.ReadTerraform(pwd, GetUnmarshallerVersion(2))
}

// AlwaysAttributeEqualsInt - requires all elements to have an exact match on attributes or it fails
func (m *Match) AlwaysAttributeEqualsInt(searchKey string, searchValue int) error {
	startingEntryCount := len(m.MatchingEntries)
	err := m.AttributeEqualsInt(searchKey, searchValue)
	if err != nil {
		return err
	}

	if startingEntryCount != len(m.MatchingEntries) {
		return fmt.Errorf("not all entries match on %v : %v", searchKey, searchValue)
	}
	return nil
}

// AlwaysAttributeEquals - requires all elements to have an exact match on attributes or it fails
func (m *Match) AlwaysAttributeEquals(searchKey, searchValue string) error {
	startingEntryCount := len(m.MatchingEntries)
	err := m.AttributeEquals(searchKey, searchValue)
	if err != nil {
		return err
	}

	if startingEntryCount != len(m.MatchingEntries) {
		return fmt.Errorf("not all entries match on %v : %v", searchKey, searchValue)
	}
	return nil
}

// AlwaysAttributeDoesNotEqualInt - all elements' attr of search key must not match the search value
func (m *Match) AlwaysAttributeDoesNotEqualInt(searchKey string, searchValue int) error {
	startingEntryCount := len(m.MatchingEntries)
	err := m.AttributeDoesNotEqualInt(searchKey, searchValue)
	if err != nil {
		return err
	}

	if startingEntryCount != len(m.MatchingEntries) {
		return fmt.Errorf("we found un-expected matches on %v : %v", searchKey, searchValue)
	}
	return nil
}

// AlwaysAttributeDoesNotEqual - all elements' attr of search key must not match the search value
func (m *Match) AlwaysAttributeDoesNotEqual(searchKey string, searchValue string) error {
	startingEntryCount := len(m.MatchingEntries)
	err := m.AttributeDoesNotEqual(searchKey, searchValue)
	if err != nil {
		return err
	}

	if startingEntryCount != len(m.MatchingEntries) {
		return fmt.Errorf("we found un-expected matches on %v : %v", searchKey, searchValue)
	}
	return nil
}

// AlwaysAttributeRegex - Check if all elements match regex
func (m *Match) AlwaysAttributeRegex(attributeName, regexString string) error {
	startingEntryCount := len(m.MatchingEntries)
	err := m.AttributeRegex(attributeName, regexString)
	if err != nil {
		return err
	}

	if startingEntryCount != len(m.MatchingEntries) {
		return fmt.Errorf("not every entry matches on %v : %v", attributeName, regexString)
	}
	return nil
}

func (m *Match) AlwaysAttributeGreaterThan(searchKey string, searchValue int) error {
	startingEntryCount := len(m.MatchingEntries)
	err := m.AttributeGreaterThan(searchKey, searchValue)
	if err != nil {
		return err
	}

	if startingEntryCount != len(m.MatchingEntries) {
		return fmt.Errorf("not every entry matches on %v : %v", searchKey, searchValue)
	}
	return nil
}

func (m *Match) AlwaysAttributeLessThan(searchKey string, searchValue int) error {
	startingEntryCount := len(m.MatchingEntries)
	err := m.AttributeLessThan(searchKey, searchValue)
	if err != nil {
		return err
	}

	if startingEntryCount != len(m.MatchingEntries) {
		return fmt.Errorf("not every entry matches on %v : %v", searchKey, searchValue)
	}
	return nil
}

type Unmarshaller func(p []byte, v interface{}) error

func GetUnmarshallerVersion(version int) Unmarshaller {
	if version == 1 {
		return hcl.Unmarshal
	}

	p := new(hcl2.Parser)
	return p.Unmarshal
}

// ReadTerrraform a simple matcher to init from terraform in
// a given directory
func (m *Match) ReadTerraform(tpath string, unmarshal Unmarshaller) error {
	baseHCL := make(map[string]interface{}, 0)
	dirContents := ""
	files, err := ioutil.ReadDir(tpath)
	if err != nil {
		return fmt.Errorf("could not read dir: %v", err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".tf") {
			contents, err := ioutil.ReadFile(tpath + "/" + file.Name())
			if err != nil {
				return fmt.Errorf("couldnt read file: %v", err)
			}

			dirContents += string(contents)
		}
	}

	err = unmarshal([]byte(dirContents), &baseHCL)
	if err != nil {
		return fmt.Errorf("hcl Unmarshal failed: %v", err)
	}

	m.HCLEntries = parseHCL2(baseHCL)
	return nil
}

func flattenArrayMap(listOfMaps []map[string]interface{}) map[string]interface{} {
	flattenedMap := make(map[string]interface{}, 0)
	for _, elementMap := range listOfMaps {
		for k, v := range elementMap {
			flattenedMap[k] = v
		}
	}
	return flattenedMap
}

// AOfTypeNamed will match a feature type, feature name and instance name
// exactly from the terraform
func (m *Match) AOfTypeNamed(providerFeature, providerFeatureType, instanceName string) error {
	matchingFeature, err := matchingFeatures(providerFeature, providerFeatureType, m.HCLEntries)
	if err != nil {
		entries, _ := json.Marshal(m.HCLEntries)
		return fmt.Errorf("no matches found for %v \n %v", providerFeature, string(entries))
	}

	matchingFeatureByInstanceName := make([]HCLEntry, 0)
	for _, feature := range matchingFeature {
		if feature.InstanceName == instanceName {
			matchingFeatureByInstanceName = append(matchingFeatureByInstanceName, feature)
		}
	}

	m.MatchingEntries = append(m.MatchingEntries, matchingFeatureByInstanceName...)
	return nil
}

// AOfTypeNamed will match a feature type and feature name exactly from the
// terraform. all instance names will be included
func (m *Match) AOfType(providerFeature, providerFeatureType string) error {
	matchingFeature, err := matchingFeatures(providerFeature, providerFeatureType, m.HCLEntries)
	if err != nil {
		entries, _ := json.Marshal(m.HCLEntries)
		return fmt.Errorf("no matches found for attribute %v \n %v", providerFeature, string(entries))
	}
	m.MatchingEntries = append(m.MatchingEntries, matchingFeature...)
	return nil
}

func matchingFeatures(providerFeature, providerFeatureType string, baseHCL []HCLEntry) ([]HCLEntry, error) {
	matchingEntries := make([]HCLEntry, 0)
	for _, entry := range baseHCL {
		if entry.HCLType == providerFeatureType && entry.ComponentName == providerFeature {
			matchingEntries = append(matchingEntries, entry)
		}
	}
	if len(matchingEntries) > 0 {
		return matchingEntries, nil
	}
	return nil, fmt.Errorf("no matches on %s %s found", providerFeatureType, providerFeature)
}

// AttributeExists will filter the set to only matching elements or it will
// return an error
func (m *Match) AttributeExists(searchKey string) error {
	var tmpEntries []HCLEntry
	if len(m.MatchingEntries) < 1 {
		return fmt.Errorf("no references to find attributes in")
	}

	for _, entry := range m.MatchingEntries {
		exists, _ := attributeExists(searchKey, entry.Attributes)
		if exists {
			tmpEntries = append(tmpEntries, entry)
		}
	}

	if len(tmpEntries) < 1 {
		entries, _ := json.Marshal(m.MatchingEntries)
		return fmt.Errorf("no matches found for attribute %v \n %v", searchKey, string(entries))
	}
	m.MatchingEntries = tmpEntries
	return nil
}

func attributeExists(searchName string, attributes interface{}) (bool, interface{}) {
	switch attributesCast := attributes.(type) {
	case []map[string]interface{}:
		for _, attributeMap := range attributesCast {
			for k, v := range attributeMap {
				if k == searchName {
					return true, v
				}
			}
		}
	case map[string]interface{}:
		for k, v := range attributesCast {
			if k == searchName {
				return true, v
			}
		}
	}
	return false, nil
}

// AttributeEqualsInt will filter on a full match of key value or it will
// return an error
func (m *Match) AttributeEqualsInt(searchKey string, searchValue int) error {
	var tmpEntries []HCLEntry
	if len(m.MatchingEntries) < 1 {
		return fmt.Errorf("no references to find attributes in")
	}

	for _, entry := range m.MatchingEntries {
		exists, attributeValue := attributeExists(searchKey, entry.Attributes)
		var attributeInteger int
		switch v := attributeValue.(type) {
		case int:
			attributeInteger = v
		case float64:
			attributeInteger = int(v)
		}

		if exists && attributeInteger == searchValue {
			tmpEntries = append(tmpEntries, entry)
		}
	}

	if len(tmpEntries) < 1 {
		entries, _ := json.Marshal(m.MatchingEntries)
		return fmt.Errorf("no matches found for attribute %v \n %v", searchValue, string(entries))
	}
	m.MatchingEntries = tmpEntries
	return nil
}

// AttributeDoesNotEqualInt will filter on a full match of key value or it will
// return an error
func (m *Match) AttributeDoesNotEqualInt(searchKey string, searchValue int) error {
	var tmpEntries []HCLEntry
	if len(m.MatchingEntries) < 1 {
		return fmt.Errorf("no references to find attributes in")
	}

	for _, entry := range m.MatchingEntries {
		exists, attributeValue := attributeExists(searchKey, entry.Attributes)
		var attributeInteger int
		switch v := attributeValue.(type) {
		case int:
			attributeInteger = v
		case float64:
			attributeInteger = int(v)
		}

		if exists && attributeInteger != searchValue {
			tmpEntries = append(tmpEntries, entry)
		}
	}

	if len(tmpEntries) < 1 {
		entries, _ := json.Marshal(m.MatchingEntries)
		return fmt.Errorf("no matches found for attribute %v \n %v", searchValue, string(entries))
	}
	m.MatchingEntries = tmpEntries
	return nil
}

// AttributeEquals will filter on a full match of key value or it will
// return an error
func (m *Match) AttributeEquals(searchKey, searchValue string) error {
	var tmpEntries []HCLEntry
	if len(m.MatchingEntries) < 1 {
		return fmt.Errorf("no references to find attributes in")
	}

	for _, entry := range m.MatchingEntries {
		exists, attributeValue := attributeExists(searchKey, entry.Attributes)
		if exists && attributeValue == searchValue {
			tmpEntries = append(tmpEntries, entry)
		}
	}

	if len(tmpEntries) < 1 {
		entries, _ := json.Marshal(m.MatchingEntries)
		return fmt.Errorf("no matches found for attribute %v \n %v", searchValue, string(entries))
	}
	m.MatchingEntries = tmpEntries
	return nil
}

// AttributeEquals will filter on a match of key where value is not a match
// or it will return an error
func (m *Match) AttributeDoesNotEqual(searchKey, searchValue string) error {
	var tmpEntries []HCLEntry
	if len(m.MatchingEntries) < 1 {
		return fmt.Errorf("no references to find attributes in")
	}

	for _, entry := range m.MatchingEntries {
		exists, attributeValue := attributeExists(searchKey, entry.Attributes)
		if exists && attributeValue != searchValue {
			tmpEntries = append(tmpEntries, entry)
		}
	}

	if len(tmpEntries) < 1 {
		entries, _ := json.Marshal(m.MatchingEntries)
		return fmt.Errorf("no matches found for attribute %v \n %v", searchValue, string(entries))
	}
	m.MatchingEntries = tmpEntries
	return nil
}

// AttributeGreaterThan - will match (assuming int) for a greater than evaluation
func (m *Match) AttributeGreaterThan(searchKey string, searchValue int) error {
	var tmpEntries []HCLEntry
	if len(m.MatchingEntries) < 1 {
		return fmt.Errorf("no references to find attributes in")
	}

	for _, entry := range m.MatchingEntries {
		exists, attributeValue := attributeExists(searchKey, entry.Attributes)
		if !exists {
			return fmt.Errorf("no attribute found named: %s", searchKey)
		}
		var actualValue int
		switch attributeValue := attributeValue.(type) {
		case string:
			var err error
			actualValue, err = strconv.Atoi(attributeValue)
			if err != nil {
				return fmt.Errorf("could not translate to int: %v", err)
			}
		case int:
			actualValue = attributeValue
		case float64:
			actualValue = int(attributeValue)
		default:
			return fmt.Errorf("could not translate to int")
		}

		if exists && actualValue > searchValue {
			tmpEntries = append(tmpEntries, entry)
		}
	}

	if len(tmpEntries) < 1 {
		entries, _ := json.Marshal(m.MatchingEntries)
		return fmt.Errorf("no matches found for attribute %v \n %v", searchValue, string(entries))
	}
	m.MatchingEntries = tmpEntries
	return nil
}

// AttributeLessThan - will match (assuming int) for a less than evaluation
func (m *Match) AttributeLessThan(searchKey string, searchValue int) error {
	var tmpEntries []HCLEntry
	if len(m.MatchingEntries) < 1 {
		return fmt.Errorf("no references to find attributes in")
	}

	for _, entry := range m.MatchingEntries {
		exists, attributeValue := attributeExists(searchKey, entry.Attributes)
		if !exists {
			return fmt.Errorf("no attribute found named: %s", searchKey)
		}

		var actualValue int
		switch attributeValue := attributeValue.(type) {
		case string:
			var err error
			actualValue, err = strconv.Atoi(attributeValue)
			if err != nil {
				return fmt.Errorf("could not translate to int: %v", err)
			}
		case int:
			actualValue = attributeValue
		}

		if exists && actualValue < searchValue {
			tmpEntries = append(tmpEntries, entry)
		}
	}

	if len(tmpEntries) < 1 {
		entries, _ := json.Marshal(m.MatchingEntries)
		return fmt.Errorf("no matches found for attribute %v \n %v", searchValue, string(entries))
	}
	m.MatchingEntries = tmpEntries
	return nil
}

// ItOccursAtLeastTimes - will check how many results are in our filter and make sure we have at least that
func (m *Match) ItOccursAtLeastTimes(count int) error {
	if len(m.MatchingEntries) >= count {
		return nil
	}
	entries, _ := json.Marshal(m.MatchingEntries)
	return fmt.Errorf("no matches found for attribute %v \n %v", count, string(entries))
}

// ItOccursAtMostTimes - will check how many results are in our filter and make sure we have at most that
func (m *Match) ItOccursAtMostTimes(count int) error {
	if len(m.MatchingEntries) <= count {
		return nil
	}
	entries, _ := json.Marshal(m.MatchingEntries)
	return fmt.Errorf("no matches found for attribute %v \n %v", count, string(entries))
}

// ItOccursExactlyTimes - will check how many results are in our filter and make sure we have at exactly that
func (m *Match) ItOccursExactlyTimes(count int) error {
	if len(m.MatchingEntries) == count {
		return nil
	}
	entries, _ := json.Marshal(m.MatchingEntries)
	return fmt.Errorf("no matches found for attribute %v \n %v", count, string(entries))
}

// AttributeRegex - will use a regex to see if attributes value is a match
func (m *Match) AttributeRegex(attributeName, regexString string) error {
	var tmpEntries []HCLEntry
	if len(m.MatchingEntries) < 1 {
		return fmt.Errorf("no references to find attributes in")
	}

	for _, entry := range m.MatchingEntries {
		exists, attributeValue := attributeExists(attributeName, entry.Attributes)
		if !exists {
			return fmt.Errorf("no attribute found named: %s", attributeName)
		}

		actualValue, err := json.Marshal(attributeValue)
		if err != nil {
			return fmt.Errorf("unmarshaling json failed: %v", err)
		}

		var validAttributeValue = regexp.MustCompile(regexString)
		if validAttributeValue.MatchString(string(actualValue)) {
			tmpEntries = append(tmpEntries, entry)
		}
	}

	if len(tmpEntries) < 1 {
		entries, _ := json.Marshal(m.MatchingEntries)
		return fmt.Errorf("no matches found for attribute %s \n %v", attributeName, string(entries))
	}
	m.MatchingEntries = tmpEntries
	return nil
}
