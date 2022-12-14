# VHS documentation
#
# Output:
#   Output <path>.gif               Create a GIF output at the given <path>
#   Output <path>.mp4               Create an MP4 output at the given <path>
#   Output <path>.webm              Create a WebM output at the given <path>
#
# Require:
#   Require <string>                Ensure a program is on the $PATH to proceed
#
# Settings:
#   Set FontSize <number>           Set the font size of the terminal
#   Set FontFamily <string>         Set the font family of the terminal
#   Set Height <number>             Set the height of the terminal
#   Set Width <number>              Set the width of the terminal
#   Set LetterSpacing <float>       Set the font letter spacing (tracking)
#   Set LineHeight <float>          Set the font line height
#   Set LoopOffset <float>%         Set the starting frame offset for the GIF loop
#   Set Theme <json|string>         Set the theme of the terminal
#   Set Padding <number>            Set the padding of the terminal
#   Set Framerate <number>          Set the framerate of the recording
#   Set PlaybackSpeed <float>       Set the playback speed of the recording
#
# Sleep:
#   Sleep <time>                    Sleep for a set amount of <time> in seconds
#
# Type:
#   Type[@<time>] "<characters>"    Type <characters> into the terminal with a
#                                   <time> delay between each character
#
# Keys:
#   Backspace[@<time>] [number]     Press the Backspace key
#   Down[@<time>] [number]          Press the Down key
#   Enter[@<time>] [number]         Press the Enter key
#   Space[@<time>] [number]         Press the Space key
#   Tab[@<time>] [number]           Press the Tab key
#   Left[@<time>] [number]          Press the Left Arrow key
#   Right[@<time>] [number]         Press the Right Arrow key
#   Up[@<time>] [number]            Press the Up Arrow key
#   Down[@<time>] [number]          Press the Down Arrow key
#   Ctrl+<key>                      Press the Control key + <key> (e.g. Ctrl+C)
#
# Display:
#   Hide                            Hide the subsequent commands from the output
#   Show                            Show the subsequent commands in the output

Output demo.gif

Set FontSize 18
Set Width 1200
Set Height 600

Ctrl+l
Type "# Welcome to fun-run!" Sleep 500ms Enter

Sleep 5s

Type "fun-run --version" Sleep 500ms Enter

Sleep 3s

Ctrl+l
Type "# Let's create a new config file..." Sleep 500ms Enter

Sleep 1s

Type "fun-run init" Sleep 500ms Enter

Sleep 1s

Ctrl+l
Type "# And then examine the file..." Sleep 500ms Enter

Sleep 1s

Type "less fun-run.yaml" Sleep 500ms Enter

Sleep 5s

Type "q" Sleep 500ms Enter

Ctrl+l
Type "# Let's validate the file..." Sleep 500ms Enter

Sleep 1s

Type "fun-run validate fun-run.yaml" Sleep 500ms Enter

Sleep 2s

Ctrl+l
Type "# Now...let's run the file..." Sleep 500ms Enter

Sleep 1s

Type "fun-run run fun-run.yaml" Sleep 500ms Enter

Sleep 10s

Ctrl+c

Sleep 5s



