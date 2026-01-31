package utils

import "fmt"

// GenerateRoomID generates a consistent room ID for two users
// The ID is deterministic - same pair of users will always get same room ID
func GenerateRoomID(userA, userB int) string {
	if userA > userB {
		userA, userB = userB, userA
	}
	return fmt.Sprintf("room_%d_%d", userA, userB)
}
