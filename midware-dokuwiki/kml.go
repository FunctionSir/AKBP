/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2024-09-22 14:49:31
 * @LastEditTime: 2024-09-24 22:49:43
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
	Lon string
	Lat string
	Asl string
}

const (
	KML_FMT_STR string = `<?xml version="1.0" encoding="UTF-8"?>` +
		`<kml xmlns="http://www.opengis.net/kml/2.2" xmlns:gx="http://www.google.com/kml/ext/2.2" ` +
		`xmlns:kml="http://www.opengis.net/kml/2.2" xmlns:atom="http://www.w3.org/2005/Atom">` +
		`<Document><name>%s</name>%s</Document></kml>` // Gened by Google Earth Pro. %s: Name, %s: Content.
	PLACEMARK_FMT_STR string = "<Placemark><name>%s</name>%s</Placemark>"               // %s: Name, %s: Content.
	LINE_FMT_STR      string = "<LineString><coordinates>%s</coordinates></LineString>" // %s: CoordinatesList.
)

func mkCoordsListStr(coords []Coord) string {
	r := ""
	for _, x := range coords {
		r += x.Lon + "," + x.Lat + "," + x.Asl + " "
	}
	return strings.TrimSpace(r)
}

func KmlAssemble(name string, content string) string {
	return fmt.Sprintf(KML_FMT_STR, name, content)
}

func KmlLineString(name string, coords []Coord) string {
	coordsList := mkCoordsListStr(coords)
	return fmt.Sprintf(PLACEMARK_FMT_STR,
		name, fmt.Sprintf(LINE_FMT_STR, coordsList))
}
