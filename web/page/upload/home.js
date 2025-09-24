{{ define "upload/home.js" }}

var CURRENT_PATH='/';

function setFolder(folderName) {
  CURRENT_PATH=folderName;
  updatePath();
}

function enterFolder(folderName) {
  if (CURRENT_PATH == '/') {
    CURRENT_PATH += folderName;
  } else {
    CURRENT_PATH += '/';
    CURRENT_PATH += folderName;
  }
  updatePath();
}

function updatePath() {
  console.log('switching to: '+CURRENT_PATH);

  generateBreadcrumb();

  // reset table listing
  document.getElementById('triage-listing').innerHTML = '';

  loadFolders();
  loadFiles();
}

function generateBreadcrumb() {
  var path = CURRENT_PATH;
  pathAsArray = path.split('/');
  max = pathAsArray.length;
  progressivePath = '';

  newBr  = '<ol class="breadcrumb">';
  newBr += '<li class="breadcrumb-item';
  if (path == '/') {
    newBr += ' active">';
  } else {
    newBr += ' br-link" onclick="setFolder(\'/\')">';
  }
  newBr += 'Triage</li>';
  for (var i=0; i < max; i++) {
    folder = pathAsArray[i];
    if (folder == "") {
      continue;
    } else {
      progressivePath += '/'+folder;
    }

    newBr += '<li class="breadcrumb-item';
    if (i+1 == max) {
      newBr += ' active">';
    } else {
      newBr += ' br-link" onClick="setFolder(\''+progressivePath+'\')">';
    }
    newBr += folder+'</li>';
  }

  newBr += '</ol>';

  document.getElementById('path').innerHTML = newBr;
}

function loadFolders() {
  $.ajax('/api/upload/triage/folder', {
    method: 'POST',
    data: {
      'path': CURRENT_PATH,
    },
    xhr: function () {
      var xhr = new XMLHttpRequest();
      return xhr;
    },

    success: function (data) {
      showFolders(data.folders);
    },

    error: function (data) {
      sendToast(
        'Unable to show folders',
        '',
        'bg-warning',
        data.responseJSON.error,
      );
    },
  });
}

function loadFiles() {
  $.ajax('/api/upload/triage/file', {
    method: 'POST',
    data: {
      'path': CURRENT_PATH,
    },
    xhr: function () {
      var xhr = new XMLHttpRequest();
      return xhr;
    },

    success: function (data) {
      showFiles(data.files);
      updateSelects();
    },

    error: function (data) {
      sendToast(
        'Unable to show files',
        '',
        'bg-warning',
        data.responseJSON.error,
      );
    },
  });
}

function showFolders(folders) {
  folderEntries = '';
  Object.keys(folders).forEach(function(folder) {
    folderEntries += '<tr ondblclick="enterFolder(\''+folder+'\');">';
    folderEntries += '<td></td>';
    folderEntries += '<td><i class="far fa-folder"></i></td>';
    folderEntries += '<td>'+folder+'</td>';
    folderEntries += '<td>—</td>';
    folderEntries += '<td>'+folders[folder]+' items</td>';
    folderEntries += '</tr>';
  });

  document.getElementById('triage-listing').innerHTML += folderEntries;
}

function showFiles(files) {
  fileEntries = '';
  Object.keys(files).forEach(function(file) {
    // store file infos
    fileInfo = files[file];
    fileSize = fileInfo['Size'] == 0 ? 'Empty' : humanFileSize(fileInfo['Size'], true);
    const options = { year: "numeric", month: "short", day: "numeric" };
    fileDate = new Date(fileInfo['LastModification']).toLocaleDateString(undefined, options);

    // guess filetype based on the extension
    fileType = 'unknown';
    fileIcon = 'fas fa-file';
    if (file.match(/.*\.(mp4|mkv|webm)/)) {
      fileType = 'video';
      fileIcon = 'far fa-play-circle';
    } else if (file.match(/.*\.(png|jpg|jpeg)/)) {
      fileType = 'image';
      fileIcon = 'far fa-image';
    } else if (file.match(/.*\.(zip|tar)/)) {
      fileType = 'archive';
      fileIcon = 'far fa-file-archive';
    }

    // build file path
    filePath = CURRENT_PATH == '/' ? filePath = '/'+file : CURRENT_PATH+'/'+file;

    // render files
    fileEntries += '<tr data-type="file" data-filename="'+file+'" data-filepath="'+filePath+'" data-filetype="'+fileType+'" onclick="showFileDetails(event, this, \''+file+'\', \''+filePath+'\', \''+fileType+'\', \''+fileSize+'\', \''+fileDate+'\')">';
    fileEntries += '<td onclick="singleSelectToggle(this);"><i class="far fa-square checks"></i></td>';
    fileEntries += '<td><i class="'+fileIcon+'"></i></td>';
    fileEntries += '<td>'+file+'</td>';
    fileEntries += '<td>'+fileDate+'</td>';
    fileEntries += '<td>'+fileSize+'</td>';
    fileEntries += '</tr>';
  });

  document.getElementById('triage-listing').innerHTML += fileEntries;
}

