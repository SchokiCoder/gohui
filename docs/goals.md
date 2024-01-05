# Explanation: Limitations and goals 

There are limitations, primary goals and secondary goals.  

Limitations describe what is __not allowed__ in the project/program.
Examples for that could be:  

- max. 3000 sloc
- C99, Linux code style
- libc only allowed dependency
- no incompatibility with FreeBSD, OpenBSD or Linux 

All __primary goals__ are to be implemented __before the final release__, unless
the implementation of a single goal proves to be too difficult or would come
with a too high cost such as a high increase of the sloc count.
The update implementing the last primary goal, becomes the final release update,
which in semantic versioning would be 1.0.0.  
Optimally the primary goals are documented as a roadmap.  

All __secondary goals__ are to be assigned __after the final release__.  
Optimally the secondary goals are documented as a roadmap,
starting where the primary goal-roadmap ended.  

# Limitations

- written in Go, with default gofmt as code style
- max. 3000 sloc
- max. 200 sloc for build system
- strictly POSIX compatible
- this rewrite must strife to implement the same features as the original before
  it can implement new ones in itself

# Roadmap

## v0.1.0 User menu

Upon start, the user is greeted with a menu.  
A menu has entries containing commands and submenus.  
Every menu has it's own custom title, which is displayed under the header.  
Use of fore- and background-colors are planned.  
They are configurable.  
The following areas can be configured (fore- and background color):  

- header
- title
- entry
- entry on which the user's cursor is hovering
- feedback line
- internal command line

The keybinds and commands are generally inspired by Vim and Neovim.  
Application defined commands are entered after pressing colon and many commands
are 1:1 compatible with Vim.  

keys:

- h: back
- j: down
- k: up
- l: into menu
- L: execute
- _Enter_: enter command (command line)

commands:

- q, quit, exit: exit app
- NUMBER: goto entry $NUMBER (starting at 1)

## v0.2.0 Config file

- file format
- read priority "/etc" over "~/.config/$FILE" over "~/.$FILE" over "$CWD/$FILE"
- if no config found, panic
- read config file and remove source code config

## v0.3.0 Child process execution

- just like hui 1.4
- add courier and let it handle multiline feedback

## v1.0.0 Final polish

This task is for creating/polishing the end-user cmd-line interface.  
That includes:  

- mangen?
- config file manpage
- hui manpage
- courier manpage
- POSIX call options
- return values
- print messages (consistency, version information, license information)

Provide generalist standard configuration which says so itself via main menu
title.

## v1.1.0 configurable aligns and margins

- Header
- Title
- Entries
- Feedback
- Cmdline
- Multiline feedback

## v1.2.0rc Scripting interface (Experimental)

Maybe Lua?
Lua in Rust didn't work too well last time i tried.
