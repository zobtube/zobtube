# Changelog

## Version 0.3.9
### Features
* [245f316](https://github.com/zobtube/zobtube/commit/245f31673b7234d9b079f0b655308c89475a628e) feat(actor/list): allow filtering by name
### Chores
* [137247e](https://github.com/zobtube/zobtube/commit/137247e0b51ea494699276dc96aecde82d34bcf5) chore(adm): remove dead code on adm/edit
* [137247e](https://github.com/zobtube/zobtube/commit/137247e0b51ea494699276dc96aecde82d34bcf5) chore(ui): rationalize headers
* [245f316](https://github.com/zobtube/zobtube/commit/245f31673b7234d9b079f0b655308c89475a628e) chore(actor/list): display video and link count
* [94fd399](https://github.com/zobtube/zobtube/commit/94fd3991b00f07b27f650dd9245bfdc81102bf25) chore(deps): update actions/upload-artifact action to v5
### Fixes
* [23303fd](https://github.com/zobtube/zobtube/commit/23303fd79303bb488b1ae4d34da35676113cdb09) fix(video/edit): move folder instead of video (fix #153)
* [6c153c9](https://github.com/zobtube/zobtube/commit/6c153c9846bd9b32fcf4b54443379d6064642343) fix(deps): update module github.com/urfave/cli/v3 to v3.5.0

## Version 0.3.8
### Features
* [00d6039](https://github.com/zobtube/zobtube/commit/00d603982f00cdb8a6fcd3e8094d51f12a8a034e) feat(clip/view): add mobile support
* [00d6039](https://github.com/zobtube/zobtube/commit/00d603982f00cdb8a6fcd3e8094d51f12a8a034e) feat(clip/view): allow scolling with arrow up and down
* [00d6039](https://github.com/zobtube/zobtube/commit/00d603982f00cdb8a6fcd3e8094d51f12a8a034e) feat(clip/view): allow scrolling with thumbs
### Chores
* [031a84a](https://github.com/zobtube/zobtube/commit/031a84ae317091f955c94434bd0025a9282c4f95) chore(e2e): remove test on main
* [4ba3d67](https://github.com/zobtube/zobtube/commit/4ba3d67e595bf23cbce343276084079857502243) chore(deps): update alpine docker tag to v3.22.2
* [c5b408b](https://github.com/zobtube/zobtube/commit/c5b408b46b70c297f65d46d443a9010051bd4e52) chore(go): cleanup dependencies
### Fixes
* [1da145c](https://github.com/zobtube/zobtube/commit/1da145cc6c257798a2eb7168148050cf80d37938) fix(deps): update module golang.org/x/image to v0.32.0
* [3a74450](https://github.com/zobtube/zobtube/commit/3a7445088634fc5ba85ba441a7a836c76f2a26ae) fix(deps): upgrade quic to fix CVE-2025-59530
* [e979c22](https://github.com/zobtube/zobtube/commit/e979c2239514338cb7c1090026ae1fc7aa56e5fe) fix(deps): update module golang.org/x/text to v0.30.0
### Documentation
* [46d7310](https://github.com/zobtube/zobtube/commit/46d73109b3329697ff3cd761d5aea4e68dc4a92a) doc: update stash differences

## Version 0.3.7
### Features
* [9000c13](https://github.com/zobtube/zobtube/commit/9000c1370f2d8d92488422d48f60f61a38246616) feat(adm/home): report errors on adm home
* [9000c13](https://github.com/zobtube/zobtube/commit/9000c1370f2d8d92488422d48f60f61a38246616) feat(main): ensure ffmpeg and ffprobe are available
* [9000c13](https://github.com/zobtube/zobtube/commit/9000c1370f2d8d92488422d48f60f61a38246616) feat(main): load providers softly and register errors
* [9de87c2](https://github.com/zobtube/zobtube/commit/9de87c2690fe297b526a984dd69a7e0e95f8453e) feat(test): add units on actor-api
* [9de87c2](https://github.com/zobtube/zobtube/commit/9de87c2690fe297b526a984dd69a7e0e95f8453e) feat(test): add units on auth-common
* [9de87c2](https://github.com/zobtube/zobtube/commit/9de87c2690fe297b526a984dd69a7e0e95f8453e) feat(test): add units on cleanup
* [9de87c2](https://github.com/zobtube/zobtube/commit/9de87c2690fe297b526a984dd69a7e0e95f8453e) feat(test): add units on providers
* [9de87c2](https://github.com/zobtube/zobtube/commit/9de87c2690fe297b526a984dd69a7e0e95f8453e) feat(test): add units on shutdown
* [a0a3dee](https://github.com/zobtube/zobtube/commit/a0a3dee27e051306f2a39306a599dee7f10bc387) feat(provider): add iafd (fix #89)
* [ac0a566](https://github.com/zobtube/zobtube/commit/ac0a56651bd278788e0f797eb4ca75d574a82ad9) feat: add unit tests on providers and config
* [e512031](https://github.com/zobtube/zobtube/commit/e512031c9a2b1fc9bd78b70a27f6afe7f6fda17e) feat(e2e): ensure video type is correct
### Chores
* [22babeb](https://github.com/zobtube/zobtube/commit/22babeb140d5c130b7da27068eba9d925636b85a) chore: add taskfile to check quality
* [22babeb](https://github.com/zobtube/zobtube/commit/22babeb140d5c130b7da27068eba9d925636b85a) chore: check code with gocritic
* [22babeb](https://github.com/zobtube/zobtube/commit/22babeb140d5c130b7da27068eba9d925636b85a) chore: ensure mod is up-to-date
* [22babeb](https://github.com/zobtube/zobtube/commit/22babeb140d5c130b7da27068eba9d925636b85a) chore: fix security recommendation of gosec
* [22babeb](https://github.com/zobtube/zobtube/commit/22babeb140d5c130b7da27068eba9d925636b85a) chore: increase code quality
* [22babeb](https://github.com/zobtube/zobtube/commit/22babeb140d5c130b7da27068eba9d925636b85a) chore: lint with gofumpt
* [22babeb](https://github.com/zobtube/zobtube/commit/22babeb140d5c130b7da27068eba9d925636b85a) chore: perform static check analysis
* [22babeb](https://github.com/zobtube/zobtube/commit/22babeb140d5c130b7da27068eba9d925636b85a) chore: prepare go revive
* [2bd9ebc](https://github.com/zobtube/zobtube/commit/2bd9ebc1a2e504b4a114cffc4e94783eb2f6fde0) chore(ci): disable sonarqube on pr
* [332b1a8](https://github.com/zobtube/zobtube/commit/332b1a868ea93a414a257b58920264b351c62c75) chore(test/e2e): rename variables
* [4dfaa7a](https://github.com/zobtube/zobtube/commit/4dfaa7acad0c6521fba05045166523490dcd4f21) chore(ci/playwright): remove python installation
* [a340236](https://github.com/zobtube/zobtube/commit/a340236281733a469a120e52338b4d0f95ad8659) chore(actor/edit): remove card on profile picture
* [b12141f](https://github.com/zobtube/zobtube/commit/b12141f6ac9410010506dcb525ce82b46932e30f) chore(ci): split tests between pr and main
### Fixes
* [4dfaa7a](https://github.com/zobtube/zobtube/commit/4dfaa7acad0c6521fba05045166523490dcd4f21) fix(ci/playwright): move playwright to self-hosted runner to avoid flakies
* [8596257](https://github.com/zobtube/zobtube/commit/8596257d16af952a22d5c9a04e4e54914e15dd4b) fix: remove sonar code smells
* [8efbd4c](https://github.com/zobtube/zobtube/commit/8efbd4c344ce9b675122199caf698dc527ead3ec) fix(provider/iafd): add logo
* [a0a3dee](https://github.com/zobtube/zobtube/commit/a0a3dee27e051306f2a39306a599dee7f10bc387) fix(actor/edit): bring back link deletion
* [ac0a566](https://github.com/zobtube/zobtube/commit/ac0a56651bd278788e0f797eb4ca75d574a82ad9) fix(provider/babepedia): return error properly
* [b842d5d](https://github.com/zobtube/zobtube/commit/b842d5dcfb92f018ed59df8f065143a64154358c) fix(controller): remove abstract controller typo
* [d4b905d](https://github.com/zobtube/zobtube/commit/d4b905d4db71c78289c345cfa4b9aa9f4221e0f8) fix(video/rename): check database return

## Version 0.3.6
### Features
* [190f90c](https://github.com/zobtube/zobtube/commit/190f90cd13a114e26fc7abeaf186dd7c4eb2dd13) feat(adm/home): rework ui
* [3f68e9b](https://github.com/zobtube/zobtube/commit/3f68e9b54f5e47e154ee42181b8c90e033055518) feat(video/view): improve interface
* [3f68e9b](https://github.com/zobtube/zobtube/commit/3f68e9b54f5e47e154ee42181b8c90e033055518) feat(view/view): add download link (fix #56)
* [436174a](https://github.com/zobtube/zobtube/commit/436174a3616f4cc30e8d9bc896ee9c2b56905c0f) feat(healthcheck): report status after error
* [e73d01f](https://github.com/zobtube/zobtube/commit/e73d01fd96d63d8b44fe3d1c3944b1fcf26a2618) feat: add sonarqube static analyzer (fix #57)
### Fixes
* [5929dcd](https://github.com/zobtube/zobtube/commit/5929dcd55101dcc55b25a88e012ac7d4521be429) fix(upload): remove forgottent sql debug statements
* [b7f9f09](https://github.com/zobtube/zobtube/commit/b7f9f09c3efaf38a654ae4ce59c9afcd870ebc2d) fix(ci): switch tests to pull_request_target
* [cf33ba1](https://github.com/zobtube/zobtube/commit/cf33ba101ccc3ce93856528e351a5a6eed167e7c) fix(category/list): remove max-height (fix #109)
* [eee1967](https://github.com/zobtube/zobtube/commit/eee1967886998acb99985c4409e54761b1966776) fix: use id from objects instead of user ones
* [f1948b2](https://github.com/zobtube/zobtube/commit/f1948b26f3069723e4083bb9fd69c760e3475627) fix(router): handle 404 through dedicated page

## Version 0.3.5
### Features
* [554b386](https://github.com/zobtube/zobtube/commit/554b3861bf0b0168c47c1c5b1f5bd973773e12b3) feat: add configuration to disable providers
* [554b386](https://github.com/zobtube/zobtube/commit/554b3861bf0b0168c47c1c5b1f5bd973773e12b3) feat: add configuration to set offline mode
* [554b386](https://github.com/zobtube/zobtube/commit/554b3861bf0b0168c47c1c5b1f5bd973773e12b3) feat(http/router): add template error handling
* [554b386](https://github.com/zobtube/zobtube/commit/554b3861bf0b0168c47c1c5b1f5bd973773e12b3) feat: store providers in database
### Chores
* [554b386](https://github.com/zobtube/zobtube/commit/554b3861bf0b0168c47c1c5b1f5bd973773e12b3) chore(adm): mutualize tabs shards
### Fixes
* [554b386](https://github.com/zobtube/zobtube/commit/554b3861bf0b0168c47c1c5b1f5bd973773e12b3) fix(actor/edit): handle properly provider failure

## Version 0.3.4
### Fixes
* [4f5b47b](https://github.com/zobtube/zobtube/commit/4f5b47b9da1154220b5932a4fd44db14a24a3fe4) fix(actor/edit): picture correct size (fix #107)
* [70d6bfb](https://github.com/zobtube/zobtube/commit/70d6bfb9c9c24f2f6490c209424b3207b995f4a6) fix(deps): update module github.com/urfave/cli-altsrc/v3 to v3.1.0
* [ce92998](https://github.com/zobtube/zobtube/commit/ce929983d4c3d92a1bf7659d06a2375a2c1da4b0) fix(goreleaser): extract relevant changelog part

## Version 0.3.3
### Fixes
* [df9b4c8](https://github.com/zobtube/zobtube/commit/df9b4c864e39120ebd918fddf3e84cd3f03accb8) fix(adm/categories): remove max-height

## Version 0.3.2
### Fixes
* [7903d70](https://github.com/zobtube/zobtube/commit/7903d70d1711a54416d7093b3b42bd2631127fb1) fix(goreleaser): bring back windows and darwin

## Version 0.3.1
### Features
* [4a467de](https://github.com/zobtube/zobtube/commit/4a467de5c2d9212bcb2c918a4fba6833876b4d0d) feat: add tool to generate changelog
* [a77aa88](https://github.com/zobtube/zobtube/commit/a77aa88686edbd43296e8685734675e23f577a27) feat(triage): add mass deletion
* [a77aa88](https://github.com/zobtube/zobtube/commit/a77aa88686edbd43296e8685734675e23f577a27) feat(triage): add mass import (fix #80)
* [a77aa88](https://github.com/zobtube/zobtube/commit/a77aa88686edbd43296e8685734675e23f577a27) feat(ui): add custom onload element
* [a77aa88](https://github.com/zobtube/zobtube/commit/a77aa88686edbd43296e8685734675e23f577a27) feat(ui): create common async function ajax wrapping jquery ajax method
* [a77aa88](https://github.com/zobtube/zobtube/commit/a77aa88686edbd43296e8685734675e23f577a27) feat(ui): split actor selection into a dedicated shard
* [a77aa88](https://github.com/zobtube/zobtube/commit/a77aa88686edbd43296e8685734675e23f577a27) feat(ui): split category selection into a dedicated shard
### Chores
* [4a467de](https://github.com/zobtube/zobtube/commit/4a467de5c2d9212bcb2c918a4fba6833876b4d0d) chore(goreleaser): change compilation more consistent naming
* [4a467de](https://github.com/zobtube/zobtube/commit/4a467de5c2d9212bcb2c918a4fba6833876b4d0d) chore(goreleaser): change compilation target to binary
* [a77aa88](https://github.com/zobtube/zobtube/commit/a77aa88686edbd43296e8685734675e23f577a27) chore(task): increase queue size to 1000
* [a77aa88](https://github.com/zobtube/zobtube/commit/a77aa88686edbd43296e8685734675e23f577a27) chore(ui): move sendToast to common 'main.js'
### Fixes
* [a77aa88](https://github.com/zobtube/zobtube/commit/a77aa88686edbd43296e8685734675e23f577a27) fix(session-cleaner): switch to new logger
* [a77aa88](https://github.com/zobtube/zobtube/commit/a77aa88686edbd43296e8685734675e23f577a27) fix(task): retry tasks stuck in 'todo' during boot (fix #62)
* [fbe685f](https://github.com/zobtube/zobtube/commit/fbe685fc4ba169c94e2c477f78b8314744ae7974) fix(adm): remove delete rows from the counting
### Documentation
* [4a467de](https://github.com/zobtube/zobtube/commit/4a467de5c2d9212bcb2c918a4fba6833876b4d0d) doc: add changelog on existing versions
* [d30d908](https://github.com/zobtube/zobtube/commit/d30d90889ec04ccc342c0890da636bb44c351d82) doc(readme): add stash postgres support
* [dd82526](https://github.com/zobtube/zobtube/commit/dd8252677f27348009823e9247c2eae549c0666d) doc(readme): fix typo
* [f2a153d](https://github.com/zobtube/zobtube/commit/f2a153df8ea17c06989cc1ec4877aeed9d6c502e) doc: prepare change for release 0.3.1

## Version 0.3.0
### Features
* [dde5d7f](https://github.com/zobtube/zobtube/commit/dde5d7f70a71368401a7eae4ffed864c4a49eb3c) feat: add configuration to disable authentication
* [dde5d7f](https://github.com/zobtube/zobtube/commit/dde5d7f70a71368401a7eae4ffed864c4a49eb3c) feat: add conistent logging with rs/zerolog
* [dde5d7f](https://github.com/zobtube/zobtube/commit/dde5d7f70a71368401a7eae4ffed864c4a49eb3c) feat: add password reset from cli
* [dde5d7f](https://github.com/zobtube/zobtube/commit/dde5d7f70a71368401a7eae4ffed864c4a49eb3c) feat(config): store part of configuration in database
* [dde5d7f](https://github.com/zobtube/zobtube/commit/dde5d7f70a71368401a7eae4ffed864c4a49eb3c) feat: disable gin debugging mode by default
* [dde5d7f](https://github.com/zobtube/zobtube/commit/dde5d7f70a71368401a7eae4ffed864c4a49eb3c) feat: improve onboarding experience
* [dde5d7f](https://github.com/zobtube/zobtube/commit/dde5d7f70a71368401a7eae4ffed864c4a49eb3c) feat(onboarding): remove user and config failsafe by creating defaults
### Chores
* [dde5d7f](https://github.com/zobtube/zobtube/commit/dde5d7f70a71368401a7eae4ffed864c4a49eb3c) chore(adm): split tasks in a dedicated page
* [dde5d7f](https://github.com/zobtube/zobtube/commit/dde5d7f70a71368401a7eae4ffed864c4a49eb3c) chore(config): switch from kelseyhightower/envconfig to urfave/cli
* [dde5d7f](https://github.com/zobtube/zobtube/commit/dde5d7f70a71368401a7eae4ffed864c4a49eb3c) chore(controller): refactor html rendering
### Fixes
* [269d9db](https://github.com/zobtube/zobtube/commit/269d9db7547ee6332aa0de35ae13fcd6079c3a41) fix(video/view): limit video suggestion to current type and avoid current video
* [3781bbc](https://github.com/zobtube/zobtube/commit/3781bbc7a9509b6b22fc58d502ef39065b1d220d) fix(failsafe): improve reloading
* [c060a5b](https://github.com/zobtube/zobtube/commit/c060a5b9e0f3f037ac56c1f3ce2e3edafcf56c70) fix(deps): update module github.com/gin-gonic/gin to v1.11.0

## Version 0.2.0
### Chores
* [8df17b9](https://github.com/zobtube/zobtube/commit/8df17b90ec35ff8f68b0f884bbcb992891d8b479) chore: change default ui port to 8069
### Fixes
* [5606eac](https://github.com/zobtube/zobtube/commit/5606eacd519cb87486919fde0f30b769512c537a) fix(doc): provide correct configuration example
* [b32fb75](https://github.com/zobtube/zobtube/commit/b32fb75db1a8d6972aef25bc54dbedf1bcb956a3) fix(docker-compose): remove badly copy-paste comments
### Documentation
* [4ebc631](https://github.com/zobtube/zobtube/commit/4ebc631e080401304fe24477302f9764c0aff0cf) doc: add stash and whisparr comparison
* [837218f](https://github.com/zobtube/zobtube/commit/837218fbacce7827bfe56a8d04a9d03979b6e162) doc: add docker compose example
* [ea33922](https://github.com/zobtube/zobtube/commit/ea3392217c8aec6462daa9446bcc5c3fd6d7246d) doc: add screenshots

## Version 0.1.70
### Features
* [0ddc5f5](https://github.com/zobtube/zobtube/commit/0ddc5f5828ddc49a588d61820839111e5771a48b) feat(docker): add latest tag

## Version 0.1.69
### Features
* [8cee4b6](https://github.com/zobtube/zobtube/commit/8cee4b6440930a3c6fb769527b9260f44ff422e3) feat(actor/new): redirect to edition after creation
* [fbe5759](https://github.com/zobtube/zobtube/commit/fbe5759b238b39028930880a5b9322f1321f9a61) feat(failsafe): reload page smoothly after restart
### Chores
* [16e1850](https://github.com/zobtube/zobtube/commit/16e1850dd983266456d2dc23755ea59a9f8f297b) chore: change license to MIT
### Fixes
* [328a06a](https://github.com/zobtube/zobtube/commit/328a06a0eaf905320768afa98f4862cad5e2997f) fix: rename sexuality to trans women

## Version 0.1.68
### Fixes
* [f10a5f6](https://github.com/zobtube/zobtube/commit/f10a5f618e5b43d1c8781bad90fd9adabcebb53c) fix(clip/view): set title and description properly
* [f44cabb](https://github.com/zobtube/zobtube/commit/f44cabbb7fdab2820723775702a1ac7d9fc9153b) fix(clip/view): add edit button
### Documentation
* [074646c](https://github.com/zobtube/zobtube/commit/074646c14a1cfefabdf9da0bfecd398ae340ad90) doc: update status and todo list
* [0fa1003](https://github.com/zobtube/zobtube/commit/0fa10039f4c785dccab3c81bcef30677fd0ae71c) doc: add security.md
* [248f6c8](https://github.com/zobtube/zobtube/commit/248f6c89a0d5a361b61a39b6663587672ba43bf0) doc: add contributor covenant code of conduct
* [fe31c09](https://github.com/zobtube/zobtube/commit/fe31c09e42141ee91d159f7558d2d2905db9723c) doc: add docker quickstart

## Version 0.1.67
### Features
* [4d3f87a](https://github.com/zobtube/zobtube/commit/4d3f87a1bf20d4b4a4b2755e76cf24a7563ab340) feat: implement clip view
* [959269c](https://github.com/zobtube/zobtube/commit/959269c2b313bd875553ef46ac556f37a538368e) feat(actor): allow renaming actors
* [aabefae](https://github.com/zobtube/zobtube/commit/aabefae8a680c3b315ad921a91869fd72ca02b64) feat: implement categories
### Chores
* [494a3fb](https://github.com/zobtube/zobtube/commit/494a3fbefaf7e6d775ede902fa6fedce872515f5) chore: add tests for new category routes
### Fixes
* [207daae](https://github.com/zobtube/zobtube/commit/207daaefcd5809438a3f24d42a342ecc40d84a92) fix(video/edit): fix function call after an initial video renaming
* [743b549](https://github.com/zobtube/zobtube/commit/743b54988c804af809cdc35dfc6e03ea4499264c) fix(actor/edit): remove actor-link-edition
### Documentation
* [5c95c06](https://github.com/zobtube/zobtube/commit/5c95c06ef0bd81fe19683fb2994b11472c4cb734) doc: update todo list with latest update

## Version 0.1.66
### Chores
* [21d8d16](https://github.com/zobtube/zobtube/commit/21d8d16c3a6fa793126db1864ab4b5cbc389e5b7) chore(deps): update dependency pytest-playwright to v0.7.1
### Fixes
* [2761231](https://github.com/zobtube/zobtube/commit/2761231820c1c2d9e538c6ba10d543b5aac5b17c) fix(deps): update module gorm.io/gorm to v1.31.0
* [dbb608d](https://github.com/zobtube/zobtube/commit/dbb608d4c094229da60978c5dc910ef4aed09a91) fix(deps): update module golang.org/x/image to v0.31.0

## Version 0.1.65
### Fixes
* [07c063d](https://github.com/zobtube/zobtube/commit/07c063ddffb5274d9aa7d60ed6066190deb4f726) fix(deps): update module gorm.io/gorm to v1.30.3
* [34dbfe2](https://github.com/zobtube/zobtube/commit/34dbfe2187e1f928a8ea04eae2932b61ba44a037) fix(deps): update module golang.org/x/text to v0.29.0
### Documentation
* [05dc347](https://github.com/zobtube/zobtube/commit/05dc3475899d653d3f5f13230f89f80fccc87521) doc: fix metadata format
* [1a4648a](https://github.com/zobtube/zobtube/commit/1a4648a0a547bcac2496cf564c2f167ccffbb54a) doc: add issue templates
* [7fae350](https://github.com/zobtube/zobtube/commit/7fae350802eefb6ac2640512bf058dce89c59d07) doc: add issue template configuration
* [bff4e65](https://github.com/zobtube/zobtube/commit/bff4e656f56dc8d2eb09a07bf0c68fc95a68c817) doc: split contributing

## Version 0.1.64
### Chores
* [092cc39](https://github.com/zobtube/zobtube/commit/092cc390afa44ff5ad51315648e8e83bf3214437) chore(deps): update dependency playwright to v1.55.0
* [4e6ee1f](https://github.com/zobtube/zobtube/commit/4e6ee1f024c16d1a02a971fcc7cf99fbec4a4d7f) chore(deps): update actions/setup-python action to v6
* [e614525](https://github.com/zobtube/zobtube/commit/e6145253d5d0a18dce91b3b0f9ff0859b78a61a6) chore(deps): update actions/setup-go action to v6
### Fixes
* [c389fe2](https://github.com/zobtube/zobtube/commit/c389fe296cf8d4ea5b0c52fb4d2c254d1014b25b) fix(deps): update module gorm.io/gorm to v1.30.2
* [d1a71ea](https://github.com/zobtube/zobtube/commit/d1a71ea3d3f2af1e8b7f17333cfe3e680490c86a) fix(deps): update module github.com/stretchr/testify to v1.11.1

## Version 0.1.63
### Fixes
* [2116769](https://github.com/zobtube/zobtube/commit/2116769f906cb117f441b1a6b51ac627765a9c3b) fix(adm/home): add category count
* [3a4a21c](https://github.com/zobtube/zobtube/commit/3a4a21c4cfc1329cf9baedcfed2a9fa73dba6608) fix(model/category-sub): add uuid annotation needed for postgres
* [b7eb032](https://github.com/zobtube/zobtube/commit/b7eb0327971663b8397c973ed1f35817dbb5afd6) fix(adm/category): use correct ajax urls

## Version 0.1.62
### Features
* [8b30202](https://github.com/zobtube/zobtube/commit/8b30202c8c9427863809c6a48491e391fc5c7e13) feat: add live reload
* [c6c3a4a](https://github.com/zobtube/zobtube/commit/c6c3a4a171a14d5a2cc119a73d00446c4986d85b) feat: implement categories
### Chores
* [84b679d](https://github.com/zobtube/zobtube/commit/84b679d30f5c96257de784023d53b66cc128504e) chore(deps): update actions/checkout action to v5
### Documentation
* [c4c3de8](https://github.com/zobtube/zobtube/commit/c4c3de817983f1d6fd603b32a555637b4cca9378) doc: add basic contributing guide

## Version 0.1.61
### Features
* [2d3a49c](https://github.com/zobtube/zobtube/commit/2d3a49c4d261bda8dc575ddf521b9ba8a6e520c4) feat(video/edit): add button to go back to the video viewer
### Fixes
* [1c6c90f](https://github.com/zobtube/zobtube/commit/1c6c90ff2f5a53105d7af7ab974e16162eed360a) fix(deps): update module golang.org/x/text to v0.28.0
* [4e16a22](https://github.com/zobtube/zobtube/commit/4e16a22b12d2efcd8522d791e70e9d8578f69249) fix(deps): update module golang.org/x/image to v0.30.0

## Version 0.1.60
### Chores
* [e91c014](https://github.com/zobtube/zobtube/commit/e91c014676c2661b531b9a027ab4611daa4bd933) chore(deps): update alpine docker tag to v3.22.1
* [ee8db8e](https://github.com/zobtube/zobtube/commit/ee8db8e0e8966efd5c39e9655733bc32caff92da) chore(deps): update dependency playwright to v1.54.0
### Fixes
* [6ae31c0](https://github.com/zobtube/zobtube/commit/6ae31c031d5b379a4c64365cafe55963cc94299b) fix(deps): update module gorm.io/gorm to v1.30.1

## Version 0.1.59
### Chores
* [1852669](https://github.com/zobtube/zobtube/commit/1852669795e8d4438bb847cf24e1cdc9c25af79e) chore(deps): update alpine docker tag to v3.22.0
* [b1a353e](https://github.com/zobtube/zobtube/commit/b1a353e45bcb1683bb74b062d9876498dd464864) chore(deps): update dependency playwright to v1.53.0
### Fixes
* [24cdc07](https://github.com/zobtube/zobtube/commit/24cdc079ebec453ec4d9e4cb102e71e5443c7f85) fix(deps): update module golang.org/x/image to v0.28.0
* [7d733a2](https://github.com/zobtube/zobtube/commit/7d733a28a49bb0eefbc12cb32ed5254f870fa569) fix(deps): update module golang.org/x/image to v0.29.0
* [85d6451](https://github.com/zobtube/zobtube/commit/85d645176e875d9b9064b9a40e9ceb5696ef7b24) fix(deps): update module golang.org/x/text to v0.27.0
* [ba8dafb](https://github.com/zobtube/zobtube/commit/ba8dafb4dc798b0de39a775ed16e5e057e5fd21c) fix(deps): update module gorm.io/driver/postgres to v1.6.0
* [c908e1d](https://github.com/zobtube/zobtube/commit/c908e1d36789e25c38f0deb71a00a654f52bc826) fix(deps): update module golang.org/x/text to v0.26.0
### Documentation
* [324b53e](https://github.com/zobtube/zobtube/commit/324b53ec2cf94fafd2deb4825d30d472d7209c27) doc: cleanup readme

## Version 0.1.58
### Features
* [c06e826](https://github.com/zobtube/zobtube/commit/c06e826a42b452427727c6fe6effd98e4ab9395e) feat(upload): add folder creation
### Chores
* [3060cdb](https://github.com/zobtube/zobtube/commit/3060cdba017c78c1dc5ae73cdd9254481f0aa114) chore(readme): cleanup bugs
### Fixes
* [66b5b1c](https://github.com/zobtube/zobtube/commit/66b5b1c9e96de77da476706bd0a8e37faa15e772) fix(video/edit): clean dead code on channel update
* [bb7faec](https://github.com/zobtube/zobtube/commit/bb7faecba1d00cb1b609cc4d7dbf32045dba0695) fix(video/edit): update video actor list properly
* [dae72fe](https://github.com/zobtube/zobtube/commit/dae72fe99310959e89a3a02f7a0d4341b172a7e6) fix(channel): add error message
* [dae72fe](https://github.com/zobtube/zobtube/commit/dae72fe99310959e89a3a02f7a0d4341b172a7e6) fix(channel): use correct form struct

## Version 0.1.57
### Fixes
* [6b7fe8f](https://github.com/zobtube/zobtube/commit/6b7fe8f25a518c72d9ac1143e288699570a6f34b) fix(channel): dipslay creation error

## Version 0.1.56
### Features
* [f791ffd](https://github.com/zobtube/zobtube/commit/f791ffda1dc3be184760224b007579331433a2b2) feat(channels): add crud + picture
* [f791ffd](https://github.com/zobtube/zobtube/commit/f791ffda1dc3be184760224b007579331433a2b2) feat(videos): change channels
### Fixes
* [66f0d51](https://github.com/zobtube/zobtube/commit/66f0d51dca9e5b137031ab21bf2a69e5845e7ef4) fix(controller): clean controller type

## Version 0.1.55
### Features
* [aa01e27](https://github.com/zobtube/zobtube/commit/aa01e27604209aef0956011197ef710a41257f4a) feat(task): add retry on error

## Version 0.1.54
### Features
* [5b48a07](https://github.com/zobtube/zobtube/commit/5b48a07392486ee1cdb62d1bcf821c5a7da3b0f7) feat(docker): add container source label

## Version 0.1.53
### Fixes
* [7f67662](https://github.com/zobtube/zobtube/commit/7f6766271202a8121d55dccf418342c456e5bf26) fix(triage): remove bug where double whitespaces would break video import
* [b4be99c](https://github.com/zobtube/zobtube/commit/b4be99c18debea95e60346037970dea990a36a99) fix(deps): update module gorm.io/gorm to v1.30.0
* [d0098f7](https://github.com/zobtube/zobtube/commit/d0098f7f9b8fd23ad0881d583f977a38c2fcdb27) fix(deps): update module github.com/gin-gonic/gin to v1.10.1

## Version 0.1.52
### Chores
* [18efb37](https://github.com/zobtube/zobtube/commit/18efb37a3a514806ad1a60298aeeb5a0e6c102bc) chore(deps): update dependency playwright to v1.52.0
* [358b80e](https://github.com/zobtube/zobtube/commit/358b80e004a245e585399f07d4837498c43da884) chore(deps): update golangci/golangci-lint-action action to v8
### Fixes
* [22f4772](https://github.com/zobtube/zobtube/commit/22f47724a66747460c165d9459e30a157dfeab61) fix(deps): update module gorm.io/gorm to v1.26.1
* [47742ce](https://github.com/zobtube/zobtube/commit/47742ce99f3e571be2aaafde673205696618ba7b) fix(deps): update module golang.org/x/image to v0.27.0

## Version 0.1.51
### Features
* [1c024b4](https://github.com/zobtube/zobtube/commit/1c024b4a812cef10837167ecede038fcda8987b8) feat: add user crud
* [3275d05](https://github.com/zobtube/zobtube/commit/3275d05c2ca0e63856c3802589130722ba2cc6b9) feat: retrieve build details from goreleaser
* [4235422](https://github.com/zobtube/zobtube/commit/4235422c91d8b7ee291c786f01feb5cd3fdf1e8a) feat: restrict non-admin user accesses
* [e1ee632](https://github.com/zobtube/zobtube/commit/e1ee632845b6d8f5c3dca2d3664e4238e33c5b5c) feat: add authentication tests
### Chores
* [57f4bf8](https://github.com/zobtube/zobtube/commit/57f4bf85cf4602dd5492979c5fdd13b86f6ed8eb) chore(deps): update actions/setup-python action to v5
* [8237433](https://github.com/zobtube/zobtube/commit/82374336e10e354766dc9837f78424442856e198) chore(deps): update dependency python to 3.13
* [c0f3627](https://github.com/zobtube/zobtube/commit/c0f36278607893944ab6e44aeffe4d943e91c011) chore: split html and js to ease readability
### Fixes
* [89da444](https://github.com/zobtube/zobtube/commit/89da444cd85dd3288e31883960de3761d5c0f9c0) fix(ci): rename ui test
### Documentation
* [1735884](https://github.com/zobtube/zobtube/commit/1735884b9aa570a4a9eaa36952c87584ea5ad538) doc: improve readme

## Version 0.1.50
### Features
* [b02b515](https://github.com/zobtube/zobtube/commit/b02b51575b642c04b5b5bf8b56c2fd73249c0379) feat(triage): add file deletion
### Chores
* [6a57635](https://github.com/zobtube/zobtube/commit/6a57635ef242376aee46b6feb5517be3b186dade) chore: bump to 0.1.50
### Fixes
* [8f94402](https://github.com/zobtube/zobtube/commit/8f9440249b363e173ce0edaf49c89a14400da5f1) fix(tools): stop air launching latest build on error
### Documentation
* [057c379](https://github.com/zobtube/zobtube/commit/057c379c2a25b443158dfc650d3a20f246c259f3) doc: ack actor alias editing

## Version 0.1.49
### Features
* [596b91c](https://github.com/zobtube/zobtube/commit/596b91c881f0f3ac351788cd8cac78ce4cccb89d) feat(actor): add alias crud
### Chores
* [3c40170](https://github.com/zobtube/zobtube/commit/3c4017032572fd8312841aee5b0b14cf66addeee) chore: bump to 0.1.49
### Documentation
* [8b13bfd](https://github.com/zobtube/zobtube/commit/8b13bfdbf9097738c18fc3a93a8ce28072472dc7) doc: add auth improvement

## Version 0.1.48
### Features
* [942d301](https://github.com/zobtube/zobtube/commit/942d30154a0888be9f451a61c5fec042f0cfcb7c) feat(adm): add task listing and view
### Chores
* [c21071b](https://github.com/zobtube/zobtube/commit/c21071b729d769d493d2c5f05e6c29e901d83582) chore: bump to 0.1.48
### Fixes
* [942d301](https://github.com/zobtube/zobtube/commit/942d30154a0888be9f451a61c5fec042f0cfcb7c) fix(task): fix done_at field not using correct type
* [942d301](https://github.com/zobtube/zobtube/commit/942d30154a0888be9f451a61c5fec042f0cfcb7c) fix(task): fix done_at not being set after task finish

## Version 0.1.47
### Features
* [a832745](https://github.com/zobtube/zobtube/commit/a832745d50855cc0db8876d354030bcb1489d664) feat(task): add video deletion
* [f5bc7ee](https://github.com/zobtube/zobtube/commit/f5bc7eea5ebf868e3d2b6ab9a20f2a09cec1eebc) feat(video): bring back the new thumbnail generation
### Chores
* [4b1c6b7](https://github.com/zobtube/zobtube/commit/4b1c6b797e8143bd06085f1d3841fcc99ad2874e) chore(deps): update golangci/golangci-lint-action action to v7
* [c9793dc](https://github.com/zobtube/zobtube/commit/c9793dc00350b2a76d5ee95752205c536980d16e) chore: bump to 0.1.47
### Fixes
* [a909345](https://github.com/zobtube/zobtube/commit/a90934525e2fa9f5dff22899e925666c958834c1) fix(deps): update module golang.org/x/text to v0.24.0
* [e08109e](https://github.com/zobtube/zobtube/commit/e08109e9f34121e649c4c946380a852761e6d846) fix(deps): update module golang.org/x/image to v0.26.0
* [f5bc7ee](https://github.com/zobtube/zobtube/commit/f5bc7eea5ebf868e3d2b6ab9a20f2a09cec1eebc) fix(video): remove a bug where video would be stuck in creating

## Version 0.1.46
### Features
* [c77c3cd](https://github.com/zobtube/zobtube/commit/c77c3cd3a2089b811c5085e2f422bd5c4ceea6f5) feat: add async tasks to perform on-disk actions
* [fd3c607](https://github.com/zobtube/zobtube/commit/fd3c607de695f62e07c357e82f2da47c05d8214b) feat: new clip listing
### Chores
* [0a06b2b](https://github.com/zobtube/zobtube/commit/0a06b2b6f7ae165b2a5a69a166959b583b2dca36) chore: bump to 0.1.46
* [2806fbe](https://github.com/zobtube/zobtube/commit/2806fbe003d1e922804caab66513824df57a21e3) chore(controller): split web and api files
### Fixes
* [2de5681](https://github.com/zobtube/zobtube/commit/2de5681458ca2255b103e266ea5d5a72ddc486a9) fix(docker): remove unused directory
* [6d830ec](https://github.com/zobtube/zobtube/commit/6d830ecbd9277cc586f99cda9ebc1687f58f9997) fix: use own favicon
* [c4d9f68](https://github.com/zobtube/zobtube/commit/c4d9f68ffb19d426f20a408af4f23a0179e145de) fix(auth): fix session issue when session did not have any user
### Documentation
* [2ef32a7](https://github.com/zobtube/zobtube/commit/2ef32a7483bd752bca6a3b4f81ecb550a2d23785) doc: remove bootstrap help from todo

## Version 0.1.45
### Features
* [0d6ec84](https://github.com/zobtube/zobtube/commit/0d6ec848dcfb3df7d3ac67b6300d1e3ce247c880) feat: add test coverage
* [d763fae](https://github.com/zobtube/zobtube/commit/d763fae555cf1846bc2f79502b8babe565a7059b) feat: list videos/actors/channels through admin pages
* [f952995](https://github.com/zobtube/zobtube/commit/f95299559bced5b904c85c1d3c9073d0ca3ff4eb) feat: add healthcheck test
* [f952995](https://github.com/zobtube/zobtube/commit/f95299559bced5b904c85c1d3c9073d0ca3ff4eb) feat: add test through CI
### Chores
* [7edec53](https://github.com/zobtube/zobtube/commit/7edec53911d866db8870c0df9a26b4447c486fff) chore: bump to 0.1.45
* [b09d91e](https://github.com/zobtube/zobtube/commit/b09d91edb161276ec5e7dd146ea189d1ed02735b) chore(model): migrate from interface to any
* [d01c474](https://github.com/zobtube/zobtube/commit/d01c4743ea35e38f93e7b15dc1f7678bd7f4e5f9) chore: add license
### Fixes
* [0db0d33](https://github.com/zobtube/zobtube/commit/0db0d3387614d6908cbef5d2f23fd57236601b57) fix(actor): allow manual input of links
* [324f21e](https://github.com/zobtube/zobtube/commit/324f21e78712e733239e0a1b2e1a54f7442bf7d3) fix(actor/edit): remove previous picture preview on modal dismiss
* [4d08124](https://github.com/zobtube/zobtube/commit/4d0812473b7e0484ba10d05f1a6460a264b41c5a) fix(actor): auto search links add display suggestions properly
* [65e0254](https://github.com/zobtube/zobtube/commit/65e0254486a19e58ebb5f2441218e6eaa0691ef2) fix(ci): add missing write for coverage
* [676b011](https://github.com/zobtube/zobtube/commit/676b01102bf9b16ace354341ad4b839142b41409) fix(actor): bring back deletion
* [78d2cf8](https://github.com/zobtube/zobtube/commit/78d2cf846084a6b8292eb2823b2a10ee63ce0395) fix(ci): remove coverage matrix
* [953e486](https://github.com/zobtube/zobtube/commit/953e486d4948c8caf36235ef0ccfaa2b58d7b281) fix(auth): remove non-nilness check
* [cfd6e03](https://github.com/zobtube/zobtube/commit/cfd6e03987b5fd59e61b378662f9a7547a9ddcb5) fix(actor/edit): fix bad zooming on picture selection
* [d9e6b42](https://github.com/zobtube/zobtube/commit/d9e6b421d5ec9402d0db1b895cc126cb7e7e63e2) fix: support thinner thumbnail
### Documentation
* [49d6db6](https://github.com/zobtube/zobtube/commit/49d6db690dfb873977d200fd7731e06b0a27990a) doc: update latest fix
* [dac6c8b](https://github.com/zobtube/zobtube/commit/dac6c8bd3fd8cf7276ef629d1f2b97f08b5af221) doc: reorder todo

## Version 0.1.44
### Features
* [4327b6b](https://github.com/zobtube/zobtube/commit/4327b6b0d6fb547dc3697bb5977b4d0b18485f71) feat: add bootstrapping and failsafe mode
* [840c6a5](https://github.com/zobtube/zobtube/commit/840c6a541dea433915b795d7fe57ed295af21a9c) feat: restart server after failsafe
* [a95b7f3](https://github.com/zobtube/zobtube/commit/a95b7f39cbd6f4babe7b0fdcfd7eaf18ed82a418) feat: create first admin user dynamically
* [bf30926](https://github.com/zobtube/zobtube/commit/bf30926a5eea314870e86f36b2fc8d1c61373ddf) feat: create library folder if missing at boot
### Chores
* [376912f](https://github.com/zobtube/zobtube/commit/376912f2ec05a78ea6030ff4c9ab7e66e13e65e4) chore: bump to 0.1.44
### Fixes
* [22b6b49](https://github.com/zobtube/zobtube/commit/22b6b496931ea2ac9c966029b6bd7bcd69b3b446) fix(video/view): bring back suggestions
### Documentation
* [d8e82e3](https://github.com/zobtube/zobtube/commit/d8e82e36b6eac3c5f4d81ceddb81b6c6fc13cfcb) doc: remove typo

## Version 0.1.43
### Features
* [e08f630](https://github.com/zobtube/zobtube/commit/e08f630a956535b12198b1fa1d6e17ffe481137c) feat: upgrade to go 1.24.1
### Chores
* [954e753](https://github.com/zobtube/zobtube/commit/954e7534f4987ab228a360fe1fd1e883902d5884) chore: bump to 0.1.43
* [98f2185](https://github.com/zobtube/zobtube/commit/98f21854d8806f02362518c12505156f36183caf) chore(deps): update docker/login-action action to v3
* [e08f630](https://github.com/zobtube/zobtube/commit/e08f630a956535b12198b1fa1d6e17ffe481137c) chore: upgrade all dependencies
### Fixes
* [44af17c](https://github.com/zobtube/zobtube/commit/44af17c30543b6c6bbbbe0168db16464769c5c70) fix(deps): update module github.com/google/uuid to v1.6.0
* [4cf24fc](https://github.com/zobtube/zobtube/commit/4cf24fc1f9c052fa2819353ae60e328efcfed6be) fix(deps): update module gorm.io/driver/postgres to v1.5.11
* [d3d210d](https://github.com/zobtube/zobtube/commit/d3d210dd75f831c137fb55b27ef29ef114a55e3f) fix(deps): update module gorm.io/gorm to v1.25.12
* [e08f630](https://github.com/zobtube/zobtube/commit/e08f630a956535b12198b1fa1d6e17ffe481137c) fix: golang.org/x/net to resolve security issue
* [e08f630](https://github.com/zobtube/zobtube/commit/e08f630a956535b12198b1fa1d6e17ffe481137c) fix: upgrade golang.org/x/crypto to resolve security issue
### Documentation
* [29cbd0d](https://github.com/zobtube/zobtube/commit/29cbd0d5eeb3c92ff4e508ede1739377b7ac26e0) doc: add onboarding readme

## Version 0.1.42
### Features
* [8f517ab](https://github.com/zobtube/zobtube/commit/8f517ab92e425843e4135c5e14fbb8a7b291135f) feat: add ci linting check
### Chores
* [27ea5ea](https://github.com/zobtube/zobtube/commit/27ea5ea66d6dc45917ba57ad7de9fe40eaf34fa2) chore: bump to 0.1.42
* [856b3ef](https://github.com/zobtube/zobtube/commit/856b3ef1627c3d1b4f25e1936ed37c4401c22368) chore: fix linting recommandations
### Fixes
* [13b0ad4](https://github.com/zobtube/zobtube/commit/13b0ad4cbef1d4134e57191e6cfc98b8d5876199) fix(github): rename workflow to avoid confusion
### Documentation
* [77d51a5](https://github.com/zobtube/zobtube/commit/77d51a54297c17c41b9d85b1658f021bb97051db) doc: add new bug

## Version 0.1.41
### Chores
* [0acb99c](https://github.com/zobtube/zobtube/commit/0acb99c461efd6a2715bba56ba8987a15290f588) chore: bump to 0.1.41
### Fixes
* [72cb796](https://github.com/zobtube/zobtube/commit/72cb7968694d97028a6163c5b94ffd9c75e81ad0) fix(github): update removed parameter

## Version 0.1.40
### Chores
* [d2b4335](https://github.com/zobtube/zobtube/commit/d2b433524814f4bb600cf71ed956af28840d3265) chore: bump to 0.1.40
### Fixes
* [76b1ab3](https://github.com/zobtube/zobtube/commit/76b1ab3f9e51bca927e45dc5c3edd9b6462358d2) fix(github): add missing permission to push container

## Version 0.1.39
### Chores
* [a1edb2e](https://github.com/zobtube/zobtube/commit/a1edb2eae2e60b94f32f21638b94732927d72fec) chore: bump to 0.1.39
### Fixes
* [c606aea](https://github.com/zobtube/zobtube/commit/c606aea3736df65645bc5ef70914ccb76bde0c28) fix(goreleaser): fix typo on registry name again

## Version 0.1.38
### Chores
* [1fcb15b](https://github.com/zobtube/zobtube/commit/1fcb15b582998285cb50f5762bb8f09cdac2e2bc) chore: bump to 0.1.38
### Fixes
* [004f2e6](https://github.com/zobtube/zobtube/commit/004f2e6b1ff071464e5ebfb6bf74cf301ffeb5d0) fix(goreleaser): fix typo on registry name

## Version 0.1.37
### Chores
* [71fe004](https://github.com/zobtube/zobtube/commit/71fe00427a6182e85149dafd5cdaa36548eab3a4) chore: bump to 0.1.37
### Fixes
* [3b97043](https://github.com/zobtube/zobtube/commit/3b9704341dc3c84f79e5ff5c74fbbfda38f930ef) fix(github): login before pushing to the registry

## Version 0.1.36
### Chores
* [adbe477](https://github.com/zobtube/zobtube/commit/adbe477f481ad040e055471847e03f6ed61819be) chore: migrate to github
* [efacb4e](https://github.com/zobtube/zobtube/commit/efacb4ebe6e8dddce3b5f26860a5050c9acb1922) chore: bump to 0.1.36
### Documentation
* [de3fd5e](https://github.com/zobtube/zobtube/commit/de3fd5edcf96b36253019eece3e9ff8ded9f4353) doc: add more known bugs and todos

## Version 0.1.35
### Chores
* [c8f425e](https://github.com/zobtube/zobtube/commit/c8f425ee667f966e700aa6cdb7341dd3efaca077) chore: bump to 0.1.35
### Fixes
* [ff4520c](https://github.com/zobtube/zobtube/commit/ff4520c1aa7c8d7bba9418772e865c01088ef7ab) fix(docker): add ffmpeg, fixing the thumbnail generation issue
* [ff4520c](https://github.com/zobtube/zobtube/commit/ff4520c1aa7c8d7bba9418772e865c01088ef7ab) fix(docker): use alpine as base, fixing the /tmp issue
### Documentation
* [1f888cb](https://github.com/zobtube/zobtube/commit/1f888cb40f6cb106c4f2d5eceb5d5c92bb3d1ad2) doc: add readme requirement
* [1f888cb](https://github.com/zobtube/zobtube/commit/1f888cb40f6cb106c4f2d5eceb5d5c92bb3d1ad2) doc: fix hierarchy
* [41436bc](https://github.com/zobtube/zobtube/commit/41436bc9d43a13be1314a8a7c4788523165b20ba) doc: update todo with release targets

## Version 0.1.34
### Chores
* [778bbb2](https://github.com/zobtube/zobtube/commit/778bbb242618e2875a3cff3319c372c270ee85e7) chore: bump to 0.1.34
### Fixes
* [778bbb2](https://github.com/zobtube/zobtube/commit/778bbb242618e2875a3cff3319c372c270ee85e7) fix(docker): fix chmod usage

## Version 0.1.33
### Chores
* [6e0deb5](https://github.com/zobtube/zobtube/commit/6e0deb55aa884511e248f99dbf635cd0c4027567) chore: bump to 0.1.33
### Fixes
* [6e0deb5](https://github.com/zobtube/zobtube/commit/6e0deb55aa884511e248f99dbf635cd0c4027567) fix(docker): use copy with chmod

## Version 0.1.32
### Chores
* [4aa11de](https://github.com/zobtube/zobtube/commit/4aa11de4ff93d8754488c7cbc451af50df4499ee) chore: bump to 0.1.32
### Fixes
* [4aa11de](https://github.com/zobtube/zobtube/commit/4aa11de4ff93d8754488c7cbc451af50df4499ee) fix(docker): add permission on /tmp

## Version 0.1.31
### Chores
* [b536259](https://github.com/zobtube/zobtube/commit/b53625993e90df52ffecbca28027eb7d308f3ccd) chore: bump to 0.1.31
### Fixes
* [b536259](https://github.com/zobtube/zobtube/commit/b53625993e90df52ffecbca28027eb7d308f3ccd) fix(docker): add missing -p flag

## Version 0.1.30
### Fixes
* [cd807b8](https://github.com/zobtube/zobtube/commit/cd807b837e911232a142e7704a7336bf11af7f0b) fix(docker): create /tmp folder from second layer

## Version 0.1.29
### Chores
* [ff6d11d](https://github.com/zobtube/zobtube/commit/ff6d11d044fe9cbc80099b19c200c03191f01131) chore: bump to 0.1.29
### Fixes
* [7ce2304](https://github.com/zobtube/zobtube/commit/7ce23048bc922f14a52ef34df923bfd04a4daa6f) fix: remove sql debug statements
* [ff6d11d](https://github.com/zobtube/zobtube/commit/ff6d11d044fe9cbc80099b19c200c03191f01131) fix(docker): create /tmp dir another way

## Version 0.1.28
### Chores
* [b6d6012](https://github.com/zobtube/zobtube/commit/b6d6012149965b3ccdbd6465e7c5afdc2be07459) chore: bump to 0.1.28
### Fixes
* [0d4cd62](https://github.com/zobtube/zobtube/commit/0d4cd62681073cb8ca5848a3c74c28f4bc7bd9ce) fix(triage): use correct error function
* [f73161e](https://github.com/zobtube/zobtube/commit/f73161e65dc4a8d0b4fcf81932716a3e20cb18b8) fix(docker): add missing /tmp needed for upload

## Version 0.1.27
### Chores
* [55d1de9](https://github.com/zobtube/zobtube/commit/55d1de93cab142cf89a0423babff1045b8f167a7) chore(triage): add error return on upload
* [c69e7f1](https://github.com/zobtube/zobtube/commit/c69e7f19be6ffbae6542021eb06280f5c417d038) chore: bump to 0.1.27
### Fixes
* [8b67844](https://github.com/zobtube/zobtube/commit/8b67844492fe1489e4658153f7e20f63523ae5e5) fix(common): remove deprecated css entry
* [e28a442](https://github.com/zobtube/zobtube/commit/e28a4425e2f2a48c8281039b73a7d32eeb6ad94f) fix(triage): refresh properly list after upload
* [e891d1b](https://github.com/zobtube/zobtube/commit/e891d1ba559027e90c7e791dbd4ef8064c885d48) fix(triage): remove debug button

## Version 0.1.26
### Features
* [b32bcd6](https://github.com/zobtube/zobtube/commit/b32bcd644e93a9130c683c40d72f9633c2c3ce19) feat: add view count
* [f26f013](https://github.com/zobtube/zobtube/commit/f26f013b544b0ef9f8d252b5e7a43db0cd4c27dc) feat: add profile page with most viewed videos/actors
### Chores
* [cd8c913](https://github.com/zobtube/zobtube/commit/cd8c91398404f91c9d7199aedf65e469f7634734) chore: bump to 0.1.26
### Fixes
* [c6fd354](https://github.com/zobtube/zobtube/commit/c6fd354d73558c6255f86ffa8a555d119740ad8a) fix(video): fix bad return on error

## Version 0.1.25
### Features
* [0696ccc](https://github.com/zobtube/zobtube/commit/0696cccae6ea5849f2f62bd61fd13d78474783b6) feat: rework triage/upload view
### Chores
* [6584228](https://github.com/zobtube/zobtube/commit/6584228a7fece7c742a8315d5fa9469e2a1af635) chore: revert to width-fixed display
* [844f689](https://github.com/zobtube/zobtube/commit/844f6892c7279c45fb2dc42d90e5d7d04e567143) chore: bump to 0.1.25
* [a7dbfa0](https://github.com/zobtube/zobtube/commit/a7dbfa0b3b60009ec9ddaca587d48c2dd8680487) chore: cleanup css

## Version 0.1.24
### Features
* [e05db34](https://github.com/zobtube/zobtube/commit/e05db3488dca2eabdfb89a159f66ee2e0a2d5468) feat: add movies back
### Chores
* [0afd1a0](https://github.com/zobtube/zobtube/commit/0afd1a0de6b661d9f6bd97dc638da05fe20dbabd) chore: bump to 0.1.24

## Version 0.1.23
### Chores
* [b8df14c](https://github.com/zobtube/zobtube/commit/b8df14cda1d729cfb6f6c72d5c71131c3e8338b5) chore: bump to 0.1.23
### Fixes
* [b8df14c](https://github.com/zobtube/zobtube/commit/b8df14cda1d729cfb6f6c72d5c71131c3e8338b5) fix(renovate): remove bad image name

## Version 0.1.22
### Chores
* [1ff375f](https://github.com/zobtube/zobtube/commit/1ff375f554aad9bcf5fe09316d8a7db82e16c6b7) chore: bump to 0.1.22
### Fixes
* [1ff375f](https://github.com/zobtube/zobtube/commit/1ff375f554aad9bcf5fe09316d8a7db82e16c6b7) fix(ci): login to registry before pushing

## Version 0.1.21
### Chores
* [14f9849](https://github.com/zobtube/zobtube/commit/14f9849d3fca8a3d2e9bafbcbe5c4334aa15cbc7) chore: bump to 0.1.21
### Fixes
* [5094805](https://github.com/zobtube/zobtube/commit/5094805451009c5a7f2a18bf933c3f56bb8c5ec6) fix(ci): pass registry credentials

## Version 0.1.20
### Chores
* [ac9f252](https://github.com/zobtube/zobtube/commit/ac9f252fab1d9c0d17708a8d43391e264c344b70) chore: bump to 0.1.20
### Fixes
* [8f3da5d](https://github.com/zobtube/zobtube/commit/8f3da5d90718f219e037d209dd6f4b11675ae170) fix(ci): remove renovate snapshot flag

## Version 0.1.19
### Chores
* [bfcdd07](https://github.com/zobtube/zobtube/commit/bfcdd07d76c6be237e2a1b5234642394feb6f1e2) chore: bump to 0.1.19
* [ded20e0](https://github.com/zobtube/zobtube/commit/ded20e0a272d7d6edd41821da3093998a48624bf) chore(renovate): update configuration to remove deprecation
### Fixes
* [33e3548](https://github.com/zobtube/zobtube/commit/33e35487dbaccc97a1e38ca5f6a93975f05ec183) fix(ci): switch build to docker
* [966624d](https://github.com/zobtube/zobtube/commit/966624deccec57c9550f5ae4dde5b8ec564d5922) fix(renovate): typo

## Version 0.1.18
### Chores
* [fc5a87c](https://github.com/zobtube/zobtube/commit/fc5a87c499e08c9937fc776e72f6bbe1df7fa1ff) chore: bump to 0.1.18
### Fixes
* [bcabf29](https://github.com/zobtube/zobtube/commit/bcabf29f606bc01fb0c70498a33e9a4c2bb9967d) fix(ci): fix docker build

## Version 0.1.17
### Features
* [d117442](https://github.com/zobtube/zobtube/commit/d11744238b8fbee0f7c643cf9fa671cd86c40c99) feat(ci): add docker images
### Chores
* [a2ff955](https://github.com/zobtube/zobtube/commit/a2ff9554d1c3610841ee464f6d4260f715ae446a) chore: bump to 0.1.17

## Version 0.1.16
### Chores
* [45cc38d](https://github.com/zobtube/zobtube/commit/45cc38d5d099b8f8f99d62b2923428f6bc746cb2) chore: bump to 0.1.16
### Fixes
* [224cac0](https://github.com/zobtube/zobtube/commit/224cac031f5ac8899c7479cf20725fe83f9d0b86) fix(deps): using proper logger url

## Version 0.1.15
### Chores
* [2755270](https://github.com/zobtube/zobtube/commit/2755270479b7fb76545a19626645772de7965c02) chore: bump to 0.1.15
### Fixes
* [0ede51f](https://github.com/zobtube/zobtube/commit/0ede51f34d5f618cdb6cb29177c245ebb110c6a5) fix(db): remove gorm logging
* [70069cc](https://github.com/zobtube/zobtube/commit/70069ccc3531f8f47c3186dd05ee98d8a48fd6aa) fix(auth): stop middleware when non-authed
* [d3fa7bb](https://github.com/zobtube/zobtube/commit/d3fa7bb6743bdf96120d45443e9885730f544a1b) fix(actor/list): add missing icon for actor adding
* [f03c98e](https://github.com/zobtube/zobtube/commit/f03c98e6fbbf348b9ffd458b567aeec39817b0ac) fix(auth): return 401 when unauthenticated

## Version 0.1.14
### Features
* [e45db76](https://github.com/zobtube/zobtube/commit/e45db76367d12473ad6af8b4d189a9a1b461daad) feat(video): display random video really randomly
### Chores
* [102c3c6](https://github.com/zobtube/zobtube/commit/102c3c654f51ad3651da6f2f7b65e4d16efadf93) chore: bump to 0.1.14
### Fixes
* [e45db76](https://github.com/zobtube/zobtube/commit/e45db76367d12473ad6af8b4d189a9a1b461daad) fix(video): remove removed shard

## Version 0.1.13
### Features
* [811b2dc](https://github.com/zobtube/zobtube/commit/811b2dcdacb9ff20e7c02f42a11a9bae52288824) feat: rework home ui
* [d9c67d2](https://github.com/zobtube/zobtube/commit/d9c67d2043e4df56cff69f3255f3981ed12900f7) feat(rework-ui): remove counters and other video types from home
### Chores
* [56be31a](https://github.com/zobtube/zobtube/commit/56be31a9706ee03962ffeb0db97e4584e64a34bc) chore: bump to 0.1.13

## Version 0.1.12
### Chores
* [3e27f1c](https://github.com/zobtube/zobtube/commit/3e27f1cf4f843ecbafd03645b82730b5d0607436) chore: bump to 0.1.12
* [c356a7b](https://github.com/zobtube/zobtube/commit/c356a7be6022165534f23ead9b8ada0a12a07daf) chore(upload): add webm support
### Fixes
* [70cc0cf](https://github.com/zobtube/zobtube/commit/70cc0cfd6027255ea52ca3ea0b4bc846e488d548) fix(web/upload): fix video import size
* [926b169](https://github.com/zobtube/zobtube/commit/926b16913270c778400815a2bc20a93a50aa7046) fix(web/upload): ease video preloading on preview
* [d19192a](https://github.com/zobtube/zobtube/commit/d19192acb52b705a8d8875504febee2daaa980f1) fix(controller/video): ensure database save
* [f730f95](https://github.com/zobtube/zobtube/commit/f730f95aa6f09120447397891255af214e57ec3e) fix(auth): check session to avoid null pointer
* [fb93a6e](https://github.com/zobtube/zobtube/commit/fb93a6e8daad9f63f818955159e3531a57db8341) fix(web/upload): fix race condition on duration computing

## Version 0.1.11
### Features
* [e5f452d](https://github.com/zobtube/zobtube/commit/e5f452dfdadfad6561f821616f1d36918d5a0f5e) feat(upload): import video from triage
### Chores
* [0822d43](https://github.com/zobtube/zobtube/commit/0822d43a7e9b4eba1b9e062eeda3a98419a538f5) chore(upload): refactor file icons
* [d835630](https://github.com/zobtube/zobtube/commit/d835630b5db55a9562a4ac0a7c09d7b7dc580d69) chore: bump to 0.1.11
### Fixes
* [0afabdc](https://github.com/zobtube/zobtube/commit/0afabdc6ad3a4bcb563f1a5f1f14cced1d5bf45e) fix(upload): display files once uploaded
* [50873b4](https://github.com/zobtube/zobtube/commit/50873b443df41088eaa7c1b4a36234f4d787c96b) fix(upload): remove dead code

## Version 0.1.10
### Features
* [45aa3b4](https://github.com/zobtube/zobtube/commit/45aa3b4b61f28dc26985b50fe6a46b3451e700d5) feat: decrease io pressure by lazy-loading
### Chores
* [3196e33](https://github.com/zobtube/zobtube/commit/3196e33b7d26c77aeaa2514dbcbe51e69ea6211d) chore: bump to 0.1.10
* [3516550](https://github.com/zobtube/zobtube/commit/3516550681e09695449c70bea03fd4cf70dbde30) chore(static): indicate lazyload version
* [35e34dd](https://github.com/zobtube/zobtube/commit/35e34ddc3361745209f9682b45c58a00775c7f82) chore(static): upgrade to bootstrap 5.3.3
* [6af9b19](https://github.com/zobtube/zobtube/commit/6af9b19ec9be3904893f9a42c3a97da41495d331) chore(static): download cropper.js 1.5.13
* [6f43224](https://github.com/zobtube/zobtube/commit/6f43224e5b08d452ca952449be84cb2c94eb8bbe) chore(static): remove bootstrap icons as unused
* [d93dc14](https://github.com/zobtube/zobtube/commit/d93dc142890c3df29cc4e53da67dc1424f08d2ca) chore(static): remove last external dependency by downloading poppins font
### Fixes
* [1a4dd88](https://github.com/zobtube/zobtube/commit/1a4dd88cee3c28bd1282e9dddfe5bc4f7e564e1b) fix(static): remove unused dependencies
* [22069a0](https://github.com/zobtube/zobtube/commit/22069a056950d97c8bc6f524f11c94bb8da4971e) fix(static): remove unused jscolor
* [3db32a7](https://github.com/zobtube/zobtube/commit/3db32a74e93d8a95da9f6c178e769521cf22c2c4) fix(static): add back megamenu as it was needed

## Version 0.1.9
### Features
* [33c2b53](https://github.com/zobtube/zobtube/commit/33c2b536df56f30e4b144d060fb8edf3dd167d4a) feat(video/edit): allow switching video type
* [6ad3745](https://github.com/zobtube/zobtube/commit/6ad3745249966d457fa4bd5be5aadc5f4045c439) feat(triage): improve ux
### Chores
* [5ff1475](https://github.com/zobtube/zobtube/commit/5ff147531c136fa3a10a1dfd7778c6b81169b2c4) chore(actor/edit): improve ux
* [70ba39f](https://github.com/zobtube/zobtube/commit/70ba39ff2656c8881d47a42d49cbac2697a28198) chore: remove legacy triage homepage
* [7bf30b5](https://github.com/zobtube/zobtube/commit/7bf30b5ae9e8f84bcd87ec9ed9ebed5882316604) chore: bump to 0.1.9

## Version 0.1.8
### Features
* [31408fd](https://github.com/zobtube/zobtube/commit/31408fd43ab26a2c0f3ab12fa51af90d3564ad45) feat(main): ensure all folder are present during boot
* [42b1a79](https://github.com/zobtube/zobtube/commit/42b1a79b67309997313fbf85c13b3392e43346ca) feat: rework triage view
* [4b21209](https://github.com/zobtube/zobtube/commit/4b2120943113f9f3550e8f7837f8ef0e5376a51e) feat(upload): allow uploading files in triage
### Chores
* [0bd8a52](https://github.com/zobtube/zobtube/commit/0bd8a52f0eb94458eb25bf39a3b37e0bd0980cb3) chore(doc): remove deployed feature
* [64688be](https://github.com/zobtube/zobtube/commit/64688be56d978c071ead27c22beeda18aabfb70a) chore: bump to 0.1.8
### Fixes
* [6fd76ba](https://github.com/zobtube/zobtube/commit/6fd76bad8ec593c9ceb3d30ecfa346eb30245e31) fix(auth): allow session with pg by setting usersession->userid as optional
### Documentation
* [cb8c040](https://github.com/zobtube/zobtube/commit/cb8c040a6f807d94c388a33d0da71127b303279b) doc: add a knwon bug on triage

## Version 0.1.7
### Chores
* [d5c39a8](https://github.com/zobtube/zobtube/commit/d5c39a8e1fb73accbfe4e678cc1f0e0aa59aac93) chore: bump to 0.1.7
### Fixes
* [3d4cb7b](https://github.com/zobtube/zobtube/commit/3d4cb7b3f68b4d4c5c14a89e5278393ef4505943) fix(auth/login): switch to sha256 method instead of crypto api to support http-only traffic

## Version 0.1.6
### Features
* [994c9d2](https://github.com/zobtube/zobtube/commit/994c9d23ef3efaa5637a232e8369440600ed3763) feat(provider/pornhub): support new profile picture page
* [c8aca57](https://github.com/zobtube/zobtube/commit/c8aca57fd00b7512dd1af21d5ef393b61857fb2d) feat: add authentication
### Chores
* [90b212a](https://github.com/zobtube/zobtube/commit/90b212ab21080d7085f0bd7d3e09d55cac738519) chore(actor/view): removed unused extra details
* [9e49cd8](https://github.com/zobtube/zobtube/commit/9e49cd8b13627692f01e0b24f4224d69cd67b387) chore: bump to 0.1.6
### Fixes
* [4c91469](https://github.com/zobtube/zobtube/commit/4c91469bb85bf0ae405741186ec5fdbd19f31d36) fix(dev): use correct env vars
* [90b212a](https://github.com/zobtube/zobtube/commit/90b212ab21080d7085f0bd7d3e09d55cac738519) fix(actor/view): display properly links
* [9c232fc](https://github.com/zobtube/zobtube/commit/9c232fc22f45d1b7848dca1d311141f16912dd23) fix(config): bring back support of env

## Version 0.1.5
### Features
* [929c507](https://github.com/zobtube/zobtube/commit/929c5071a5910d27f785fb83f5db43e2ec5c0076) feat(video/edit): allow deleting videos
### Chores
* [1bf216b](https://github.com/zobtube/zobtube/commit/1bf216b73bdfc11d26958ffd37020a7c45494b3c) chore: bump to 0.1.5
### Fixes
* [91fbf18](https://github.com/zobtube/zobtube/commit/91fbf18b4fc1c98b3e12cb5bff11344bbbea1977) fix(home): remove video type mixup

## Version 0.1.4
### Features
* [73ae891](https://github.com/zobtube/zobtube/commit/73ae891251409b3bcf8c8afb09257e5fa48084fa) feat(actor/edit): limit number of actors displayed to avoid dos
### Chores
* [6eadfe5](https://github.com/zobtube/zobtube/commit/6eadfe5c1de4c55ced68ef970bda964629c114dc) chore: bump to 0.1.4
* [e8bd371](https://github.com/zobtube/zobtube/commit/e8bd3716859b0c858dc6fedbca4d66dc29ac82ff) chore: backlog todo
### Fixes
* [55cf8d6](https://github.com/zobtube/zobtube/commit/55cf8d6063649e73bb96cf762b2786ae04fc1fa3) fix(video/edit): remove case where title would be set to undefined on renaming
* [73ae891](https://github.com/zobtube/zobtube/commit/73ae891251409b3bcf8c8afb09257e5fa48084fa) fix(video/edit): repair actor filtering

## Version 0.1.3
### Chores
* [002c014](https://github.com/zobtube/zobtube/commit/002c0141194221a7e945674d8975d75205e3e560) chore(video/list): remove limit as long as pagination is not implemented
* [d22f87f](https://github.com/zobtube/zobtube/commit/d22f87fc97e71971413ad3215a6675472a743744) chore: bump to 0.1.3
### Fixes
* [f0b4dec](https://github.com/zobtube/zobtube/commit/f0b4dec13a65453032c4c4f2d67ddbe00940fba3) fix(home): order by desc

## Version 0.1.2
### Chores
* [99e90c1](https://github.com/zobtube/zobtube/commit/99e90c1432edea705a8256657d5e2b96dd258f57) chore: bump to 0.1.2
* [9b50ffa](https://github.com/zobtube/zobtube/commit/9b50ffad3216f14993380e702d34581c9dbe994d) chore: rationalize video types
### Fixes
* [ab8a27c](https://github.com/zobtube/zobtube/commit/ab8a27c853382d4f30262b571f3df4f4a904c6a0) fix(home): limit properly video amount
* [ce6869a](https://github.com/zobtube/zobtube/commit/ce6869a6befb23a5a7b1789fa269a61f687a1b2a) fix(video/view): limit number of proposed videos
* [ce6869a](https://github.com/zobtube/zobtube/commit/ce6869a6befb23a5a7b1789fa269a61f687a1b2a) fix(video/view): remove debug flag

## Version 0.1.1
### Fixes
* [801cae3](https://github.com/zobtube/zobtube/commit/801cae3bce8d66d2a5094f536dcd44715370540f) fix(cfg): bind according to configuration

## Version 0.1.0
### Features
* [5b5475a](https://github.com/zobtube/zobtube/commit/5b5475a2e65c11956b76ec060b10908f4d45b63d) feat(http): bind port from configuration
* [b9dfaa1](https://github.com/zobtube/zobtube/commit/b9dfaa19905b4414b3124584a010620ca319bce6) feat: allow configuration through ci and yaml
### Documentation
* [e9144b4](https://github.com/zobtube/zobtube/commit/e9144b49ffdeec1b7546d54bd13f49afa832a2fa) doc: increase todo

## Version 0.0.3
### Fixes
* [bc434c8](https://github.com/zobtube/zobtube/commit/bc434c8a35849db30e86f9d5482eda9c1309eb56) fix: embed static and template assets

## Version 0.0.2
### Features
* [eb6eb80](https://github.com/zobtube/zobtube/commit/eb6eb80a94a71e03517b3b0986ff2609c826771b) feat: add ci to release on tag

