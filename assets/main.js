
htmx.on('#form', 'htmx:xhr:progress', function(evt) {
  htmx.find('#progress').setAttribute('value', evt.detail.loaded/evt.detail.total * 100)
  if (evt.detail.loaded/evt.detail.total == 1) {
    setTimeout(() => {
      htmx.find('#progress').setAttribute('value', 0)
    }, 1000);
  }
});

// function loadAsset(path, callback=null) {
//   var xhr = new XMLHttpRequest();
//   xhr.open('GET', path, true);
//   xhr.setRequestHeader('updownin-dest', 'api');
//   xhr.onreadystatechange = function() {
//       if (xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200) {
// 	  var scriptElement = document.createElement('script');
// 	  scriptElement.type = 'text/javascript';
// 	  scriptElement.innerHTML = xhr.responseText;
// 	  document.body.appendChild(scriptElement);
// 	  if (callback) callback();
//       }
//   };
//   xhr.send();
// }
// function loadScripts() {
//   loadAsset('/assets/main.js')
// }
// loadAsset('/assets/htmx.min.js', loadScripts)


