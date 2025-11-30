package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func DispatchMqttPublisher() {

	broker := "tcp://localhost:1883"
	topic := "sensor/payload"
	clientID := "golang-random-publisher"

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	opts.OnConnect = func(c mqtt.Client) {
		fmt.Println("✅ Connected to MQTT broker:", broker)
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	rand.Seed(time.Now().UnixNano())

	for {
		// ==== struktur JSON ====
		payload := make(map[string]interface{})
		data := make(map[string]interface{})

		// ===== isi 15 property =====

		data["1_Chiller_Operating_State"] = []bool{rand.Intn(2) == 1}
		data["1_Entering_Chilled_Water_Temp"] = []int{rand.Intn(120)}
		data["1_Compressor1_Load"] = []int{rand.Intn(800)}
		data["1_Compressor2_Load"] = []int{rand.Intn(800)}
		data["1_Leaving_Chilled_Water_Temp"] = []int{rand.Intn(120)}
		data["1_Entering_Cooled_Water_Temp"] = []int{rand.Intn(400)}
		data["1_Leaving_Cooled_Water_Temp"] = []int{rand.Intn(400)}
		data["1_Leaving_Chilled_Water_Temp_Settings"] = []int{60}
		data["1_Chiller_COP"] = []float64{rand.Float64()*8 + 1}
		data["1_Comp1Power"] = []int{rand.Intn(700)}
		data["1_Comp2Power"] = []int{rand.Intn(700)}
		data["1_Comp1RunningTime"] = []int{22909}
		data["1_Comp2RunningTime"] = []int{22914}
		data["1_Comp1OperatingState"] = []bool{rand.Intn(2) == 1}
		data["1_Comp2OperatingState"] = []bool{rand.Intn(2) == 1}

		// ===== masukkan ke payload =====
		payload["d"] = data
		payload["ts"] = time.Now().Format(time.RFC3339Nano)

		jsonData, _ := json.Marshal(payload)

		token := client.Publish(topic, 0, false, jsonData)
		token.Wait()

		fmt.Println("📤 Published:", string(jsonData))

		time.Sleep(2 * time.Second)
	}
}
