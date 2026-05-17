package controller

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/storage"
)

const errFileEmpty = "file name cannot be empty"

var uploadVideoExt = regexp.MustCompile(`(?i)\.(mp4|mkv|webm)$`)
var uploadImageExt = regexp.MustCompile(`(?i)\.(png|jpg|jpeg)$`)
var uploadZipExt = regexp.MustCompile(`(?i)\.zip$`)

const (
	assignTargetVideo       = "video"
	assignTargetActor       = "actor"
	assignTargetChannel     = "channel"
	assignTargetCategorySub = "category_sub"
)

func classifyVideoBySize(sizeBytes, clipVideoMB, videoMovieMB int64) string {
	mb := sizeBytes / (1024 * 1024)
	if mb < clipVideoMB {
		return "c"
	}
	if mb < videoMovieMB {
		return "v"
	}
	return "m"
}

func scanTypeEnabled(enabled struct {
	C bool `json:"c"`
	V bool `json:"v"`
	M bool `json:"m"`
}, typ string) bool {
	switch typ {
	case "c":
		return enabled.C
	case "v":
		return enabled.V
	case "m":
		return enabled.M
	default:
		return false
	}
}

// uploadLibraryID returns the library ID to use for upload/triage: if formLibraryID is non-empty
// and a library with that ID exists, returns it; otherwise returns the default library ID.
func (c *Controller) uploadLibraryID(formLibraryID string) string {
	if formLibraryID == "" {
		return c.config.DefaultLibraryID
	}
	var lib model.Library
	if c.datastore.First(&lib, "id = ?", formLibraryID).RowsAffected < 1 {
		return c.config.DefaultLibraryID
	}
	return lib.ID
}
//
//	@Summary	Import file from triage as video
//	@Tags		upload
//	@Accept		json
//	@Param		body	body	object	true	"JSON with path, import_as"
//	@Success	201	{object}	map[string]interface{}
//	@Failure	400	{object}	map[string]interface{}
//	@Router		/upload/import [post]
func (c *Controller) UploadImport(g *gin.Context) {
	var body struct {
		Path      string `json:"path"`
		ImportAs  string `json:"import_as"`
		LibraryID string `json:"library_id"`
	}
	if err := g.ShouldBindJSON(&body); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	libID := c.uploadLibraryID(body.LibraryID)
	video := &model.Video{
		Name:          body.Path,
		Filename:      body.Path,
		Thumbnail:     false,
		ThumbnailMini: false,
		Type:          body.ImportAs,
		LibraryID:     &libID,
	}
	if err := c.datastore.Create(video).Error; err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusCreated, gin.H{"id": video.ID, "redirect": "/video/" + video.ID})
}

// UploadImportPhotoset godoc
//
//	@Summary	Import a zip archive from triage as a photoset
//	@Tags		upload
//	@Accept		json
//	@Param		body	body	object	true	"JSON with path, name, optional library_id and delete_from_triage"
//	@Success	202	{object}	map[string]interface{}
//	@Failure	400	{object}	map[string]interface{}
//	@Failure	404	{object}	map[string]interface{}
//	@Router		/upload/triage/import-photoset [post]
func (c *Controller) UploadImportPhotoset(g *gin.Context) {
	type importForm struct {
		Path             string `json:"path" binding:"required"`
		Name             string `json:"name" binding:"required"`
		LibraryID        string `json:"library_id"`
		DeleteFromTriage bool   `json:"delete_from_triage"`
	}
	form := importForm{}
	if err := g.ShouldBindJSON(&form); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file := strings.TrimPrefix(strings.TrimPrefix(form.Path, "/"), "\\")
	if file == "" {
		g.JSON(http.StatusBadRequest, gin.H{"error": errFileEmpty})
		return
	}
	if !uploadZipExt.MatchString(file) {
		g.JSON(http.StatusBadRequest, gin.H{"error": "file must be a .zip archive"})
		return
	}

	name := strings.TrimSpace(form.Name)
	if name == "" {
		g.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	libID := c.uploadLibraryID(form.LibraryID)
	store, err := c.storageResolver.Storage(libID)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	srcPath := filepath.Join("triage", file)
	exists, err := store.Exists(srcPath)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !exists {
		g.JSON(http.StatusNotFound, gin.H{"error": "triage file not found"})
		return
	}

	ps := &model.Photoset{
		Name:      name,
		LibraryID: &libID,
		Status:    model.PhotosetStatusCreating,
	}
	if err := c.datastore.Create(ps).Error; err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	archivePath := filepath.Join(ps.FolderRelativePath(), "_upload.zip")
	if err := store.MkdirAll(ps.FolderRelativePath()); err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := storage.CopyObject(store, store, srcPath, archivePath); err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if form.DeleteFromTriage {
		if err := store.Delete(srcPath); err != nil {
			g.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
	}

	if c.runner == nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": "task runner not available"})
		return
	}
	if err := c.runner.NewTask("photoset/unzip", map[string]string{
		"photosetID":  ps.ID,
		"archivePath": archivePath,
	}); err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusAccepted, gin.H{
		"photoset_id": ps.ID,
		"id":          ps.ID,
		"redirect":    "/adm/tasks",
	})
}

