package classifier

import (
	"fmt"
	"log"
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

func Linkage(groups []*Group, mode int) []*Group {

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

	return newGroups

}

func Merge(matrix map[string]map[string]float64) {
	groups := convertToGroups(matrix)

	for len(groups) > 1 {
		groups = Linkage(groups, completeLinkageMode)
		log.Println(groups)
	}

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
