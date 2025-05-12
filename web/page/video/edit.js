{{ define "video/edit.js" }}

/* -- start globals set at boot -- */

// store video id
const video_id = '{{ .Video.ID }}';
const url_video_actor_edit = '/api/video/{{ .Video.ID }}/actor/00000000-0000-0000-0000-000000000000';
var csrf_token = '{% csrf_token %}';

/*
    'aliases': [
      {% for alias in actor.alias %}
      '{% alias.name %}',
      {% endfor %}
    ],
*/

// store all actors at start
var actors_all = {
  {{ range $actor := .Actors }}
  '{{ $actor.ID }}': {
    'name': '{{ $actor.Name }}',
    'aliases': [],
  },
  {{ end }}
};

// store all actors in the video at start
var actors_in_video = {
  {{ range $actor := .Video.Actors }}
  '{{ $actor.ID }}': undefined,
  {{ end }}
};

/* -- end globals set at boot -- */

// globals
var bootstrap = true;
var actors_complete;

/* videoAddActorInVideo - FilterPresentActors */
function videoAddActorInVideoFilterPresentActors() {
  console.log('update actors presents in video');
  var actors_chips = document.getElementsByClassName('add-actor-list');
  for (const actor_chip of actors_chips) {
    actor_id = actor_chip.getAttribute('actor-id');
    if (actor_id in actors_in_video) {
      actor_chip.querySelector('.btn-success').style.display = 'none';
      actor_chip.querySelector('.btn-danger').style.display = '';
    } else {
      actor_chip.querySelector('.btn-success').style.display = '';
      actor_chip.querySelector('.btn-danger').style.display = 'none';
    }
  }
}

/* videoAddActorInVideo - FilterPresentActors */
function videoAddActorInVideoFilterActorsFromInput(filter = '') {
  var re = new RegExp(filter, 'i');
  var actors_chips = document.getElementsByClassName('add-actor-list');
  found_count = 0;
  found_count_max = 15;
  for (const actor_chip of actors_chips) {
    actor_id = actor_chip.getAttribute('actor-id');
    found = re.test(actors_all[actor_id]['name']);
    if (found) {
      found_count++;
    } else {
      for (const a of actors_all[actor_id]['aliases']) {
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
}

/* videoRemoveActor */
function videoRemoveActor(actor_id) {
  console.log('remove actor '+actor_id+' from video '+video_id);
  url = url_video_actor_edit.replace('00000000-0000-0000-0000-000000000000', actor_id);
  $.ajax(url, {
    method: 'DELETE',
    headers: {
      'X-CSRFToken': '{% csrf_token %}',
    },

    xhr: function () {
      var xhr = new XMLHttpRequest();
      return xhr;
    },

    success: function (res) {
      console.debug('success, got', res);
      sendToast('Actor removed', '', 'bg-success', actors_all[actor_id]['name']+' removed.');
      delete actors_in_video[actor_id];
      videoAddActorInVideoFilterPresentActors();
    },

    error: function () {
      console.debug('failed');
      sendToast('Actor not removed', '', 'bg-danger', actors_all[actor_id]['name']+' not removed, call failed.');
    },

    complete: function () {
    },
  });
}

/* videoRemoveActor */
function videoAddActor(actor_id) {
  console.log('add actor '+actor_id+' in video '+video_id);
  url = url_video_actor_edit.replace('00000000-0000-0000-0000-000000000000', actor_id);
  $.ajax(url, {
    method: 'PUT',
    headers: {
      'X-CSRFToken': '{% csrf_token %}',
    },

    xhr: function () {
      var xhr = new XMLHttpRequest();
      return xhr;
    },

    success: function (res) {
      console.debug('success, got', res);
      sendToast('Actor added', '', 'bg-success', actors_all[actor_id]['name']+' added.');
      actors_in_video[actor_id] = undefined;
      videoAddActorInVideoFilterPresentActors();
    },

    error: function () {
      console.debug('failed');
      sendToast('Actor not added', '', 'bg-danger', actors_all[actor_id]['name']+' not added, call failed.');
    },

    complete: function () {
    },
  });
}

/* filter actor modal */
function actorInputUpdate(e) {
  videoAddActorInVideoFilterActorsFromInput(e.target.value);
}

function updateActorEntries(filter = '') {
  var content = '';
  var re = new RegExp(filter, 'i');
  var alreaySelectedEntries = $('#actors')[0].value.split(',');
  for (actor in actors_complete['name']) {
    if (!re.test(actor)) {
      continue;
    }

    if (alreaySelectedEntries.includes(actors_complete['name'][actor])) {
      continue;
    }

    content += '<div class="cs-entry" onClick="addEntry(\''+actor+'\')">'+actor+'</div>';
  }
  $(".cs-dataset").html(content)
  if (filter == '') {
    $(".cs-menu").hide();
  } else {
    $(".cs-menu").show();
  }
}

window.onload = function() {
  // toggle right button in modal
  videoAddActorInVideoFilterPresentActors();

  // add event on modal input
  document.getElementById('addActorInput').addEventListener('input', actorInputUpdate);
  videoAddActorInVideoFilterActorsFromInput('');

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
      btn.onClick = video_title_edit;
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

{{ end }}
