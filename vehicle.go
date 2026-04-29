package gomavlinkdroneapi

import (
	"log"
	"time"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

func (s *DroneAPI) GetConnectedVehicle() (targetSystem byte, targetComponent byte) {
	log.Println("Waiting for vehicle heartbeat...")

	// 1. Initialize a 10-second timeout timer
	timeoutDuration := 10 * time.Second
	timer := time.NewTimer(timeoutDuration)
	defer timer.Stop()

	// 2. Continuous loop to process events
	for {
		select {
		// Handle incoming events from gomavlib
		case evt, ok := <-s.drone.node.Events():
			if !ok {
				log.Println("Event channel closed.")
				return 0, 0
			}

			// Filter strictly for incoming packet frames
			if frm, ok := evt.(*gomavlib.EventFrame); ok {

				// Narrow down to Heartbeat messages to identify the target
				if _, ok := frm.Message().(*common.MessageHeartbeat); ok {
					targetSystem = frm.SystemID()
					targetComponent = frm.ComponentID()

					log.Printf("Connected to Vehicle!\n")
					log.Printf("TargetSystem: %d\n", targetSystem)
					log.Printf("TargetComponent: %d\n", targetComponent)

					// Successfully found the heartbeat, return the IDs
					return targetSystem, targetComponent
				}
			}

		// Trigger if no heartbeat is received within the duration
		case <-timer.C:
			log.Println("Timeout: No vehicle heartbeat received.")
			return 0, 0
		}
	}
}
