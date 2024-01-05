```
j6      j6
QQ      4P
QQ
QQ
QQQQQ6  jg
QQQQQQ  QQ
QQ  QQ  QQ
QQ  QQ  QQ
QQ  QQQQQQ
4P  4QQQQP
```

# What is this

A customizable terminal user-interface for common tasks and personal tastes,
inspired by suckless software.  
You can create TUI menus in a config file and then deploy it to your user.  
Then set hui as their default shell, to chain them into specific tasks :D  
A scripting interface allows you to tack logic onto the menus.  
With it you can even create entire menus at runtime.  

# HUI

This project is a Go rewrite of
[the original hui](https://github.com/SchokiCoder/hui).  
The purpose of this is to figure out if Go can be a 100% replacement for C in
the context of hui.  
So this will reimplement all of the original features plus a config file system.  
Depending on how well this works, this repo may cause the original to become
obsolete.  

# Install (no scripts yet, won't work)

Follow these steps:  

- NA

This will install two binaries "hui" and "courier".  
Courier is the pager that also lives here, because they share a lot of code so
they can look and feel similar.  
If you don't wish to have "courier", edit
[NA](https://github.com/SchokiCoder/gohui/blob/main/NA_E404).  

# Contributing

If you wish to do that, follow these steps to make it convenient for both of us:  

- have a look at docs/todo.md to see short term goals
- have a look at docs/goals.md to see long term goals
- create a fork on GitHub
- commit your changes
- create a pull request on GitHub

Thank you and have fun.  
