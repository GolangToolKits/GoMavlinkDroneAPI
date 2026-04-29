package gomavlinkdroneapi

import (
	"context"
	"log"
	"time"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/ardupilotmega"
)

func (s *DroneAPI) GetDroneChannel() *gomavlib.Channel {
	// Wait for the drone to send a heartbeat so we can capture its Channel
	log.Println("Waiting for heartbeat...")
	// var droneChannel *gomavlib.Channel

	// 1. Create a context that expires after 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Always clean up the context

	// 2. Read from the gomavlib Events channel
	eventsCh := s.drone.node.Events()

	for {
		select {
		// Case A: A new event is received from gomavlib
		case evt, ok := <-eventsCh:
			if !ok {
				log.Println("Event channel closed.")
				return nil
			}

			if frm, ok := evt.(*gomavlib.EventFrame); ok {
				if _, isHeartbeat := frm.Message().(*ardupilotmega.MessageHeartbeat); isHeartbeat {
					log.Println("Heartbeat received, channel captured!")
					return frm.Channel
				}
			}

		// Case B: The 10-second timer expires first
		case <-ctx.Done():
			log.Println("Error: Timed out waiting for drone heartbeat.")
			return nil
		}
	}
}
