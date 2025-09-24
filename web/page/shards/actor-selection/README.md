# Actor Selection

This shards exposes two HTML parts for the actor selection.

## UI element: listing

Defined as `shards/actor-selection/list.html`.

Usable with:

```gotemplate
{{ template "shards/actor-selection/list.html" . }}
```

Will require:
- `.Actors` to be defined with a list of actors already selected

## UI element: modal for edition

Defined as `shards/actor-selection/edit-modal.html`.

Usable with:

```gotemplate
{{ template "shards/actor-selection/edit-modal.html" . }}
```

Will require:
- `.Actors` to be defined with a list of selectable actors

## JS element

Defined as `shards/actor-selection/main.js`.

Usage with:

```gotemplate
{{ template "shards/actor-selection/main.js" . }}
```

Exposed variables:
- `actorSelectable` containing all selectable actors
- `actorSelected` containing all selected actors

Exposed methods:
- `onActorSelectBefore(actor_id)` called before an actor selection is accepted
- `onActorSelectAfter(actor_id)` called after an actor selection is accepted
- `onActorDeselectBefore(actor_id)` called before an actor deselection is accepted
- `onActorDeselectAfter(actor_id)` called after an actor deselection is accepted
- `onModalHide(calling_event)` called on modal closing
