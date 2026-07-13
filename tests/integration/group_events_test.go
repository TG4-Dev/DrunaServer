package integration

import (
	"druna_server/pkg/model"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type groupFixture struct {
	ownerName    string
	ownerAccess  string
	memberName   string
	memberAccess string
	groupID      int
}

// setupGroupWithMember creates an owner and a member who are friends and both
// belong to a freshly created group.
func setupGroupWithMember(t *testing.T, router http.Handler) groupFixture {
	t.Helper()
	ownerName := fmt.Sprintf("ge_owner_%d", time.Now().UnixNano())
	memberName := fmt.Sprintf("ge_member_%d", time.Now().UnixNano())

	ownerAccess, _ := signUpAndSignIn(t, router, ownerName, "secret12345")
	memberAccess, _ := signUpAndSignIn(t, router, memberName, "secret12345")

	doJSON(t, router, http.MethodPost, "/api/v1/friends/request", map[string]string{"username": memberName}, ownerAccess)
	doJSON(t, router, http.MethodPost, "/api/v1/friends/accept", map[string]string{"username": ownerName}, memberAccess)

	createRec := doJSON(t, router, http.MethodPost, "/api/v1/groups/create", map[string]string{"name": "Events Group"}, ownerAccess)
	if createRec.Code != http.StatusOK {
		t.Fatalf("create group failed: %d %s", createRec.Code, createRec.Body.String())
	}
	groupData := decodeData[map[string]interface{}](t, createRec)
	gid, ok := groupData["groupId"].(float64)
	if !ok {
		t.Fatalf("expected groupId in response: %+v", groupData)
	}
	groupID := int(gid)

	addRec := doJSON(t, router, http.MethodPost, fmt.Sprintf("/api/v1/groups/%d/members", groupID), map[string]string{"username": memberName}, ownerAccess)
	if addRec.Code != http.StatusOK {
		t.Fatalf("add member failed: %d %s", addRec.Code, addRec.Body.String())
	}

	return groupFixture{
		ownerName:    ownerName,
		ownerAccess:  ownerAccess,
		memberName:   memberName,
		memberAccess: memberAccess,
		groupID:      groupID,
	}
}

func createGroupEvent(t *testing.T, router http.Handler, groupID int, access, title, start, end string) *httptest.ResponseRecorder {
	t.Helper()
	return doJSON(t, router, http.MethodPost, fmt.Sprintf("/api/v1/groups/%d/events", groupID), map[string]string{
		"title":     title,
		"startTime": start,
		"endTime":   end,
		"type":      "meeting",
	}, access)
}

func TestGroupEventMemberCreateAndList(t *testing.T) {
	router, cleanup := setupIntegration(t)
	defer cleanup()

	f := setupGroupWithMember(t, router)

	createRec := createGroupEvent(t, router, f.groupID, f.memberAccess, "Standup", "2026-06-17T10:00:00Z", "2026-06-17T11:00:00Z")
	if createRec.Code != http.StatusOK {
		t.Fatalf("member create group event failed: %d %s", createRec.Code, createRec.Body.String())
	}

	listRec := doJSON(t, router, http.MethodGet, fmt.Sprintf("/api/v1/groups/%d/events", f.groupID), nil, f.ownerAccess)
	if listRec.Code != http.StatusOK {
		t.Fatalf("list group events failed: %d %s", listRec.Code, listRec.Body.String())
	}
	result := decodeData[model.EventListResponse](t, listRec)
	if result.Total != 1 || len(result.Events) != 1 {
		t.Fatalf("expected 1 group event, got total=%d len=%d", result.Total, len(result.Events))
	}
	ev := result.Events[0]
	if ev.GroupID == nil || *ev.GroupID != f.groupID {
		t.Fatalf("expected groupID %d in event, got %+v", f.groupID, ev.GroupID)
	}

	// personal event list of the creator must NOT include the group event
	personalRec := doJSON(t, router, http.MethodGet, "/api/v1/events/", nil, f.memberAccess)
	if personalRec.Code != http.StatusOK {
		t.Fatalf("personal list failed: %d", personalRec.Code)
	}
	personal := decodeData[model.EventListResponse](t, personalRec)
	if personal.Total != 0 {
		t.Fatalf("expected personal list to exclude group event, got total=%d", personal.Total)
	}
}

func TestGroupEventNonMemberForbidden(t *testing.T) {
	router, cleanup := setupIntegration(t)
	defer cleanup()

	f := setupGroupWithMember(t, router)

	strangerName := fmt.Sprintf("ge_stranger_%d", time.Now().UnixNano())
	strangerAccess, _ := signUpAndSignIn(t, router, strangerName, "secret12345")

	createRec := createGroupEvent(t, router, f.groupID, strangerAccess, "Sneaky", "2026-06-17T10:00:00Z", "2026-06-17T11:00:00Z")
	if createRec.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for non-member create, got %d %s", createRec.Code, createRec.Body.String())
	}

	listRec := doJSON(t, router, http.MethodGet, fmt.Sprintf("/api/v1/groups/%d/events", f.groupID), nil, strangerAccess)
	if listRec.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for non-member list, got %d %s", listRec.Code, listRec.Body.String())
	}
}

