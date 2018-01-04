# hakase

Calculate ownership-ness of code from how many authors committed.

```
✘╹◡╹✘ < hakase -repo ~/devel/src/github.com/plack/plack -max_commits 100 -file lib/Plack/Component.pm | jq .
{
  "files": {
    "lib/Plack/Component.pm": {
      "Jay Hannah": {
        "count": 1,
        "score": 0.05
      },
      "Karen Etheridge": {
        "count": 1,
        "score": 0.05
      },
      "Shawn M Moore": {
        "count": 1,
        "score": 0.05
      },
      "Stevan Little": {
        "count": 5,
        "score": 0.25
      },
      "Tatsuhiko Miyagawa": {
        "count": 11,
        "score": 0.55
      },
      "hiratara": {
        "count": 1,
        "score": 0.05
      }
    }
  }
}
```

It is useful for building a tool such as facebook/mention-bot but more configurable one.
