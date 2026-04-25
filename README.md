# vint

Very fast linter for Golang. Designed to be MUCH FASTER replacement of golangci-lint and support as much rules as possible.

<p align="center">
  <img src="./assets/logo.png" alt="" width="300">
</p>

Here's how `vint` is different from `golangci-lint`:

- More than 7x faster running the same rules as golangci-lint.
- Clear centralized unified configuration.

This is still in progress, docs are not updated and in need of some cleanup, but vint already has close to 500 rules supported and benchmarks show quite well performance.

## Install and use

```bash
go get -tool github.com/vint-go/vint # need go 1.25+
```

```bash
go tool vint ./...
```

## Credits

This project is a fork of [revive](https://github.com/mgechev/revive) and contains all its original code and rules, with the addition of rules reproducing the behavior of these linters:

| | | | | |
| --- | --- | --- | --- | --- |
| [asasalint](https://github.com/alingse/asasalint) | [asciicheck](https://github.com/golangci/asciicheck) | [bidichk](https://github.com/breml/bidichk) | [canonicalheader](https://github.com/lasiar/canonicalheader) | [containedctx](https://github.com/sivchari/containedctx) |
| [contextcheck](https://github.com/kkHAIKE/contextcheck) | [cyclop](https://github.com/bkielbasa/cyclop) | [decorder](https://gitlab.com/bosi/decorder) | [dupword](https://github.com/Abirdcfly/dupword) | [durationcheck](https://github.com/charithe/durationcheck) |
| [errchkjson](https://github.com/breml/errchkjson) | [errname](https://github.com/Antonboom/errname) | [exhaustive](https://github.com/nishanths/exhaustive) | [exptostd](https://github.com/ldez/exptostd) | [fatcontext](https://github.com/Crocmagnon/fatcontext) |
| [forbidigo](https://github.com/ashanbrown/forbidigo) | [forcetypeassert](https://github.com/gostaticanalysis/forcetypeassert) | [gci](https://github.com/daixiang0/gci) | [ginkgolinter](https://github.com/nunnatsa/ginkgolinter) | [gochecknoglobals](https://github.com/leighmcculloch/gochecknoglobals) |
| [gochecksumtype](https://github.com/alecthomas/go-check-sumtype) | [gocognit](https://github.com/uudashr/gocognit) | [godot](https://github.com/tetafro/godot) | [gofmt](https://github.com/golangci/gofmt) | [gofumpt](https://github.com/mvdan/gofumpt) |
| [goheader](https://github.com/denis-tingaikin/go-header) | [goimports](https://github.com/golang/tools) | [golint](https://github.com/golang/lint) | [gomoddirectives](https://github.com/ldez/gomoddirectives) | [gomodguard](https://github.com/ryancurrah/gomodguard) |
| [gosimple](https://github.com/dominikh/go-tools) | [gosmopolitan](https://github.com/xen0n/gosmopolitan) | [grouper](https://github.com/leonklingele/grouper) | [iface](https://github.com/uudashr/iface) | [importas](https://github.com/julz/importas) |
| [inamedparam](https://github.com/macabu/inamedparam) | [interfacebloat](https://github.com/sashamelentyev/interfacebloat) | [kubeapilinter](https://github.com/kubernetes-sigs/kube-api-linter) | [logcheck](https://github.com/timonwong/loggercheck) | [makezero](https://github.com/ashanbrown/makezero) |
| [mirror](https://github.com/butuzov/mirror) | [modernize](https://github.com/golang/tools) | [musttag](https://github.com/go-simpler/musttag) | [nestif](https://github.com/nakabonne/nestif) | [nilerr](https://github.com/gostaticanalysis/nilerr) |
| [nilnesserr](https://github.com/alingse/nilnesserr) | [nilnil](https://github.com/Antonboom/nilnil) | [nonamedreturns](https://github.com/firefart/nonamedreturns) | [nosprintfhostport](https://github.com/stbenjam/no-sprintf-host-port) | [paralleltest](https://github.com/kunwardeep/paralleltest) |
| [perfsprint](https://github.com/catenacyber/perfsprint) | [prealloc](https://github.com/alexkohler/prealloc) | [predeclared](https://github.com/nishanths/predeclared) | [promlinter](https://github.com/yeya24/promlinter) | [protogetter](https://github.com/ghostiam/protogetter) |
| [reassign](https://github.com/curioswitch/go-reassign) | [recvcheck](https://github.com/raeperd/recvcheck) | [rowserrcheck](https://github.com/jingyugao/rowserrcheck) | [scopelint](https://github.com/kyoh86/scopelint) | [sorted](https://github.com/ravsii/sorted) |
| [spancheck](https://github.com/jjti/go-spancheck) | [sqlclosecheck](https://github.com/ryanrolds/sqlclosecheck) | [structcheck](https://github.com/golangci/check) | [stylecheck](https://github.com/dominikh/go-tools) | [tagliatelle](https://github.com/ldez/tagliatelle) |
| [testableexamples](https://github.com/maratori/testableexamples) | [thelper](https://github.com/kulti/thelper) | [tparallel](https://github.com/moricho/tparallel) | [typecheck](https://github.com/golangci/golangci-lint) | [usestdlibvars](https://github.com/sashamelentyev/usestdlibvars) |
| [usetesting](https://github.com/ldez/usetesting) | [varcheck](https://github.com/golangci/check) | [wastedassign](https://github.com/sanposhiho/wastedassign) | [wrapcheck](https://github.com/tomarrell/wrapcheck) | [wsl_v5](https://github.com/bombsimon/wsl) |
| [zerologlint](https://github.com/ykadowak/zerologlint) | | | | |


## License

MIT
