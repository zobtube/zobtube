package controller

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/zobtube/zobtube/internal/model"
)

type admActorDuplicateEntry struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	VideoCount int64  `json:"video_count"`
	Thumbnail  bool   `json:"thumbnail"`
	Sex        string `json:"sex"`
}

type admActorDuplicateGroup struct {
	Name   string                   `json:"name"`
	Actors []admActorDuplicateEntry `json:"actors"`
}

func actorPairKey(id1, id2 string) string {
	a, b := model.NormalizeActorPair(id1, id2)
	return a + "|" + b
}

func loadDismissedActorPairKeys(db *gorm.DB) (map[string]struct{}, error) {
	var dismissed []model.ActorDismissedDuplicate
	if err := db.Find(&dismissed).Error; err != nil {
		return nil, err
	}
	keys := make(map[string]struct{}, len(dismissed))
	for _, d := range dismissed {
		keys[actorPairKey(d.ActorID1, d.ActorID2)] = struct{}{}
	}
	return keys, nil
}

func groupHasNonDismissedPair(actorIDs []string, dismissed map[string]struct{}) bool {
	for i := 0; i < len(actorIDs); i++ {
		for j := i + 1; j < len(actorIDs); j++ {
			if _, ok := dismissed[actorPairKey(actorIDs[i], actorIDs[j])]; !ok {
				return true
			}
		}
	}
	return false
}

func actorVideoCounts(db *gorm.DB, actorIDs []string) (map[string]int64, error) {
	counts := make(map[string]int64, len(actorIDs))
	if len(actorIDs) == 0 {
		return counts, nil
	}
	type row struct {
		ActorID string
		Count   int64
	}
	var rows []row
	err := db.Table("video_actors").
		Select("actor_id, COUNT(*) as count").
		Where("actor_id IN ?", actorIDs).
		Group("actor_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		counts[r.ActorID] = r.Count
	}
	return counts, nil
}

// AdmActorDuplicates godoc
//
//	@Summary	List actor name duplicate groups (admin)
//	@Tags		admin
//	@Produce	json
//	@Success	200	{object}	map[string]interface{}
//	@Router		/adm/actor/duplicates [get]
func (c *Controller) AdmActorDuplicates(g *gin.Context) {
	dismissed, err := loadDismissedActorPairKeys(c.datastore)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var actors []model.Actor
	if err := c.datastore.Where(NOT_DELETED).Order("LOWER(name), name").Find(&actors).Error; err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	byName := make(map[string][]model.Actor)
	nameOrder := make([]string, 0)
	for _, a := range actors {
		key := strings.ToLower(a.Name)
		if _, ok := byName[key]; !ok {
			nameOrder = append(nameOrder, key)
		}
		byName[key] = append(byName[key], a)
	}

	allIDs := make([]string, 0, len(actors))
	for _, a := range actors {
		allIDs = append(allIDs, a.ID)
	}
	videoCounts, err := actorVideoCounts(c.datastore, allIDs)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	groups := make([]admActorDuplicateGroup, 0)
	for _, nameKey := range nameOrder {
		members := byName[nameKey]
		if len(members) < 2 {
			continue
		}
		ids := make([]string, len(members))
		for i, m := range members {
			ids[i] = m.ID
		}
		if !groupHasNonDismissedPair(ids, dismissed) {
			continue
		}
		displayName := members[0].Name
		entries := make([]admActorDuplicateEntry, 0, len(members))
		for _, m := range members {
			entries = append(entries, admActorDuplicateEntry{
				ID:         m.ID,
				Name:       m.Name,
				VideoCount: videoCounts[m.ID],
				Thumbnail:  m.Thumbnail,
				Sex:        m.Sex,
			})
		}
		groups = append(groups, admActorDuplicateGroup{
			Name:   displayName,
			Actors: entries,
		})
	}

	g.JSON(http.StatusOK, gin.H{"groups": groups, "total": len(groups)})
}

