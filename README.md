<p align="center">
  <img src="https://github.com/EgeBalci/sgn/raw/master/img/banner.png">
  </br>
  <a href="https://github.com/EgeBalci/sgn">
    <img src="https://img.shields.io/badge/version-2.0.0-green.svg?style=flat-square">
  </a>
  <a href="https://goreportcard.com/report/github.com/egebalci/sgn">
    <img src="https://goreportcard.com/badge/github.com/egebalci/sgn?style=flat-square">
  </a>
  <a href="https://github.com/EgeBalci/sgn/issues">
    <img src="https://img.shields.io/github/issues/egebalci/sgn?style=flat-square&color=red">
  </a>
  <a href="https://raw.githubusercontent.com/EgeBalci/sgn/master/LICENSE">
    <img src="https://img.shields.io/github/license/egebalci/sgn.svg?style=flat-square">
  </a>
  <a href="https://twitter.com/egeblc">
    <img src="https://img.shields.io/badge/twitter-@egeblc-55acee.svg?style=flat-square">
  </a>
</p>

SGN is a polymorphic binary encoder for offensive security purposes such as generating statically undetecable binary payloads. It uses a additive feedback loop to encode given binary instructions similar to [LFSR](https://en.wikipedia.org/wiki/Linear-feedback_shift_register). This project is the reimplementation of the [original Shikata ga nai](https://github.com/rapid7/metasploit-framework/blob/master/modules/encoders/x86/shikata_ga_nai.rb) in golang with many improvements. 

This Repo is a fork version of [github.com/EgeBalci/sgn](https://github.com/EgeBalci/sgn) to bring sgn to javascript.

## Build

this repo replace all cgo([github.com/EgeBalci/keystone-go](https://github.com/EgeBalci/keystone-go)) to [github.com/AlexAltea/keystone.js](https://github.com/AlexAltea/keystone.js). Then it is compiled to wasm.

```
$ git clone github.com/akkuman/sgn
$ cd sgn
$ GOARCH=wasm GOOS=js go build -ldflags="-s -w" -trimpath
```

## Demo

Please see [github.com/akkuman/sgn-html](https://github.com/akkuman/sgn-html)

## Why not PR to upstream

Because I think this repo has a different purpose. I will communicate with the author([EgeBalci](https://github.com/EgeBalci)) to see if it can be used as a third party link for sgn.
