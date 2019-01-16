package main

import (
	"io/ioutil"
	"log"
	"time"

	resources "github.com/kyeett/platformer/assets"

	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/mp3"
	"github.com/hajimehoshi/ebiten/audio/wav"
)

const (
	sampleRate = 44100
)

var jumpSound []byte
var bounceSound []byte

// Player represents the current audio state.
type Player struct {
	audioContext *audio.Context
	audioPlayer  *audio.Player
	current      time.Duration
	total        time.Duration
	seBytes      []byte
	seCh         chan []byte
	volume128    int
}

var globalAudioContext *audio.Context

func NewPlayer() (*Player, error) {

	// return &Player{}, nil

	var err error
	globalAudioContext, err = audio.NewContext(sampleRate)
	if err != nil {
		log.Fatal(err)
	}

	// f, err := ioutil.ReadFile(resources.Lookup["assets/audio/getting-it-done-medium.mp3"])
	// if err != nil {
	// log.Fatal(err)
	// }

	const bytesPerSample = 4 // TODO: This should be defined in audio package
	path := "assets/audio/getting-it-done-medium.mp3"
	s, err := mp3.Decode(globalAudioContext, audio.BytesReadSeekCloser(resources.Lookup[path]))
	if err != nil {
		return nil, err
	}
	p, err := audio.NewPlayer(globalAudioContext, s)
	if err != nil {
		return nil, err
	}
	player := &Player{
		audioContext: globalAudioContext,
		audioPlayer:  p,
		total:        time.Second * time.Duration(s.Length()) / bytesPerSample / sampleRate,
		volume128:    128,
		seCh:         make(chan []byte),
	}
	if player.total == 0 {
		player.total = 1
	}

	path = "assets/audio/jump_01.wav"
	s2, err := wav.Decode(globalAudioContext, audio.BytesReadSeekCloser(resources.Lookup[path]))
	if err != nil {
		log.Fatal("failed to load", err)
	}
	b, err := ioutil.ReadAll(s2)
	if err != nil {
		log.Fatal("failed to read", err)
	}
	jumpSound = b

	path = "assets/audio/SFX_Jump_07.wav"
	s2, err = wav.Decode(globalAudioContext, audio.BytesReadSeekCloser(resources.Lookup[path]))
	if err != nil {
		log.Fatal("failed to load", err)
	}
	b, err = ioutil.ReadAll(s2)
	if err != nil {
		log.Fatal("failed to read", err)
	}
	bounceSound = b

	player.audioPlayer.Play()

	return player, nil
}

func (p *Player) PlayAudio(b []byte) {
	if b == nil {
		log.Println("tried to play empty bytes")
		return
	}
	tmpP, err := audio.NewPlayerFromBytes(globalAudioContext, b)
	if err != nil {
		log.Fatal(err)
	}
	tmpP.Play()
}
