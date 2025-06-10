package resources

import (
	"github.com/ouijan/ingenuity/pkg/client/input"
	"github.com/ouijan/ingenuity/pkg/core/utils"
)

type UserInputStore struct {
	Values  utils.InputActionValues
	Manager input.InputManager
}
