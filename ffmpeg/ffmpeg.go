package ffmpeg

import (
	"fmt"
	"os/exec"
)

func Encode(ffmpegLocation string, input string, output string, resolution string) {
	out, err := exec.Command(ffmpegLocation, "-y", "-i", input, "-s:v",
		resolution, output).Output()
	if err != nil {
		fmt.Println("error occured")
		fmt.Printf("%s", err)
	}
	fmt.Printf("%s", out)

}
