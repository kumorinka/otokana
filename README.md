# OtoKana

A desktop app for practicing hiragana and katakana: type the romaji of each kana that shows up. With a UTAU voicebank, each kana is also sung.

Built in Go with [mygo](https://github.com/egoist/mygo)'s native UI. It runs on Windows, macOS and Linux, needs no cgo, and keeps everything on your device.

## Features

- Practice hiragana, katakana or both, with basic kana, dakuten and yōon.
- Click a card to see its romaji as a hint. The score counts the cards you answered without one.
- Import a UTAU voicebank from a folder or a zip file to hear each kana. Set its base note to hear cards a few semitones around it, with the note's name on the card.
- Choose when sounds play: when you answer right, when a card shows, or both.

## Run

You need Go 1.27 or newer.

```sh
go tool mygo dev     # run in development
go tool mygo build   # build the app into build/
```

## Voicebanks

No voicebank or audio comes with the app. Use one you have the right to use. Practice works without one, just silently.

## Fonts

Klee One and Zen Kaku Gothic New, under the SIL Open Font License: see `assets/fonts/`.
