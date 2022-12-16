package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
)

func calcstats(csvFile string) []string {
	rows := readLog(csvFile)
	var output []string
	// check data is received before doing anything, else will crash due to no data in the csv file
	if len(rows) > 1 {
		fmt.Printf("Successfully processed %d data points!\n", len(rows))
		rows, output = calculate(rows)
	}
	return output
}

func readLog(name string) [][]string {

	f, err := os.Open(name)
	// Usually we would return the error to the caller and handle
	// all errors in function `main()`. However, this is just a
	// small command-line tool, and so we use `log.Fatal()`
	// instead, in order to write the error message to the
	// terminal and exit immediately.
	if err != nil {
		log.Fatalf("Cannot open '%s': %s\n", name, err.Error())
	}

	// After this point, the file has been successfully opened,
	// and we want to ensure that it gets closed when no longer
	// needed, so we add a deferred call to `f.Close()`.
	defer f.Close()

	// To read in the CSV data, we create a new CSV reader that
	// reads from the input file.
	//
	// The CSV reader is aware of the CSV data format. It
	// separates the input stream into rows and columns,
	// and returns a slice of slices of strings.
	r := csv.NewReader(f)

	// We can even adjust the reader to recognize a semicolon,
	// rather than a comma, as the column separator.
	r.Comma = ','

	// Read the whole file at once. (We don't expect large files.)
	rows, err := r.ReadAll()

	// Again, we check for any error,
	if err != nil {
		log.Fatalln("Cannot read CSV data:", err.Error())
	}

	// and finally we can return the rows.
	return rows
}

func check(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}

// Calculate statistics during a race: Best lap time, track top speed, and lap sector times for La Selva Circuit
// Returns Best Lap Time, followed by Track Top Speed, then an array of times for Sectors 1-4
func calcRaceStats(csvFile string) (bestLapTime string, trackTopSpeed string, times []string) {
	rows := readLog(csvFile)
	// check data is received before doing anything, else will crash due to no data in the csv file
	if len(rows) < 1 {
		log.Fatalf("CSV File is empty!")
	}

	// Find row numbers based on column header names (row 0)
	bestLapRow := 0
	currentLapRow := 0
	distanceRow := 0
	lapNumRow := 0
	speedRow := 0

	for k, v := range rows[0] {
		if v == "BestLap" {
			bestLapRow = k
		} else if v == "CurrentLap" {
			currentLapRow = k
		} else if v == "DistanceTraveled" {
			distanceRow = k
		} else if v == "LapNumber" {
			lapNumRow = k
		} else if v == "Speed" {
			speedRow = k
		}
	}

	var t []float64  // array of current lap time values
	var d []float64  // array of distance values
	var l []float64  // array of lap number values
	var s []float64  // array of speed values
	var bl []float64 // array of BestLap time values

	for i := range rows {
		if i == 0 { // skip first row (header/column names)
			continue
		}

		// Add current lap times from row to array
		time, err := strconv.ParseFloat(rows[i][currentLapRow], 32)
		check(err)
		t = append(t, time)

		// Add distance values from row to array
		dist, err := strconv.ParseFloat(rows[i][distanceRow], 32)
		check(err)
		d = append(d, dist)

		// Add current lap numbers from row to array
		lap, err := strconv.ParseFloat(rows[i][lapNumRow], 32)
		check(err)
		l = append(l, lap)

		// Add speed value from row to array
		speed, err := strconv.ParseFloat(rows[i][speedRow], 32)
		check(err)
		s = append(s, (speed * 2.237)) // convert to MPH

		// Add Best lap times from row to array
		best, err := strconv.ParseFloat(rows[i][bestLapRow], 32)
		check(err)
		bl = append(bl, best)
	}

	// Find the best lap time
	bestLap := bl[len(bl)-1]
	// Check time at the end of the race if you're at the finish line
	if d[len(d)-1]-(5951*l[len(l)-1]) > 5951 {
		if t[len(t)-1]+0.0125 < bestLap {
			bestLap = t[len(t)-1] + 0.0125 // The game seems to take this much extra time when finishing the race
		}
	}
	// Convert to 00:00.000 time format
	min := strconv.FormatFloat(math.Floor(bestLap/60), 'f', 0, 32)
	if len(min) == 1 {
		min = "0" + min
	}
	sec := strconv.FormatFloat(math.Mod(bestLap, 60), 'f', 3, 32)
	if len(sec) == 1 {
		sec = "0" + sec
	}
	bestLapStr := min + ":" + sec

	// Find the track top speed
	sort.Float64s(s)
	topSpeed := s[len(s)-1]
	topSpeedStr := strconv.FormatFloat(topSpeed, 'f', 2, 32)

	// Calculate Track Sector Times
	s1Time := 0.0
	s2Time := 0.0
	s3Time := 0.0
	s4Time := 0.0
	s1TimeTmp := 0.0
	s2TimeTmp := 0.0
	s3TimeTmp := 0.0
	s4TimeTmp := 0.0
	s1End := 0.0
	s2End := 0.0
	s3End := 0.0
	dist := 0.0
	endLapDist := 0.0

	for i, val := range d {
		if i == 0 || l[i] != l[i-1] {
			endLapDist = val
		}
		dist = val - endLapDist

		if dist > 1878 && dist < 1879 {
			s1End = t[i]
			s1TimeTmp = s1End
		} else if dist > 3184 && dist < 3185 {
			s2End = t[i]
			s2TimeTmp = s2End - s1End
		} else if dist > 4311 && dist < 4312 {
			s3End = t[i]
			s3TimeTmp = s3End - s2End
		} else if i == len(l)-1 || l[i] != l[i+1] {
			s4TimeTmp = (t[i] - s3End) + 0.0125
			if t[i] <= bestLap {
				s1Time = s1TimeTmp
				s2Time = s2TimeTmp
				s3Time = s3TimeTmp
				s4Time = s4TimeTmp
			}
		}
	}

	// Convert to hh:mm:ss.000 duration format
	minutes := strconv.FormatFloat(math.Floor(s1Time/60), 'f', 0, 32)
	seconds := strconv.FormatFloat(math.Mod(s1Time, 60), 'f', 3, 32)
	if math.Mod(s1Time, 60) < 1 {
		seconds = "0" + seconds
	}
	s1TimeStr := "00:" + minutes + ":" + seconds

	minutes = strconv.FormatFloat(math.Floor(s2Time/60), 'f', 0, 32)
	seconds = strconv.FormatFloat(math.Mod(s2Time, 60), 'f', 3, 32)
	if math.Mod(s2Time, 60) < 1 {
		seconds = "0" + seconds
	}
	s2TimeStr := "00:" + minutes + ":" + seconds

	minutes = strconv.FormatFloat(math.Floor(s3Time/60), 'f', 0, 32)
	seconds = strconv.FormatFloat(math.Mod(s3Time, 60), 'f', 3, 32)
	if math.Mod(s3Time, 60) < 1 {
		seconds = "0" + seconds
	}
	s3TimeStr := "00:" + minutes + ":" + seconds

	minutes = strconv.FormatFloat(math.Floor(s4Time/60), 'f', 0, 32)
	seconds = strconv.FormatFloat(math.Mod(s4Time, 60), 'f', 3, 32)
	if math.Mod(s4Time, 60) < 1 {
		seconds = "0" + seconds
	}
	s4TimeStr := "00:" + minutes + ":" + seconds

	var sectorTimes []string
	sectorTimes = append(sectorTimes, s1TimeStr, s2TimeStr, s3TimeStr, s4TimeStr)

	return bestLapStr, topSpeedStr, sectorTimes
}

