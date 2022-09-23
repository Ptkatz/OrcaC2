package common

import (
	"image/color"
	"math"
)

/*
 	Black 	#000000 	(0,0,0)
  	White 	#FFFFFF 	(255,255,255)
  	Red 	#FF0000 	(255,0,0)
  	Lime 	#00FF00 	(0,255,0)
  	Blue 	#0000FF 	(0,0,255)
  	Yellow 	#FFFF00 	(255,255,0)
  	Cyan 	#00FFFF 	(0,255,255)
  	Magenta	#FF00FF 	(255,0,255)
  	Silver 	#C0C0C0 	(192,192,192)
  	Gray 	#808080 	(128,128,128)
  	Maroon 	#800000 	(128,0,0)
  	Olive 	#808000 	(128,128,0)
  	Green 	#008000 	(0,128,0)
  	Purple 	#800080 	(128,0,128)
  	Teal 	#008080 	(0,128,128)
  	Navy 	#000080 	(0,0,128)
*/

var Black = color.RGBA{0, 0, 0, 0}
var White = color.RGBA{255, 255, 255, 0}
var Red = color.RGBA{255, 0, 0, 0}
var Lime = color.RGBA{0, 255, 0, 0}
var Blue = color.RGBA{0, 0, 255, 0}
var Yellow = color.RGBA{255, 255, 0, 0}
var Cyan = color.RGBA{0, 255, 255, 0}
var Magenta = color.RGBA{255, 0, 255, 0}
var Silver = color.RGBA{192, 192, 192, 0}
var Gray = color.RGBA{128, 128, 128, 0}
var Maroon = color.RGBA{128, 0, 0, 0}
var Olive = color.RGBA{128, 128, 0, 0}
var Green = color.RGBA{0, 128, 0, 0}
var Purple = color.RGBA{128, 0, 128, 0}
var Teal = color.RGBA{0, 128, 128, 0}
var Navy = color.RGBA{0, 0, 128, 0}

func ColorDistance(c1 color.RGBA, c2 color.RGBA) float64 {
	return math.Sqrt((float64(c1.R)-float64(c2.R))*(float64(c1.R)-float64(c2.R)) +
		(float64(c1.G)-float64(c2.G))*(float64(c1.G)-float64(c2.G)) +
		(float64(c1.B)-float64(c2.B))*(float64(c1.B)-float64(c2.B)))
}
