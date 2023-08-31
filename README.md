# fred

FRiendlier ED is a line based text editor for terminals very similar to `ed` (The Standard Unix Editor).

`fred` is the author's fourth attempt at an `ed` clone. Instead of starting with "Software Tools"[^1] or it's sequel "Software Tools for Pascal"[^2], the author started with "Lexical Scanning in Go - Rob Pike"[^3]. This effort is nowhere near the beauty of Thompson's recursive global state machine but it was written to be read.

In other words, the hope is that the effort's source code is readable. To that end, "brute force" was used copiously in an attempt to stay true to the axiom "clear is better than clever"--because the author is not clever[^4]. There are plenty of leftover comments throughout the codebase to remind the future of what the past was thinking.

The end result is not pretty, certainly not optimized or cleverly crafted, but it is testable and after letting it sit, still comprehendible to the author.

The differences between `fred` and `ed` reflect the authors usage patterns when working with `ed`.

## Installation

`$ make install`

By default `fred` uses a tmp file as scratch space. This reduces the memory footprint. If it's more desireable to keep all the scratch space in memory, build `fred` in memory mode:

`$ FREDMODE=memory make install`

[^1]: [Software Tools](https://a.co/d/57j2eG0)
[^2]: [Software Tools for Pascal](https://a.co/d/jllgMxg)
[^3]: [Lexical Scanning in Go - Rob Pike](https://www.youtube.com/watch?v=HxaD_trXwRE)
[^4]: "Everyone knows that debugging is twice as hard as writing a program in the first place. So if you're as clever as you can be when you write it, how will you ever debug it? --Brian Kernighan; The Elements of Programming Style, 2nd ed. Page 10"
