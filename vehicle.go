package gomavlinkdroneapi

import (
	"log"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

func (s *DroneAPI) GetConnectedVehicle() (targetSystem byte, targetComponent byte) {
	log.Println("Waiting for vehicle heartbeat...")
	// var targetSystem byte
	// var targetComponent byte

	// 2. Loop through incoming events
	for evt := range s.drone.node.Events() {
		// Filter strictly for incoming packet frames
		if frm, ok := evt.(*gomavlib.EventFrame); ok {

			// Narrow down to Heartbeat messages to identify the target
			if _, ok := frm.Message().(*common.MessageHeartbeat); ok {

				// 3. Extract your TargetSystem and TargetComponent
				targetSystem = frm.SystemID()
				targetComponent = frm.ComponentID()

				log.Printf("🎯 Connected to Vehicle!\n")
				log.Printf("TargetSystem: %d\n", targetSystem)
				log.Printf("TargetComponent: %d\n", targetComponent)

				// Break
				break
			}
		}
	}
	return targetSystem, targetComponent
}
