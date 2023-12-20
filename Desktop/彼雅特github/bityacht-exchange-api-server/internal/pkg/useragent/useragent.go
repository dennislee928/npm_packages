package useragent

import (
	"github.com/mileusna/useragent"
)

type UserAgent struct {
	useragent.UserAgent
}

func Parse(userAgent string) UserAgent {
	return UserAgent{
		UserAgent: useragent.Parse(userAgent),
	}
}

func (ua UserAgent) GetBrowserString() string {
	browser := "Unknown"

	if ua.Name != "" {
		browser = ua.Name
	}
	if ua.Version != "" {
		browser += " " + ua.Version
	}

	return browser
}

func (ua UserAgent) GetDeviceString() string {
	var device string

	switch {
	case ua.Device != "":
		device = ua.Device
	case ua.Mobile:
		device = "Mobile"
	case ua.Tablet:
		device = "Tablet"
	case ua.Desktop:
		device = "Desktop"
	case ua.Bot:
		device = "Bot"
	default:
		device = "Unknown"
	}

	if ua.OS != "" {
		if ua.OSVersion == "" {
			device += " (" + ua.OS + ")"
		} else {
			device += " (" + ua.OS + " " + ua.OSVersion + ")"
		}
	}

	return device
}
