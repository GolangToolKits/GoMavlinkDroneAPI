package gomavlinkdroneapi

import (
	"context"
	"log"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

func (s *DroneAPI) GetConnectedVehicle(ctx context.Context) (targetSystem byte, targetComponent byte) {
	log.Println("Waiting for vehicle heartbeat...")

	for {
		select {
		case <-ctx.Done():
			log.Println("Search aborted or timed out.")
			return 0, 0

		case evt, ok := <-s.drone.node.Events():
			if !ok {
				return 0, 0
			}

			if frm, ok := evt.(*gomavlib.EventFrame); ok {
				if msg, ok := frm.Message().(*common.MessageHeartbeat); ok {

					// 1. Filter for actual vehicles (Type < 5 usually covers fixed-wing/quad/sub)
					// This avoids accidentally connecting to a GCS (Type 6)
					if msg.Type == common.MAV_TYPE_GCS || msg.Type == common.MAV_TYPE_ONBOARD_CONTROLLER {
						continue
					}

					log.Printf("Connected to Vehicle on %s!", frm.Channel.String())
					log.Printf("TargetSystem: %d, TargetComponent: %d", frm.SystemID(), frm.ComponentID())

					return frm.SystemID(), frm.ComponentID()
				}
			}
		}
	}
}
