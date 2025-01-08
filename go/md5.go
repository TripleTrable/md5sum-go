package main

import (
	"encoding/binary"
	"fmt"
	"math/bits"
)

func tableK (index int) uint32 {
	var k = [64]uint32{
        0xd76aa478, 0xe8c7b756, 0x242070db, 0xc1bdceee ,
        0xf57c0faf, 0x4787c62a, 0xa8304613, 0xfd469501 ,
        0x698098d8, 0x8b44f7af, 0xffff5bb1, 0x895cd7be ,
        0x6b901122, 0xfd987193, 0xa679438e, 0x49b40821 ,
        0xf61e2562, 0xc040b340, 0x265e5a51, 0xe9b6c7aa ,
        0xd62f105d, 0x02441453, 0xd8a1e681, 0xe7d3fbc8 ,
        0x21e1cde6, 0xc33707d6, 0xf4d50d87, 0x455a14ed ,
        0xa9e3e905, 0xfcefa3f8, 0x676f02d9, 0x8d2a4c8a ,
        0xfffa3942, 0x8771f681, 0x6d9d6122, 0xfde5380c ,
        0xa4beea44, 0x4bdecfa9, 0xf6bb4b60, 0xbebfbc70 ,
        0x289b7ec6, 0xeaa127fa, 0xd4ef3085, 0x04881d05 ,
        0xd9d4d039, 0xe6db99e5, 0x1fa27cf8, 0xc4ac5665 ,
        0xf4292244, 0x432aff97, 0xab9423a7, 0xfc93a039 ,
        0x655b59c3, 0x8f0ccc92, 0xffeff47d, 0x85845dd1 ,
        0x6fa87e4f, 0xfe2ce6e0, 0xa3014314, 0x4e0811a1 ,
        0xf7537e82, 0xbd3af235, 0x2ad7d2bb, 0xeb86d391 }

    return k[index]
}

func tableS (index int) uint32 {

	var s = [64]uint32{
		7, 12, 17, 22, 7, 12, 17, 22, 7, 12, 17, 22, 7, 12, 17, 22,
		5, 9, 14, 20, 5, 9, 14, 20, 5, 9, 14, 20, 5, 9, 14, 20,
		4, 11, 16, 23, 4, 11, 16, 23, 4, 11, 16, 23, 4, 11, 16, 23,
		6, 10, 15, 21, 6, 10, 15, 21, 6, 10, 15, 21, 6, 10, 15, 21,
	}

    return s[index]
}


func funcF (b,c,d uint32) uint32 { return (b & c) | (^b & d) }
func funcG (b,c,d uint32) uint32 { return (b & d) | (c & ^d) }
func funcH (b,c,d uint32) uint32 { return b ^ c ^ d }
func funcI (b,c,d uint32) uint32 { return c ^ (b | ^d) }

func funcg (i uint32) uint32 {
    if i < 16 {
        return i
    } else if i < 32 {
        return (5* i + 1) % 16
    } else if i < 48 {
        return (3* i + 5) % 16
    } else {
        return (7* i) % 16
    }
}

type md5Data struct {
    a uint32
    b uint32
    c uint32
    d uint32
}

func (data *md5Data)calc (f func(b uint32,c uint32,d uint32) uint32, messageWord uint32,loopIndex int) {
        a_ :=  f(data.b,data.c,data.d) + data.a + tableK(loopIndex) + messageWord
        a_ = data.b + bits.RotateLeft32(a_,int(tableS(loopIndex))) 
        data.a, data.b, data.c, data.d = data.d, a_, data.b, data.c
        
}


func (data *md5Data) digest(chunk []byte) {
    // Break chunk into sixteen 32-bit words
    var m [16]uint32
    for j := 0; j < 16; j++ {
        m[j] = uint32(chunk[j*4]) | uint32(chunk[j*4+1])<<8 | uint32(chunk[j*4+2])<<16 | uint32(chunk[j*4+3])<<24
    }

    originalA, originalB, originalC, originalD := data.a, data.b, data.c, data.d


    // Main loop
    for j := 0; j < 64; j++ {
        var f func(b uint32,c uint32,d uint32) uint32
        var g = funcg(uint32(j))
        if j < 16 {
            f = funcF
        } else if j < 32 {
            f = funcG
        } else if j < 48 {
            f = funcH
        } else {
            f = funcI
        }
        data.calc(f,m[g],j)
    }

    // Add this chunk's hash to result so far
    data.a += originalA
    data.b += originalB
    data.c += originalC
    data.d += originalD
}

func swapEndianess (word uint32) uint32 {
    data := make([]byte,4)
    binary.NativeEndian.PutUint32(data[:], word)
    return binary.BigEndian.Uint32(data)
}

func MD5Sum (data string) string{
    // Constants for MD5
    md5dat := md5Data{
		a: 0x67452301,
		b: 0xefcdab89,
		c: 0x98badcfe,
		d: 0x10325476,
	}

	// Padding the input
	msg := []byte(data)
	originalLen := uint64(len(msg) * 8)
	msg = append(msg, 0x80)
	for len(msg)%64 != 56 {
		msg = append(msg, 0)
	}

	// Append original length (before padding) as a 64-bit little-endian integer
    for i := uint(0); i < 8; i++ {
		msg = append(msg, byte(originalLen>>(8*i)))
	}


	// Process the message in 512-bit chunks
	for i := 0; i < len(msg); i += 64 {
		chunk := msg[i : i+64]

        md5dat.digest(chunk)
	}

	// Produce the final hash value (little-endian)
	result := fmt.Sprintf("%08x%08x%08x%08x", swapEndianess(md5dat.a), swapEndianess(md5dat.b), swapEndianess(md5dat.c), swapEndianess(md5dat.d))
	return result
}
