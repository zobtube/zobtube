{{ define "actor/edit.js" }}

// store all categories at start
var categories_all = {
  {{ range $category := .Categories }}
  {{ range $sub:= $category.Sub }}
  '{{ $sub.ID }}': {
    'name': '{{ $sub.Name }}',
  },
  {{ end }}
  {{ end }}
};

// store all categories in the video at start
var categories_in_actor = {
  {{ range $sub := .Actor.Categories }}
  '{{ $sub.ID }}': undefined,
  {{ end }}
};

const G_OfflineMode = {{ .OfflineMode }};

// list all known links
G_links = {
{{ range $link := .Actor.Links }}
  '{{ $link.Provider }}': {
    'link_url': '{{ $link.URL }}',
    'link_id': '{{ $link.ID }}',
  },
{{ end }}
};

// list all providers
G_providers = {
  {{ range $provider := .Providers }}
  '{{ $provider.ID }}': {
    'NiceName': '{{ $provider.NiceName }}',
    'Enabled': {{ $provider.Enabled }},
  },
{{ end }}
};

G_first_time = {{ if eq ( len .Actor.Links ) 0 }}true{{ else }}false{{ end }};

G_actor_id = "{{ .Actor.ID }}";
/* -- end globals set at boot -- */

/* -- start of categories scripts -- */
function actorUpdateCategoryStatus() {
  console.log('update categories presents for actor');
  // categories presents for the actor
  var categoriess_chips = document.getElementsByClassName('actor-category-list');
  for (const category_chip of categoriess_chips) {
    category_id = category_chip.getAttribute('category-id');
    if (category_id in categories_in_actor) {
      category_chip.style.display = '';
    } else {
      category_chip.style.display = 'none';
    }
  }

  // categories available for category add/removal modal
  var categories_chips = document.getElementsByClassName('add-category-list');
  for (const category_chip of categories_chips) {
    category_id = category_chip.getAttribute('category-id');
    if (category_id in categories_in_actor) {
      category_chip.querySelector('.btn-success').style.display = 'none';
      category_chip.querySelector('.btn-danger').style.display = '';
    } else {
      category_chip.querySelector('.btn-success').style.display = '';
      category_chip.querySelector('.btn-danger').style.display = 'none';
    }
  }
}

function actorAddCategory(category_id) {
  console.log('add category'+category_id+' for actor '+actor_id);
  url = url_actor_category_edit.replace('00000000-0000-0000-0000-000000000000', category_id);
  $.ajax(url, {
    method: 'PUT',

    xhr: function () {
      var xhr = new XMLHttpRequest();
      return xhr;
    },

    success: function (res) {
      console.debug('success, got', res);
      sendToast('Category added', '', 'bg-success', categories_all[category_id]['name']+' added.');
      categories_in_actor[category_id] = undefined;
      actorUpdateCategoryStatus();
    },

    error: function () {
      console.debug('failed');
      sendToast('Category not added', '', 'bg-danger', categories_all[category_id]['name']+' not added, call failed.');
    },
  });
}

function actorRemoveCategory(category_id) {
  console.log('remove category '+category_id+' from actor '+actor_id);
  url = url_actor_category_edit.replace('00000000-0000-0000-0000-000000000000', category_id);
  $.ajax(url, {
    method: 'DELETE',

    xhr: function () {
      var xhr = new XMLHttpRequest();
      return xhr;
    },

    success: function (res) {
      console.debug('success, got', res);
      sendToast('Category removed', '', 'bg-success', categories_all[category_id]['name']+' removed.');
      delete categories_in_actor[category_id];
      actorUpdateCategoryStatus();
    },

    error: function () {
      console.debug('failed');
      sendToast('Category not removed', '', 'bg-danger', categories_all[category_id]['name']+' not removed, call failed.');
    },
  });
}


/* -- end of categories scripts -- */

function actorLinkAutomaticSearch(actorName, providerName, providerSlug, url) {
  sendToast(
    'Automatic actor search',
    '',
    'bg-info',
    'Searching for '+actorName+' on '+providerName,
  );

  $.ajax(url, {
    method: 'GET',
    xhr: function () {
      var xhr = new XMLHttpRequest();
      return xhr;
    },
    success: function (result) {
      sendToast(
        'Automatic actor search',
        '',
        'bg-success',
        actorName+' found on '+providerName+'!',
      );
      G_links[providerSlug] = {
        "link_url": result.link_url,
        "link_id": result.link_id,
      };
      updateProviders();
    },
    error: function () {
      sendToast(
        'Automatic actor search',
        '',
        'bg-warning',
        actorName+' not found on '+providerName,
      );
    },
    complete: function() {
      document.providerInit[providerSlug] = true;
    },
  });

}

