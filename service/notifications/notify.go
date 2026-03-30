package notifications

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"service/internal/domain"
	"service/internal/system/notes"

	"github.com/gen2brain/beeep"
)

type NotificationArgs struct {
	DeviceId   string
	DeviceName string
	All        bool
}

func Notify(args NotificationArgs, notification Notification) {
	defer resolveNotification(args, notification)
	fmt.Println("Received notification from ", notification.App, " at ", notification.CreatedAt.Format("Mon, 02 Jan 2006 15:04:05 MST"), " saying ", notification.Message)
}

func resolveNotification(args NotificationArgs, notification Notification) {
	var device domain.Device
	if args.All {
		notifyAll(notification)
		return
	}
	fmt.Println(args)
	if args.DeviceId != "" {
		device = GetDeviceById(device.ID)
	} else if device.Name != "" {
		device = GetDeviceByName(device.Name)
	} else {
		return
	}

	device.OS = "ios"

	switch device.OS {
	case "ios":
		sendToPushover(notification)
	case "linux":
		sendToLinux(notification)
	case "macos":
		sendToWebSocket(notification)
	}
}

func notifyAll(notification Notification) error {
	sendToLinux(notification)
	sendToWebSocket(notification)
	sendToPushover(notification)

	return nil
}

func GetDeviceById(deviceId string) domain.Device {
	return domain.Device{Name: "hume", ID: "1", OS: "linux", Token: ""}
}

func GetDeviceByName(deviceName string) domain.Device {
	return domain.Device{Name: "iphone12", ID: "2", OS: "iOS", Token: ""}
}

func sendToPushover(notification Notification) {
	_, err := http.PostForm("https://api.pushover.net/1/messages.json", url.Values{
		"token":   {""},
		"user":    {""},
		"message": {notification.Message},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Notification sent!")
}

func sendToLinux(notification Notification) {
	beeep.Notify(notification.App, notification.Message, "")
	beeep.Beep(notes.C5, beeep.DefaultDuration)
	beeep.Beep(notes.G4, beeep.DefaultDuration)
	beeep.Beep(notes.E4, beeep.DefaultDuration)
}

func sendToWebSocket(notification Notification) {

}