function clickInputButton() {
  document.getElementById('upload-input-path').value = CURRENT_PATH;
  document.getElementById('upload-input-file').click();
}

var MD5 = function(d){var r = M(V(Y(X(d),8*d.length)));return r.toLowerCase()};function M(d){for(var _,m="0123456789ABCDEF",f="",r=0;r<d.length;r++)_=d.charCodeAt(r),f+=m.charAt(_>>>4&15)+m.charAt(15&_);return f}function X(d){for(var _=Array(d.length>>2),m=0;m<_.length;m++)_[m]=0;for(m=0;m<8*d.length;m+=8)_[m>>5]|=(255&d.charCodeAt(m/8))<<m%32;return _}function V(d){for(var _="",m=0;m<32*d.length;m+=8)_+=String.fromCharCode(d[m>>5]>>>m%32&255);return _}function Y(d,_){d[_>>5]|=128<<_%32,d[14+(_+64>>>9<<4)]=_;for(var m=1732584193,f=-271733879,r=-1732584194,i=271733878,n=0;n<d.length;n+=16){var h=m,t=f,g=r,e=i;f=md5_ii(f=md5_ii(f=md5_ii(f=md5_ii(f=md5_hh(f=md5_hh(f=md5_hh(f=md5_hh(f=md5_gg(f=md5_gg(f=md5_gg(f=md5_gg(f=md5_ff(f=md5_ff(f=md5_ff(f=md5_ff(f,r=md5_ff(r,i=md5_ff(i,m=md5_ff(m,f,r,i,d[n+0],7,-680876936),f,r,d[n+1],12,-389564586),m,f,d[n+2],17,606105819),i,m,d[n+3],22,-1044525330),r=md5_ff(r,i=md5_ff(i,m=md5_ff(m,f,r,i,d[n+4],7,-176418897),f,r,d[n+5],12,1200080426),m,f,d[n+6],17,-1473231341),i,m,d[n+7],22,-45705983),r=md5_ff(r,i=md5_ff(i,m=md5_ff(m,f,r,i,d[n+8],7,1770035416),f,r,d[n+9],12,-1958414417),m,f,d[n+10],17,-42063),i,m,d[n+11],22,-1990404162),r=md5_ff(r,i=md5_ff(i,m=md5_ff(m,f,r,i,d[n+12],7,1804603682),f,r,d[n+13],12,-40341101),m,f,d[n+14],17,-1502002290),i,m,d[n+15],22,1236535329),r=md5_gg(r,i=md5_gg(i,m=md5_gg(m,f,r,i,d[n+1],5,-165796510),f,r,d[n+6],9,-1069501632),m,f,d[n+11],14,643717713),i,m,d[n+0],20,-373897302),r=md5_gg(r,i=md5_gg(i,m=md5_gg(m,f,r,i,d[n+5],5,-701558691),f,r,d[n+10],9,38016083),m,f,d[n+15],14,-660478335),i,m,d[n+4],20,-405537848),r=md5_gg(r,i=md5_gg(i,m=md5_gg(m,f,r,i,d[n+9],5,568446438),f,r,d[n+14],9,-1019803690),m,f,d[n+3],14,-187363961),i,m,d[n+8],20,1163531501),r=md5_gg(r,i=md5_gg(i,m=md5_gg(m,f,r,i,d[n+13],5,-1444681467),f,r,d[n+2],9,-51403784),m,f,d[n+7],14,1735328473),i,m,d[n+12],20,-1926607734),r=md5_hh(r,i=md5_hh(i,m=md5_hh(m,f,r,i,d[n+5],4,-378558),f,r,d[n+8],11,-2022574463),m,f,d[n+11],16,1839030562),i,m,d[n+14],23,-35309556),r=md5_hh(r,i=md5_hh(i,m=md5_hh(m,f,r,i,d[n+1],4,-1530992060),f,r,d[n+4],11,1272893353),m,f,d[n+7],16,-155497632),i,m,d[n+10],23,-1094730640),r=md5_hh(r,i=md5_hh(i,m=md5_hh(m,f,r,i,d[n+13],4,681279174),f,r,d[n+0],11,-358537222),m,f,d[n+3],16,-722521979),i,m,d[n+6],23,76029189),r=md5_hh(r,i=md5_hh(i,m=md5_hh(m,f,r,i,d[n+9],4,-640364487),f,r,d[n+12],11,-421815835),m,f,d[n+15],16,530742520),i,m,d[n+2],23,-995338651),r=md5_ii(r,i=md5_ii(i,m=md5_ii(m,f,r,i,d[n+0],6,-198630844),f,r,d[n+7],10,1126891415),m,f,d[n+14],15,-1416354905),i,m,d[n+5],21,-57434055),r=md5_ii(r,i=md5_ii(i,m=md5_ii(m,f,r,i,d[n+12],6,1700485571),f,r,d[n+3],10,-1894986606),m,f,d[n+10],15,-1051523),i,m,d[n+1],21,-2054922799),r=md5_ii(r,i=md5_ii(i,m=md5_ii(m,f,r,i,d[n+8],6,1873313359),f,r,d[n+15],10,-30611744),m,f,d[n+6],15,-1560198380),i,m,d[n+13],21,1309151649),r=md5_ii(r,i=md5_ii(i,m=md5_ii(m,f,r,i,d[n+4],6,-145523070),f,r,d[n+11],10,-1120210379),m,f,d[n+2],15,718787259),i,m,d[n+9],21,-343485551),m=safe_add(m,h),f=safe_add(f,t),r=safe_add(r,g),i=safe_add(i,e)}return Array(m,f,r,i)}function md5_cmn(d,_,m,f,r,i){return safe_add(bit_rol(safe_add(safe_add(_,d),safe_add(f,i)),r),m)}function md5_ff(d,_,m,f,r,i,n){return md5_cmn(_&m|~_&f,d,_,r,i,n)}function md5_gg(d,_,m,f,r,i,n){return md5_cmn(_&f|m&~f,d,_,r,i,n)}function md5_hh(d,_,m,f,r,i,n){return md5_cmn(_^m^f,d,_,r,i,n)}function md5_ii(d,_,m,f,r,i,n){return md5_cmn(m^(_|~f),d,_,r,i,n)}function safe_add(d,_){var m=(65535&d)+(65535&_);return(d>>16)+(_>>16)+(m>>16)<<16|65535&m}function bit_rol(d,_){return d<<_|d>>>32-_}

