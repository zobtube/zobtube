# Category Selection

This shards exposes two HTML parts for the category selection.

## UI element: listing

Defined as `shards/category-selection/list.html`.

Usable with:

```gotemplate
{{ template "shards/category-selection/list.html" . }}
```

Will require:
- `.Categories` to be defined with a list of categories already selected

## UI element: modal for edition

Defined as `shards/category-selection/edit-modal.html`.

Usable with:

```gotemplate
{{ template "shards/category-selection/edit-modal.html" . }}
```

Will require:
- `.Categories` to be defined with a list of selectable categories

## JS element

Defined as `shards/category-selection/main.js`.

Usage with:

```gotemplate
{{ template "shards/category-selection/main.js" . }}
```

Exposed variables:
- `categorySelectable` containing all selectable categories
- `categorySelected` containing all selected categories

Exposed methods:
- `onCategorySelectBefore(category_id)` called before an category selection is accepted
- `onCategorySelectAfter(category_id)` called after an category selection is accepted
- `onCategoryDeselectBefore(category_id)` called before an category deselection is accepted
- `onCategoryDeselectAfter(category_id)` called after an category deselection is accepted
- `onModalHide(calling_event)` called on modal closing
