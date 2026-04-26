
Yes, you absolutely need to check that the connection worked. [1] 
In Go, calling gomavlib.NewNode() only checks that your local network adapter or serial port is accessible. It does not verify that a physical drone is actively listening, powered on, or communicating back to you. [2, 3] 
To confirm a successful connection, you must wait to receive a HEARTBEAT message from the drone's System ID. [4, 5] 
------------------------------
## 💓 The MAVLink Heartbeat Rule [6] 
All MAVLink-compatible flight controllers (like PX4 and ArduPilot) broadcast a HEARTBEAT message once every second. If you aren't receiving heartbeats, you are not connected. [4, 5, 6, 7, 8] 
Here is how you can perform a connection check in gomavlib with a timeout fallback:
```
package main
import (
	"context"
	"log"
	"time"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)
func main() {
	node, err := gomavlib.NewNode(gomavlib.NodeConf{
		Endpoints: []gomavlib.EndpointConf{
			gomavlib.EndpointUDPClient{Address: "127.0.0.1:14550"},
		},
		Dialect:     common.Dialect,
		OutVersion:  gomavlib.V2,
		OutSystemID: 255,
	})
	if err != nil {
		log.Fatalf("Failed to create local node: %v", err)
	}
	defer node.Close()

	log.Println("Node started. Waiting for drone heartbeat...")

	// Create a context that times out after 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	connected := false

	// Loop through events until we find a heartbeat or the timeout expires
	for {
		select {
		case <-ctx.Done():
			log.Fatal("Connection failed: No heartbeat received from drone.")
		case evt, ok := <-node.Events():
			if !ok {
				log.Fatal("Event channel closed.")
			}

			if frm, ok := evt.(*gomavlib.EventFrame); ok {
				// Check if the message is a HEARTBEAT (Message ID 0)
				if frm.Message().GetID() == 0 {
					log.Printf("🎉 Success! Drone connected. System ID: %d", frm.SystemID())
					connected = true
					break
				}
			}
		}
		if connected {
			break
		}
	}

	// Your drone control logic goes here!
}
```
------------------------------
## 📋 Best Practices for Connection Management

* Track the System ID: Drones usually broadcast as System ID: 1. Ground stations use 255. Catching the drone's specific ID prevents you from accidentally targeting a component like a gimbal or a local simulation bridge.
* Handle In-Flight Disconnections: You shouldn't check for a heartbeat just once. In a production application, you should record the timestamp of the last received heartbeat. If more than 5 seconds pass without a heartbeat, your code should trigger a failsafe or pause operations.
* Double Check UDP Ports: If you are using UDP and don't see a heartbeat, double check that your drone and computer aren't both set to "Client" mode, leaving no one acting as the "Server" listening for packets. [4, 9, 10, 11] 

Would you like an example of how to continuously track the heartbeat in a background routine to detect mid-flight disconnections?

[1] [https://www.youtube.com](https://www.youtube.com/watch?v=esNTt3939BA&t=24)
[2] [https://mavsdk.mavlink.io](https://mavsdk.mavlink.io/main/en/cpp/guide/connections.html)
[3] [https://stackoverflow.com](https://stackoverflow.com/questions/79220849/gomavlib-arming-command-not-working-with-qgroundcontrol)
[4] [https://mavlink.io](https://mavlink.io/en/services/heartbeat.html)
[5] [https://discuss.bluerobotics.com](https://discuss.bluerobotics.com/t/sending-and-recieving-heartbeats-to-bluerov-from-surface-computer-using-pymavlink/10940#:~:text=Likely%20the%20most%20intended%20way%20%28including%20how,master.recv_match%28type=%27ATTITUDE%27%2C%20blocking=True%29%20if%20display.trigger%28%29:%20print%28att.roll%2C%20att.pitch%2C%20att.yaw%29)
[6] [https://www.youtube.com](https://www.youtube.com/watch?v=PTd5ViDsH1Q&t=3)
[7] [https://discuss.bluerobotics.com](https://discuss.bluerobotics.com/t/heartbeat-not-detected-on-blueos/12570#:~:text=Here%20are%20some%20things%20you%20can%20try:,available.%20Heartbeats%20should%20come%20roughly%20every%20second.)
[8] [https://github.com](https://github.com/bluenviron/gomavlib/blob/main/README.md)
[9] [https://github.com](https://github.com/bluenviron/gomavlib/discussions/233)
[10] [https://ardupilot.org](https://ardupilot.org/dev/docs/mavlink-basics.html)
[11] [https://github.com](https://github.com/mavlink/mavros/issues/1768)
