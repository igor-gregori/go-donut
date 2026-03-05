package main

import (
	"fmt"
	"math"
	"time"
)

const screen_width float64 = 40
const screen_height float64 = 60

const theta_spacing float64 = 0.07
const phi_spacing float64 = 0.02

const R1 float64 = 1
const R2 float64 = 2
const K2 float64 = 5
const K1 float64 = screen_width * K2 * 3 / (8 * (R1 + R2))

func render(A float64, B float64) {
	var cosA float64 = math.Cos(A)
	var sinA float64 = math.Sin(A)
	var cosB float64 = math.Cos(B)
	var sinB float64 = math.Sin(B)

	output := make([][]byte, int64(screen_width))
	for i := range output {
		output[i] = make([]byte, int64(screen_height))
		for j := range output[i] {
			output[i][j] = ' '
		}
	}

	zbuffer := make([][]float64, int64(screen_width))
	for i := range zbuffer {
		zbuffer[i] = make([]float64, int64(screen_height))
		for j := range zbuffer[i] {
			zbuffer[i][j] = 0
		}
	}

	for theta := 0.0; theta < 2*math.Pi; theta += theta_spacing {
		var costheta float64 = math.Cos(theta)
		var sintheta float64 = math.Sin(theta)

		for phi := 0.0; phi < 2*math.Pi; phi += phi_spacing {
			var cosphi = math.Cos(phi)
			var sinphi = math.Sin(phi)

			var circlex float64 = R2 + R1*costheta
			var circley float64 = R1 * sintheta

			// With the math above, we have our final 3D (x, y, z) coordinates after rotations
			var x float64 = circlex*(cosB*cosphi+sinA*sinB*sinphi) - circley*cosA*sinB
			var y float64 = circlex*(sinB*cosphi-sinA*cosB*sinphi) + circley*cosA*cosB
			var z float64 = K2 + cosA*circlex*sinphi + circley*sinA
			var ooz float64 = 1 / z

			var xp int64 = int64(screen_width/2 + K1*ooz*x)
			var yp int64 = int64(screen_height/2 - K1*ooz*y)
			if yp < 0 {
				yp = 0
			}

			// L is luminance
			var L float64 = cosphi*costheta*sinB - cosA*costheta*sinphi - sinA*sintheta + cosB*(cosA*sintheta-costheta*sinA*sinphi)

			// If L is less than zero, we dont need plot something
			if L > 0 {
				if ooz > zbuffer[xp][yp] {
					zbuffer[xp][yp] = ooz
					var luminance_index int64 = int64(L * 8)
					output[xp][yp] = ".,-~:;=!*#$@"[luminance_index]
				}
			}
		}
	}

	fmt.Printf("\x1b[H")

	for i := 0; i < int(screen_width); i++ {
		for j := 0; j < int(screen_height); j++ {
			fmt.Printf("%c", output[i][j])
		}
		fmt.Println("")
	}
}

func main() {
	var A float64 = 1
	var B float64 = 1
	for {
		render(A, B)

		time.Sleep(100 * time.Millisecond)

		A += 0.07
		B += 0.03
	}
}
