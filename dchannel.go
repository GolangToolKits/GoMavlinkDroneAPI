package gomavlinkdroneapi

import (
	"log"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/ardupilotmega"
)

func (s *DroneAPI) GetDroneChannel() *gomavlib.Channel {
	// Wait for the drone to send a heartbeat so we can capture its Channel
	log.Println("Waiting for heartbeat indefinitely...")

	// 1. Read directly from the gomavlib Events channel using range.
	// This will safely process events forever until a heartbeat is found
	// OR until the library automatically closes the channel.
	for evt := range s.drone.node.Events() {

		if frm, ok := evt.(*gomavlib.EventFrame); ok {
			if _, isHeartbeat := frm.Message().(*ardupilotmega.MessageHeartbeat); isHeartbeat {
				log.Println("Heartbeat received, channel captured!")
				return frm.Channel
			}
		}
	}

	// If the range loop finishes, it means node.Events() was closed.
	log.Println("Event channel closed.")
	return nil
}