func TestGroupEventOverlapRejected(t *testing.T) {
	router, cleanup := setupIntegration(t)
	defer cleanup()

	f := setupGroupWithMember(t, router)

	firstRec := createGroupEvent(t, router, f.groupID, f.ownerAccess, "Block", "2026-06-17T10:00:00Z", "2026-06-17T11:00:00Z")
	if firstRec.Code != http.StatusOK {
		t.Fatalf("first group event failed: %d %s", firstRec.Code, firstRec.Body.String())
	}

	overlapRec := createGroupEvent(t, router, f.groupID, f.memberAccess, "Overlap", "2026-06-17T10:30:00Z", "2026-06-17T11:30:00Z")
	if overlapRec.Code == http.StatusOK {
		t.Fatal("expected overlapping group event to be rejected")
	}
}

func TestGroupEventReducesFreeTime(t *testing.T) {
	router, cleanup := setupIntegration(t)
	defer cleanup()

	f := setupGroupWithMember(t, router)

	createRec := createGroupEvent(t, router, f.groupID, f.ownerAccess, "Meeting", "2026-06-17T14:00:00Z", "2026-06-17T16:00:00Z")
	if createRec.Code != http.StatusOK {
		t.Fatalf("create group event failed: %d %s", createRec.Code, createRec.Body.String())
	}

	freeRec := doJSON(t, router, http.MethodPost, fmt.Sprintf("/api/v1/groups/%d/free-time", f.groupID), map[string]string{
		"date": "2026-06-17",
	}, f.ownerAccess)
	if freeRec.Code != http.StatusOK {
		t.Fatalf("group free-time failed: %d %s", freeRec.Code, freeRec.Body.String())
	}
	slots := decodeData[map[string][]model.TimeSlot](t, freeRec)
	busyStart, _ := time.Parse(time.RFC3339, "2026-06-17T14:00:00Z")
	midpoint := busyStart.Add(time.Hour) // 15:00, inside the group event
	for _, slot := range slots["freeSlots"] {
		if !slot.Start.After(midpoint) && slot.End.After(midpoint) {
			t.Fatalf("free slot overlaps the group event window: %+v", slot)
		}
	}
}

func TestGroupEventUpdateDeletePermissions(t *testing.T) {
	router, cleanup := setupIntegration(t)
	defer cleanup()

	f := setupGroupWithMember(t, router)

	// member creates an event
	createRec := createGroupEvent(t, router, f.groupID, f.memberAccess, "MemberEvent", "2026-06-17T10:00:00Z", "2026-06-17T11:00:00Z")
	if createRec.Code != http.StatusOK {
		t.Fatalf("member create failed: %d %s", createRec.Code, createRec.Body.String())
	}
	created := decodeData[map[string]interface{}](t, createRec)
	eventID := int(created["eventId"].(float64))

	// creator can update
	updateRec := doJSON(t, router, http.MethodPatch, fmt.Sprintf("/api/v1/groups/%d/events/%d", f.groupID, eventID), map[string]string{
		"title":     "MemberEvent v2",
		"startTime": "2026-06-17T10:00:00Z",
		"endTime":   "2026-06-17T11:30:00Z",
		"type":      "meeting",
	}, f.memberAccess)
	if updateRec.Code != http.StatusOK {
		t.Fatalf("creator update failed: %d %s", updateRec.Code, updateRec.Body.String())
	}

	// third member (not creator, not owner) is forbidden to delete
	otherName := fmt.Sprintf("ge_other_%d", time.Now().UnixNano())
	otherAccess, _ := signUpAndSignIn(t, router, otherName, "secret12345")
	doJSON(t, router, http.MethodPost, "/api/v1/friends/request", map[string]string{"username": otherName}, f.ownerAccess)
	doJSON(t, router, http.MethodPost, "/api/v1/friends/accept", map[string]string{"username": f.ownerName}, otherAccess)
	addOtherRec := doJSON(t, router, http.MethodPost, fmt.Sprintf("/api/v1/groups/%d/members", f.groupID), map[string]string{"username": otherName}, f.ownerAccess)
	if addOtherRec.Code != http.StatusOK {
		t.Fatalf("add third member failed: %d %s", addOtherRec.Code, addOtherRec.Body.String())
	}

	forbiddenRec := doJSON(t, router, http.MethodDelete, fmt.Sprintf("/api/v1/groups/%d/events/%d", f.groupID, eventID), nil, otherAccess)
	if forbiddenRec.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for non-creator non-owner delete, got %d %s", forbiddenRec.Code, forbiddenRec.Body.String())
	}

	// owner can delete someone else's event
	ownerDeleteRec := doJSON(t, router, http.MethodDelete, fmt.Sprintf("/api/v1/groups/%d/events/%d", f.groupID, eventID), nil, f.ownerAccess)
	if ownerDeleteRec.Code != http.StatusOK {
		t.Fatalf("owner delete failed: %d %s", ownerDeleteRec.Code, ownerDeleteRec.Body.String())
	}
}
