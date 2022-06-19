# What is Countube?
In a nutshell, it's a simple and very specific video generator.
<br/>

Many years ago Gerd Jansen showed me one of this art projects, **_Countune_**, which is an online mega-picture which follows some defined pattern with configurable constraints. People all over the world contribute by creating their own picture, which is appended to the whole picture.
I was fascinated by the patterns and the colors (it reminded me a lot of color patterns
in old computers), and I liked the idea of having it in my videos.

Therefore, I thought of writing a script to generate videos that can scroll multiple Countune pictures.
I wrote the first one many years ago in Python but, unfortunately, lost the sources. I decided to rewrite it in Golang this time, to play around with that language.
In complete lack of originality, for now I decided to name it Countube.

# How does it work?
You basically give Countube a bunch of parameters, mainly the duration of the video and how fast you want to scroll the picture. Then it will pick random pictures from the whole Countune set and generate the video for you.
For now, this is configurable in `main.go`

# Status of the project
Even though it is usable, it's not "production ready", and it's lacking at least a command line client. Therefore, for now I'm not providing instructions for non Golang devs.
I will work on this and some other improvements soon.

# I created a .frames file. Now how do I make a video out of it?
The `.frames` file is a concatenation of jpegs, each corresponding to each frame in the video. **_ffmpeg_** tool can help in creating a video out of it. Example:
`$ ffmpeg -y -framerate 60 -i test.frames out.mp4`
_Important_:  Make sure that the frame rate provided to **_ffmpeg_** matches the one used to generate the frames file.

# How to add music to the video
One possibility is to create a single audio file concatenating all audio tracks.
You can have a look at the [ffmpeg documentation on concatenation](https://trac.ffmpeg.org/wiki/Concatenate).
Example:
1. Create `track_list.txt`
```
# this is a comment
file '/path/to/file1.wav'
file '/path/to/file2.wav'
file '/path/to/file3.wav'
```
2. Concatenate audio using the list file:
`$ ffmpeg -f concat -safe 0 -i track_list.txt -c copy sample.wav`
3. Create a video using the .frames file and the audio file:
`$ ffmpeg -y -framerate 60 -i sample.frames -i sample.wav sample.mp4`
