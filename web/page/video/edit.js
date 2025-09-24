{{ define "video/edit.js" }}

/* -- start globals set at boot -- */

// store video id
const video_id = '{{ .Video.ID }}';
const url_video_actor_edit = '/api/video/{{ .Video.ID }}/actor/00000000-0000-0000-0000-000000000000';
const url_video_category_edit = '/api/video/{{ .Video.ID }}/category/00000000-0000-0000-0000-000000000000';
/* -- end globals set at boot -- */

// globals
var bootstrap = true;
var actors_complete;

function generateThumbnail(video_id) {
  url = "/api/video/00000000-0000-0000-0000-000000000000/generate-thumbnail/00:00:00";
  url = url.replace('00000000-0000-0000-0000-000000000000', video_id);
  url = url.replace('00:00:00', video_timing);
  $.ajax(url, {
    method: 'POST',
    xhr: function () {
      var xhr = new XMLHttpRequest();
      return xhr;
    },

    success: function (e) {
      console.debug('thumbnail successfully generated');
      ext = document.getElementById('video-thumb');
      ext.innerText = 'Generated';
      ext.classList.remove('bg-warning');
      ext.classList.add('bg-success');
    },

    error: function () {
      console.debug('failed');
    },
  });

  return false;
}

function deleteVideo(video_id) {
  url = "/api/video/"+video_id;
  $.ajax(url, {
    method: 'DELETE',
    xhr: function () {
      var xhr = new XMLHttpRequest();
      return xhr;
    },

    success: function (e) {
      window.location.replace('/');
    },

    error: function () {
      console.debug('failed');
    },
  });

  return false;
}

function changeTo(new_type, video_id) {
  url = "/api/video/"+video_id+"/migrate";
  $.ajax(url, {
    method: 'POST',
    data: {
      'new_type': new_type,
    },
    xhr: function () {
      var xhr = new XMLHttpRequest();
      return xhr;
    },

    success: function (e) {
      sendToast('Change video type', '', 'bg-success', 'Successfull');
      console.debug('success');
    },

    error: function (data) {
      sendToast('Change video type', 'failed', 'bg-danger', data.responseJSON.error);
    },
  });

  return false;
}



/*
 * timing tracking
 */
var video_timing = '00:00:00';
function onTrackedVideoFrame(currentTime){
    date = new Date(null);
    date.setMilliseconds(currentTime*1000);
    video_timing = date.toISOString().substr(11,12);
}

function video_title_edit() {
  title = document.getElementById('video-title');
  title.disabled = false;

  btn = document.getElementById('video-title-edit');
  btn.classList.remove('btn-outline-warning');
  btn.classList.add('btn-outline-success');
  btn.innerText = 'Send';
  btn.onclick = video_title_send;
}

function video_title_send() {
  url = "/api/video/"+video_id+"/rename";
  $.ajax(url, {
    method: 'POST',
    data: {
      'name': document.getElementById('video-title').value,
    },
    xhr: function () {
      var xhr = new XMLHttpRequest();
      return xhr;
    },

    success: function (e) {
      console.debug('import successfull');

      title = document.getElementById('video-title');
      title.disabled = true;

      btn = document.getElementById('video-title-edit');
      btn.classList.add('btn-outline-warning');
      btn.classList.remove('btn-outline-success');
      btn.innerText = 'Edit';
      btn.onclick = video_title_edit;
    },

    error: function () {
      console.debug('failed');
    },
  });

  return false;
}

function video_channel_edit() {
  // get channel list
  $.ajax("/api/channels", {
    method: 'GET',

    xhr: function () {
      var xhr = new XMLHttpRequest();

      return xhr;
    },

    success: function (result) {
      // prepare channel list for the select
      selectChannelList = '<option value="x">None</option>';
      for (const [channelID, channelName] of Object.entries(result.channels)) {
        selectChannelList += '<option value="'+channelID+'">'+channelName+'</option>';
      }
      document.getElementById('channel-list').innerHTML = selectChannelList;

      // get channel edition modal and display it
      document.modalChannelEditModal.show();
    },

    error: function () {
      console.debug('failed');
    },
  });
}

function video_channel_send() {
  // send channel
  $.ajax("/api/video/"+video_id+"/channel", {
    method: 'POST',
    data: {
      'channelID': document.getElementById('channel-list').value,
    },
    xhr: function () {
      var xhr = new XMLHttpRequest();
      return xhr;
    },

    success: function () {
      // get channel edition modal and hide it
      document.modalChannelEditModal.hide();

      // to improve later: reload to show new channel
      window.location.reload();
    },

    error: function () {
      console.debug('failed');
    },
  });
}