// UploadPreview godoc
//
//	@Summary	Preview file from triage folder
//	@Tags		upload
//	@Param		filepath	path	string	true	"URL-encoded file path"
//	@Success	200	file	bytes
//	@Failure	500	{object}	map[string]interface{}
//	@Router		/upload/preview/{filepath} [get]
func (c *Controller) UploadPreview(g *gin.Context) {
	filePathEncoded := g.Param("filepath")
	filePath, err := url.QueryUnescape(filePathEncoded)
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	libID := c.uploadLibraryID(g.Query("library_id"))
	store, err := c.storageResolver.Storage(libID)
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	path := filepath.Join("triage", filePath)
	c.serveFromStorage(g, store, path)
}

// UploadTriageFolder godoc
//
//	@Summary	List folders in triage path with file counts
//	@Tags		upload
//	@Accept		x-www-form-urlencoded
//	@Param		path	formData	string	true	"Path in triage"
//	@Success	200	{object}	map[string]interface{}
//	@Failure	500	{object}	map[string]interface{}
//	@Router		/upload/triage/folder [post]
func (c *Controller) UploadTriageFolder(g *gin.Context) {
	path := g.PostForm("path")
	libID := c.uploadLibraryID(g.PostForm("library_id"))
	store, err := c.storageResolver.Storage(libID)
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	prefix := filepath.Join("triage", path)
	entries, err := store.List(prefix)
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	items := make(map[string]int)
	for _, e := range entries {
		if !e.IsDir {
			continue
		}
		subPrefix := filepath.Join(prefix, e.Name)
		sub, err := store.List(subPrefix)
		if err != nil {
			continue
		}
		count := 0
		for _, s := range sub {
			if !s.IsDir {
				count++
			}
		}
		items[e.Name] = count
	}
	g.JSON(http.StatusOK, gin.H{"folders": items})
}

type FileInfo struct {
	Size             int64
	LastModification time.Time
}

// UploadTriageFile godoc
//
//	@Summary	List files in triage path
//	@Tags		upload
//	@Accept		x-www-form-urlencoded
//	@Param		path	formData	string	true	"Path in triage"
//	@Success	200	{object}	map[string]interface{}
//	@Failure	500	{object}	map[string]interface{}
//	@Router		/upload/triage/file [post]
func (c *Controller) UploadTriageFile(g *gin.Context) {
	path := g.PostForm("path")
	libID := c.uploadLibraryID(g.PostForm("library_id"))
	store, err := c.storageResolver.Storage(libID)
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	prefix := filepath.Join("triage", path)
	entries, err := store.List(prefix)
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	items := make(map[string]FileInfo)
	for _, e := range entries {
		if e.IsDir {
			continue
		}
		items[e.Name] = FileInfo{
			Size:             e.Size,
			LastModification: e.ModTime,
		}
	}
	g.JSON(http.StatusOK, gin.H{"files": items})
}

