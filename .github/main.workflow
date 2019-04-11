workflow "bug label" {
  on = "label"
  resolves = [
    "GitHub Action for Slack"
  ]
}

action "Filters for GitHub Actions" {
  uses = "actions/bin/filter@3c98a2679187369a2116d4f311568596d3725740"
  args = "label urgent"
}

action "GitHub Action for Slack" {
  uses = "Ilshidur/action-slack@ab5f0955362cfdff2e0f0990f0272624e8cb5d13"
  secrets = ["SLACK_WEBHOOK"]
}
