package main

import (
				"os"
				"os/exec"
				"log"
				"strings"
				"time"
			 )

type StringArray []string
type ResolutionData map[string]StringArray

func XrandrRead() (ResolutionData, error) {
	out, err := exec.Command("xrandr").Output()
	if err != nil {
		return nil, err
	}
	sOut := strings.TrimSpace(string(out))
	lines := strings.Split(sOut, "\n")
	var lastName string = ""
	var values = make(ResolutionData)
	for _, line := range lines {
		connIndex := strings.Index(line, " connected")
		if connIndex != -1 {
			lastName = line[:connIndex]
			values[lastName] = []string{}
		} else if line[:3] == "   " {
			endOfNewRes := strings.Index(line[3:], " ")
			newRes := line[3:3+endOfNewRes]
			values[lastName] = append(values[lastName], newRes)
		}
	}
	return values, nil
}

func FindCommonResolution(resolutions ResolutionData, devices []string) (string, bool) {
	baseDevice := devices[0]
	log.Println("selecting ", baseDevice, " as base device")
	i := 0
	iChances := len(resolutions[baseDevice])
	for i < iChances {
		log.Println("trying to match resolution ", resolutions[baseDevice][i])
		goodOne := true
		j := 1
		jChances := len(devices)
		for goodOne && j < jChances {
			device := devices[j]
			log.Println("let's see if it works for device ", device)
			exists := false
			k := 0
			kChances := len(resolutions[device])
			for !exists && k < kChances {

				exists = resolutions[device][k] == resolutions[baseDevice][i]
				k++
			}
			goodOne = exists
			j++
		}
		if goodOne {
			return resolutions[devices[0]][i], true
		}
		i++
	}
	return "", false
}

func main() {
	res, err := XrandrRead()
	if err != nil {
		log.Fatal(err)
	}
	var allDevices []string
	allDevices = make([]string, len(res))
	i := 0
	for key := range res {
		allDevices[i] = key
		i++
	}
	var devices []string
	if len(os.Args) < 2 {
		log.Println("devices not specified via commandline, defaulting to all existing devices: ", allDevices)
		devices = allDevices
	} else {
		devices = os.Args[1:]
	}
	resolution, ok := FindCommonResolution(res, devices)
	if !ok {
		log.Fatal("common resolution not found for ", devices, " out of ", res)
	}
	log.Println("looks like", resolution, "is a way to go")
	setCmd := []string{"xrandr"}
	modeMap := make(map[string]bool)
	for _, device := range allDevices {
		modeMap[device] = false
	}
	for _, device := range devices {
		modeMap[device] = true
	}
	for _, device := range allDevices {
		setCmd = append(setCmd, "--output")
		setCmd = append(setCmd, device)
		if modeMap[device] {
			setCmd = append(setCmd, "--mode")
			setCmd = append(setCmd, resolution)
		} else {
			setCmd = append(setCmd, "--off")
		}
	}
	log.Println(setCmd)
	log.Println("applying in 3 seconds if not opposed")
	time.Sleep(3*time.Second)
	err = exec.Command("xrandr", setCmd[1:]...).Run()
	if err != nil {
		log.Fatal(err)
	}
}
