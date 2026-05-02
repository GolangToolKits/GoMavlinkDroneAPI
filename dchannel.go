package gomavlinkdroneapi

import (
	"context"
	"log"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/ardupilotmega"
)

func (s *DroneAPI) GetDroneChannel(ctx context.Context) *gomavlib.Channel {
	log.Println("Waiting for heartbeat...")

	for {
		select {
		case <-ctx.Done():
			log.Printf("Wait cancelled or timed out: %v", ctx.Err())
			return nil
		case evt, ok := <-s.drone.node.Events():
			if !ok {
				log.Println("Node event channel closed")
				return nil
			}

			// Use a type switch for cleaner message handling
			frm, ok := evt.(*gomavlib.EventFrame)
			if !ok {
				continue
			}

			if _, isHB := frm.Message().(*ardupilotmega.MessageHeartbeat); isHB {
				log.Printf("Heartbeat detected from %v", frm.Channel)
				return frm.Channel
			}
		}
	}
}
