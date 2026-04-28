package gomavlinkdroneapi

import (
	"log"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/ardupilotmega"
)

func (s *DroneAPI) GetDroneChannel() *gomavlib.Channel {
	// Wait for the drone to send a heartbeat so we can capture its Channel
	log.Println("Waiting for heartbeat...")
	var droneChannel *gomavlib.Channel

	for evt := range s.drone.node.Events() {
		if frm, ok := evt.(*gomavlib.EventFrame); ok {
			if _, isHeartbeat := frm.Message().(*ardupilotmega.MessageHeartbeat); isHeartbeat {
				droneChannel = frm.Channel // Captured the actual communication channel!
				break
			}
		}
	}
	return droneChannel
}
