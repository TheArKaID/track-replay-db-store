package main

// // Log to a file
// func logToFile(message string) {
// 	// Open a file
// 	f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Close the file
// 	defer f.Close()

// 	// Write to the file, and add new line
// 	if _, err := f.WriteString(message + "\n"); err != nil {
// 		log.Fatal(err)
// 	}
// }

// // function to format JSON data
// func formatJSON(data []byte) (string, error) {
// 	var out bytes.Buffer
// 	err := json.Indent(&out, data, "", " ")

// 	if err != nil {
// 		return "", err
// 	}

// 	d := out.Bytes()
// 	return string(d), nil
// }
