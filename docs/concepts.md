# Concepts

## Generics

ZobTube allows sorting videos into several kinds:

- **Videos**: Medium-length videos
- **Movies**: Long videos, that usually contains several scenes
- **Clips**: Short videos, easily scrollable, with a TikTok-like viewer

Each **video** can be associated to **actors**, **channels**, and **categories**.

**Channels** can be used to match production companies (like __Brazzers__) or creator feed (like a content creator on __OnlyFans__).

**Categories** are fully customizable, allowing them to reflect personal preferences.

**Actors** informations can be retrieved from well-known platforms:

- Babepedia
- Babes Directory
- Boobpedia
- PornHub
- IAFD

__If you think a provider is missing and could be integrated, feel free to [create a feature request](https://github.com/zobtube/zobtube/issues/new/choose)__

## Triage


ZobTube is a tool aiming to help sorting videos. To achieve this, all videos are initially stored in the `/triage` folder. You can place them directly inside it or upload them through the upload window.

A video in triage can then be imported (by selected its kind: video, movie or clip).

Once imported, the video will be moved a new dedicated folder (depending on its kind), the duration will be computed (using `ffmpeg`) and a thumbnail will be generated (also using `ffmpeg`).

The video will then be editable, allowing to add **channels**, **actors** and **categories**.