window.onload = function() {
  // for thumbnail generation
  video = document.getElementById("video")
  video.addEventListener(
    "timeupdate",
    function(event){
      onTrackedVideoFrame(this.currentTime);
    }
  );

  document.modalChannelEdit = document.getElementById("editChannelModal");
  document.modalChannelEditModal = new bootstrap.Modal(document.modalChannelEdit);
}

//
// actor selection
//
{{ template "shards/actor-selection/main.js" . }}

async function videoEditActorAdd(actor_id) {
  console.log('add actor '+actor_id+' in video '+video_id);
  url = url_video_actor_edit.replace('00000000-0000-0000-0000-000000000000', actor_id);
  try {
    var res = await ajax(url, 'PUT');
    sendToast('Actor added', '', 'bg-success', this.actorSelectable[actor_id]['name']+' added.');
  } catch(e) {
    console.debug(e);
    sendToast('Actor not added', '', 'bg-danger', this.actorSelectable[actor_id]['name']+' not added, call failed.');
  }
}

async function videoEditActorRemove(actor_id) {
  console.log('remove actor '+actor_id+' from video '+video_id);
  url = url_video_actor_edit.replace('00000000-0000-0000-0000-000000000000', actor_id);
  try {
    var res = await ajax(url, 'DELETE');
    sendToast('Actor removed', '', 'bg-success', this.actorSelectable[actor_id]['name']+' removed.');
  } catch(e) {
    console.debug(e);
    sendToast('Actor not removed', '', 'bg-danger', this.actorSelectable[actor_id]['name']+' not removed, call failed.');
  }
}

window.zt.onload.unshift(function videoEditActorSelectableConfiguration() {
  // all actors
  window.zt.actorSelection.actorSelectable = {
    {{ range $actor := .Actors }}
    '{{ $actor.ID }}': {
      'name': '{{ $actor.Name }}',
      'aliases': [],
    },
    {{ end }}
  };

  // add actors in video
  window.zt.actorSelection.actorSelected = {
    {{ range $actor := .Video.Actors }}
    '{{ $actor.ID }}': undefined,
    {{ end }}
  };

  // call api on actor add
  window.zt.actorSelection.onActorSelectBefore = videoEditActorAdd;

  // call api on actor removal
  window.zt.actorSelection.onActorDeselectBefore = videoEditActorRemove;
});
//
// actor selection end
//

//
// category selection
//
{{ template "shards/category-selection/main.js" . }}

async function videoEditCategoryAdd(category_id) {
  console.log('add category'+category_id+' in video '+video_id);
  url = url_video_category_edit.replace('00000000-0000-0000-0000-000000000000', category_id);
  try {
    var res = await ajax(url, 'PUT');
    sendToast('Category added', '', 'bg-success', this.categorySelectable[category_id]['name']+' added.');
  } catch(e) {
    console.debug(e);
    sendToast('Category not added', '', 'bg-danger', this.categorySelectable[category_id]['name']+' not added, call failed.');
  }
}

async function videoEditCategoryRemove(category_id) {
  console.log('remove category '+category_id+' from video '+video_id);
  url = url_video_category_edit.replace('00000000-0000-0000-0000-000000000000', category_id);
  try {
    var res = await ajax(url, 'DELETE');
    sendToast('Category removed', '', 'bg-success', this.categorySelectable[category_id]['name']+' removed.');
  } catch(e) {
    console.debug(e);
    sendToast('Category not removed', '', 'bg-danger', this.cateogrySelectable[category_id]['name']+' not removed, call failed.');
  }
}

window.zt.onload.unshift(function videoEditCategorySelectableConfiguration() {
  // store all categories at start
  window.zt.categorySelection.categorySelectable = {
    {{ range $category := .Categories }}
    {{ range $sub:= $category.Sub }}
    '{{ $sub.ID }}': {
      'name': '{{ $sub.Name }}',
    },
    {{ end }}
    {{ end }}
  };

// store all categories in the video at start
  window.zt.categorySelection.categorySelected = {
    {{ range $sub := .Video.Categories }}
    '{{ $sub.ID }}': undefined,
    {{ end }}
  };

  // call api on actor add
  window.zt.categorySelection.onCategorySelectBefore = videoEditCategoryAdd;

  // call api on actor removal
  window.zt.categorySelection.onCategoryDeselectBefore = videoEditCategoryRemove;
});
//
// category selection end
//

{{ end }}
