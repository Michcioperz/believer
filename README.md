a quick solution for attaching another screen with xrandr

- first it parses xrandr output in a way that cries to heaven for vengeance
- next it figures out the best common resolution (it assumes xrandr had put them in a best-to-worst order)
- then it says what it's about to do
- it waits 3 seconds just in case
- finally it executes an xrandr command that disables all screens that aren't to be set and sets all the other screens to the calculated resolution

if you give it any commandline arguments it assumes they are device names that you want set up (and it assumes you want the rest of devices off)

if you don't give it any commandline arguments it assumes you want to configure all devices

it has no Go dependencies outside the standard library