function createToast(path, file) {

  toastTemplate = document.getElementById('toastFileUploadTemplate');

  // create new toast
  newtoast = toastTemplate.cloneNode(true);

  div_id = 'toast_file_'+MD5(file+'/'+path);
  newtoast.id = div_id;

  // set title
  //n_title = newtoast.getElementsByClassName('zt-toast-title');
  //n_title[0].innerText = title;

  // set subtitle
  //n_subtitle = newtoast.getElementsByClassName('zt-toast-subtitle');
  //n_subtitle[0].innerText = subtitle;

  // set color
  //n_header = newtoast.getElementsByClassName('toast-header');
  //n_header[0].className += " "+title_color

  // set content
  n_content = newtoast.getElementsByClassName('zt-toast-body');
  n_content[0].innerHTML = 'Uploading of '+file;

 
  toastContainer = document.getElementById('zt-toast-container');
  toastContainer.appendChild(newtoast);
  toast = new bootstrap.Toast(newtoast, {'autohide': false});
  toast.show();
 
}

function updateToastProgress(path, file, current, max) {
  div_id = 'toast_file_'+MD5(file+'/'+path);
  toast = document.getElementById(div_id);

  value = Math.floor(current*100/max);

  // set subtitle
  n_subtitle = toast.getElementsByClassName('zt-toast-subtitle');
  n_subtitle[0].innerText = value + ' %';

  if (value == 100) {
    // set color
    n_header = toast.getElementsByClassName('toast-header');
    n_header[0].classList.remove('bg-primary');
    n_header[0].classList.add('bg-success');

    // set title
    n_title = toast.getElementsByClassName('zt-toast-title');
    n_title[0].innerText = 'Upload finished';

    updatePath();

    // close a little bit later
    setTimeout(function() {
      $(toast).hide();
    }, 5000);
  }
}