// calculate stats
func calculate(rows [][]string) (r [][]string, out []string) {
	// Find row numbers based on column header names (row 0)
	powerRow := 0
	torqueRow := 0
	speedRow := 0
	boostRow := 0
	classPIRow := 0
	driveTrainRow := 0
	timeRow := 0
	gearRow := 0
	classLetterRow := 0

	for k, v := range rows[0] {
		if v == "Speed" {
			speedRow = k
		} else if v == "Boost" {
			boostRow = k
		} else if v == "CarPerformanceIndex" {
			classPIRow = k
		} else if v == "DrivetrainType" {
			driveTrainRow = k
		} else if v == "Power" {
			powerRow = k
		} else if v == "Torque" {
			torqueRow = k
		} else if v == "TimestampMS" {
			timeRow = k
		} else if v == "Gear" {
			gearRow = k
		} else if v == "CarClass" {
			classLetterRow = k
		}
	}

	var output []string
	carClass := rows[1][classLetterRow]
	var t []float64  // array of timestamp values
	var s []float64  // array of speed values
	var b []float64  // array of boost values
	var p []float64  // array of power values
	var tq []float64 // array of torque values
	var g []float64  // array of gear values
	// fmt.Printf("%T\n", s) // print type

	for i := range rows {

		if i == 0 { // skip first row (header/column names)
			continue
		}

		// Add timestamps from row to array
		time, err := strconv.ParseFloat(rows[i][timeRow], 32)
		check(err)
		t = append(t, (time / 1000)) // convert from milliseconds to seconds

		// Add power value from row to array
		power, err := strconv.ParseFloat(rows[i][powerRow], 32) // convert power string to int
		check(err)
		p = append(p, (power * 0.0013410220888)) // convert from Watts to Mechanical Horsepower

		// Add torque value from row to array
		torque, err := strconv.ParseFloat(rows[i][torqueRow], 32) // convert torque string to int
		check(err)
		tq = append(tq, (torque * 0.7375621493)) // convert from nm to ft-lb

		// Add speed value from row to array
		speed, err := strconv.ParseFloat(rows[i][speedRow], 32) // convert speed string to int
		check(err)
		s = append(s, (speed * 2.237)) // convert to MPH

		// Add boost value from row to array
		boost, err := strconv.ParseFloat(rows[i][boostRow], 32) // convert boost string to int
		check(err)
		b = append(b, boost) // convert to PSI, not 100% sure what value this is natively?

		// Add gear value from row to array
		gear, err := strconv.ParseFloat(rows[i][gearRow], 32) // convert speed string to int
		check(err)
		g = append(g, (gear))
	}

	var totalSpeed float64
	for _, value := range s {
		totalSpeed += value // add all speed values together for getting average later
	}

	//fmt.Printf("\nRace statistics:\n")

	// Get average speed
	//fmt.Printf("Average speed: %.2f MPH \n", totalSpeed/float64(len(s))) // truncate to 2 decimal places

	//Get PI Index Number
	pINum := rows[1][classPIRow]
	output = append(output, pINum)

	//Get Drivetrain Type
	drivetrainStr := ""
	drivetrainNum, err := strconv.ParseInt(rows[1][driveTrainRow], 10, 32)
	check(err)
	if drivetrainNum == 0 {
		drivetrainStr = "FWD"
	} else if drivetrainNum == 1 {
		drivetrainStr = "RWD"
	} else if drivetrainNum == 2 {
		drivetrainStr = "AWD"
	}
	output = append(output, drivetrainStr)

	// Get peak horsepower
	// Only looks at power numbers when the car is in 2nd gear or higher,
	// because when bouncing off the rev limiter during a launch the game
	// will output higher horsepower numbers than the car actually has.
	var adjustedPowers []float64
	var gearTotal float64
	for _, value := range g {
		gearTotal += value
	}
	if gearTotal == float64(len(g)) { // If the car only has 1 gear (i.e. electric cars) then don't adjust
		adjustedPowers = p
	}

	for i, value := range p {
		if g[i] != 1 { // If the car is not in 1st gear,
			adjustedPowers = append(adjustedPowers, value) // add that power number to the adjusted list
		}
	}
	sort.Float64s(adjustedPowers)
	topPower := adjustedPowers[len(adjustedPowers)-1]
	output = append(output, strconv.FormatFloat(topPower, 'f', 0, 32))

	// Get peak torque
	sort.Float64s(tq)
	topTorque := tq[len(tq)-1]
	output = append(output, strconv.FormatFloat(topTorque, 'f', 0, 32))

	// Get 0-60mph time
	zeroTo60, err := getTimeBetween(0, 60, t, s)
	if err != nil {
		output = append(output, "Failed!")
	} else {
		output = append(output, strconv.FormatFloat(zeroTo60, 'f', 3, 32))
	}

	// Get 0-100mph time
	zeroTo100, err := getTimeBetween(0, 100, t, s)
	if err != nil {
		output = append(output, "Failed!")
	} else {
		output = append(output, strconv.FormatFloat(zeroTo100, 'f', 3, 32))
	}

	if carClass == "0" || carClass == "1" { // D Class or C Class
		// Get 25-75mph time
		twentyfiveTo75, err := getTimeBetween(25, 75, t, s)
		if err != nil {
			output = append(output, "Failed!")
		} else {
			output = append(output, strconv.FormatFloat(twentyfiveTo75, 'f', 3, 32))
		}

		// Get 50-100mph time
		fiftyTo100, err := getTimeBetween(50, 100, t, s)
		if err != nil {
			output = append(output, "Failed!")
		} else {
			output = append(output, strconv.FormatFloat(fiftyTo100, 'f', 3, 32))
		}
	} else if carClass == "2" || carClass == "3" { // B Class or A Class
		// Get 50-100mph time
		fiftyTo100, err := getTimeBetween(50, 100, t, s)
		if err != nil {
			output = append(output, "Failed!")
		} else {
			output = append(output, strconv.FormatFloat(fiftyTo100, 'f', 3, 32))
		}

		// Get 60-150mph time
		sixtyTo150, err := getTimeBetween(60, 150, t, s)
		if err != nil {
			output = append(output, "Failed!")
		} else {
			output = append(output, strconv.FormatFloat(sixtyTo150, 'f', 3, 32))
		}
	} else { // S1, S2 and X Class
		// Get 60-150mph time
		sixtyTo150, err := getTimeBetween(60, 150, t, s)
		if err != nil {
			output = append(output, "Failed!")
		} else {
			output = append(output, strconv.FormatFloat(sixtyTo150, 'f', 3, 32))
		}

		// Get 100-200mph time
		hundredTo200, err := getTimeBetween(100, 200, t, s)
		if err != nil {
			output = append(output, "Failed!")
		} else {
			output = append(output, strconv.FormatFloat(hundredTo200, 'f', 3, 32))
		}
	}

	// Get 60-0mph time
	sixtytoZero, err := getTimeBetween(60, 0, t, s)
	if err != nil {
		output = append(output, "Failed!")
	} else {
		output = append(output, strconv.FormatFloat(sixtytoZero, 'f', 3, 32))
	}

	// Get 100-0mph time
	hundredToZero, err := getTimeBetween(100, 0, t, s)
	if err != nil {
		output = append(output, "Failed!")
	} else {
		output = append(output, strconv.FormatFloat(hundredToZero, 'f', 3, 32))
	}

	// Get top speed
	// fmt.Println(s)
	sort.Float64s(s)
	topSpeed := s[len(s)-1]
	//fmt.Printf("Top speed: %.2f MPH \n", topSpeed)
	output = append(output, strconv.FormatFloat(topSpeed, 'f', 2, 32))

	// Get peak boost
	sort.Float64s(b)
	topBoost := b[len(b)-1]
	//fmt.Printf("Peak boost: %.2f PSI \n", topBoost)
	output = append(output, strconv.FormatFloat(topBoost, 'f', 1, 32))

	return rows, output
}

