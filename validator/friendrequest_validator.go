package validator

import (
    "fmt"
)

func ValidateFriendRequest(fromUserID, toUserID string) error {
    validations := []func() error{
        func() error { return validateRequired(fromUserID, "sender user ID is required") },
        func() error { return validateRequired(toUserID, "recipient user ID is required") },
        func() error { return validateDifferentUsers(fromUserID, toUserID) },
    }

    for _, validate := range validations {
        if err := validate(); err != nil {
            return err
        }
    }
    return nil
}

func validateDifferentUsers(fromUserID, toUserID string) error {
    if fromUserID == toUserID {
        return fmt.Errorf("cannot send friend request to self")
    }
    return nil
}