function fileUpload(e) {
  var file = e.files[0];

  createToast(CURRENT_PATH,file.name);
  $.ajax({
    url: '/api/upload/file',
    type: 'POST',

    data: new FormData(e.form),

    cache: false,
    contentType: false,
    processData: false,

    // Custom XMLHttpRequest
    xhr: function () {
      var myXhr = $.ajaxSettings.xhr();
      if (myXhr.upload) {
        // For handling the progress of the upload
        myXhr.upload.addEventListener('progress', function (e) {
          if (e.lengthComputable) {
            updateToastProgress(CURRENT_PATH, file.name, e.loaded, e.total);
          }
        }, false);
      }
      return myXhr;
    },
    done: function() {
      updateToastProgress(CURRENT_PATH, file.name, e.loaded, e.total);
    },
    fail: function(errorThrown) {
      console.error('upload failed');
      console.error(errorThrown);
    },
  });
}

function importVideo(filename, filepath) {
  const filenameInput = window.VideoImport.querySelector('#videoImportFileName');
  filenameInput.innerText = filename;

  const filepathInput = window.VideoImport.querySelector('#videoImportFilePath');
  filepathInput.innerText = filepath;

  // store filepath as global to avoid html whitespace stripping
  window.G_filepath = filepath;

  // display modal
  window.VideoImportModal.show();
}

function showFileDetails(e, element, fileName, filePath, fileType, fileSize, fileDate) {
  // prevent the tick from showing the off-canvas
  e = e || event;
  var target = e.target || e.srcElement;
  if (element.childNodes[0].childNodes[0] == target)
    return;
  if (element.childNodes[0].childNodes[0] == target.childNodes[0])
    return;

  // Nice file types
  const niceFileTypes = {
    'unknown': 'Unrecognized type',
    'video': 'Video',
    'image': 'Picture',
    'archive': 'Compressed archive',
  };

  // generate canvas content - common atrributes
  canvasHTML = '<h3>File details</h3><hr />';
  canvasHTML += '<b>Name</b><p>'+fileName+'</p>';
  canvasHTML += '<b>Type</b><p>'+niceFileTypes[fileType]+'</p>';
  canvasHTML += '<b>Size</b><p>'+fileSize+'</p>';
  canvasHTML += '<b>Last modification</b><p>'+fileDate+'</p>';

  // generate canvas content - video specifics
  if (fileType == 'video') {
    canvasHTML += '<h5 style="padding-top: 50px;">Video preview</h5><hr />';
    canvasHTML += '<video controls class="previewVideo" style="width: 100%"';
    canvasHTML += 'src="/upload/preview/'+encodeURIComponent(filePath.substring(1))+'" >';
    canvasHTML += '</video>';
  }
  // generate canvas content - picture specifics
  else if (fileType == 'image') {
    canvasHTML += '<h5 style="padding-top: 50px;">Picture preview</h5><hr />';
    canvasHTML += '<img style="width: 100%"';
    canvasHTML += 'src="/upload/preview/'+encodeURIComponent(filePath.substring(1))+'" ';
    canvasHTML += '/>';
  }

  // generate canvas content - actions
  canvasHTML += '<h5 style="padding: 50px 0px 15px 0px;">Actions</h5>';

  // generate canvas content - video specific actions
  if (fileType == 'video') {
    canvasHTML += '<button class="btn btn-outline-primary" onclick="importVideo(\''+fileName+'\', \''+filePath+'\')">';
    canvasHTML += '<i class="fas fa-file-import"></i> Import</button>';
  }

  // generate canvas content - common actions
  canvasHTML += '<a class="btn btn-outline-primary" target="_blank" rel="noopener noreferrer" ';
  canvasHTML += 'href="/upload/preview/'+encodeURIComponent(filePath.substring(1))+'">';
  canvasHTML += '<i class="fas fa-download"></i> Download</a>';

  canvasHTML += '<button class="btn btn-outline-danger" onclick="deleteFile(\''+filePath.substring(1)+'\');"><i class="far fa-trash-alt"></i> Delete</button>';

  // set content
  window.DetailPanelContent.innerHTML = canvasHTML;

  // display canvas
  window.DetailPanelOffCanvas.show();
}

