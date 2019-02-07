# easytile
Adds keyboard-shortcut driven tiling to any floating WM (Only tested on openbox).

Made this tool to learn Golang, inspired from https://github.com/ssokolow/quicktile. This has only been tested on Openbox, however, I don't see it having problems on other window managers. The layout of the program looks very monolithic to me and I wouldn't call myself a Go expert at all.

To run it, simply clone this repository and run `./easytile` (You might need to make this executable though).

## shortcuts
These shortcuts will move and resize the active window according to the final key in the combination:

 - Control + Shift + Numpad 6 - Right half
 - Control + Shift + Numpad 4 - Left half
 - Control + Shift + Numpad 8 - Top half
 - Control + Shift + Numpad 2 - Bottom half
 - Control + Shift + Numpad 7 - Top left
 - Control + Shift + Numpad 1 - Bottom left
 - Control + Shift + Numpad 9 - Top right
 - Control + Shift + Numpad 3 - Bottom Right
 - Control + Shift + Numpad 5 - Maximize Window
 - Control + Shift + Numpad Enter - Move window to the next head (Tested on 3 monitors, I don't know how it'll fair on a single monitor)
