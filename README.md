# lcpre
##### Determine longest common prefix among multiple strings

`lcpre` finds the longest common prefix of a given set of strings. 

```
$ lcpre superman supercalifragilistic
super
```

The components of each string may be partitioned into substrings ("words") separated by a common delimeter substring (flag `-s`).

For example, if the given strings are file paths and the delimeter is the file path separator (e.g., `/`), it will find the nearest common parent directory:

```
$ lcpre -s "/" "/root/path/a/a1" "/root/path/b" "/root/path/c/c1/c1.1"
/root/path

$ lcpre -s "/" "a/b/c" "d/e/f"


$ lcpre -s "/" "a/b/c" "a/d"
a
```

If the longest common prefix is simply the delimeter itself, the flag `-d` will control whether or not the delimeter is printed. This could be useful if the delimeter has special meaning as in the following example:

```
$ lcpre -s "/" "/A/B/C" "/D/E/F"


$ lcpre -d -s "/" "/A/B/C" "/D/E/F"
/
```

If a single string is input, it is returned verbatim. If there does not exist a common prefix, the empty string is returned.

