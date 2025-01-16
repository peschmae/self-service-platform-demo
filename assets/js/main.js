function addLabel() {
    const fieldset = document.querySelector('fieldset#label');
    const count = document.querySelectorAll('fieldset#label label').length;
    const label = document.createElement('label');
    label.innerHTML = 'Labels <input type="text" name="labels['+count+'][key]" placeholder="sub.example.biz/placeholder" /> <input type="text" name="labels['+count+'][value]" placeholder="value" />';
    fieldset.appendChild(label);
}
function deleteLabel() {
    const fieldset = document.querySelector('fieldset#label');
    fieldset.removeChild(fieldset.lastChild);
}

function addEgress() {
    const fieldset = document.querySelector('fieldset#egress');
    const count = document.querySelectorAll('fieldset#egress label').length;
    const label = document.createElement('label');
    label.innerHTML = 'Egress endpoints <input type="text" name="egress['+count+'][cidr]" placeholder="IP/CIDR" /> <input type="number" name="egress['+count+'][port]" placeholder="Port" />';
    fieldset.appendChild(label);
}
function deleteEgress() {
    const fieldset = document.querySelector('fieldset#egress');
    fieldset.removeChild(fieldset.lastChild);
}

function enableCheck() {
    const enableChecks = document.querySelector('input#enableChecks');
    const addCheck = document.querySelector('button[onclick="addCheck()"]');
    const deleteCheck = document.querySelector('button[onclick="deleteCheck()"]');
    if (enableChecks.checked) {
        addCheck.disabled = false;
        deleteCheck.disabled = false;
    } else {
        addCheck.disabled = true;
        deleteCheck.disabled = true;
    }
}

function addCheck() {
    const fieldset = document.querySelector('fieldset#checks');
    const count = document.querySelectorAll('fieldset#checks label').length;

    const label = document.createElement('label');

    label.innerHTML = 'Checks <input type="text" name="checks['+count+']" placeholder="https://example.com" />';
    fieldset.appendChild(label);
}

function deleteCheck() {
    const fieldset = document.querySelector('fieldset#checks');
    fieldset.removeChild(fieldset.lastChild);
}

// This function is taken from https://stackoverflow.com/a/54733755/5459839
function deepSet(obj, path, value) {
    if (Object(obj) !== obj) return obj; // When obj is not an object
    // If not yet an array, get the keys from the string-path
    if (!Array.isArray(path)) path = path.toString().match(/[^.[\]]+/g) || [];
    path.slice(0,-1).reduce((a, c, i) => // Iterate all of them except the last one
        Object(a[c]) === a[c] // Does the key exist and is its value an object?
            // Yes: then follow that path
            ? a[c]
            // No: create the key. Is the next key a potential array-index?
            : a[c] = Math.abs(path[i+1])>>0 === +path[i+1]
                ? [] // Yes: assign a new array object
                : {}, // No: assign a new plain object
        obj)[path[path.length-1]] = value; // Finally assign the value to the last key
    return obj; // Return the top-level object to allow chaining
}

// Use it for formData:
function formDataObject(form) {
    const formData = new FormData(form);
    const root = {};
    for (const [path, value] of formData) {
        deepSet(root, path, value);
    }
    return root;
}

function displayResponse(message) {
    const waiting = document.querySelector('#response');
    waiting.innerHTML = "<span>" + message + "</span>"
}

function submitJson() {
    const form = document.querySelector('form');

    form.style.display = 'none';
    const waiting = document.querySelector('#response');
    waiting.style.display = 'block';

    obj= formDataObject(form);

    console.log(obj);

    fetch(form.attributes['action'].value, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(obj),
    })
        .then((response) => response.json())
        .then((data) => {
            // should check status code and display message
            displayResponse(data['error'])
            console.log('Success:', data['error']);
        })
        .catch((error) => {
            console.error('Error:', error);
        });
}
