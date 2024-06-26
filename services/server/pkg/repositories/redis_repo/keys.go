package redis_repo

import (
	"spaces-p/pkg/models"
	"spaces-p/pkg/uuid"
)

// ---- USER ----

var userFields = struct {
	userFirstNameField string
	userLastNameField  string
	userUsernameField  string
	userAvatarUrlField string
}{userFirstNameField: "first_name", userLastNameField: "last_name", userUsernameField: "username", userAvatarUrlField: "avatar_url"}

// getUserKey returns a redis key: users:[user_uid]
//
// The keys holds a HASH value with the following fields: "is_signed_up", "first_name", "last_name" "username", "avatar_url"
func getUserKey(userId models.UserUid) string {
	return "users:" + string(userId)
}

// getUserSpacesKey returns a redis key: users:[user_uid]:spaces
//
// The keys hold SORTED SET values with the space ids as MEMBERS and joining time as SCORES
func getUserSpacesKey(userId models.UserUid) string {
	return getUserKey(userId) + ":spaces"
}

// ---- SPACE COORDINATES ----

// space_coords
func getSpaceCoordinatesKey() string {
	return "space_coords"
}

// ---- SPACE ----

var spaceFields = struct {
	nameField               string
	themeColorHexaCodeField string
	radiusField             string
	locationField           string
	createdAtField          string
	adminIdField            string
}{
	nameField:               "name",
	themeColorHexaCodeField: "color",
	radiusField:             "radius",
	locationField:           "location",
	createdAtField:          "created_at",
	adminIdField:            "admin",
}

// spaces:[spaceid] hash of space data
func getSpaceKey(spaceId uuid.Uuid) string {
	return "spaces:" + spaceId.String()
}

// must be subset of users:[user_uid]
//
// spaces:[spaceid]:subscribers
//
// The keys hold SORTED SET values with the user ids as MEMBERS and joining time as SCORES
func getSpaceSubscribersKey(spaceId uuid.Uuid) string {
	return getSpaceKey(spaceId) + ":subscribers"
}

// spaces:[spaceid]:subscribers[userid]:sessions
//
// The keys holds a SORTED SET value with the session ids as MEMBERS and starting session time as SCORES
func getSpaceActiveSubscriberSessionsKey(spaceId uuid.Uuid, userId models.UserUid) string {
	return getSpaceSubscribersKey(spaceId) + ":" + string(userId) + ":sessions"
}

// must be subset of spaces:[spaceid]:subscribers
//
// spaces:[spaceid]:active_subscribers
//
// The keys hold SORTED SET values with the user ids as MEMBERS and joining times as SCORES
func getSpaceActiveSubscribersKey(spaceId uuid.Uuid) string {
	return getSpaceKey(spaceId) + ":active_subscribers"
}

// spaces:[spaceid]:toplevel_threads_by_time
func getSpaceToplevelThreadsByTimeKey(spaceId uuid.Uuid) string {
	return getSpaceKey(spaceId) + ":toplevel_threads_by_time"
}

// spaces:[spaceid]:toplevel_threads_by_popularity
func getSpaceToplevelThreadsByPopularityKey(spaceId uuid.Uuid) string {
	return getSpaceKey(spaceId) + ":toplevel_threads_by_popularity"
}

// ---- THREAD ----

var threadFields = struct {
	likesField           string
	messagesCountField   string
	spaceIdField         string
	createdAtField       string
	parentMessageIdField string
	firstMessageIdField  string
}{
	likesField:           "likes",
	messagesCountField:   "messages_count",
	spaceIdField:         "space_id",
	parentMessageIdField: "parent_message_id", // only for sublevel threads
	firstMessageIdField:  "first_message_id",  // only for toplevel threads
	createdAtField:       "created_at",
}

// threads:[threadid]
func getThreadKey(threadId uuid.Uuid) string {
	return "threads:" + threadId.String()
}

// threads:[threadid]:messages_by_time
func getThreadMessagesByTimeKey(threadId uuid.Uuid) string {
	return getThreadKey(threadId) + ":messages_by_time"
}

// threads:[threadid]:messages_by_popularity
func getThreadMessagesByPopularityKey(threadId uuid.Uuid) string {
	return getThreadKey(threadId) + ":messages_by_popularity"
}

// ---- MESSAGE ----

var messageFields = struct {
	contentField       string
	senderIdField      string
	typeField          string
	createdAtField     string
	childThreadIdField string
	threadIdField      string
	likesField         string
}{
	contentField:       "content",
	senderIdField:      "sender_id",
	typeField:          "type",
	createdAtField:     "created_at",
	childThreadIdField: "child_thread_id",
	threadIdField:      "thread_id",
	likesField:         "likes",
}

// messages:[messageid]
func getMessageKey(messageId uuid.Uuid) string {
	return "messages:" + messageId.String()
}

// ---- ADDRESS ----

// getAddressKey returns a redis key: addresses:[geohash]
//
// The keys holds a STRINGIFIED JSON value which represent an address
func getAddressKey(geohash string) string {
	return "addresses:" + geohash
}
