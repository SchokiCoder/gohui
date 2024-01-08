# v1.0.0

- the license into binary at compile time thing can be easily done with include
  macro?
  Do i even still need that with GPL2?

- mangen?
- config file manpage
- hui manpage
- courier manpage
- POSIX call options
- return values
- print messages (consistency, version information, license information)
- generalist standard configuration which says so itself via main menu title
- install scripts
- update README.md#Install

# v0.3.0

- child process execution test vim (maybe it already works)
- add courier
- give hui multiline feedback to courier

# v0.2.0

- find config file format
- add config path priority
  "/etc" over "~/.config/$FILE" over "~/.$FILE" over "$CWD/$FILE"
  if no config found, panic

- use new config structs
- read config file
- remove source code config

# v0.1.0

+ add mainloop
+ add header
+ add menu
+ add title draw
- add menu draw
- add keyboard input and close via ctl + 'c' and 'q'
- add cursor hide/show
- add terminal clear
- add raw terminal mode
- add menu entry cursor (up and down)
- add basic shell command execution
- add feedback line
- add feedback color
- add menu navigation (left, right)
- add entry prefix and postfix for each entry type
- add command line enter via ':'
- add command line leave via ctl + 'c'
- add command line typing and display
- add command interpretation via enter
- add config sys for text fore- and background
- add config values for key binds
- set version to 0.1.0
