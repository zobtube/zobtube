{{ define "shards/category-selection/main.js" }}
window.zt.categorySelection = Object({
  categorySelectable: {},
  categorySelected: {},

  onCategorySelectBefore() { return true; },
  onCategorySelectAfter() { return true; },

  onCategoryDeselectBefore() { return true; },
  onCategoryDeselectAfter() { return true; },

  onModalHide() { return; },

  _onModalHide(e) {
    return this.onModalHide(e);
  },

  categorySelect(category_id) {
    // before
    res = this.onCategorySelectBefore(category_id);
    if (!res) {
      console.error('failed to call to function onCategorySelectBefore');
      return
    }

    // perform selection
    this.categorySelected[category_id] = undefined;
    this._updateSelectedCategories();

    // after
    res = this.onCategorySelectAfter(category_id);
    if (!res) {
      console.error('failed to call to function onCategorySelectAfter');
      return
    }
  },

  categoryDeselect(category_id) {
    // before
    res = this.onCategoryDeselectBefore(category_id);
    if (!res) {
      console.error('failed to call to function onCategoryDeselectBefore');
      return
    }

    // perform deselection
    delete this.categorySelected[category_id];
    this._updateSelectedCategories();

    // after
    res = this.onCategoryDeselectAfter(category_id);
    if (!res) {
      console.error('failed to call to function onCategoryDeselectAfter');
      return
    }
  },


  // will update the list and the edit modal with selected categories
  _updateSelectedCategories() {
    console.debug('update categories presents in video');
    // categories presents in the video
    var categories_chips = document.getElementsByClassName('video-category-list');
    for (const category_chip of categories_chips) {
      category_id = category_chip.getAttribute('category-id');
      if (category_id in this.categorySelected) {
        category_chip.style.display = '';
      } else {
        category_chip.style.display = 'none';
      }
    }

    // categories available for category add/removal modal
    var categories_chips = document.getElementsByClassName('add-category-list');
    for (const category_chip of categories_chips) {
      category_id = category_chip.getAttribute('category-id');
      if (category_id in this.categorySelected) {
        category_chip.querySelector('.btn-success').style.display = 'none';
        category_chip.querySelector('.btn-danger').style.display = '';
      } else {
        category_chip.querySelector('.btn-success').style.display = '';
        category_chip.querySelector('.btn-danger').style.display = 'none';
      }
    }
  },

});

window.zt.onload.push(function() {
  // create modal object
  this.Modal = document.getElementById('categorySelectionModal');
  this.ModalBS = new bootstrap.Modal(this.Modal);

  // handle modal closing
  this.Modal.addEventListener('hidden.bs.modal', window.zt.categorySelection._onModalHide.bind(window.zt.categorySelection));

  // toggle right button in modal
  window.zt.categorySelection._updateSelectedCategories();
});

{{ end }}