function updateProviders() {
  for (const provider_slug in G_providers) {
    if (provider_slug in G_links) {
      link_url = G_links[provider_slug].link_url;
      link_id = G_links[provider_slug].link_id;
      providerLinkPresent(provider_slug, link_url, link_id);
      continue;
    }

    providerLinkAbsent(provider_slug);
  }
}

function providerFirstTime(provider_slug) {
  // hide buttons
  document.getElementById('provider-action-'+provider_slug+'-view').style.display = 'none';
  document.getElementById('provider-action-'+provider_slug+'-delete').style.display = 'none';
  document.getElementById('provider-action-'+provider_slug+'-search').style.display = 'none';
  document.getElementById('provider-action-'+provider_slug+'-add').style.display = 'none';

  // show first time text
  document.getElementById('provider-text-'+provider_slug+'-first-time').style.display = '';

}

function providerLinkPresent(provider_slug, link_url, link_id) {
  // hide search buttons
  document.getElementById('provider-action-'+provider_slug+'-search').style.display = 'none';
  document.getElementById('provider-action-'+provider_slug+'-add').style.display = 'none';
  document.getElementById('provider-text-'+provider_slug+'-first-time').style.display = 'none';
  // show edit ones
  btnView = document.getElementById('provider-action-'+provider_slug+'-view');
  btnView.style.display = '';
  btnDelete = document.getElementById('provider-action-'+provider_slug+'-delete');
  btnDelete.style.display = '';

  // update urls
  btnView.href = link_url;
}

function linkDelete(provider_slug) {
  link_id = G_links[provider_slug].link_id;
  url = "/api/actor/link/"+link_id;
  $.ajax(url, {
    method: 'DELETE',
    xhr: function () {
      var xhr = new XMLHttpRequest();
      return xhr;
    },
    success: function (result) {
      sendToast(
        'Actor\'s link deleted',
        '',
        'bg-success',
        'Success',
      );
      delete G_links[provider_slug];
      updateProviders();
    },
    error: function () {
      sendToast(
        'Actor\'s link deletion',
        '',
        'bg-warning',
        'Unable to delete',
      );
    },
  });
}

function providerLinkAbsent(provider_slug) {
  // hide buttons
  document.getElementById('provider-action-'+provider_slug+'-view').style.display = 'none';
  document.getElementById('provider-action-'+provider_slug+'-delete').style.display = 'none';
  document.getElementById('provider-text-'+provider_slug+'-first-time').style.display = 'none';
  // show search ones
  document.getElementById('provider-action-'+provider_slug+'-search').style.display = '';
  document.getElementById('provider-action-'+provider_slug+'-add').style.display = '';
}

function firstTimeSearch() {
  if (!G_first_time) {
    return
  }

  console.debug('First time running, discoverying providers');

  for (const provider_slug in G_providers) {
    if (!(provider_slug in G_links) || G_links[provider_slug].link_url == undefined) {
      if (!G_providers[provider_slug]['Enabled']) {
        continue;
      }
      document.getElementById('provider-action-'+provider_slug+'-search').click();
    }
  }
}

updateProviders();

window.onload = function() {
  // init vars
  document.providerInit = {};

  // store modals
  document.modalNewAlias = document.getElementById("addActorAliasModal");
  document.modalNewAliasModal = new bootstrap.Modal(document.modalNewAlias);

  const previewModal = document.getElementById('profilePictureModal');
  previewModal.addEventListener('shown.bs.modal', previewModalShown);
  previewModal.addEventListener('hide.bs.modal', event => {
    cropper.destroy();
  });

  // register image onload
  document.getElementById('image').onload = function() {
    gen_new_crop();
  }

  // register calls on modals
  $("#addActorAliasInput").on('keyup', function (e) {
      if (e.key === 'Enter' || e.keyCode === 13) {
        actorAliasAdd();
      }
});

  // calls
  if (!G_OfflineMode) {
    firstTimeSearch();
    suggestProfilePicture();
  }

  // update categories
  actorUpdateCategoryStatus();
}

function handlePictureError(event) {
  if (G_OfflineMode) {
    event.target.src = '/static/images/actor-picture-error-offline-mode.svg';
  } else {
    event.target.src = '/static/images/actor-picture-error-not-found.svg';
  }
}

