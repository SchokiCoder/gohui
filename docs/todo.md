# niche bugs

- courier: issues with neofetch
  it's not cuz: hui:handleShell:ReadAll, SplitByLines:trimming, drawContent
  Also neofetch does something weird with lines in general; try scrolling down

  `less` seems to filter (some) CSI Sequences out, maybe try that

# later

- add configurable padding
  (lPadding and rPadding)
	- how does padding interact with alignment?
	- Header
	- Title
	- Entries
	- Feedback
	- Cmdline
	- courier: content

- update demo config
- set version

# v1.3 niceties

- add arrow key support for navigation
- add arrow key support for cmdline
- add del key support for cmdline
- add return key support for cmdline
- pgup and pgdown support
- add specific feedback color for errors
- add feedback for when hitting "right" on a shell entry or "execute" on a menu
  entry
  "Entry type is \"menu\", cannot execute."
  "Entry type is \"shell\", cannot enter."
- add cursor for each menu in menu_path
  Thus a "left" key press would send you to the menu entry that you entered.
- set version to 1.3

# bleh

- look for original features that need to be implemented before implementing
  next

# v1.2

+ hui: add go scripting interface for entries

For that, runtime variables have been moved from local main function variables
to a struct defined in common.
I first wanted to make it global variables in hui but... eh,
plus I would have had to import hui in scripts, which means circular inclusion.

- hui: add go-entry config-values

- hui: add scripting interface for start and quit
- courier: add scripting interface for start and quit

- add config validation for if a script function actually exists in func map
- update demo configs
- set version

-----

# v1.1

+ add configurable aligns for Header and Title

Also change shell scripts to encourage devs to use the new build.sh rather than
the single-target build scripts.
This is to prevent forgetfulness about checking wether common code compiles for
all targets.

+ rework common.SplitByLines to be simpler and more reliable
  May the merciful lord allow this function to just work without ever being
  touched again.

+ courier: fix content not adjusting to term width changes

+ hui: add configurable aligns for Entries

+ fix alignment stretching background color along the entire left padding

+ add configurable aligns for Feedback and Cmdline

Also add missing cfg validity check for EntryAlignment.

+ courier: add configurable aligns for content

_Padding postponed because it would be a limitation violation...
just like alignments._ 

+ update demo config

Also set version to 1.1.

# v1.0

+ UNNECESSARY: return values
  (no shell scripting with those lol)

+ courier: add info args
+ hui: add info args

+ add demo configuration
+ hui: add config validity check to call out empty menus
+ fix my skill issue
  hui: Shell entries can't handle `cat` or `neofetch` (child returns 127)
  (cfg just didn't point towards an existing pager (forgot "./" for local test))
+ fix handle shell session: some apps messing up first draw after return

+ add help args
+ add docs/goals_omissions.md

+ fix panic message consistency
+ add install scripts
+ enable install scripts to do user local installs
+ common cfg: add validity check for if Pager can be found
+ update README.md

Also set version to 1.0.

# v0.3.0

+ add shell session execution

Aka add support for child processes with their own mainloop.

+ FAILED: maybe unify shell and shellsession
	- use handleShellSession as base
	  (we hand our own stdout to child)
	- record record our own stdout
	  (while child runs)
	- once child is done, do the normal decision making of return feedback
	  from stdout (recorded) vs stderr as in handleShell
	- could handing over our own stdin mess with piping?
	  "cat myfile | idklol"
	  maybe not since it's encapsulated by "sh -c %v"
	- remove ShellSession value
	- discard feedback that came from a session's Stdout...
	  (withoput explicit knowledge gained from the cfg, this is impossible,
	   i am afrad)
	- how do you determine if a shell command runs a binary with a mainloop?
	  You don't.
	  Execution time is not reliable.

+ seperate common code from hui
+ add courier base
+ courier: add file read from arg

Also add a missing file close when config reading.

+ courier: add scroll

Also fix last line being omitted by common.SplitByLines.

+ FAILED: try rune for keys in configs
+ fix colored prints not resetting themselves
+ courier: add scroll via cmdline number

+ courier: add optional title arg
+ hui: add termH awareness for drawMenu

+ add hui giving big feedback to courier
  try temp files first this time
	+ hui: fix not passing correct string as feedback to pager
	+ fix: add feedback reset after pager call
	+ test compat with other pager

+ remove compile flags for gdb because it's cumbersome af

Use delve from now on for debugging.
Also improve build scripts a bit.
Also set version to 0.3.0.

# v0.2.0

+ add basic toml config file reading
+ fix toml reading
+ config: change keys to be strings

Otherwise the toml umarshal would expect literal integers in the toml file.

+ add XDG config env var to config paths

+ remove EntryContent from Entry

This is then replaced with the values itself.
Sanity checks are implemented right after the config unmarshal.
Thanks to this the config doesn't need to contain arbitrary integers anymore.
Also set version to 0.2.0.

# v0.1.0

+ add mainloop
+ add header
+ add menu
+ add title draw
+ add menu draw
+ add raw terminal mode
	+ fix stdin read and permanent redraw
	+ fix draw magic tab characters

+ can EntryContent be implemented as empty interface and used via RTTI?
  Yes but it seems more prone to failure due to requiring/having a default case.
```Go
type EntryContentMenu string
type EntryContentShell string

type Entry struct {
	caption string
	content interface{}
}



switch cur_menu.entries[i].content.(type) {
		case EntryContentMenu:
			...

		default:
			panic("unknown entry content type")
		
```

+ add keyboard input and close via ctl + 'c' and 'd'
+ add menu entry cursor (up and down)
+ add menu navigation (left, right)
+ add basic shell command execution
+ add feedback line
+ add command line enter and quit command
+ add command line display
+ fix command line not getting emptied
+ add command line leave via SIGINT and SIGTSTP
+ add command line number parsing
+ fix successful command not clearing feedback
+ add cursor hide/show
+ add config sys for text fore- and background
  
Also fix default-color-sequences

+ add config values for key binds
+ add feedback trim
+ add print prevention for feedback longer than 1 line

Also add configurable cmdline/feedback prefix
to prevent a temporary hack solution.
We need the prefix for detecting needed lines
for feedback print.

+ set version to 0.1.0
