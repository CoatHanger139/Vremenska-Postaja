// Package vreme provides a very simple "simulation" of weather stations

package vreme

import ( 
	"fmt"
)

// Holds all measurements and data on where they were taken.
// Also hold min. and max. allowed measure temp.
type Postaja struct {
	ime     string
	kraj    string
	meritve []float64
	minTemp float64
	maxTemp float64
}

// Creates Postaja with name, location and min. and max. allowd temp.
// Returns a new Postaja with empty measurements array.
func UstvariPostajo( p_ime, p_kraj string, p_mintemp, p_maxTemp float64 ) Postaja {
	return Postaja{ ime: p_ime, kraj: p_kraj, minTemp: p_mintemp, maxTemp: p_maxTemp }
}

// Adds a measurement to a specified station.
// Checks if the measurement is valid against min. and max. allowed temp.
// if not valid it is rejected or station not found in map (prints error).
func DodajMeritev( p_stationMap map[ string ]Postaja, p_key string, p_temp float64 ) {
	station, ok := p_stationMap[ p_key ]
	if !ok {
		fmt.Printf( "Station sign %s does not exist.\n", p_key )
		return
	}
	if p_temp < station.minTemp || p_temp > station.maxTemp {
		fmt.Printf( "Invalid temperature (%.2f°C).\n", p_temp )
		return
	}

	station.meritve = append( station.meritve, p_temp )
	p_stationMap[ p_key ] = station
}

// Calculates max. min. and avg. measured temp. on given sation
// Returnes max. min. avg. temp and ok if station was found 
// and -1, -1, -1 and !ok if station not found
func statistika( p_stationMap map[ string ]Postaja, p_key string ) ( high, low, avg float64, ok bool ) {
	station, found := p_stationMap[ p_key ]
    if !found || len( station.meritve ) < 3 {
        return -1, -1, -1, false
    }

	ok = true
	high = station.meritve[ 0 ]
	low = station.meritve[ 0 ]
	sum := 0.0
	for _, entry := range station.meritve {
		high = max( entry, high )
		low  = min( entry, low )

		sum += entry
	}
	avg = sum / float64( len( station.meritve ) )
	return
}

// Prints map key, name, location and all measurements for all stations.
func IzpisMeritev( p_stationMap map[ string ]Postaja ) {
	for key, station := range p_stationMap {
		fmt.Printf( "%s - %s (%s): %.2f\n", key, station.ime, station.kraj, station.meritve )
	}
}

// Prints a report for every sation.
// A report contains the station name, its min. max. and avg temps 
// and a description of the avg. temperature
// >25 --> Vroče
// <25 --> Toplo
// <15 --> Hladno
// <5  --> Mrzlo 
func IzpisPorocila( p_stationMap map[ string ]Postaja ) {
	for key, station := range p_stationMap {
		high, low, avg, ok := statistika( p_stationMap, key )

		if !ok {
			fmt.Printf( "%s: premalo meritev\n", station.ime )
			continue
		}

		var desc string
		switch {
			case avg >= 25:
				desc = "Vroče"
			case avg >= 15:
				desc = "Toplo"
			case avg >= 5:
				desc = "Hladno"
			default:
				desc = "Mrzlo"
		}

		fmt.Printf( "%s: min %.2f °C, max %.2f°C, avg %.2f°C -> %s\n", station.ime, low, high, avg, desc )
	}
}