function showActorPictures() {
  const profilePictureURL = "/api/actor/link/00000000-0000-0000-0000-000000000000/thumb";
  template = document.getElementById('profile-picture-template');

  document.getElementById('profile-picture-propositions-row').style.display = '';

  for (const provider_slug in G_links) {
    if (!G_providers[provider_slug]['Enabled']) {
      continue;
    }
    link_id = G_links[provider_slug].link_id;

    document.getElementById('profile-picture-suggestion-welcome').style.display = 'none';
    document.getElementById('profile-picture-suggestion-click-info').style.display = '';

    url = profilePictureURL.replace('00000000-0000-0000-0000-000000000000', link_id);

    newPicture = template.cloneNode(true);
    newPicture.style.display = '';
    newPicture.id = 'profile-picture-provider-'+provider_slug;
    newPicture.querySelector('img').addEventListener('error', handlePictureError);
    newPicture.querySelector('img').src = url;
    newPicture.querySelector('.card-footer').innerText = G_providers[provider_slug]['NiceName'];
    newPicture.onclick = function() {
      const providerSlug = provider_slug;
      console.log('picture from '+provider_slug);
      console.log(this);
    }

    propositionRoot = document.getElementById('profile-picture-propositions');
    propositionRoot.appendChild(newPicture);
  }
}

function suggestProfilePicture() {
  if ({{ if .Actor.Thumbnail }}true{{ else }}false{{ end }}) {
    return;
  }

  // if first time, wait for init
  if (G_first_time) {
    for (provider_slug in G_providers) {
      if (!(provider_slug in document.providerInit)) {
        console.debug('provider '+provider_slug+' not init, waiting');
        setTimeout(suggestProfilePicture, "1000");
        return
      }
    }
  }

  showActorPictures();
}

// set profile picture modal logic
function previewModalShown(event) {
  // update status
  document.getElementById('cropper-status').innerText = 'Cropper library loading...';
  // set image src to caller one
  document.getElementById('image').src = event.relatedTarget.src;
}

var cropper;
function gen_new_crop() {
  var image = document.querySelector('#image');
  var minAspectRatio = 1;
  var maxAspectRatio = 1;
  cropper = new Cropper(image, {
    viewMode: 2,
    initialAspectRatio: 1,
    ready: function () {
      var cropper = this.cropper;
      var containerData = cropper.getContainerData();
      var cropBoxData = cropper.getCropBoxData();
      var size = Math.min(image.width, image.height);

      cropper.setCropBoxData({
        top: 0,
        left: 0,
        width: size,
        height: size,
      });
      cropper.moveTo(0);
      document.getElementById('cropper-status').innerText = 'Cropper library ready';
    },

    cropmove: function () {
      var cropper = this.cropper;
      var cropBoxData = cropper.getCropBoxData();
      var aspectRatio = cropBoxData.width / cropBoxData.height;

      if (aspectRatio < minAspectRatio) {
        cropper.setCropBoxData({
            width: cropBoxData.height * minAspectRatio
        });
      } else if (aspectRatio > maxAspectRatio) {
        cropper.setCropBoxData({
          width: cropBoxData.height * maxAspectRatio
        });
      }
    },
  });
}

function send_crop() {
  cropper.getCroppedCanvas().toBlob(function (blob) {
    var formData = new FormData();

    formData.append('pp', blob);
    formData.append('csrfmiddlewaretoken', '{% csrf_token %}');
    $.ajax("/api/actor/"+G_actor_id+"/thumb", {
      method: 'POST',
      data: formData,
      processData: false,
      contentType: false,

      xhr: function () {
        var xhr = new XMLHttpRequest();

        return xhr;
      },

      success: function () {
        console.debug('success, redirecting');
        window.location.href = "/actor/"+G_actor_id+"/edit"; //TODO: add in-page reload
      },

      error: function () {
        console.debug('failed');
      },

      complete: function () {
      },
    });
  });
}

function addNewLinkModal(actorID, providerSlug) {
  console.debug("show modal to add link for actor "+actorID+" on "+providerSlug)
  if (document.modalNewLink == undefined) {
    document.modalNewLink = document.getElementById("manualLinkModal")
    document.modalNewLinkModal = new bootstrap.Modal(document.modalNewLink)
  }
  document.newLinkActorID = actorID;
  document.newLinkProviderSlug = providerSlug;
  document.getElementById("newLinkInput").value = '';
  document.modalNewLinkModal.show();
}

