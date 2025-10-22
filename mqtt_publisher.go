package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// daftar parameter sesuai data kamu
var parameters = []string{
	"Entering_Chilled_Water_Temp",
	"Leaving_Chilled_Water_Temp",
	"Entering_Cooled_Water_Temp",
	"Leaving_Cooled_Water_Temp",
	"Leaving_Chilled_Water_Temp_Settings",
	"C1_Suction_Temp",
	"C1_Discharge_Temp",
	"C1_Suction_Pressure",
	"C1_Discharge_Pressure",
	"C1_Load",
	"C1_Current",
	"C1_Actual_Speed",
	"C1_Voltage",
	"C1_Power",
	"C2_Suction_Temp",
	"C2_Discharge_Temp",
	"C2_Suction_Pressure",
	"C2_Discharge_Pressure",
	"C2_Load",
	"C2_Current",
	"C2_Actual_Speed",
	"C2_Voltage",
	"C2_Power",
}

func DispatchMqttPublisher() {
	// ====== Konfigurasi MQTT ======
	broker := "tcp://localhost:1883"
	topic := "sensor/data"
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

	// ====== Loop kirim data setiap 2 detik ======
	for {
		data := make(map[string]map[string]float64)
		data["data"] = make(map[string]float64)
		for _, param := range parameters {
			data["data"][param] = float64(rand.Intn(100) + 1) // nilai acak 1–100
		}

		jsonData, _ := json.Marshal(data)
		token := client.Publish(topic, 0, false, jsonData)
		token.Wait()

		fmt.Println("📤 Published:", string(jsonData))
		time.Sleep(2 * time.Second)
	}
}
