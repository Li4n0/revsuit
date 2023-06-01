# Changelog

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.

### [0.7.1](https://github.com/Li4n0/revsuit/compare/v0.7.0...v0.7.1) (2023-06-01)


### Bug Fixes

* fix the error that the maximum length of http rule response body is 191 when using mysql ([bf8a43f](https://github.com/Li4n0/revsuit/commit/bf8a43fa875327c3191499a214c420bfa78334d0))

## [0.7.0](https://github.com/Li4n0/revsuit/compare/v0.6.0...v0.7.0) (2023-05-23)


### Features

* add `value` field to the dns log ([#67](https://github.com/Li4n0/revsuit/issues/67)) ([223c14a](https://github.com/Li4n0/revsuit/commit/223c14ab33757e06b0d8177091982582e08e70c4))

## [0.6.0](https://github.com/Li4n0/revsuit/compare/v0.5.2...v0.6.0) (2023-03-24)


### Features

* support for using `Flag-Filter` request headers to filter the flags that the client expects to receive ([#64](https://github.com/Li4n0/revsuit/issues/64)) ([b33de11](https://github.com/Li4n0/revsuit/commit/b33de11ba8f126b03a45ad15cae69331ad59b865)) (Thanks to @whwlsfb) 


### Bug Fixes

* fix static resource loading failure after modifying `admin_path_prefix` ([b76757d](https://github.com/Li4n0/revsuit/commit/b76757da331b56f240ec2a0c69d479f3184bc4c2)), closes [#63](https://github.com/Li4n0/revsuit/issues/63)

### [0.5.2](https://github.com/Li4n0/revsuit/compare/v0.5.1...v0.5.2) (2022-11-16)


### Bug Fixes

* fix bug that the files associated with records still exist after deleting mysql and ftp records ([5e80591](https://github.com/Li4n0/revsuit/commit/5e80591c856a63c851ba5e17812af89373d76cf5))
* fix panic when update config ([#62](https://github.com/Li4n0/revsuit/issues/62)) ([48f55c5](https://github.com/Li4n0/revsuit/commit/48f55c53273fcea579c9fdfaf2939c90df3588b1))

### [0.5.1](https://github.com/Li4n0/revsuit/compare/v0.5.0...v0.5.1) (2022-09-13)


### Bug Fixes

* fix the bug that authentication always fails after token change ([#59](https://github.com/Li4n0/revsuit/issues/59)) ([3fa2cae](https://github.com/Li4n0/revsuit/commit/3fa2cae86260d39fa2ffd43d59b2754327f6169a))
* fix the error of `database is locked` when use sqlite ([#60](https://github.com/Li4n0/revsuit/issues/60)) ([d29b767](https://github.com/Li4n0/revsuit/commit/d29b767d9cb2b70f365d0830c4bbeaae09524a5a))

## [0.5.0](https://github.com/Li4n0/revsuit/compare/v0.4.0...v0.5.0) (2022-01-16)


### Features

* **build:** replace sqlite driver with which does not depend cgo ([#52](https://github.com/Li4n0/revsuit/issues/52)) ([82a2023](https://github.com/Li4n0/revsuit/commit/82a20236e1b1a092277b645554a38513e3bb4087))
* **dns:** support use multiple root domains ([#51](https://github.com/Li4n0/revsuit/issues/51)) ([913e0f7](https://github.com/Li4n0/revsuit/commit/913e0f78db380cb1beeb88b57dc30f47e7c26c1c))
* **ldap:** add ldap protocol support ([#50](https://github.com/Li4n0/revsuit/issues/50)) ([fc505b0](https://github.com/Li4n0/revsuit/commit/fc505b0733a8375d29b910aa897c6e0b642ab34b))
* support auto check upgrade ([#53](https://github.com/Li4n0/revsuit/issues/53)) ([f6afe0f](https://github.com/Li4n0/revsuit/commit/f6afe0fa12188f17b74548129035f7b89a149660))


### Bug Fixes

* **settings:** fix the bug of update rmi & mysql addr settings not working ([#49](https://github.com/Li4n0/revsuit/issues/49)) ([06ae3d0](https://github.com/Li4n0/revsuit/commit/06ae3d0555eff3bf2d6a53b2c405d1dd50beea87))

## [0.4.0](https://github.com/Li4n0/revsuit/compare/v0.3.0...v0.4.0) (2021-12-26)


### Features

* **database:** support postgres ([#43](https://github.com/Li4n0/revsuit/issues/43)) ([8b30f83](https://github.com/Li4n0/revsuit/commit/8b30f83075bab3f62546793c471b74075745d8d0))


### Bug Fixes

* **rmi:** fix the problem that rmi has multiple rules matching in one record ([6f921f6](https://github.com/Li4n0/revsuit/commit/6f921f6c97f39638c73de13d87879514485125ef))
* **server:** fix auth failure when complex characters in token ([#42](https://github.com/Li4n0/revsuit/issues/42)) ([3242a49](https://github.com/Li4n0/revsuit/commit/3242a49b2731392125e5c22375549586495ee6fc))

## [0.3.0](https://github.com/Li4n0/revsuit/compare/v0.2.1...v0.3.0) (2021-12-22)


### Features

* **database:** add mysql support ([#39](https://github.com/Li4n0/revsuit/issues/39)) ([9d7a5b4](https://github.com/Li4n0/revsuit/commit/9d7a5b45984bb3fee187146b0f27e77d6ec0ea0a))


### Bug Fixes

* **cli:** fix typos ([#41](https://github.com/Li4n0/revsuit/issues/41)) ([8e96962](https://github.com/Li4n0/revsuit/commit/8e969627917fda9ce220283eb34599b781f03ae3))
* **database:** fix the error when load package ([3ca242d](https://github.com/Li4n0/revsuit/commit/3ca242d0866f5f10aade3557d80b8252be7ef938))
* **mysql/record:** fix the bug when searching mysql record with username as keyword ([ce8d760](https://github.com/Li4n0/revsuit/commit/ce8d76056faec45d3e68164b3802ceaaec0468b9))
* **pkg:** fix the conflict between rank field and built-in function ([#38](https://github.com/Li4n0/revsuit/issues/38)) ([a59540d](https://github.com/Li4n0/revsuit/commit/a59540d42ddf820e445a8395c2995fa89fc14323))

### [0.2.1](https://github.com/Li4n0/revsuit/compare/v0.2.0...v0.2.1) (2021-07-24)


### Features

* **rule:** add `flag` as a built-in template variable ([ee36205](https://github.com/Li4n0/revsuit/commit/ee36205ac160a75d23dbb244c72f3bb2d185ab80))

## [0.2.0](https://github.com/Li4n0/revsuit/compare/v0.1.7...v0.2.0) (2021-07-22)


### Features

* support deleting logs based on filtering results ([#35](https://github.com/Li4n0/revsuit/issues/35)) ([9863821](https://github.com/Li4n0/revsuit/commit/9863821dbc351734569b7ec1eda2b6712bd02bb0))

### [0.1.7](https://github.com/Li4n0/revsuit/compare/v0.1.6...v0.1.7) (2021-07-11)


### Bug Fixes

* **http:** fix the bug that http record url search does not work ([ffc4318](https://github.com/Li4n0/revsuit/commit/ffc4318725017c7eda7ec11ca67d5de16ad1e8a8))
* **mysql:** fix the bug that crashes when mysql port conflicts ([c0bda70](https://github.com/Li4n0/revsuit/commit/c0bda707b750b7858fdd4efc414b726ee37e71da))

### [0.1.6](https://github.com/Li4n0/revsuit/compare/v0.1.5...v0.1.6) (2021-06-25)


### Bug Fixes

* **cli:** fix the bug that the prompt may not be output when generating config file ([5615351](https://github.com/Li4n0/revsuit/commit/56153510e11783c08843cbb8b0f827d5f8567842))
* **ipinfo:** fix a potential panic when downloading ip location database ([#33](https://github.com/Li4n0/revsuit/issues/33)) ([b1773ee](https://github.com/Li4n0/revsuit/commit/b1773ee5e9ebd8c8be72f315981268992c08b3f1))

### [0.1.5](https://github.com/Li4n0/revsuit/compare/v0.1.5-beta...v0.1.5) (2021-06-22 公开该项目以庆祝毕业，愿工作多年后的我，归来时仍是少年)


### Bug Fixes

* **notice:** fix the buf of notice with null record field ([#31](https://github.com/Li4n0/revsuit/issues/31)) ([ef38cbc](https://github.com/Li4n0/revsuit/commit/ef38cbc790f69716a335167321e1c7c8bdee2e41))

### [0.1.5-beta](https://github.com/Li4n0/revsuit/compare/v0.1.4-beta-fix-1...v0.1.5-beta) (2021-06-13)


### Features

* **ipinfo:** support GeoIP ([#26](https://github.com/Li4n0/revsuit/issues/26)) ([dc43b69](https://github.com/Li4n0/revsuit/commit/dc43b6973a5ac98e439d1353102ade2029b1d382))
* **platform:** support custom admin path prefix ([#25](https://github.com/Li4n0/revsuit/issues/25)) ([f7be2bc](https://github.com/Li4n0/revsuit/commit/f7be2bc2e67841178e9316995999b20f80a49df7))


### Bug Fixes

* **frontend:** fix the bug of autoRefresh ([#27](https://github.com/Li4n0/revsuit/issues/27)) ([d1013c2](https://github.com/Li4n0/revsuit/commit/d1013c2416c04bd72b9cbb0c7e6b4b3f2e717837))
* **frontend:** fix the frontend display bugs ([#28](https://github.com/Li4n0/revsuit/issues/28)) ([ddeced3](https://github.com/Li4n0/revsuit/commit/ddeced354042e0be6db5dd4feffa13ef22db231f))

### [0.1.4-beta-fix-1](https://github.com/Li4n0/revsuit/compare/v0.1.4-beta...v0.1.4-beta-fix-1) (2021-05-28)


### Bug Fixes

* **frontend:** fix the problem of autoRefresh failure ([#24](https://github.com/Li4n0/revsuit/issues/24)) ([b945069](https://github.com/Li4n0/revsuit/commit/b945069117fec6eb5de88557da5d4c2c996cfd90))
* **http:** fix the problem that ip_header setting does not take effect ([#23](https://github.com/Li4n0/revsuit/issues/23)) ([74f512d](https://github.com/Li4n0/revsuit/commit/74f512d2140fb97128acf56be803d1bd3b888fa3))

### [0.1.4-beta](https://github.com/Li4n0/revsuit/compare/v0.1.3-beta-fix1...v0.1.4-beta) (2021-05-27)

### Features

* **dns:** support custom dns service
  port ([7fb4c97](https://github.com/Li4n0/revsuit/commit/7fb4c97279e57d120a4e4aef5dda5c8f3c024835))
* **frontend:** auto refresh data when switching back to the revsuit
  console ([9b31406](https://github.com/Li4n0/revsuit/commit/9b314062a39ddc7acf7a7eab3570b24b9bb9d122))

### Bug Fixes

* **http:** fix the problem that ip_header configuration does not take
  effect ([b76c39e](https://github.com/Li4n0/revsuit/commit/b76c39e2fc1ada189feb783fdec76daffa11d1c7))