function setNewLink() {
  newURL = document.getElementById("newLinkInput").value;

  $.ajax('/api/actor/'+document.newLinkActorID+'/link', {
    method: 'POST',
    data: {
      'url': newURL,
      'provider': document.newLinkProviderSlug,
    },
    xhr: function () {
      var xhr = new XMLHttpRequest();
      return xhr;
    },
    success: function (result) {
      sendToast(
        'New link added for {{ .Actor.Name }}',
        '',
        'bg-success',
        'Success',
      );
      G_links[document.newLinkProviderSlug] = {
        "link_url": result.link_url,
        "link_id": result.link_id,
      };
      updateProviders();
      document.modalNewLinkModal.hide();
    },
    error: function () {
      sendToast(
        'Automatic actor search',
        '',
        'bg-warning',
        'Adding failed',
      );
    },
  });
}

function actorAliasAdd() {
  console.log('add new alias for actor');

  // retrieve button
  button = document.getElementById('modalActorAliasButton');

  // set button as loader
  button.classList.add('spinner-border', 'text-secondary');
  button.classList.remove('btn', 'btn-success');
  button.innerHTML = '<span class="visually-hidden">Loading...</span>';

  // store input
  aliasInput = document.getElementById('addActorAliasInput');

  // perform http call
  url = '/api/actor/'+G_actor_id+'/alias';
  $.ajax(url, {
    method: 'POST',
    data: {
      'alias': aliasInput.value,
    },
    xhr: function () {
      var xhr = new XMLHttpRequest();
      return xhr;
    },
    success: function (result) {
      sendToast(
        'Actor\'s new alias',
        '',
        'bg-success',
        'New alias recorded!',
      );

      // add new alias
      newAliasChip  = document.createElement('div');
      newAliasChip.classList.add('chip');

      newAliasChipContent  = aliasInput.value;
      newAliasChipContent += '<button class="btn btn-danger" onclick="actorAliasRemove(\'';
      newAliasChipContent += result.id;
      newAliasChipContent += '\');"><i class="fa fa-trash-alt"></i></button></div>';
      newAliasChip.innerHTML = newAliasChipContent;

      newChipButton = document.getElementById('aliasChipsNew');
      newChipParent = newChipButton.parentNode;
      newChipParent.insertBefore(newAliasChip, newChipButton);

      // reset input field
      aliasInput.value = '';

      // revert loader button
      button.classList.remove('spinner-border', 'text-secondary');
      button.classList.add('btn', 'btn-success');
      button.innerHTML = 'Create new alias';

      // hide modal
      document.modalNewAliasModal.hide();
    },
    error: function () {
      sendToast(
        'Actor\'s new alias',
        '',
        'bg-warning',
        'Unable to create this new alias',
      );
    },
  });
}

function actorAliasRemove(aliasID) {
  console.log('remove alias for actor');

  // perform http call
  url = '/api/actor/alias/'+aliasID;
  $.ajax(url, {
    method: 'DELETE',
    xhr: function () {
      var xhr = new XMLHttpRequest();
      return xhr;
    },
    success: function (result) {
      // send notification
      sendToast(
        'Actor\'s alias removal',
        '',
        'bg-success',
        'Successfully removed!',
      );

      // remove node
      chips = document.getElementById('aliasChips');
      for (chip of chips.children) {
        if (chip.getAttribute('alias-id') == aliasID) {
          chip.remove();
        }
      }
    },
    error: function () {
      sendToast(
        'Actor\'s alias removal',
        '',
        'bg-warning',
        'Unable to remove alias',
      );
    },
  });
}

function actorNameEdit() {
  title = document.getElementById('actor-name');
  title.disabled = false;

  btn = document.getElementById('actor-name-edit');
  btn.classList.remove('btn-outline-warning');
  btn.classList.add('btn-outline-success');
  btn.innerText = 'Send';
  btn.onclick = actorNameSend;
}

function actorNameSend() {
  url = "/api/actor/"+actor_id+"/rename";
  $.ajax(url, {
    method: 'POST',
    data: {
      'name': document.getElementById('actor-name').value,
    },
    xhr: function () {
      var xhr = new XMLHttpRequest();
      return xhr;
    },

    success: function (e) {
      console.debug('actor renaming successful');

      title = document.getElementById('actor-name');
      title.disabled = true;

      btn = document.getElementById('actor-name-edit');
      btn.classList.add('btn-outline-warning');
      btn.classList.remove('btn-outline-success');
      btn.innerText = 'Edit';
      btn.onclick = actorNameEdit;
    },

    error: function () {
      console.debug('failed');
    },
  });

  return false;
}

{{ end }}
