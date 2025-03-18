package controller

const ACTOR_FILEPATH = "/actors"
const ACTOR_PROFILE_PICTURE_MISSING = "/static/images/actor-missing-profile-picture.png"

var fileTypeToPath = map[string]string{
	"clip":  "/clips",
	"movie": "/movies",
	"video": "/videos",
}

const VIDEO_THUMB_NOT_GENERATED = "/static/images/video-thumb-not-generated.png"

const TRIAGE_FILEPATH = "/triage"

const ZT_VERSION = "0.1.30"
