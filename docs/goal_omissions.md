# v1.0.0 Final polish

## everything man-related

Roff sucks.  
See [rants](https://github.com/SchokiCoder/gohui/blob/main/docs/rants.md#manpages).  

## return values

...why?  
It's unnecessary.  

# v1.2

## Scripting explicitly in Lua

Don't want to.  
I tried but doing that is annoying.  
Here are some thoughts:  

> Do you wish to add functionality to gohui?
> Do you want it to have a special quirk?
> Good, you are probably the kind that can handle having to have a compiler around.
> If not, get comfortable.

I do the suckless way in this case and i think that is the opitmal solution.  
Getting into hui is buttery smooth with the toml config but getting deep into it
then requires using a compiler (indirectly since my `./build.sh` is an easy
option).  

So the goal is not really omitted but the planned implementation thereof.  
