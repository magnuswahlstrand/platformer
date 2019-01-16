package main

import (
	"io/ioutil"
	"log"
	"time"

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
	var err error
	globalAudioContext, err = audio.NewContext(sampleRate)
	if err != nil {
		log.Fatal(err)
	}

	f, err := ioutil.ReadFile("assets/audio/getting-it-done-medium.mp3")
	if err != nil {
		log.Fatal(err)
	}

	const bytesPerSample = 4 // TODO: This should be defined in audio package
	s, err := mp3.Decode(globalAudioContext, audio.BytesReadSeekCloser(f))
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

	// Jumping sound
	bytes, err := ioutil.ReadFile("assets/audio/jump_01.wav")
	if err != nil {
		log.Fatal(err)
	}

	s2, err := wav.Decode(globalAudioContext, audio.BytesReadSeekCloser(bytes))
	if err != nil {
		log.Fatal(err)
	}
	b, err := ioutil.ReadAll(s2)
	if err != nil {
		log.Fatal(err)
	}
	jumpSound = b

	// Jumping sound
	bytes, err = ioutil.ReadFile("assets/audio/SFX_Jump_07.wav")
	if err != nil {
		log.Fatal(err)
	}

	s2, err = wav.Decode(globalAudioContext, audio.BytesReadSeekCloser(bytes))
	if err != nil {
		log.Fatal(err)
	}
	b, err = ioutil.ReadAll(s2)
	if err != nil {
		log.Fatal(err)
	}

	bounceSound = b

	player.audioPlayer.Play()

	return player, nil
}

func playAudio(b []byte) {
	if b == nil {
		log.Fatal("noo?")
	}
	p, err := audio.NewPlayerFromBytes(globalAudioContext, b)
	if err != nil {
		log.Fatal(err)
	}
	p.Play()
}
