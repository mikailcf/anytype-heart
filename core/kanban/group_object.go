package kanban

import (
	"fmt"
	"sort"
	"strings"

	"github.com/anyproto/anytype-heart/core/domain"
	"github.com/anyproto/anytype-heart/pkg/lib/database"
	"github.com/anyproto/anytype-heart/pkg/lib/localstore/objectstore"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/anyproto/anytype-heart/util/slice"
)

type GroupObject struct {
	Key     domain.RelationKey
	store   objectstore.ObjectStore
	Records []database.Record
}

func (g *GroupObject) InitGroups(spaceID string, f *database.Filters) error {
	if spaceID == "" {
		return fmt.Errorf("spaceId is required")
	}

	// Filter records that have a non-empty value for this object relation
	filterObject := database.FiltersAnd{
		database.FilterNot{Filter: database.FilterEmpty{Key: g.Key}},
	}

	if f == nil {
		f = &database.Filters{FilterObj: filterObject}
	} else {
		f.FilterObj = database.FiltersAnd{f.FilterObj, filterObject}
	}

	records, err := g.store.SpaceIndex(spaceID).QueryRaw(f, 0, 0)
	if err != nil {
		return fmt.Errorf("init kanban by object, objectStore query error: %w", err)
	}

	g.Records = records

	return nil
}

func (g *GroupObject) MakeGroups() (GroupSlice, error) {
	var groups GroupSlice

	uniqMap := make(map[string]bool)

	// Collect all unique object IDs from the relation values
	allObjectIds := make(map[string]bool)
	for _, rec := range g.Records {
		objectIDs := rec.Details.GetStringList(g.Key)
		for _, id := range objectIDs {
			allObjectIds[id] = true
		}
	}

	// Single object groups - one group per unique object ID
	for objectID := range allObjectIds {
		if !uniqMap[objectID] {
			uniqMap[objectID] = true
			groups = append(groups, Group{
				Id:   objectID,
				Data: GroupData{Ids: []string{objectID}},
			})
		}
	}

	// Multiple object groups - for records with multiple objects
	for _, rec := range g.Records {
		objectIDs := slice.Filter(rec.Details.GetStringList(g.Key), func(objectID string) bool {
			return allObjectIds[objectID]
		})

		if len(objectIDs) > 1 {
			sort.Strings(objectIDs)
			hash := strings.Join(objectIDs, "")
			if !uniqMap[hash] {
				uniqMap[hash] = true
				groups = append(groups, Group{
					Id:   hash,
					Data: GroupData{Ids: objectIDs},
				})
			}
		}
	}

	return groups, nil
}

func (g *GroupObject) MakeDataViewGroups() ([]*model.BlockContentDataviewGroup, error) {
	var result []*model.BlockContentDataviewGroup

	groups, err := g.MakeGroups()
	if err != nil {
		return nil, err
	}

	sort.Sort(groups)

	// Use Tag value type since it has the same structure (Ids []string)
	for _, grp := range groups {
		result = append(result, &model.BlockContentDataviewGroup{
			Id: Hash(grp.Id),
			Value: &model.BlockContentDataviewGroupValueOfTag{
				Tag: &model.BlockContentDataviewTag{
					Ids: grp.Data.Ids,
				}},
		})
	}

	// Add empty group at the beginning
	result = append([]*model.BlockContentDataviewGroup{{
		Id: "empty",
		Value: &model.BlockContentDataviewGroupValueOfTag{
			Tag: &model.BlockContentDataviewTag{
				Ids: make([]string, 0),
			}},
	}}, result...)

	return result, nil
}