// Returns a floating point number representing how fast (in seconds)
// the car traveled from the start speed to the end speed. (ex: 0-60 mph or 100-0 mph)
// Returns an error if car does not reach either speed.
func getTimeBetween(startSpeed float64, endSpeed float64, timeValues []float64, speedValues []float64) (float64, error) {
	var startTime float64
	var endTime float64
	var speeds []float64

	if len(timeValues) != len(speedValues) {
		return 0, errors.New("Time Array and Speed Array lengths do not match")
	}

	for _, value := range speedValues {
		speeds = append(speeds, value)
	}
	sort.Float64s(speeds)
	minSpeed := speeds[0]
	maxSpeed := speeds[len(speeds)-1]

	// If either the Start or End speed is beyond what the car actually performed, returns an Error.
	if startSpeed > maxSpeed+0.1 || endSpeed > maxSpeed+0.1 || (minSpeed > 0.1 && startSpeed < minSpeed) || (minSpeed > 0.1 && endSpeed < minSpeed) {
		return 0, errors.New("Start or End Speed is outside of data range.")
	}

	if startSpeed > endSpeed { // If startSpeed > endSpeed, then we're calculating deceleration (ex: 60-0mph)
		for i := 0; i < len(speedValues); i++ {
			if speedValues[i] > startSpeed+0.1 {
				startTime = timeValues[i]
			} else if speedValues[i] > endSpeed+0.1 {
				endTime = timeValues[i]
			}
		}
	} else { // Otherwise, we're calculating acceleration (ex: 0-60mph)
		// This solution stops after finding the first moment the car passes the start speed and ignores
		// possible low speeds later on, so that the startTime isn't greater than the endTime.
		// (For best usage this assumes a test run accelerating from 0mph to Top Speed, then back to 0mph)
		i := 0
		for speedValues[i] < startSpeed+0.1 {
			startTime = timeValues[i]
			i++
		}
		j := 0
		for speedValues[j] < endSpeed+0.1 {
			endTime = timeValues[j]
			j++
		}
	}

	if startTime > endTime {
		return 0, errors.New("Calculated negative time, something went wrong with the data input.")
	}
	return (endTime - startTime), nil
}

// Returns the Ordinal Number of the current car, or the first car used during data collection.
func getOrdinalNumber(csvFile string) (string, error) {
	rows := readLog(csvFile)
	if len(rows) < 1 {
		return "", errors.New("csvFile is empty.")
	}
	return rows[1][53], nil
}

func getAllOrdinalNumbers(csvFile string) ([]string, error) {
	rows := readLog(csvFile)
	if len(rows) < 1 {
		return nil, errors.New("csvFile is empty.")
	}
	var o []string
	lastVal := ""
	for i := range rows {
		if i == 0 { // skip first row (header/column names)
			continue
		}
		// Add Ordinals from row to array
		num := rows[i][53]
		if lastVal == "" || num != lastVal {
			o = append(o, num)
			lastVal = num
		}
	}
	return o, nil
}

// `intToFloatString` takes an integer `n` and calculates the floating point value representing `n/100` as a string.
// func intToFloatString(n int) string {
// 	intgr := n / 100
// 	frac := n - intgr*100
// 	return fmt.Sprintf("%d.%d", intgr, frac)
// }

// func printSlice(s []int) {
// 	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
// }