/* startup */
window.onload = function() {
  generateBreadcrumb();
  setFolder(CURRENT_PATH);

  // initialize detail panel
  window.DetailPanel = document.getElementById('item-details');
  window.DetailPanelOffCanvas = new bootstrap.Offcanvas(window.DetailPanel);
  window.DetailPanelContent = document.getElementById('item-details-content');

  // initialize mass action panel
  window.MassActionPanel = document.getElementById('mass-action-details');
  window.MassActionPanelModal = new bootstrap.Modal(window.MassActionPanel);

  // initialize mass import panel
  window.MassImportPanel = document.getElementById('mass-import-panel');
  window.MassImportPanelModal = new bootstrap.Modal(window.MassImportPanel);

  // pause video when hidden
  window.DetailPanel.addEventListener('hide.bs.offcanvas', event => {
    const video = DetailPanelContent.querySelector('.previewVideo');
    if (video != null) {
      video.pause();
    }
  });

  // initialize video import modal
  window.VideoImport = document.getElementById('video-import-modal');
  window.VideoImportModal = new bootstrap.Modal(window.VideoImport);

  // "enter" handling on folder creation
  document.getElementById('folder-new').addEventListener("keydown", function (e) {
    if (e.code === "Enter") {
        new_folder_send();
    }
  });
}

function importAsVideoButtonUpdate(sel) {
  const btn = document.getElementById('importAsVideoButton');
  btn.disabled = (sel.value == 'x');
}

function importVideoAjax(filetype) {
  oc = document.getElementById('importAsVideoOffCanvas');
  const title = window.VideoImport.querySelector('#videoImportFileName').innerText;

  $.ajax('/api/video', {
    method: 'POST',
    data: {
      'name': title,
      'filename': window.G_filepath,
      'type': filetype,
    },
    xhr: function () {
      var xhr = new XMLHttpRequest();
      return xhr;
    },

    success: function (data) {
      // get video id from response
      video_id = data.video_id;

      // send notification
      sendToast(
        'Import a video',
        '',
        'bg-success',
        '<a href="/video/'+video_id+'/edit" target="_blank">You can edit the video more here</a>',
      );

      // close canvas

      // remove video from view

      // close modal
      window.VideoImportModal.hide();

      // close off canvas
      window.DetailPanelOffCanvas.hide();

      // refresh view
      updatePath();

    },
    error: function (data) {
      sendToast(
        'Unable to import video',
        '',
        'bg-warning',
        data.responseJSON.error,
      );
    },
  });
}

function deleteFile(file) {
  oc = document.getElementById('importAsVideoOffCanvas');
  const filename = window.VideoImport.querySelector('#videoImportFilePath').innerText;
  const title = window.VideoImport.querySelector('#videoImportFileName').innerText;

  console.log('delete file:', file);

  $.ajax('/api/upload/file', {
    method: 'DELETE',
    data: JSON.stringify({
      'file': file,
    }),
    contentType: 'application/json; charset=utf-8',
    dataType: 'json',
    xhr: function () {
      var xhr = new XMLHttpRequest();
      return xhr;
    },

    success: function (data) {
      // get video id from response
      video_id = data.video_id;

      // send notification
      sendToast(
        'Delete a file',
        '',
        'bg-success',
        'File deleted successfully',
      );

      // close modal
      window.VideoImportModal.hide();

      // close off canvas
      window.DetailPanelOffCanvas.hide();

      // refresh view
      updatePath();

    },
    error: function (data) {
      sendToast(
        'Unable to delete file',
        '',
        'bg-warning',
        data.responseJSON.error,
      );
    },
  });
}

