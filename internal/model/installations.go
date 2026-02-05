package model

// InstallationPlatform represents the push notification platform
type InstallationPlatform string

const (
	PlatformFCM  InstallationPlatform = "FCM"
	PlatformAPNS InstallationPlatform = "APNS"
)

// DeviceInstallationRequest represents the request to register/update a device installation
type DeviceInstallationRequest struct {
	Platform    InstallationPlatform `json:"platform"`
	PushChannel string               `json:"pushChannel"`
	Language    string               `json:"language"`
}
