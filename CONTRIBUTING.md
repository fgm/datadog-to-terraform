Thanks for your interest in contributing! If anything in here is unclear or outdated, 
please update it too 😄

# Development setup

1. Fork this repo
1. Run `go get ./...` in the project directory

## Testing locally

You can run the tests with

```
go test ./...
```

We use [goldie v2](https://github.com/sebdah/goldie) for golden files testing.
To update the snapshots after making changes, you can run:

```
go test -update
```

Don't forget to review the updates and make sure only intended changes are included!

# Pull Requests

- _Document any change in behaviour._ Make sure the README and any other relevant documentation are kept up-to-date.
- _Create topic branches._ Don't ask us to pull from your master branch.
- _One pull request per feature._ If you want to do more than one thing, send multiple pull requests.
- _Send coherent history._ Make sure each individual commit in your pull request is meaningful. If you had to make multiple intermediate commits while developing, please squash them before submitting them.