function new_folder_send() {
  url = "/api/upload/folder";
  folder_name = CURRENT_PATH;
  if (folder_name != '/') {
    folder_name += '/';
  }
  folder_name += document.getElementById('folder-new').value;

  $.ajax(url, {
    method: 'POST',
    data: {
      'name': folder_name,
    },
    xhr: function () {
      var xhr = new XMLHttpRequest();
      return xhr;
    },

    success: function (e) {
      // notify success
      sendToast(
        'New folder created',
        '',
        'bg-success',
        'Folder '+folder_name+' created',
      );

      // hide modal
      $('#newFolderModal').modal('hide');

      // cleanup input
      document.getElementById('folder-new').value = '';

      // update folder list
      updatePath();
    },

    error: function (data) {
      sendToast(
        'Unable to create folder',
        '',
        'bg-warning',
        data.responseJSON.error,
      );
    },
  });

  return false;
}

/**
 * mass import related functions
 */
allEnabled = false;
function multiSelectToggle() {
  if (allEnabled) {
    document.getElementById('triage-listing').childNodes.forEach(function(item){
      if (item.childNodes[0] == undefined || item.childNodes[0].childNodes[0] == undefined)
        return;
      item.childNodes[0].childNodes[0].classList.add('fa-square');
      item.childNodes[0].childNodes[0].classList.remove('fa-check-square');
    });
    allEnabled = false;
  } else {
    document.getElementById('triage-listing').childNodes.forEach(function(item){
      if (item.childNodes[0] == undefined || item.childNodes[0].childNodes[0] == undefined)
        return;
      item.childNodes[0].childNodes[0].classList.remove('fa-square');
      item.childNodes[0].childNodes[0].classList.add('fa-check-square');
    });
    allEnabled = true;
  }
  updateSelects();
}

function singleSelectToggle(caller) {
  tick = caller.childNodes[0];
  if (tick.classList.contains('fa-check-square')) {
    tick.classList.remove('fa-check-square');
    tick.classList.add('fa-square');
  } else {
    tick.classList.remove('fa-square');
    tick.classList.add('fa-check-square');
  }
  updateSelects();
}

function updateSelects() {
  selected = 0;
  unselected = 0;
  document.getElementById('triage-listing').childNodes.forEach(function(item){
      if (item.childNodes[0] == undefined || item.childNodes[0].childNodes[0] == undefined)
        return;
    if (item.childNodes[0].childNodes[0].classList.contains('fa-square'))
      unselected++;
    if (item.childNodes[0].childNodes[0].classList.contains('fa-check-square'))
      selected++;
  });

  selectAllTick = document.getElementById('selectAllTick');
  if (unselected == 0 && selected > 0) {
    selectAllTick.classList.remove('fa-square');
    selectAllTick.classList.add('fa-check-square');
  } else {
    selectAllTick.classList.remove('fa-check-square');
    selectAllTick.classList.add('fa-square');
  }

  massImportButton = document.getElementById('button-mass-action');
  if (selected == 0) {
    massImportButton.classList.add('disabled');
  } else {
    massImportButton.classList.remove('disabled');
  }
}

function massActionGetFileList() {
  fileList = [];

  document.getElementById('triage-listing').childNodes.forEach(function(item){
    // skip folders
    if (item.childNodes[0] == undefined || item.childNodes[0].childNodes[0] == undefined)
      return;

    // skip all unchecked items
    if (item.childNodes[0].childNodes[0].classList.contains('fa-square')) {
      return;
    }

    fileList.push(item);
  });

  return fileList;
}

function showMassActionModal() {
  fileList = massActionGetFileList();

  massActionItemListTable = '';
  fileList.forEach(function(item){
    massActionItemListTable += '<tr>';
    massActionItemListTable += '<td>'+item.dataset.filepath+'</td>';
    massActionItemListTable += '<td>'+item.dataset.filetype+'</td>';
    massActionItemListTable += '</tr>';
  });

  document.getElementById('mass-action-item-list').innerHTML = massActionItemListTable;

  window.MassActionPanelModal.show();

  window.zt.actorSelection.onModalHide = function() {
    window.MassImportPanelModal.show();
  };

  window.zt.categorySelection.onModalHide = function() {
    window.MassImportPanelModal.show();
  };
}

