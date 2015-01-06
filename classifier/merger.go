package classifier

import (
	"fmt"
	"math"
	"strings"
)

type Group struct {
	Proximites map[*Group]float64
	Objects    []string
}

func NewGroup(firstObject string) *Group {
	return &Group{
		Objects:    []string{firstObject},
		Proximites: map[*Group]float64{},
	}
}

func (g *Group) String() string {
	return fmt.Sprintf("[%s]", strings.Join(g.Objects, ","))
}

func (g *Group) MinProximity() (*Group, float64) {

	var nearestGroup *Group
	var minProximity float64

	for group, proximity := range g.Proximites {
		if nearestGroup == nil || proximity < minProximity {
			nearestGroup = group
			minProximity = proximity
		}

	}

	return nearestGroup, minProximity
}

const (
	completeLinkageMode = iota
	singleLinkageMode   = iota
)

type mergeStep struct {
	level     float64
	normLevel float64
	groups    []*Group
}

func newMergeStep(groups []*Group, level float64) *mergeStep {

	//copy groups (to avoid rearrangement through pointer)
	copiedGroups := []*Group{}
	for _, group := range groups {
		newGroup := Group{Objects: group.Objects}
		copiedGroups = append(copiedGroups, &newGroup)
	}

	return &mergeStep{
		level:  level,
		groups: copiedGroups,
	}
}

func Linkage(groups []*Group, mode int) ([]*Group, *Group) {

	//get the two groups with the smallest proximity
	var minProximity float64
	var groupA, groupB *Group
	for _, group := range groups {
		targetGroup, proximity := group.MinProximity()

		if groupA == nil || proximity < minProximity {
			groupA = group
			groupB = targetGroup
			minProximity = proximity
		}
	}

	newGroups := []*Group{}

	groupA.Objects = append(groupA.Objects, groupB.Objects...)

	for _, group := range groups {
		if group != groupB {
			newGroups = append(newGroups, group)
		}

		if group != groupA && group != groupB {

			proximityA := group.Proximites[groupA]
			proximityB := group.Proximites[groupB]

			var newProximity float64
			switch mode {
			case singleLinkageMode:
				newProximity = math.Min(proximityA, proximityB)
			case completeLinkageMode:
				newProximity = math.Max(proximityA, proximityB)
			}

			group.Proximites[groupA] = newProximity
			groupA.Proximites[group] = newProximity
		}

		//delete refs to removed group
		delete(group.Proximites, groupB)
	}

	return newGroups, groupA

}

func Merge(matrix map[string]map[string]float64) []*Group {
	groups := convertToGroups(matrix)

	var newGroup *Group
	steps := map[int]*mergeStep{}

	steps[len(groups)] = &mergeStep{
		level: 0.0, //TODO replace
	}

	minLevel := 10000.0  //TODO replace
	maxLevel := -10000.0 //TODO replace

	for len(groups) > 1 {
		groups, newGroup = Linkage(groups, completeLinkageMode)

		var level float64 = -10000 //TODO replace
		for _, objectA := range newGroup.Objects {
			for _, objectB := range newGroup.Objects {
				if objectA == objectB {
					continue
				}
				proximity := matrix[objectA][objectB]
				if proximity > level {
					level = proximity
				}
			}
		}

		if level > maxLevel {
			maxLevel = level
		}
		if level < minLevel {
			minLevel = level
		}

		step := newMergeStep(groups, level)
		steps[len(groups)] = step
	}

	//rescale levels to 1-7
	for _, step := range steps {
		step.normLevel = 1 + (step.level-minLevel)*(7-1)/(maxLevel-minLevel)
		//log.Println(len(step.groups), step.level, step.normLevel)
	}

	clusterCount := testMojena(steps)
	return steps[clusterCount].groups
}

func testMojena(steps map[int]*mergeStep) int {

	for k := len(steps) - 1; k > 1; k-- {
		sum := 0.0
		for i := len(steps) - 1; i >= k; i-- {
			sum += steps[i].normLevel
		}
		vk := sum / float64(len(steps)-k)

		sum = 0.0
		for i := len(steps) - 1; i >= k; i-- {
			sum += math.Pow(steps[i].normLevel-vk, 2)
		}
		sk := math.Sqrt(sum * 1.0 / (float64(len(steps)-k) - 1.0))

		next := steps[k-1].normLevel
		test := (next - vk) / sk

		if math.IsNaN(test) == false && test > 2.75 {
			return k
		}

		/*log.Println("step", len(steps)-k)
		log.Println("cluster", k)
		log.Println("durchschnitt", vk)
		log.Println("standardabweichung", sk)
		log.Println("vk+1", next)
		log.Println("test", test)
		log.Println("---------")*/

	}

	return 0
}

//convert a matrix of objects into a 'one group one object' structure
func convertToGroups(matrix map[string]map[string]float64) []*Group {
	groups := []*Group{}
	tmpGroups := map[string]*Group{}

	//create group vor every object
	for name, proximities := range matrix {

		var group *Group
		var targetGroup *Group
		var exists bool
		if group, exists = tmpGroups[name]; exists == false {
			group = NewGroup(name)
			tmpGroups[name] = group
		}

		for target, proximity := range proximities {

			if target == name {
				continue
			}

			if targetGroup, exists = tmpGroups[target]; exists == false {
				targetGroup = NewGroup(target)
				tmpGroups[target] = targetGroup
			}
			group.Proximites[targetGroup] = proximity
		}

		groups = append(groups, group)
	}

	return groups
}