// AdmActorDuplicateDismiss godoc
//
//	@Summary	Mark two actors as not duplicates (admin)
//	@Tags		admin
//	@Accept		json
//	@Param		body	body	object	true	"JSON with actor_id_1, actor_id_2"
//	@Success	201	{object}	map[string]interface{}
//	@Failure	400	{object}	map[string]interface{}
//	@Failure	404	{object}	map[string]interface{}
//	@Failure	409	{object}	map[string]interface{}
//	@Router		/adm/actor/duplicates/dismiss [post]
func (c *Controller) AdmActorDuplicateDismiss(g *gin.Context) {
	var body struct {
		ActorID1 string `json:"actor_id_1"`
		ActorID2 string `json:"actor_id_2"`
	}
	if err := g.ShouldBindJSON(&body); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id1 := strings.TrimSpace(body.ActorID1)
	id2 := strings.TrimSpace(body.ActorID2)
	if id1 == "" || id2 == "" {
		g.JSON(http.StatusBadRequest, gin.H{"error": "actor_id_1 and actor_id_2 are required"})
		return
	}
	if id1 == id2 {
		g.JSON(http.StatusBadRequest, gin.H{"error": "actor_id_1 and actor_id_2 must be different"})
		return
	}

	id1, id2 = model.NormalizeActorPair(id1, id2)

	var a1, a2 model.Actor
	if res := c.datastore.Where(NOT_DELETED).First(&a1, "id = ?", id1); res.RowsAffected < 1 {
		g.JSON(http.StatusNotFound, gin.H{"error": "actor not found"})
		return
	}
	if res := c.datastore.Where(NOT_DELETED).First(&a2, "id = ?", id2); res.RowsAffected < 1 {
		g.JSON(http.StatusNotFound, gin.H{"error": "actor not found"})
		return
	}

	var existing model.ActorDismissedDuplicate
	if res := c.datastore.Where("actor_id1 = ? AND actor_id2 = ?", id1, id2).First(&existing); res.RowsAffected > 0 {
		g.JSON(http.StatusConflict, gin.H{"error": "pair already dismissed", "id": existing.ID})
		return
	}

	record := &model.ActorDismissedDuplicate{
		ActorID1: id1,
		ActorID2: id2,
	}
	if err := c.datastore.Create(record).Error; err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusCreated, gin.H{"id": record.ID})
}

// AdmActorDuplicateDismissRemove godoc
//
//	@Summary	Remove dismissed duplicate pair (admin)
//	@Tags		admin
//	@Param		id	path	string	true	"Dismissed duplicate record ID"
//	@Success	204
//	@Failure	404	{object}	map[string]interface{}
//	@Router		/adm/actor/duplicates/dismiss/{id} [delete]
func (c *Controller) AdmActorDuplicateDismissRemove(g *gin.Context) {
	id := g.Param("id")
	record := &model.ActorDismissedDuplicate{ID: id}
	if res := c.datastore.First(record); res.RowsAffected < 1 {
		g.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if err := c.datastore.Delete(record).Error; err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusNoContent, gin.H{})
}

type admDismissedDuplicateItem struct {
	ID         string `json:"id"`
	ActorID1   string `json:"actor_id_1"`
	ActorID2   string `json:"actor_id_2"`
	ActorName1 string `json:"actor_name_1"`
	ActorName2 string `json:"actor_name_2"`
	CreatedAt  string `json:"created_at"`
}

// AdmActorDuplicateDismissedList godoc
//
//	@Summary	List dismissed actor duplicate pairs (admin)
//	@Tags		admin
//	@Produce	json
//	@Success	200	{object}	map[string]interface{}
//	@Router		/adm/actor/duplicates/dismissed [get]
func (c *Controller) AdmActorDuplicateDismissedList(g *gin.Context) {
	var records []model.ActorDismissedDuplicate
	if err := c.datastore.Order("created_at DESC").Find(&records).Error; err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	actorIDs := make(map[string]struct{})
	for _, r := range records {
		actorIDs[r.ActorID1] = struct{}{}
		actorIDs[r.ActorID2] = struct{}{}
	}
	idList := make([]string, 0, len(actorIDs))
	for id := range actorIDs {
		idList = append(idList, id)
	}

	names := make(map[string]string)
	if len(idList) > 0 {
		var actors []model.Actor
		_ = c.datastore.Unscoped().Where("id IN ?", idList).Find(&actors).Error
		for _, a := range actors {
			names[a.ID] = a.Name
		}
	}

	items := make([]admDismissedDuplicateItem, 0, len(records))
	for _, r := range records {
		items = append(items, admDismissedDuplicateItem{
			ID:         r.ID,
			ActorID1:   r.ActorID1,
			ActorID2:   r.ActorID2,
			ActorName1: names[r.ActorID1],
			ActorName2: names[r.ActorID2],
			CreatedAt:  r.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	g.JSON(http.StatusOK, gin.H{"items": items, "total": len(items)})
}
