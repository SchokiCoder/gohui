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
- add courier
- give hui multiline feedback to courier

-----

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
