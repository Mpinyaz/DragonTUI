package pages

import (
	"sync"
)

var (
	menuModelOnce sync.Once
	menuModel     *MenuModel

	aboutModelOnce sync.Once
	aboutModel     *AboutModel

	contactModelOnce sync.Once
	contactModel     *ContactModel
)

func GetMenuModel(width, height int) *MenuModel {
	menuModelOnce.Do(func() {
		menuModel = NewMenuModel(width, height)
	})
	return menuModel
}

func GetAboutModel(width, height int) *AboutModel {
	aboutModelOnce.Do(func() {
		aboutModel = NewAboutModel(width, height)
	})
	return aboutModel
}

func GetContactModel(width, height int) *ContactModel {
	contactModelOnce.Do(func() {
		contactModel = NewContactModel(width, height)
	})
	return contactModel
}
