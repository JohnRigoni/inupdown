//
// htmx.on('#form', 'htmx:xhr:progress', function(evt) {
//   htmx.find('#progress').setAttribute('value', evt.detail.loaded/evt.detail.total * 100)
//   if (evt.detail.loaded/evt.detail.total == 1) {
//     setTimeout(() => {
//       htmx.find('#progress').setAttribute('value', 0)
//     }, 1000);
//   }
// });
//
let global_is_uploading = false;
async function handleSubmit(_) {
    if (global_is_uploading) {
        return;
    }
    let input = document.createElement('input');
    input.type = 'file';
    input.multiple = true;
    input.onchange = async _ => await input_auto_submit(input);
    input.click();
}

async function input_auto_submit(input) {
    global_is_uploading = true;

    for (let value of input.files) {
        if (value.name === "") {
            break;
        };
        const formDataIN = new FormData();
        formDataIN.append("file", value);
        const fetchOptions = {
            method: "post",
            body: formDataIN,
        };
        let response;
        try {
            response = await fetch("/?api=upload", fetchOptions);
        }
        catch (e) {
            window.alert("Error: " + e);
            continue;
        }
        if (response.status != '200' && response.status != '201') {
            let e_str = "Error uploading: " + value.name
                + "\n\n" + response.status + " " + response.statusText.toLowerCase();
            window.alert(e_str);
            continue;
        }
    }

    global_is_uploading = false;
    htmx.trigger('#flist', 'globup')
}

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


