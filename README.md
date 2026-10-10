# YADMG - yet another DMG (Gameboy emulator) written in Go

This isn't exactly a novel project, but it's one that I wanted to build to gain a little more familiarity with how systems actually work.
Currently, this emulator passes blargg's cpu_instrs (the individual tests) and instr_timing.

### Build
```sh
go build -o yadmg ./cmd/emu/
```

### Usage
No boot ROM required, just pass in a cartridge ROM like so:
```sh
yadmg rom.gb
```
Debug/other options:
```sh
-debug
    show debug output to track PC, SP and register values
```

### Configuration
yadmg supports configuring controls via a YAML file `keymap.yml`, located either in the emulator's directory or in `$HOME/.config/yadmg`. An example configuration is included in the repository. All letter keys (make sure they're upper case in the config file) are mappable, as well as left shift, right shift, enter, and arrow keys. 

### Goals
- [x] Configurable controls with YAML config file
- [ ] Working APU
- [ ] Fix PPU rendering bugs that appear with certain games (Link's Awakening, Pac-Man known so far)

### Non-goals
- CGB support. At least for now, I'm not particularly interested in implementing all of the disparate state and branching logic that comes with it.

### Acknowledgements 
Thanks to these projects for making this much easier to develop:
- [goboy](https://github.com/Humpheh/goboy): very valuable second opinion on Go implementations
- [mooneye-gb](https://github.com/Gekkio/mooneye-gb): excellent reference for hardware-accurate implementations
- [DMG-01](https://github.com/rylev/DMG-01): easy to read and understand to get an idea of what to do 
- [gameboy-doctor](https://github.com/robert/gameboy-doctor): super useful debugging tool