// UploadFile godoc
//
//	@Summary	Upload file to triage folder
//	@Tags		upload
//	@Accept		multipart/form-data
//	@Param		file	formData	file	true	"File to upload"
//	@Param		path	formData	string	true	"Destination path"
//	@Success	200
//	@Failure	500	{object}	map[string]interface{}
//	@Router		/upload/file [post]
func (c *Controller) UploadFile(g *gin.Context) {
	file, err := g.FormFile("file")
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	_path := g.PostForm("path")
	path := filepath.Join("triage", _path, file.Filename)
	libID := c.uploadLibraryID(g.PostForm("library_id"))
	store, err := c.storageResolver.Storage(libID)
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if err := store.MkdirAll(filepath.Dir(path)); err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	src, err := file.Open()
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer src.Close()
	dst, err := store.Create(path)
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer dst.Close()
	if _, err := io.Copy(dst, src); err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	g.JSON(200, gin.H{})
}

// UploadDeleteFile godoc
//
//	@Summary	Delete file from triage
//	@Tags		upload
//	@Param		File	formData	string	true	"File path in triage"
//	@Success	200
//	@Failure	400	{object}	map[string]interface{}
//	@Router		/upload/file [delete]
func (c *Controller) UploadDeleteFile(g *gin.Context) {
	type fileDeleteForm struct {
		File      string `json:"File"`
		LibraryID string `json:"library_id"`
	}
	form := fileDeleteForm{}
	if err := g.ShouldBindJSON(&form); err != nil {
		g.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if form.File == "" {
		g.JSON(400, gin.H{"error": errFileEmpty})
		return
	}
	libID := c.uploadLibraryID(form.LibraryID)
	store, err := c.storageResolver.Storage(libID)
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	path := filepath.Join("triage", form.File)
	if err := store.Delete(path); err != nil {
		g.JSON(422, gin.H{"error": err.Error()})
		return
	}
	g.JSON(200, gin.H{})
}

// UploadMassDelete godoc
//
//	@Summary	Delete multiple files from triage
//	@Tags		upload
//	@Accept		json
//	@Param		body	body	object	true	"JSON with files array"
//	@Success	200
//	@Failure	400	{object}	map[string]interface{}
//	@Router		/upload/triage/mass-action [delete]
func (c *Controller) UploadMassDelete(g *gin.Context) {
	// get file list from request
	type fileDeleteForm struct {
		Files      []string `json:"files" binding:"required"`
		LibraryID  string   `json:"library_id"`
	}

	form := fileDeleteForm{}
	err := g.ShouldBindJSON(&form)
	if err != nil {
		g.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	// ensure not empty
	files := form.Files
	if len(files) == 0 {
		g.JSON(400, gin.H{
			"error": "mass deletion requested without any files",
		})
		return
	}
	libID := c.uploadLibraryID(form.LibraryID)
	for _, file := range files {
		c.logger.Debug().Str("file", file).Send()
		if file == "" {
			g.JSON(400, gin.H{"error": errFileEmpty})
			return
		}
		store, err := c.storageResolver.Storage(libID)
		if err != nil {
			g.JSON(500, gin.H{"error": err.Error()})
			return
		}
		path := filepath.Join("triage", file)
		if err := store.Delete(path); err != nil {
			g.JSON(422, gin.H{"error": err.Error()})
			return
		}
	}
	g.JSON(200, gin.H{})
}

// UploadMassImport godoc
//
//	@Summary	Mass import files from triage as videos
//	@Tags		upload
//	@Accept		json
//	@Param		body	body	object	true	"JSON with files, type, actors, categories, channel"
//	@Success	200
//	@Failure	400	{object}	map[string]interface{}
//	@Router		/upload/triage/mass-action [post]
func (c *Controller) UploadMassImport(g *gin.Context) {
	// get file list from request
	type fileImportForm struct {
		Files              []string `json:"files" binding:"required"`
		Actors             []string `json:"actors"`
		Categories         []string `json:"categories"`
		TypeEnum           string   `json:"type" binding:"required"`
		Channel            string   `json:"channel"`
		LibraryID          string   `json:"library_id"`
		SkipReorganization *bool    `json:"skip_reorganization"`
	}

	form := fileImportForm{}
	err := g.ShouldBindJSON(&form)
	if err != nil {
		g.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	// ensure type is valid
	if form.TypeEnum != "c" && form.TypeEnum != "v" && form.TypeEnum != "m" {
		g.JSON(400, gin.H{
			"error": "type of video is invalid",
		})
	}

	// ensure not empty
	files := form.Files
	if len(files) == 0 {
		g.JSON(400, gin.H{
			"error": "mass import requested without any files",
		})
		return
	}

	// pre-check: ensure actors exists
	var actors []*model.Actor
	for _, actorID := range form.Actors {
		actor := &model.Actor{
			ID: actorID,
		}
		result := c.datastore.First(actor)

		// check result
		if result.RowsAffected < 1 {
			g.JSON(400, gin.H{
				"error": fmt.Sprintf("actor id %s does not exist", actorID),
			})
			return
		}

		actors = append(actors, actor)
	}

	// pre-check: ensure categories exists
	var categories []*model.CategorySub
	for _, subCategoryID := range form.Categories {
		subCategory := &model.CategorySub{
			ID: subCategoryID,
		}
		result := c.datastore.First(subCategory)

		// check result
		if result.RowsAffected < 1 {
			g.JSON(400, gin.H{
				"error": fmt.Sprintf("category id %s does not exist", subCategoryID),
			})
			return
		}

		categories = append(categories, subCategory)
	}

	// pre-check: ensure channel exists
	var channel *model.Channel
	channel = nil
	if form.Channel != "" {
		channel = &model.Channel{
			ID: form.Channel,
		}
		result := c.datastore.First(channel)

		// check result
		if result.RowsAffected < 1 {
			g.JSON(400, gin.H{
				"error": fmt.Sprintf("channel id %s does not exist", form.Channel),
			})
			return
		}
	}

	// pre-check: ensure files are not empty
	for _, file := range files {
		c.logger.Debug().Str("file", file).Send()
		if file == "" {
			g.JSON(400, gin.H{
				"error": errFileEmpty,
			})
			return
		}
	}

	// prepare transaction for the whole import
	tx := c.datastore.Begin()
	var videos []*model.Video
	libID := c.uploadLibraryID(form.LibraryID)

	// now perform file import
	for _, file := range files {
		video := &model.Video{
			Name:      file,
			Filename:  file,
			Type:      form.TypeEnum,
			Imported:  false,
			Thumbnail: false,
			LibraryID: &libID,
		}

		if channel != nil {
			video.Channel = channel
		}

		// save object in db
		err = tx.Create(&video).Error
		if err != nil {
			tx.Rollback()
			g.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		err = tx.Model(video).Association("Actors").Append(actors)
		if err != nil {
			tx.Rollback()
			g.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		err = tx.Model(video).Association("Categories").Append(categories)
		if err != nil {
			tx.Rollback()
			g.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		videos = append(videos, video)
	}

	// validate transaction
	tx.Commit()

	// now create task for the import
	for _, video := range videos {
		params := map[string]string{
			"videoID":         video.ID,
			"thumbnailTiming": "0",
		}
		if form.SkipReorganization != nil {
			if *form.SkipReorganization {
				params["skipReorganization"] = "true"
			} else {
				params["skipReorganization"] = "false"
			}
		}
		err = c.runner.NewTask("video/create", params)
		if err != nil {
			g.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	g.JSON(200, gin.H{})
}

// UploadTriageScan godoc
//
//	@Summary	Scan triage folder and import videos by size-based type
//	@Tags		upload
//	@Accept		json
//	@Param		body	body	object	true	"JSON with path, recursive, thresholds, enabled types, metadata"
//	@Success	200	{object}	map[string]interface{}
//	@Failure	400	{object}	map[string]interface{}
//	@Router		/upload/triage/scan [post]
func (c *Controller) UploadTriageScan(g *gin.Context) {
	type scanForm struct {
		Path               string   `json:"path"`
		Recursive          bool     `json:"recursive"`
		LibraryID          string   `json:"library_id"`
		SkipReorganization *bool    `json:"skip_reorganization"`
		Channel            string   `json:"channel"`
		Actors             []string `json:"actors"`
		Categories         []string `json:"categories"`
		Enabled            struct {
			C bool `json:"c"`
			V bool `json:"v"`
			M bool `json:"m"`
		} `json:"enabled"`
		Thresholds struct {
			ClipVideoMB  int64 `json:"clip_video_mb"`
			VideoMovieMB int64 `json:"video_movie_mb"`
		} `json:"thresholds"`
	}

	form := scanForm{}
	if err := g.ShouldBindJSON(&form); err != nil {
		g.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if !form.Enabled.C && !form.Enabled.V && !form.Enabled.M {
		g.JSON(400, gin.H{"error": "at least one video type must be enabled"})
		return
	}
	if form.Thresholds.ClipVideoMB <= 0 || form.Thresholds.VideoMovieMB <= 0 {
		g.JSON(400, gin.H{"error": "size thresholds must be positive"})
		return
	}
	if form.Thresholds.ClipVideoMB > form.Thresholds.VideoMovieMB {
		g.JSON(400, gin.H{"error": "clip threshold must not exceed video threshold"})
		return
	}

	var actors []*model.Actor
	for _, actorID := range form.Actors {
		actor := &model.Actor{ID: actorID}
		if c.datastore.First(actor).RowsAffected < 1 {
			g.JSON(400, gin.H{"error": fmt.Sprintf("actor id %s does not exist", actorID)})
			return
		}
		actors = append(actors, actor)
	}

	var categories []*model.CategorySub
	for _, subCategoryID := range form.Categories {
		subCategory := &model.CategorySub{ID: subCategoryID}
		if c.datastore.First(subCategory).RowsAffected < 1 {
			g.JSON(400, gin.H{"error": fmt.Sprintf("category id %s does not exist", subCategoryID)})
			return
		}
		categories = append(categories, subCategory)
	}

	var channel *model.Channel
	if form.Channel != "" {
		channel = &model.Channel{ID: form.Channel}
		if c.datastore.First(channel).RowsAffected < 1 {
			g.JSON(400, gin.H{"error": fmt.Sprintf("channel id %s does not exist", form.Channel)})
			return
		}
	}

	libID := c.uploadLibraryID(form.LibraryID)
	store, err := c.storageResolver.Storage(libID)
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}

	scanPath := strings.TrimPrefix(strings.TrimPrefix(form.Path, "/"), "\\")
	prefix := filepath.Join("triage", scanPath)

	type scanCandidate struct {
		rel  string
		size int64
	}
	var candidates []scanCandidate
	var skippedNonVideo int

	err = storage.WalkFiles(store, prefix, form.Recursive, func(p string, e storage.Entry) error {
		if !uploadVideoExt.MatchString(e.Name) {
			skippedNonVideo++
			return nil
		}
		rel, relErr := filepath.Rel("triage", p)
		if relErr != nil {
			return relErr
		}
		rel = filepath.ToSlash(rel)
		candidates = append(candidates, scanCandidate{rel: rel, size: e.Size})
		return nil
	})
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var skippedExisting int
	var skippedDisabled int
	var toImport []struct {
		rel  string
		typ  string
		size int64
	}

	existingFilenames := make(map[string]struct{})
	var existingRows []struct {
		Filename string
	}
	if err := c.datastore.Model(&model.Video{}).Select("filename").Where("library_id = ?", libID).Find(&existingRows).Error; err != nil {
		g.JSON(500, gin.H{"error": "unable to check existing videos"})
		return
	}
	for _, row := range existingRows {
		existingFilenames[row.Filename] = struct{}{}
	}

	for _, cand := range candidates {
		vtyp := classifyVideoBySize(cand.size, form.Thresholds.ClipVideoMB, form.Thresholds.VideoMovieMB)
		if !scanTypeEnabled(form.Enabled, vtyp) {
			skippedDisabled++
			continue
		}
		if _, ok := existingFilenames[cand.rel]; ok {
			skippedExisting++
			continue
		}
		toImport = append(toImport, struct {
			rel  string
			typ  string
			size int64
		}{cand.rel, vtyp, cand.size})
	}

	if len(toImport) == 0 {
		g.JSON(200, gin.H{
			"imported":          0,
			"skipped_existing":  skippedExisting,
			"skipped_disabled":  skippedDisabled,
			"skipped_non_video": skippedNonVideo,
		})
		return
	}

	tx := c.datastore.Begin()
	var videos []*model.Video
	for _, item := range toImport {
		video := &model.Video{
			Name:      item.rel,
			Filename:  item.rel,
			Type:      item.typ,
			Imported:  false,
			Thumbnail: false,
			LibraryID: &libID,
		}
		if channel != nil {
			video.Channel = channel
		}
		if err := tx.Create(video).Error; err != nil {
			tx.Rollback()
			g.JSON(500, gin.H{"error": err.Error()})
			return
		}
		if err := tx.Model(video).Association("Actors").Append(actors); err != nil {
			tx.Rollback()
			g.JSON(500, gin.H{"error": err.Error()})
			return
		}
		if err := tx.Model(video).Association("Categories").Append(categories); err != nil {
			tx.Rollback()
			g.JSON(500, gin.H{"error": err.Error()})
			return
		}
		videos = append(videos, video)
	}
	tx.Commit()

	for _, video := range videos {
		params := map[string]string{
			"videoID":         video.ID,
			"thumbnailTiming": "0",
		}
		if form.SkipReorganization != nil {
			if *form.SkipReorganization {
				params["skipReorganization"] = "true"
			} else {
				params["skipReorganization"] = "false"
			}
		}
		if err := c.runner.NewTask("video/create", params); err != nil {
			g.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	g.JSON(200, gin.H{
		"imported":          len(videos),
		"skipped_existing":  skippedExisting,
		"skipped_disabled":  skippedDisabled,
		"skipped_non_video": skippedNonVideo,
	})
}

// UploadFolderCreate godoc
//
//	@Summary	Create folder in triage
//	@Tags		upload
//	@Accept		x-www-form-urlencoded
//	@Param		name	formData	string	true	"Folder name"
//	@Success	200
//	@Failure	409	{object}	map[string]interface{}
//	@Router		/upload/folder [post]
func (c *Controller) UploadFolderCreate(g *gin.Context) {
	name := g.PostForm("name")
	libID := c.uploadLibraryID(g.PostForm("library_id"))
	store, err := c.storageResolver.Storage(libID)
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	path := filepath.Join("triage", name)
	exists, err := store.Exists(path)
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if exists {
		g.JSON(409, gin.H{"error": "Folder already exists"})
		return
	}
	if err := store.MkdirAll(path); err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	g.JSON(200, gin.H{})
}

// UploadAssignImage godoc
//
//	@Summary	Copy a triage image onto an entity thumbnail
//	@Tags		upload
//	@Accept		json
//	@Param		body	body	object	true	"JSON with file, target_type, target_id, optional library_id and delete_from_triage"
//	@Success	200	{object}	map[string]interface{}
//	@Failure	400	{object}	map[string]interface{}
//	@Failure	404	{object}	map[string]interface{}
//	@Router		/upload/triage/assign-image [post]
func (c *Controller) UploadAssignImage(g *gin.Context) {
	type assignForm struct {
		File             string `json:"file" binding:"required"`
		LibraryID        string `json:"library_id"`
		TargetType       string `json:"target_type" binding:"required"`
		TargetID         string `json:"target_id" binding:"required"`
		DeleteFromTriage bool   `json:"delete_from_triage"`
	}

	form := assignForm{}
	if err := g.ShouldBindJSON(&form); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file := strings.TrimPrefix(strings.TrimPrefix(form.File, "/"), "\\")
	if file == "" {
		g.JSON(http.StatusBadRequest, gin.H{"error": errFileEmpty})
		return
	}
	if !uploadImageExt.MatchString(file) {
		g.JSON(http.StatusBadRequest, gin.H{"error": "file must be a png, jpg, or jpeg image"})
		return
	}

	metaStore, err := c.metadataStoreForWrite()
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	libID := c.uploadLibraryID(form.LibraryID)
	libStore, err := c.storageResolver.Storage(libID)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	srcPath := filepath.Join("triage", file)
	exists, err := libStore.Exists(srcPath)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !exists {
		g.JSON(http.StatusNotFound, gin.H{"error": "triage file not found"})
		return
	}

	var redirect string
	var videoID string

	switch form.TargetType {
	case assignTargetVideo:
		video := &model.Video{ID: form.TargetID}
		if c.datastore.First(video).RowsAffected < 1 {
			g.JSON(http.StatusNotFound, gin.H{"error": "video not found"})
			return
		}
		dstPath := video.ThumbnailRelativePath()
		redirect = video.URLAdmEdit()
		videoID = video.ID
		if err := storage.CopyObject(libStore, metaStore, srcPath, dstPath); err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		video.Thumbnail = true
		video.Migrated = true
		if err := c.datastore.Save(video).Error; err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

	case assignTargetActor:
		actor := &model.Actor{ID: form.TargetID}
		if c.datastore.First(actor).RowsAffected < 1 {
			g.JSON(http.StatusNotFound, gin.H{"error": "actor not found"})
			return
		}
		dstPath := filepath.Join("actors", actor.ID, "thumb.jpg")
		redirect = actor.URLAdmEdit()
		if err := storage.CopyObject(libStore, metaStore, srcPath, dstPath); err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		actor.Thumbnail = true
		actor.Migrated = true
		if err := c.datastore.Save(actor).Error; err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

	case assignTargetChannel:
		channel := &model.Channel{ID: form.TargetID}
		if c.datastore.First(channel).RowsAffected < 1 {
			g.JSON(http.StatusNotFound, gin.H{"error": "channel not found"})
			return
		}
		dstPath := filepath.Join("channels", channel.ID, "thumb.jpg")
		redirect = channel.URLView()
		if err := storage.CopyObject(libStore, metaStore, srcPath, dstPath); err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		channel.Thumbnail = true
		channel.Migrated = true
		if err := c.datastore.Save(channel).Error; err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

	case assignTargetCategorySub:
		sub := &model.CategorySub{ID: form.TargetID}
		if c.datastore.First(sub).RowsAffected < 1 {
			g.JSON(http.StatusNotFound, gin.H{"error": "category sub not found"})
			return
		}
		dstPath := filepath.Join("categories", sub.ID+".jpg")
		redirect = sub.URLThumb()
		if err := storage.CopyObject(libStore, metaStore, srcPath, dstPath); err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		sub.Thumbnail = true
		sub.Migrated = true
		if err := c.datastore.Save(sub).Error; err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

	default:
		g.JSON(http.StatusBadRequest, gin.H{"error": "target_type must be video, actor, channel, or category_sub"})
		return
	}

	if videoID != "" {
		if err := c.runner.NewTask("video/mini-thumb", map[string]string{"videoID": videoID}); err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if form.DeleteFromTriage {
		if err := libStore.Delete(srcPath); err != nil {
			g.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
	}

	g.JSON(http.StatusOK, gin.H{"redirect": redirect})
}
