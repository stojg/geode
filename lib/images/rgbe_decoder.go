package images

// RGBE or Radiance HDR is an image format invented by Gregory Ward Larson for the Radiance
// rendering system. It stores pixels as one byte each for RGB (red, green, and blue) values with
// a one byte shared exponent. Thus it stores four bytes per pixel.
//
// 1) Scale the three floating point values to share a common 8-bit exponent, taken from
// the brightest of the three. Each value is then truncated to an 8-bit mantissa (fractional part).
// The result is four bytes, 32 bits, for each pixel. This results in a 6:1 compression, at the
// expense of reduced colour fidelity.
//
// 2) The second stage performs run length encoding on the 32-bit pixel values. This has a limited
// impact on the size of most rendered images, but it is fast and simple.
//
// - https://en.wikipedia.org/wiki/RGBE_image_format
// - https://en.wikipedia.org/wiki/Radiance_(software)#HDR_image_format
// - https://en.wikipedia.org/wiki/Run-length_encoding
//
// http://www.graphics.cornell.edu/~bjw/rgbe
// https://github.com/Opioid/rgbe/blob/master/decode.go

import (
	"bufio"
	"fmt"
	"io"
	"math"
)

// DecodeRGBE reads from input and returns the width, height, image data in float32 and an error
func DecodeRGBE(r io.Reader) (int, int, []float32, error) {
	br := bufio.NewReader(r)

	width, height, err := readHeader(br)
	if err != nil {
		return 0, 0, nil, err
	}
	data := make([]float32, width*height*3)
	if err := readPixelsRLE(br, width, height, data); err != nil {
		return 0, 0, nil, err
	}

	return width, height, data, nil
}

func readHeader(r *bufio.Reader) (int, int, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return 0, 0, newDecoderErr(readError, "header #1 %s", err)
	}

	if line[0] != '#' || line[1] != '?' {
		return 0, 0, newDecoderErr(formatError, "bad initial token '%d'", line[1])
	}

	formatSpecifier := false

	for {
		line, err = r.ReadString('\n')
		if err != nil {
			return 0, 0, newDecoderErr(readError, "header #2 %s", err)
		}

		if line[0] == 0 || line[0] == '\n' {
			// blank lines signifies end of meta data header
			break
		} else if line == "FORMAT=32-bit_rle_rgbe\n" {
			formatSpecifier = true
		}
	}

	if !formatSpecifier {
		return 0, 0, newDecoderErr(formatError, "no FORMAT specifier found")
	}

	line, err = r.ReadString('\n')
	if err != nil {
		return 0, 0, newDecoderErr(readError, "header #3 %s", err)
	}

	width, height := 0, 0
	if n, err := fmt.Sscanf(line, "-Y %d +X %d", &height, &width); n < 2 || err != nil {
		return 0, 0, newDecoderErr(formatError, "missing image size specifier")
	}

	return width, height, nil
}

func readPixelsRLE(r io.Reader, width, height int, result []float32) error {
	if len(result) < width*height*3 {
		return newDecoderErr(memoryError, "requires %d floats but only got %d floats available", width*height*3, len(result))
	}

	// Run Length Encoding is not allowed so read flat
	if width < 8 || width > 0x7fff {
		return readPixels(r, width*height, result)
	}

	offset := 0
	rgbe := make([]byte, 4)
	scanlineBuffer := make([]byte, 4*width)
	buf := make([]byte, 2)

	for ; height > 0; height-- {
		if _, err := io.ReadFull(r, rgbe); err != nil {
			return newDecoderErr(readError, "readfull #1 %s", err)
		}

		// this line is not RLE so read the rest of the file as flat
		if rgbe[0] != 2 || rgbe[1] != 2 || (rgbe[2]&0x80) != 0 {
			result[0], result[1], result[2] = rgbeToFloat(rgbe[0], rgbe[1], rgbe[2], rgbe[3])
			return readPixels(r, width*height-1, result[3:])
		}

		// the 3rd and 4th byte should match the width
		if int(rgbe[2])<<8|int(rgbe[3]) != width {
			return newDecoderErr(formatError, "wrong scanline width, got %d, but expected %d", int(rgbe[2])<<8|int(rgbe[3]), width)
		}

		// read each of the four channels for the scanline into the buffer
		index := 0
		for i := 0; i < 4; i++ {
			end := (i + 1) * width

			for index < end {
				if _, err := io.ReadFull(r, buf); err != nil {
					return newDecoderErr(readError, "readfull #2 %s", err)
				}

				if buf[0] > 128 {
					count := int(buf[0]) - 128

					if count == 0 || count > end-index {
						return newDecoderErr(formatError, "not enough data in scanline for rle")
					}

					for ; count > 0; count-- {
						scanlineBuffer[index] = buf[1]
						index++
					}
				} else {
					// a non-run
					count := int(buf[0])

					if count == 0 || count > end-index {
						return newDecoderErr(formatError, "not enough data in scanline")
					}

					scanlineBuffer[index] = buf[1]
					index++

					count--
					if count > 0 {
						if _, err := io.ReadFull(r, scanlineBuffer[index:index+count]); err != nil {
							return newDecoderErr(readError, "readfull #3 %s", err)
						}

						index += count
					}
				}
			}
		}

		// now convert data from buffer into floats
		for i := 0; i < width; i++ {
			r := scanlineBuffer[i]
			g := scanlineBuffer[i+width]
			b := scanlineBuffer[i+2*width]
			e := scanlineBuffer[i+3*width]

			result[offset], result[offset+1], result[offset+2] = rgbeToFloat(r, g, b, e)
			offset += 3
		}
	}

	return nil
}

func readPixels(r io.Reader, numPixels int, data []float32) error {
	rgbe := make([]byte, 4)
	offset := 0

	for ; numPixels > 0; numPixels-- {
		if _, err := io.ReadFull(r, rgbe); err != nil {
			return newDecoderErr(readError, "read pixels %s", err)
		}
		data[offset], data[offset+1], data[offset+2] = rgbeToFloat(rgbe[0], rgbe[1], rgbe[2], rgbe[3])
		offset += 3
	}
	return nil
}

// standard conversion from rgbe to float pixels, ldexp(col+0.5,exp-(128+8))
// note: 128 (0x80) was chose as the offset so that:
// > In order to cover negative exponents as well, some offset should be added to the unsigned values.
// > In this case 128 was chosen, which reserves the same range for values greater than 1 and less than 1.
// > It is possible to adjust this offset value if necessary, but since 2^127 ~= 10^38 it will
// > rarely be the necessary. This file format covers about 76 orders of magnitude with 1% relative accuracy.
// > - https://www.cg.tuwien.ac.at/research/theses/matkovic/node84.html
func rgbeToFloat(r, g, b, e byte) (float32, float32, float32) {
	if e == 0 {
		return 0, 0, 0
	}

	exp := int(e) - (128 + 8)

	// if the values needs to be normalised to [0,1] range instead of being in float32
	// f := float32(math.Ldexp(1, exp))
	// return float32(r) * f, float32(g) * f, float32(b) * f
	red := math.Ldexp(float64(r)+0.5, exp)
	green := math.Ldexp(float64(g)+0.5, exp)
	blue := math.Ldexp(float64(b)+0.5, exp)
	return float32(red), float32(green), float32(blue)
}
