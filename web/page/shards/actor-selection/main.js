{{ define "shards/actor-selection/main.js" }}
window.zt.actorSelection = Object({
  actorSelectable: {},
  actorSelected: {},

  onActorSelectBefore() { return true; },
  onActorSelectAfter() { return true; },

  onActorDeselectBefore() { return true; },
  onActorDeselectAfter() { return true; },

  onModalHide() { return; },

  _onModalHide(e) {
    return this.onModalHide(e);
  },

  actorSelect(actor_id) {
    // before
    res = this.onActorSelectBefore(actor_id);
    if (!res) {
      console.error('failed to call to function onActorSelectBefore');
      return
    }

    // perform selection
    this.actorSelected[actor_id] = undefined;
    this._updateSelectedActors();

    // after
    res = this.onActorSelectAfter(actor_id);
    if (!res) {
      console.error('failed to call to function onActorSelectAfter');
      return
    }
  },

  actorDeselect(actor_id) {
    // before
    res = this.onActorDeselectBefore(actor_id);
    if (!res) {
      console.error('failed to call to function onActorDeselectBefore');
      return
    }

    // perform deselection
    delete this.actorSelected[actor_id];
    this._updateSelectedActors();

    // after
    res = this.onActorDeselectAfter(actor_id);
    if (!res) {
      console.error('failed to call to function onActorDeselectAfter');
      return
    }
  },


  // will update the list and the edit modal with selected actors
  _updateSelectedActors() {
    console.debug('update actors presents in video');
    // actors presents in the video
    var actors_chips = document.getElementsByClassName('video-actor-list');
    for (const actor_chip of actors_chips) {
      actor_id = actor_chip.getAttribute('actor-id');
      if (actor_id in this.actorSelected) {
        actor_chip.style.display = '';
      } else {
        actor_chip.style.display = 'none';
      }
    }

    // actors available for actor add/removal modal
    var actors_chips = document.getElementsByClassName('add-actor-list');
    for (const actor_chip of actors_chips) {
      actor_id = actor_chip.getAttribute('actor-id');
      if (actor_id in this.actorSelected) {
        actor_chip.querySelector('.btn-success').style.display = 'none';
        actor_chip.querySelector('.btn-danger').style.display = '';
      } else {
        actor_chip.querySelector('.btn-success').style.display = '';
        actor_chip.querySelector('.btn-danger').style.display = 'none';
      }
    }
  },

  _updateModalDisplayedActors(filter = '') {
    var re = new RegExp(filter, 'i');
    var actors_chips = document.getElementsByClassName('add-actor-list');
    found_count = 0;
    found_count_max = 15;
    for (const actor_chip of actors_chips) {
      actor_id = actor_chip.getAttribute('actor-id');
      found = re.test(this.actorSelectable[actor_id]['name']);
      if (found) {
        found_count++;
      } else {
        for (const a of this.actorSelectable[actor_id]['aliases']) {
          if (re.test(a)) {
            found = true;
            found_count++;
            break;
          }
        }
      }
      if (found && found_count < found_count_max) {
        actor_chip.style.display = 'none';
        actor_chip.style.display = '';
      } else {
        actor_chip.style.display = '';
        actor_chip.style.display = 'none';
      }
    }
  },
});

window.zt.onload.push(function() {
  // create modal object
  this.Modal = document.getElementById('actorSelectionModal');
  this.ModalBS = new bootstrap.Modal(this.Modal);

  // handle modal closing
  this.Modal.addEventListener('hidden.bs.modal', window.zt.actorSelection._onModalHide.bind(window.zt.actorSelection));

  // toggle right button in modal
  window.zt.actorSelection._updateSelectedActors();

  // add event on modal input
  document.getElementById('actorSelectionModalInput').addEventListener('input', function(e) {
    window.zt.actorSelection._updateModalDisplayedActors(e.target.value);
  });

  // display actors according to filter
  window.zt.actorSelection._updateModalDisplayedActors('');
});

{{ end }}
