package main

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
	"unicode"

	sgn "github.com/EgeBalci/sgn/pkg"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"

	"syscall/js"
)

// Verbose output mode
var Verbose bool
var spinr = spinner.New(spinner.CharSets[9], 50*time.Millisecond)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	printBanner()
}

func sgnFunc(this js.Value, args []js.Value) interface{} {
	arch := args[0].Int()
	encCount := args[1].Int()
	obsLevel := args[2].Int()
	encDecoder := args[3].Bool()
	asciPayload := args[4].Bool()
	saveRegisters := args[5].Bool()
	badChars := args[6].String()
	input := args[7].String()
	return js.ValueOf(sgnExec(arch, encCount, obsLevel, encDecoder, asciPayload, saveRegisters, badChars, input))
}

func sgnExec(arch, encCount, obsLevel int, encDecoder, asciPayload, saveRegisters bool, badChars, input string) map[string]interface{} {
	var res = map[string]interface{}{
		"err": nil,
		"result": nil,
	}
	source, err := hex.DecodeString(strings.ReplaceAll(input, `\x`, ""))
	if err != nil {
		res["err"] = err
		return res
	}
	payload := []byte{}
	encoder := sgn.NewEncoder()
	encoder.ObfuscationLimit = obsLevel
	encoder.PlainDecoder = encDecoder
	encoder.EncodingCount = encCount
	encoder.SaveRegisters = saveRegisters
	eror(encoder.SetArchitecture(arch))

	if badChars != "" || asciPayload {
		badBytes, err := hex.DecodeString(strings.ReplaceAll(badChars, `\x`, ""))
		eror(err)

		for {
			p, err := encode(encoder, source)
			eror(err)

			if (asciPayload && isASCIIPrintable(string(p))) || (len(badBytes) > 0 && !containsBytes(p, badBytes)) {
				payload = p
				break
			}
			encoder.Seed = (encoder.Seed + 1) % 255
		}
	} else {
		payload, err = encode(encoder, source)
		eror(err)
	}
	res["result"] = hex.EncodeToString(payload)
	return res
}

func main() {
	done := make(chan int, 0)
	js.Global().Set("sgnFunc", js.FuncOf(sgnFunc))
	<-done
}

// Encode function is the primary encode method for SGN
func encode(encoder *sgn.Encoder, payload []byte) ([]byte, error) {
	var final []byte

	if encoder.SaveRegisters {
		printLog("Adding safe register suffix...")
		final = append(final, sgn.SafeRegisterSuffix[encoder.GetArchitecture()]...)
	}

	// Add garbage instrctions before the ciphered decoder stub
	garbage, err := encoder.GenerateGarbageInstructions()
	if err != nil {
		return nil, err
	}
	payload = append(garbage, payload...)
	encoder.ObfuscationLimit -= len(garbage)

	printLog("Ciphering payload...")
	ciperedPayload := sgn.CipherADFL(payload, encoder.Seed)
	decoderAssembly := encoder.NewDecoderAssembly(len(ciperedPayload))
	printLog("Selected decoder: %s", decoderAssembly))
	decoder, ok := encoder.Assemble(decoderAssembly)
	if !ok {
		return nil, errors.New("decoder assembly failed")
	}

	encodedPayload := append(decoder, ciperedPayload...)
	if encoder.PlainDecoder {
		final = encodedPayload
	} else {
		schemaSize := ((len(encodedPayload) - len(ciperedPayload)) / (encoder.GetArchitecture() / 8)) + 1
		randomSchema := encoder.NewCipherSchema(schemaSize)
		printLog("Cipher schema: %s", sgn.GetSchemaTable(randomSchema)))
		obfuscatedEncodedPayload := encoder.SchemaCipher(encodedPayload, 0, randomSchema)
		final, err = encoder.AddSchemaDecoder(obfuscatedEncodedPayload, randomSchema)
		if err != nil {
			return nil, err
		}

	}

	if encoder.SaveRegisters {
		printLog("Adding safe register prefix...")
		final = append(sgn.SafeRegisterPrefix[encoder.GetArchitecture()], final...)
	}

	if encoder.EncodingCount > 1 {
		encoder.EncodingCount--
		encoder.Seed = sgn.GetRandomByte()
		final, err = encode(encoder, final)
		if err != nil {
			return nil, err
		}
	}

	return final, nil
}

// checks if a byte array contains any element of another byte array
func containsBytes(data, any []byte) bool {
	for _, b := range any {
		if bytes.Contains(data, []byte{b}) {
			return true
		}
	}
	return false
}

// checks if s is ascii and printable, aka doesn't include tab, backspace, etc.
func isASCIIPrintable(s string) bool {
	for _, r := range s {
		if r > unicode.MaxASCII || !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}

func eror(err error) {
	if err != nil {
		pc, _, _, ok := runtime.Caller(1)
		details := runtime.FuncForPC(pc)
		if ok && details != nil {
			printLog("[%s] ERROR: %s\n", strings.ToUpper(strings.Split(details.Name(), ".")[1]), err)
		} else {
			printLog("[UNKNOWN] ERROR: %s\n", err)
		}
		os.Exit(1)
	}
}

func printLog(format string, a ...interface{}) {
	js.Global().Get("console").Call("log", fmt.Sprintf(format, a...))
}

func printBanner() {
	banner, _ := base64.StdEncoding.DecodeString("ICAgICAgIF9fICAgXyBfXyAgICAgICAgX18gICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgXyAKICBfX18gLyAvICAoXykgL19fX19fIF8vIC9fX19fIF8gIF9fXyBfX19fIF8gIF9fXyAgX19fIF8oXykKIChfLTwvIF8gXC8gLyAgJ18vIF8gYC8gX18vIF8gYC8gLyBfIGAvIF8gYC8gLyBfIFwvIF8gYC8gLyAKL19fXy9fLy9fL18vXy9cX1xcXyxfL1xfXy9cXyxfLyAgXF8sIC9cXyxfLyAvXy8vXy9cXyxfL18vICAKPT09PT09PT1bQXV0aG9yOi1FZ2UtQmFsY8SxLV09PT09L19fXy89PT09PT09djIuMC4wPT09PT09PT09ICAKICAgIOKUu+KUgeKUuyDvuLXjg70oYNCUwrQp776J77i1IOKUu+KUgeKUuyAgICAgICAgICAgKOODjiDjgpzQlOOCnCnjg44g77i1IOS7leaWueOBjOOBquOBhAo=")
	printLog(string(banner))
}
