/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2024-09-22 14:49:31
 * @LastEditTime: 2024-10-04 23:22:56
 * @LastEditors: FunctionSir
 * @Description: -
 * @FilePath: /AKBP/midware-dokuwiki/kml.go
 */

package main

import (
	"fmt"
	"strings"
)

type Coord struct {
	Lat string
	Lon string
	Asl string
}

const (
	KML_FMT_STR string = `<?xml version="1.0" encoding="UTF-8"?>` +
		`<kml xmlns="http://www.opengis.net/kml/2.2" xmlns:gx="http://www.google.com/kml/ext/2.2" ` +
		`xmlns:kml="http://www.opengis.net/kml/2.2" xmlns:atom="http://www.w3.org/2005/Atom">` +
		`<Document><name>%s</name>%s</Document></kml>` // Gened by Google Earth Pro. %s: Name, %s: Content.
	PLACEMARK_FMT_STR string = "<Placemark><name>%s</name><description>%s</description>%s</Placemark>" // %s: Name, %s: Description %s: Content.
	LINE_FMT_STR      string = "<LineString><coordinates>%s</coordinates></LineString>"                // %s: CoordinatesList.
	POINT_FMT_STR     string = "<Point><coordinates>%s</coordinates></Point>"                          // %s: Coord of the point.
)

func mkCoordsListStr(coords []Coord) string {
	r := ""
	for _, x := range coords {
		asl := x.Asl
		if x.Asl == "nil" || x.Asl == "" {
			asl = "0"
		}
		r += x.Lon + "," + x.Lat + "," + asl + " "
	}
	return strings.TrimSpace(r)
}

func coordToStr(coord *Coord) string {
	asl := coord.Asl
	if coord.Asl == "nil" || coord.Asl == "" {
		asl = "0"
	}
	return coord.Lon + "," + coord.Lat + "," + asl
}

func KmlAssemble(name string, content string) string {
	return fmt.Sprintf(KML_FMT_STR, name, content)
}

func KmlLineString(name string, description string, coords []Coord) string {
	coordsList := mkCoordsListStr(coords)
	return fmt.Sprintf(PLACEMARK_FMT_STR,
		name, description, fmt.Sprintf(LINE_FMT_STR, coordsList))
}

func KmlPoint(name string, description string, coord Coord) string {
	return fmt.Sprintf(PLACEMARK_FMT_STR,
		name, description, fmt.Sprintf(POINT_FMT_STR, coordToStr(&coord)))
}

func GenRoute(name string, route []Coord, elemName string, description string) string {
	content := KmlLineString(elemName, description, route)
	content += KmlPoint("Begin", description, route[0])
	content += KmlPoint("End", description, route[len(route)-1])
	content = KmlAssemble(name, content)
	return content
}
