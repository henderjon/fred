# fred

FRiendlier ED is a line based text editor for terminals very similar to ed(1).

`fred` is the author's fourth attempt at an ed(1) clone. Instead of starting with "Software Tools"[^1] or it's sequel "Software Tools for Pascal"[^2], the author started with "Lexical Scanning in Go - Rob Pike"[^3]. This effort is nowhere near the beauty of Thompsons's recursive global state machine but it was written to be read. In other words, in using as elementary Go as possible, the author's hope is that the effort's source code is readable even to a future version of himself. To that end, "brute force" was used copiously in an attempt to stay true to the axiom "clear is better than clever"--because the author is not clever[^4]. The end result is not pretty, certainly not optimized or cleverly crafted, but it is testable and after letting it sit, still comprehendible to the author.

`fred` has also been tweaked to the author's preferences vis-a-vis how they usually use `ed`. For more details `$ fred -h`

## Installation

`$ make install`

By default `fred` uses a tmp file as scratch space. This reduces the memory footprint. If it's more desireable to keep all the scratch space in memory, build `fred` in memory mode:

`$ FREDMODE=memory make install`

--
[^1]: [Software Tools](https://a.co/d/57j2eG0)
[^2]: [Software Tools for Pascal](https://a.co/d/jllgMxg)
[^3]: [Lexical Scanning in Go - Rob Pike](https://www.youtube.com/watch?v=HxaD_trXwRE)
[^4]: "Debugging is twice as hard as writing the code in the first place. Therefore, if you write the code as cleverly as possible, you are, by definition, not smart enough to debug it."
