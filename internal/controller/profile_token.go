package controller

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zobtube/zobtube/internal/model"
)

// ResolveUserByApiTokenHash looks up an API token by the hash of the raw token and returns the associated user if found.
// The tokenArg is the raw Bearer token; it will be hashed (SHA256) and looked up.
func (c *Controller) ResolveUserByApiTokenHash(tokenArg string) (*model.User, bool) {
	if tokenArg == "" {
		return nil, false
	}
	sum := sha256.Sum256([]byte(tokenArg))
	hash := hex.EncodeToString(sum[:])
	var apiToken model.ApiToken
	if c.datastore.Where("token_hash = ?", hash).First(&apiToken).RowsAffected < 1 {
		return nil, false
	}
	user := &model.User{ID: apiToken.UserID}
	if c.datastore.First(user).RowsAffected < 1 {
		return nil, false
	}
	return user, true
}

// ProfileTokenList returns the list of API tokens for the current user (id, name, created_at only).
func (c *Controller) ProfileTokenList(g *gin.Context) {
	user := g.MustGet("user").(*model.User)
	var tokens []model.ApiToken
	c.datastore.Where("user_id = ?", user.ID).Order("created_at desc").Find(&tokens)
	list := make([]gin.H, 0, len(tokens))
	for _, t := range tokens {
		list = append(list, gin.H{
			"id":         t.ID,
			"name":       t.Name,
			"created_at": t.CreatedAt,
		})
	}
	g.JSON(http.StatusOK, gin.H{"tokens": list})
}

// ProfileTokenCreate creates a new API token for the current user. Body: { "name": "label" }. Returns the raw token only in this response.
func (c *Controller) ProfileTokenCreate(g *gin.Context) {
	user := g.MustGet("user").(*model.User)
	var body struct {
		Name string `json:"name"`
	}
	if err := g.ShouldBindJSON(&body); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	name := body.Name
	if name == "" {
		g.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}
	rawToken := hex.EncodeToString(raw)
	sum := sha256.Sum256([]byte(rawToken))
	tokenHash := hex.EncodeToString(sum[:])
	apiToken := &model.ApiToken{
		UserID:    user.ID,
		Name:      name,
		TokenHash: tokenHash,
		CreatedAt: time.Now(),
	}
	if err := c.datastore.Create(apiToken).Error; err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save"})
		return
	}
	g.JSON(http.StatusCreated, gin.H{
		"id":         apiToken.ID,
		"name":       apiToken.Name,
		"token":      rawToken,
		"created_at": apiToken.CreatedAt,
	})
}

// ProfileTokenDelete deletes an API token by ID if it belongs to the current user.
func (c *Controller) ProfileTokenDelete(g *gin.Context) {
	user := g.MustGet("user").(*model.User)
	id := g.Param("id")
	if id == "" {
		g.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
		return
	}
	var apiToken model.ApiToken
	if c.datastore.Where("id = ? AND user_id = ?", id, user.ID).First(&apiToken).RowsAffected < 1 {
		g.JSON(http.StatusNotFound, gin.H{"error": "token not found"})
		return
	}
	if err := c.datastore.Delete(&apiToken).Error; err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete"})
		return
	}
	g.AbortWithStatus(http.StatusNoContent)
}
