# LC-3 VM

![LC-3 VM by Rob Morgan](docs/lc3-vm.png)

This project is a really simple [LC-3](https://en.wikipedia.org/wiki/LC-3) VM written in Go. LC-3 or Little Computer 3 is
a fictional computer system that is designed to teach students how to code in assembly language. My VM is originally [based
on an article](https://justinmeiners.github.io/lc3-vm/) written by Justin Meiners and Ryan Pendleton.

## Current Status

The VM is now able to execute the following programs:

- [2048](https://github.com/rpendleton/lc3-2048) by Ryan Pendleton
- [Rogue](https://github.com/justinmeiners/lc3-rogue) by Justin Meiners

## TODO

- [ ] Fix 199% CPU issue when running Rogue
- [ ] Fix failing unit tests

## Changelog

- Fixed Trap Routines for displaying output.
- Fixed the STI Op Code.
- Migrated to Termbox for display and key input
- Added a flag for verbose debug output
- Added a flag to specify the target program
- Added support for CPU profiling
- Fixed an issue where GETC was not waiting for input

## Development Resources

These resources came in handy when developing this VM:

- https://wchargin.github.io/lc3web/: I could load the programs and compare the expected behaviour against my VM.

## LICENSE

MIT, See the LICENSE file.
