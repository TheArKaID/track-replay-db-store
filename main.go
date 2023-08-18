package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ClickHouse/clickhouse-go/v2"
)

func main() {
	http.HandleFunc("/track-replay", trackReplay)

	createTable()

	fmt.Println("starting web server at ::1234")
	http.ListenAndServe(":1234", nil)
}

// Log to a file
func logToFile(message string) {
	// Open a file
	f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// Close the file
	defer f.Close()

	// Write to the file, and add new line
	if _, err := f.WriteString(message + "\n"); err != nil {
		log.Fatal(err)
	}
}

func trackReplay(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "POST" {

		// x, _ := io.ReadAll(r.Body)
		// logToFile(string(x))

		decoder := json.NewDecoder(r.Body)

		var dummy map[string]map[string]interface{}

		err := decoder.Decode(&dummy)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		delete(dummy["device"], "id")
		delete(dummy["position"], "id")

		// Stringify attributes device
		a, err := json.Marshal(dummy["device"]["attributes"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Convert attributes to stringify version
		dummy["device"]["attributes"] = string(a)

		// Stringify attributes position
		a, err = json.Marshal(dummy["position"]["attributes"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Convert attributes to stringify version
		dummy["position"]["attributes"] = string(a)

		// Parse Dummy to Data
		// Convert to JSON
		b, err := json.Marshal(dummy)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// // Convert JSON to Data
		var data Data
		err = json.Unmarshal(b, &data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Print(data.Device.Phone)
		fmt.Print("\n")
		fmt.Print(data.Device.Status)
		fmt.Print("\n")
		fmt.Print(data.Device.Name)
		fmt.Print("\n")
		conn, _ := connect()
		ctx := clickhouse.Context(context.Background(), clickhouse.WithStdAsync(false))
		{
			conn.Ping(ctx)

			query := fmt.Sprintf(`
			INSERT INTO devices (name, model, phone, status, contact, category, disabled, unique_id, attributes, expiration_time)
			VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%t', '%s', '%s', '%s')`, data.Device.Name, data.Device.Model, data.Device.Phone, data.Device.Status, data.Device.Contact, data.Device.Category, data.Device.Disabled, data.Device.UniqueId, data.Device.Attributes, data.Device.ExpirationTime)

			err := conn.AsyncInsert(ctx, query, true)
			if err != nil {
				panic(err)
			}

			query = fmt.Sprintf(`
			INSERT INTO positions (speed, valid, course, address, fix_time, network, accuracy, altitude, device_id, latitude, outdated, protocol, longitude, attributes, device_time, server_time)
			VALUES ('%d', '%t', '%d', '%s', '%s', '%s', '%d', '%d', '%s', '%f', '%t', '%s', '%f', '%s', '%s', '%s')`, data.Position.Speed, data.Position.Valid, data.Position.Course, data.Position.Address, data.Position.FixTime, data.Position.Network, data.Position.Accuracy, data.Position.Altitude, data.Device.UniqueId, data.Position.Latitude, data.Position.Outdated, data.Position.Protocol, data.Position.Longitude, data.Position.Attributes, data.Position.DeviceTime, data.Position.ServerTime)

			err = conn.AsyncInsert(ctx, query, true)
			if err != nil {
				panic(err)
			}
		}

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
	rows, err := conn.Query(context.Background(), `SELECT * FROM devices order by last_update desc limit 5`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var device Device = Device{}
	for rows.Next() {
		// if err := rows.Scan(
		// 	&d.Id,
		// 	&d.Name,
		// 	&d.Phone,
		// 	&d.LastUpdate); err != nil {
		// 	panic(err)
		// }
		if err := rows.ScanStruct(&device); err != nil {
			panic(err)
		}
		fmt.Printf("Row: %v\n", device)

	}
}

// function to format JSON data
// func formatJSON(data []byte) (string, error) {
// 	var out bytes.Buffer
// 	err := json.Indent(&out, data, "", " ")

// 	if err != nil {
// 		return "", err
// 	}

// 	d := out.Bytes()
// 	return string(d), nil
// }
