# Contributing

## Branches

Branch from `dev` and open the pull request against `dev`.

`main` is a CI-owned mirror of the latest release: the Release workflow commits
the built web assets there. A branch cut from `main` carries those commits, so
a pull request into `dev` shows them as yours. CI also runs only on pull
requests into `dev`.

```shell
git remote add upstream https://github.com/voilelab/toolgui.git
git fetch upstream
git switch -c my-change upstream/dev
```

Already branched from `main`? Move your commits onto `dev`, then change the
pull request's base to `dev`:

```shell
git fetch upstream
git rebase --onto upstream/dev upstream/main my-change
git push --force-with-lease
```

## Before opening a pull request

```shell
task stub_assets   # once, so Go builds without the frontend
go vet ./...
go test ./...
```

- Adding or changing a `//tgcomp:export` component: run `go generate ./toolgui/tgcomp`.
- An API change updates its page under `docs/src` in the same pull request.

See [README.md](README.md#for-dev) for running the demo and the E2E tests.
