// Credit to https://github.com/BenLubar for original code :)
package ffprobe


import (
	"time"
	"os/exec"
	"encoding/json"
)

//Struct containing metadata about a video file
type ProbeFormat struct {
	Filename         string            `json:"filename"`
	NBStreams        int               `json:"nb_streams"`
	NBPrograms       int               `json:"nb_programs"`
	FormatName       string            `json:"format_name"`
	FormatLongName   string            `json:"format_long_name"`
	StartTimeSeconds float64           `json:"start_time,string"`
	DurationSeconds  float64           `json:"duration,string"`
	Size             uint64            `json:"size,string"`
	BitRate          uint64            `json:"bit_rate,string"`
	ProbeScore       float64           `json:"probe_score"`
	Tags             map[string]string `json:"tags"`
}

func (f ProbeFormat) StartTime() time.Duration {
	return time.Duration(f.StartTimeSeconds * float64(time.Second))
}


func (f ProbeFormat) Duration() time.Duration {
	return time.Duration(f.DurationSeconds * float64(time.Second))
}

//Struct containing stream info and metadata
type ProbeData struct {
	Streams []Stream   `json:"streams,omitempty"`
	Format ProbeFormat `json:"format,omitempty"`
}

//Contains the height and width of a video in pixels
type Stream struct {
	Width uint64 `json:"width"`
	Height uint64 `json:"height"`
}

//returns the video height from a ProbeData object
func (f ProbeData) Height() uint64 {
	return f.Streams[0].Height
}

//returns the video height from a ProbeData object
func (f ProbeData) Width() uint64 {
	return f.Streams[0].Width
}

//Returns the ProbeData of a video file located at the location of the filename
func Probe(filename string) (*ProbeData, error) {
	cmd := exec.Command("ffprobe", "-show_format", filename,
		"-show_entries", "stream=height,width", "-print_format", "json")
	//cmd.Stderr = os.Stderr

	r, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	var v ProbeData
	err = json.NewDecoder(r).Decode(&v)
	if err != nil {
		return nil, err
	}

	err = cmd.Wait()
	if err != nil {
		return nil, err
	}

	return &v, nil
}
