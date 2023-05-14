package oval

import (
	"alt-oval-scanner/pkg/rpm"
	"alt-oval-scanner/pkg/utils"
	"regexp"
	"strconv"
)

func (o *Manager) OvalCheck(vulns map[string][]Vuln, oval OVAL, packs []rpm.Package, release utils.OsRelease) error {
	for _, p := range packs {
		if o.BaseCheck(oval, release) && o.Test(oval, p) {
			vuln := Vuln{Title: oval.Definitions.Definition[0].Metadata.Title}
			vuln.InstalledVersion = strconv.Itoa(p.Epoch) + ":" + p.Version + "-" + p.Release
			vuln.References = oval.Definitions.Definition[0].Metadata.References
			vuln.CVEs = oval.Definitions.Definition[0].Metadata.Advisory.Cves
			vuln.Severity = oval.Definitions.Definition[0].Metadata.Advisory.Severity
			vulns[p.Name] = append(vulns[p.Name], vuln)
		}
	}
	return nil
}

func (o *Manager) Test(oval OVAL, pack rpm.Package) bool {
	objets := parseObjets(oval.Objects)
	states := parseStates(oval.States)
	tests := parseTests(oval.Tests, objets, states)
	for _, d := range oval.Definitions.Definition {
		for _, c := range d.Criteria.Criterias {
			if o.testWalk(c, tests, pack) {
				return true
			}
		}
	}
	return false
}
func (o *Manager) BaseCheck(oval OVAL, release utils.OsRelease) bool {
	objects := parseBaseObjects(oval.Objects)
	states := parseBaseStates(oval.States)
	tests := parseBaseTest(oval.Tests, objects, states)
	for _, d := range oval.Definitions.Definition {
		//if d.Criteria.Operator
		for _, c := range d.Criteria.Criterions {
			test := tests[c.TestRef]
			if test.Object.Pattern.Operation == "pattern match" {
				exp := regexp.MustCompile(test.Object.Pattern.Text)
				if release.Check(exp) {
					return true
				}
			}
		}
	}
	return false
}

func parseBaseTest(tests Tests, objects map[string]TextFileContent54Object, states map[string]TextFileContent54State) map[string]AdjTestBase {
	var adjTests = map[string]AdjTestBase{}
	for _, t := range tests.TextFileContent54Tests {
		var adjTest AdjTestBase
		adjTest.ID = t.ID
		adjTest.Version = t.Version
		adjTest.Check = t.Check
		adjTest.Comment = t.Comment
		adjTest.Object = objects[t.Object.ObjectRef]
		adjTest.State = states[t.State.StateRef]
		adjTests[t.ID] = adjTest
	}
	return adjTests
}
func parseBaseObjects(objects Objects) map[string]TextFileContent54Object {
	var adjObjects = map[string]TextFileContent54Object{}
	for _, o := range objects.TextFileContent54Objects {
		adjObjects[o.ID] = o
	}
	return adjObjects
}
func parseBaseStates(states States) map[string]TextFileContent54State {
	var adjStates = map[string]TextFileContent54State{}
	for _, s := range states.TextFileContent54State {
		adjStates[s.ID] = s
	}
	return adjStates
}

func parseTests(tests Tests, objects map[string]RpmInfoObject, states map[string]RpmInfoState) map[string]AdjTest {
	var adjTests = map[string]AdjTest{}
	for _, t := range tests.RPMInfoTests {
		var adjTest AdjTest
		adjTest.ID = t.ID
		adjTest.Version = t.Version
		adjTest.Check = t.Check
		adjTest.Comment = t.Comment
		adjTest.Object = objects[t.Object.ObjectRef]
		adjTest.State = states[t.State.StateRef]
		adjTests[t.ID] = adjTest
	}
	return adjTests
}
func parseObjets(objects Objects) map[string]RpmInfoObject {
	var adjObjects = map[string]RpmInfoObject{}
	for _, o := range objects.RpmInfoObjects {
		adjObjects[o.ID] = o
	}
	return adjObjects
}
func parseStates(states States) map[string]RpmInfoState {
	var adjStates = map[string]RpmInfoState{}
	for _, s := range states.RpmInfoState {
		adjStates[s.ID] = s
	}
	return adjStates
}

func (o *Manager) testWalk(criteria Criteria, tests map[string]AdjTest, pack rpm.Package) bool {
	for _, c := range criteria.Criterias {
		res := o.testWalk(c, tests, pack)
		if res && criteria.Operator != "AND" {
			return res
		}
	}

	for _, c := range criteria.Criterions {
		test := tests[c.TestRef]
		res := checkObject(test.Object, pack) && checkState(test.State, pack)
		if criteria.Operator == "OR" && res {
			return res
		}
	}

	return false
}

func checkObject(object RpmInfoObject, pack rpm.Package) bool {
	if object.Name == pack.Name {
		return true
	}
	return false
}
func checkState(state RpmInfoState, pack rpm.Package) bool {
	//if state.Arch == pack.Arch
	version := rpm.NewVersion(state.Evr.Text)
	installedVersion := rpm.NewVersion(strconv.Itoa(pack.Epoch) + ":" + pack.Version + "-" + pack.Release)
	if state.Evr.Operation == "less than" && state.Evr.Datatype == "evr_string" {
		if installedVersion.LessThan(version) {
			return true
		}
	}
	return false
}