function showMassImportModal() {
  window.MassActionPanelModal.hide();
  window.MassImportPanelModal.show();
}

function massDelete() {
  fileList = massActionGetFileList();
  curatedFileList = [];
  fileList.forEach(function(item){
    curatedFileList.push(item.dataset.filepath);
  });

  $.ajax('/api/upload/triage/mass-action', {
    method: 'DELETE',
    contentType: 'application/json; charset=utf-8',
    dataType: 'json',
    data: JSON.stringify({
      'files': curatedFileList,
    }),
    xhr: function () {
      var xhr = new XMLHttpRequest();
      return xhr;
    },

    success: function (data) {
      window.MassActionPanelModal.hide();
      sendToast(
        'Deletion successful',
        '',
        'bg-success',
        'This folder should feel lighter now',
      );
      updatePath();
    },

    error: function (data) {
      sendToast(
        'Unable to delete selected files',
        '',
        'bg-warning',
        data.responseJSON.error,
      );
    },
  });
}

function massImportSend(videoType) {
  fileList = massActionGetFileList();
  curatedFileList = [];
  fileList.forEach(function(item){
    curatedFileList.push(item.dataset.filepath);
  });

  actorList = [];
  for (const actor_id of Object.keys(window.zt.actorSelection.actorSelected)) {
    actorList.push(actor_id);
  }
  categoryList = [];
  for (const category_id of Object.keys(window.zt.categorySelection.categorySelected)) {
    categoryList.push(category_id);
  }

  channel = document.getElementById('channel-list').value;
  if (channel == 'None')
    channel = undefined;

  $.ajax('/api/upload/triage/mass-action', {
    method: 'POST',
    contentType: 'application/json; charset=utf-8',
    dataType: 'json',
    data: JSON.stringify({
      'files': curatedFileList,
      'actors': actorList,
      'categories': categoryList,
      'type': videoType,
      'channel': channel,
    }),
    xhr: function () {
      var xhr = new XMLHttpRequest();
      return xhr;
    },

    success: function (data) {
      updatePath();
      window.MassImportPanelModal.hide();
      sendToast(
        'Import successful',
        '',
        'bg-success',
        'Import tasks will run in background',
      );
    },

    error: function (data) {
      updatePath();
      window.MassImportPanelModal.hide();
      sendToast(
        'Unable to import selected files',
        '',
        'bg-warning',
        data.responseJSON.error,
      );
    },
  });

}

// From: https://stackoverflow.com/a/14919494

/**
 * Format bytes as human-readable text.
 *
 * @param bytes Number of bytes.
 * @param si True to use metric (SI) units, aka powers of 1000. False to use
 *           binary (IEC), aka powers of 1024.
 * @param dp Number of decimal places to display.
 *
 * @return Formatted string.
 */
 function humanFileSize(bytes, si=false, dp=1) {
  const thresh = si ? 1000 : 1024;

  if (Math.abs(bytes) < thresh) {
    return bytes + ' B';
  }

  const units = si
    ? ['kB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']
    : ['KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB'];
  let u = -1;
  const r = 10**dp;

  do {
    bytes /= thresh;
    ++u;
  } while (Math.round(Math.abs(bytes) * r) / r >= thresh && u < units.length - 1);


  return bytes.toFixed(dp) + ' ' + units[u];
}

//
// actor selection
//
{{ template "shards/actor-selection/main.js" . }}

window.zt.onload.unshift(function UploadHomeActorSelectableConfiguration() {
  // all actors
  window.zt.actorSelection.actorSelectable = {
    {{ range $actor := .Actors }}
    '{{ $actor.ID }}': {
      'name': '{{ $actor.Name }}',
      'aliases': [],
    },
    {{ end }}
  };
});
//
// actor selection end
//

//
// category selection
//
{{ template "shards/category-selection/main.js" . }}

window.zt.onload.unshift(function UploadHomeCategorySelectableConfiguration() {
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
});
//
// category selection end
//


{{ end }}
