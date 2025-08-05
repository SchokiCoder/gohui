# os.Open() and io.ReadAll()

os.Open() happiliy ignores "etc/hui/hui.toml" not existing and doesn't return an
error, WTF. So it just goes with that.  
Then io.ReadAll(f) on that *ghost file* just returns an error of
"invalid argument".  
This is the most retarted waste of time i have seen in quite a while.  
  
Keep in mind that os.Open() **just** opens for reading.  
There is no point in trying to open a non existing file for a read, unlike
when writing.  
So what is this? Why is this a thing? And why is io.ReadAll()'s error message so
useless?  

## go-toml v2

Plus "go-toml/v2"'s Decode() just returned the "invalid argument" error of the
io.Reader it itself used, which at first gaslit me into thinking some toml
variable is wrong.  

# manpages

Why no manpages?  
Roff is garbage and i don't know why it is not criticized more  
but on the other hand we have:  

 > No DeVeLePoR wAnT's To WrItE dOcUmEnTaTiOn waaaah

So i protest, until the man command supports markdown.  
Instead just use -h or --help... just like with the go compiler.  
