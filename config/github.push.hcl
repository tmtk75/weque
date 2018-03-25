watches = [
  #{
  #  type = "keyprefix"
  #  prefix = "weque/github"
  #  args = ["./github.push.sh"]
  #}
  {
    type = "event"
    name = "github.push"
    args = ["./github.push.sh"]
  }
]
