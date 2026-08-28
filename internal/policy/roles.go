package policy

import (
	"fmt"
	"strings"
)

type Role string

const (
	RoleClerk    Role = "clerk"
	RoleReviewer Role = "reviewer"
	RoleManager  Role = "manager"
)

type Actor struct {
	ID   string
	Role Role
}

func (a Actor) Valid() bool {
	return strings.TrimSpace(a.ID) != "" && (a.Role == RoleClerk || a.Role == RoleReviewer || a.Role == RoleManager)
}

func Can(actor Actor, action string) bool {
	if !actor.Valid() {
		return false
	}
	switch action {
	case "register", "update", "submit":
		return actor.Role == RoleClerk || actor.Role == RoleManager
	case "approve", "reject":
		return actor.Role == RoleReviewer || actor.Role == RoleManager
	case "publish", "archive":
		return actor.Role == RoleManager
	default:
		return false
	}
}

func Require(actor Actor, action string) error {
	if Can(actor, action) {
		return nil
	}
	return fmt.Errorf("actor %s cannot %s", actor.ID, action)
}

func ParseRole(value string) Role {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "clerk":
		return RoleClerk
	case "reviewer":
		return RoleReviewer
	case "manager":
		return RoleManager
	default:
		return Role("")
	}
}
