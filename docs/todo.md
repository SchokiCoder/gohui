# niche bugs

- courier: issues with neofetch
  it's not cuz: hui:handleShell:ReadAll, SplitByLines:trimming, drawContent
  Also neofetch does something weird with lines in general; try scrolling down

  `less` seems to filter (some) CSI Sequences out, maybe try that

# later

- add csi 4 bit colors?
- add specific feedback color for errors

- add numerical modificator for key commands?
  ("2j" goes down twice)

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

# The Idiomatic Update ?

+ make Feedback its own type

Why Feedback?
Because it's often returned by functions.
Why at all?
The book told me to do so.

+ move cmdline data to common into its own struct

This simplifies things a bit.

- make runtime changing functions more idiomatic
  (pass by value, methods?)

- is App a better name for (hui/cou)Runtime ?

- shadowing via `:=` is a thing
  Ban?

- drop build via Shell for Makefile ?

- finish the book !

- set version to 1.5

-----

# The Nice Update

+ update README and add comments to help with config customization
+ hui: add arrow key support for navigation
+ hui: add arrow key support for cmdline
+ hui: add del key support for cmdline
+ hui: add backspace key support for cmdline
+ hui: add home key support for navigation and cmdline
+ hui: add end key support for navigation and cmdline
+ hui: add insert key support for cmdline

Also fix unsupported csi's being added to cmdline as input.

+ hui: add pgup and pgdown support for navigation

+ courier: add arrow key support for navigation and cmdline
+ courier: add del and backspace key support for cmdline
+ courier: add home and end key support for navigation and cmdline

Also add missing cmdline cursor reset.

+ courier: add insert key support for cmdline

Also fix unsupported csi's being added to cmdline as input.

+ courier: add pgup and pgdown support for navigation

+ rename "runtime" variables into "rt" to solve line len issues

+ hui: add cmdLineHistory, which is navigated by up and down arrow
+ courier: add cmdLineHistory, which is navigated by up and down arrow

+ gofmt

Also make hui's and courier's cmdline stuff more alike.

+ unify cmdLine logic
  (Only difference is that custom commands should be able to alter the different
   runtime structs, which is no worry anymore thanks to closures.
   Seriously they came in clutch, as I tried for weeks to elegantly solve this
   without butchering the workspace.)

+ fix const and var names

+ unify HandleKeyCmdLine functions

+ fix cmdline cursor ignoring any alignment that is not left

There is still some weirdness about right alignment in the cmdline.
I would need to add an artificial space to the end of the cmdline, to properly
set the cursor after the last character.

+ fix cursor may not being hidden after child invocation
+ fix cmdline handling running with empty input
+ fix cmdLineHistory going back even when there is nothing yet
+ fix cmdline int as cmd stepping out of range

This reveals a weird courier inconsistency,
which also caused it to be unaffected.
There is likely some other courier bug with how the Scroll value is being used.

+ UNNECESSARY: fix courier header and title breaking when aligned to the right
  Gone.
  What did you say?
  Gone.
  Are you telling me there is no bug?
  Yes.
  Oh.

+ hui: add cursor for each menu in menuPath

Thus leaving a menu brings you back to where the cursor was before.

+ hui: add feedback for useless interactions

Such as trying to open a go/shell entry or executing a menu entry.

+ hui: add entry prefix and postfix variations for on cursor hover

This is for colorless environments.

+ add chmod to install for configs to make them explicitly readable for users

+ courier: fix scrolling going past the last line

This also caused the scroll vs int-cmd inconsistency.

+ hui: add config validation against entries that lead to non-existent menus

+ set version to 1.4

# The White Magic Update

+ configs: add most variables to some struct

This is for position independency within the toml and it looks better.

+ hui: add custom commands via scripts
+ courier: add custom commands via scripts

+ BSD compat tests
	+ FreeBSD
	+ OpenBSD

+ bypass a bug with entry fg color inconsistency before cursor
compared with after

This happens once text with FG of #000000 has been printed.
After that the default FG color will change, at least on Fedora 39 and
Ubuntu 23.10.
So instead just explicitly print white instead of default color.
Also set version to 1.3.

# v1.2 The non-brazilian Update

+ hui: add go scripting interface for entries

For that, runtime variables have been moved from local main function variables
to a struct defined in common.
I first wanted to make it global variables in hui but... eh,
plus I would have had to import hui in scripts, which means circular inclusion.

+ hui: add go-entry config-values

Also add PagerTitle to hui and courier config.
Also also fix pagerTitle not being given to dev-courier at all.

+ hui: add scripting interface for start and quit
+ courier: add scripting interface for start and quit
+ move all main func variables to respective runtime structs

This was painful.

+ add config validation for if a script function actually exists in func map

Also undo me chopping everything into tiny modules.

+ update demo configs

+ fix event scripts interfering with flags (-v, -a and -h)

Also set version to 1.2.

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
