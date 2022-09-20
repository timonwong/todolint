# todolint

## Description

Requires TODO comments to be in the form of "TODO(author) ...

## Badges

![Build Status](https://github.com/timonwong/todolint/workflows/CI/badge.svg)
[![Coverage](https://img.shields.io/codecov/c/github/timonwong/todolint?token=Nutf41gwoG)](https://app.codecov.io/gh/timonwong/todolint)
[![License](https://img.shields.io/github/license/timonwong/todolint.svg)](/LICENSE)
[![Release](https://img.shields.io/github/release/timonwong/todolint.svg)](https://github.com/timonwong/todolint/releases/latest)

## Install

```shel
go install github.com/timonwong/loggercheck/cmd/loggercheck
```

## Usage

```
todolint: Requires TODO comments to be in the form of "TODO(author) ...

Usage: todolint [-flag] [package]
```

## Example

```go
package a

import "fmt"

// TODO: This is not ok // want `TODO comment should be in the form TODO\(author\)`
func NotOkFunc() {
}

// TODO(author1): This is ok
func OkFunc() {
}

type ABC struct {
	A int    // @FIXME: This field comment is not ok // want `TODO comment should be in the form FIXME\(author\)`
	B string // FIXME(author2): This field comment is ok
}

func Example() {
	// TODO(timonwong): This is ok
	//

	// 🚀🚀🚀 FixMe: 你好世界 // want `TODO comment should be in the form FIXME\(author\)`
	fmt.Println("Hello")

	fmt.Println("你好，世界") // fixme: more languages // want `TODO comment should be in the form FIXME\(author\)`
}
```