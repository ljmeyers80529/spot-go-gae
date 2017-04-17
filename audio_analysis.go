package spotify

import (
	"encoding/json"
	"net/http"
)

// AudioAnalysis contains audio information and metadata for the specified track
type AudioAnalysis struct {
	Bars      []BeatBar `json:"bars"`
	Beats     []BeatBar `json:"beats"`
	Meta      Meta      `json:"meta"`
	Sections  []Section `json:"sections"`
	Segments  []Segment `json:"segments"`
	Tatums    []Tatum   `json:"tatums"`
	TrackInfo TrackInfo `json:"track"`
}

// BeatBar contains position and duration information for a given bar or beat in a track
type BeatBar struct {
	Start      float64 `json:"start"`
	Duration   float64 `json:"duration"`
	Confidence float64 `json:"confidence"`
}

// Meta contains metadata on the analysis process for the track
type Meta struct {
	Analyzer     string  `json:"analyzer_version"`
	Platform     string  `json:"platform"`
	Status       string  `json:"detailed_status"`
	StatusCode   int     `json:"status_code"`
	Timestamp    int     `json:"timestamp"`
	AnalysisTime float64 `json:"analysis_time"`
	InputProcess string  `json:"input_process"`
}

// Section contains audio information for specific sections of a track defined by large
// variations in rhythm or timbre
type Section struct {
	Start             float64 `json:"start"`
	Duration          float64 `json:"duration"`
	Confidence        float64 `json:"confidence"`
	Loudness          float64 `json:"loudness"`
	Tempo             float64 `json:"tempo"`
	TempoConfidence   float64 `json:"tempo_confidence"`
	Key               int     `json:"key"`
	KeyConfidence     float64 `json:"key_confidence"`
	Mode              int     `json:"mode"`
	ModeConfidence    float64 `json:"mode_confidence"`
	TimeSignature     int     `json:"time_signature"`
	TimeSigConfidence float64 `json:"time_signature_confidence"`
}

// Segment contains audio information for specific segments of a track defined by their
// relatively uniform duration, loudness, pitch, and timbre.
type Segment struct {
	Start           float64   `json:"start"`
	Duration        float64   `json:"duration"`
	Confidence      float64   `json:"confidence"`
	LoudnessStart   float64   `json:"loudness_start"`
	LoudnessMaxTime float64   `json:"loudness_max_time"`
	LoudnessMax     float64   `json:"loudness_max"`
	LoudnessEnd     float64   `json:"loudness_end"`
	Pitches         []float64 `json:"pitches"`
	Timbre          []float64 `json:"timbre"`
}

// Tatum contains position and duration information on the lowest regular pulse that
// a listener will intuitively infer from the timing of musical events in a track
type Tatum struct {
	Start      float64 `json:"start"`
	Duration   float64 `json:"duration"`
	Confidence float64 `json:"confidence"`
}

// TrackInfo contains metadata on the entire track and its analysis along with additional
// information that can be used for rhythm matching and sychronization
type TrackInfo struct {
	NumSamples         int     `json:"num_samples,omitempty"`
	Duration           float64 `json:"duration"`
	SampleMD5          string  `json:"sample_md5,omitempty"`
	OffsetSeconds      float64 `json:"offset_seconds"`
	WindowSeconds      float64 `json:"window_seconds"`
	AnalysisSampleRate int     `json:"analysis_sample_rate"`
	AnalysisChannels   int     `json:"analysis_channels"`
	EndFadeIn          float64 `json:"end_of_fade_in"`
	StartFadeOut       float64 `json:"start_of_fade_out"`
	Loudness           float64 `json:"loudness"`
	Tempo              float64 `json:"tempo"`
	TempoConfidence    float64 `json:"tempo_confidence"`
	TimeSignature      int     `json:"time_signature"`
	TimeSigConfidence  float64 `json:"time_signature_confidence"`
	Key                int     `json:"key"`
	KeyConfidence      float64 `json:"key_confidence"`
	Mode               int     `json:"mode"`
	ModeConfidence     float64 `json:"mode_confidence"`
	Codestring         string  `json:"codestring"`
	CodeVersion        float64 `json:"code_version"`
	EchoPrintString    string  `json:"echoprintstring"`
	EchoPrintVersion   float64 `json:"echoprint_version"`
	SynchString        string  `json:"synchstring"`
	SynchVersion       float64 `json:"synch_version"`
	RhythmString       string  `json:"rhythmstring"`
	RhythmVersion      float64 `json:"rhythm_version"`
}

// GetAudioAnalysis takes a track ID and returns the audio analysis information for
// the associated track including loudness, tempo, key, pitch, and timbre for denoted
// sections of the track. For a full outline of the output, see: https://developer.spotify.com/web-api/get-audio-analysis/
func (c *Client) GetAudioAnalysis(id ID) (*AudioAnalysis, error) {
	spotifyURL := baseAddress + "audio-analysis/" + id.String()
	resp, err := c.http.Get(spotifyURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, decodeError(resp.Body)
	}

	var a AudioAnalysis
	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}
