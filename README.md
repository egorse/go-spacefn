# The spacefn  alike keyboard layout for linux

Having difficult time with 60% keyboard and burden of xmodmap, xkb or others solutions
made this small tool to hookup the keyboard input events and remap keys.

The fn key is "space" by default. 
The fn key is bypassed on repeteation.
You have to press mapped key prior the repeteation starts.
When mapping is activated, all non-mapped keys are bypassed.

Default mapping is:
| Input      |  Result    |
|------------|------------|
| fn + W     |  Up        |
| fn + A     |  Left      |
| fn + S     |  Down      |
| fn + D     |  Right     |
| fn + R     |  Page Up   |
| fn + F     |  Page Down | 
| fn + Q     |  Home      |
| fn + E     |  End       |
| fn + 1-9,0 |  F1-F9,F10 |
| fn + ESC   |  `         |
| fn + `     |  ESC       |

# How to use

Ensure you have permissions to access to /dev/input/by-id/* (i.e. you are in group input, or whatever your distro uses) and /dev/uinput (see i.e. https://github.com/aksommerville/wiimote-uinput how to handle it); or use sudo. Then:

```#./go-spacefn```

# Current known/expected limitations

Basic hot plug.
No command line args.
No map customization.

# TODO
- customize fn key
- customize the map
