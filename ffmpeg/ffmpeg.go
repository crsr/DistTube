package ffmpeg

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/pwed/DistTube/ffprobe"
)

//Encodes a single file to a given resolution and output file
func Encode(ffmpegLocation string, input string, resolution string, output string) error {
	out, err := exec.Command(ffmpegLocation, "-y", "-i", input, "-s:v",
		resolution, output).Output()
	if err != nil {
		return err
	}
	fmt.Printf("%s", out)
	return nil
}

func BatchEncode(ffmpegLocation string, input string, outputs ...string) error {
	if len(outputs)%2 != 0 {
		return errors.New("ffmpeg: number of arguments was incorrect")
	}

	args := []string{"-y", "-i", input}
	for n := 0; n < len(outputs); n += 2 {
		args = append(args, "-s:v", outputs[n], outputs[n+1])
	}

	out, err := exec.Command(ffmpegLocation, args...).Output()
	if err != nil {
		return err
	}

	fmt.Printf("%s", out)
	return nil
}

func Ingest(f string, i string, s string, o string) {
	err := Encode(f, i, o+".mp4", s)
	if err != nil {
		log.Fatal(err)
	}
	os.Remove(i)
}

func BulkBatchIngest(f string, i string, o string) {
	BatchEncode(f, i,
		"480x234", o+"_234p.mp4",
		"640x360", o+"_360p.mp4",
		"1280x720", o+"_720p.mp4",
		"1920x1080", o+"_1080p.mp4",
		"2560x1440", o+"_1440p.mp4",
		"3840x2160", o+"_2160p.mp4")
	os.Remove(i)
}

func GetResolution(i string) (uint64, uint64) {
	p, err := ffprobe.Probe(i)
	if err != nil {

	}
	w := p.Width()
	h := p.Height()

	return w, h
}

func SequentialBatchIngest(f string, i string, o string) {
	Encode(f, i, "480x234", o+"_234p.mp4")
	Encode(f, i, "640x360", o+"_360p.mp4")
	Encode(f, i, "1280x720", o+"_720p.mp4")
	Encode(f, i, "1920x1080", o+"_1080p.mp4")
	Encode(f, i, "2560x1440", o+"_1440p.mp4")
	Encode(f, i, "3840x2160", o+"_2160p.mp4")
	os.Remove(i)
}
