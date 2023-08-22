package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ClickHouse/clickhouse-go/v2"
)

func main() {
	http.HandleFunc("/track-replay", trackReplay)

	createTable()

	fmt.Println("starting web server at ::1234")
	http.ListenAndServe(":1234", nil)
}

func trackReplay(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "POST" {
		var data Data
		err := dataParser(&data, r.Body)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		conn, _ := connect()
		// defer conn.Close()
		ctx := clickhouse.Context(r.Context(), clickhouse.WithStdAsync(false))

		{
			conn.Ping(ctx)

			deviceCheckQuery := fmt.Sprintf(`SELECT * FROM devices WHERE unique_id = '%s'`, data.Device.UniqueId)
			rows, err := conn.Query(ctx, deviceCheckQuery)
			if err != nil {
				panic(err)
			}

			if !rows.Next() {
				newDeviceQuery := fmt.Sprintf(`
				INSERT INTO devices (name, model, phone, status, contact, category, disabled, unique_id, attributes, expiration_time)
				VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%t', '%s', '%s', '%s')`, data.Device.Name, data.Device.Model, data.Device.Phone, data.Device.Status, data.Device.Contact, data.Device.Category, data.Device.Disabled, data.Device.UniqueId, data.Device.Attributes, data.Device.ExpirationTime)

				err = conn.Exec(ctx, newDeviceQuery)
				if err != nil {
					panic(err)
				}
			}
			rows.Close()

			newPositionQuery := fmt.Sprintf(`
			INSERT INTO positions (speed, valid, course, address, fix_time, network, accuracy, altitude, device_id, latitude, outdated, protocol, longitude, attributes, device_time, server_time)
			VALUES ('%d', '%t', '%d', '%s', '%s', '%v', '%d', '%d', '%s', '%f', '%t', '%s', '%f', '%s', '%s', '%s')`, data.Position.Speed, data.Position.Valid, data.Position.Course, data.Position.Address, data.Position.FixTime, data.Position.Network, data.Position.Accuracy, data.Position.Altitude, data.Device.UniqueId, data.Position.Latitude, data.Position.Outdated, data.Position.Protocol, data.Position.Longitude, data.Position.Attributes, data.Position.DeviceTime, data.Position.ServerTime)

			err = conn.Exec(ctx, newPositionQuery)
			if err != nil {
				panic(err)
			}
		}
		// fmt.Println(conn.Stats())
		// readFromDB()

		// w.WriteHeader(http.StatusOK)
		// Send json OK
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		return
	}

	// Send json error
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]string{"error": "bad request"})
}

func readFromDB() {
	conn, err := connect()
	if err != nil {
		panic(err)
	}
	// Read from database
	rows, err := conn.Query(context.Background(), `SELECT * FROM positions order by server_time desc limit 5`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	// var device Device = Device{}
	var position Position = Position{}
	for rows.Next() {
		// if err := rows.ScanStruct(&device); err != nil {
		// 	panic(err)
		// }
		// fmt.Printf("Row: %v\n", device)

		if err := rows.Scan(
			&position.Id,
			&position.Speed,
			&position.Valid,
			&position.Course,
			&position.Address,
			&position.FixTime,
			&position.Network,
			&position.Accuracy,
			&position.Altitude,
			&position.DeviceId,
			&position.Latitude,
			&position.Outdated,
			&position.Protocol,
			&position.Longitude,
			&position.Attributes,
			&position.DeviceTime,
			&position.ServerTime); err != nil {
			panic(err)
		}
		if err := rows.ScanStruct(&position); err != nil {
			panic(err)
		}
		fmt.Printf("Row: %v\n", position)

	}
}

func dataParser(data *Data, body io.ReadCloser) error {
	// x, _ := io.ReadAll(r.Body)
	// logToFile(string(x))

	decoder := json.NewDecoder(body)

	var dummy map[string]map[string]interface{}

	err := decoder.Decode(&dummy)
	if err != nil {
		return err
	}

	delete(dummy["device"], "id")
	delete(dummy["position"], "id")

	// Stringify attributes device
	a, err := json.Marshal(dummy["device"]["attributes"])
	if err != nil {
		return err
	}
	// Convert attributes to stringify version
	dummy["device"]["attributes"] = string(a)

	// Stringify attributes position
	a, err = json.Marshal(dummy["position"]["attributes"])
	if err != nil {
		return err
	}
	// Convert attributes to stringify version
	dummy["position"]["attributes"] = string(a)

	// Stringify position network
	a, err = json.Marshal(dummy["position"]["network"])
	if err != nil {
		return err
	}
	// Convert network to stringify version
	dummy["position"]["network"] = string(a)

	// Parse Dummy to Data
	// Convert to JSON
	b, err := json.Marshal(dummy)
	if err != nil {
		return err
	}

	// // Convert JSON to Data
	// var data Data
	err = json.Unmarshal(b, data)
	if err != nil {
		return err
	}

	return nil
}
