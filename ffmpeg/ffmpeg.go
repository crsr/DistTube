package ffmpeg

import (
	"fmt"
	"os/exec"
	"errors"
)

func Encode(ffmpegLocation string, input string, output string, resolution string) error {
	out, err := exec.Command(ffmpegLocation, "-y", "-i", input, "-s:v",
		resolution, output).Output()
	if err != nil {
		return err
	}
	fmt.Printf("%s", out)
	return nil
}

func BatchEncode(ffmpegLocation string, input string, outputs ...string) error {
	if len(outputs) % 2 != 0 {
		return errors.New("ffmpeg: number of arguments was incorrect")
	}

	args := []string{"-y", "-i", input}
	for n := 0; n < len(outputs); n += 2 {
		args = append(args, "-s:v")
		args = append(args, outputs[n])
		args = append(args, outputs[n+1])
		fmt.Println(n)
	}


	fmt.Println(args)

	out, err := exec.Command(ffmpegLocation, args...).Output()
	if err != nil {
		return err
	}

	fmt.Printf("%s", out)
	return nil
}
