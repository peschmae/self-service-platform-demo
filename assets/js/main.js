function addLabel() {
    const fieldset = document.querySelector('fieldset#label');
    const label = document.createElement('label');
    label.innerHTML = 'Labels <input type="text" name="labels[]" placeholder="key=value" />';
    fieldset.appendChild(label);
}
function deleteLabel() {
    const fieldset = document.querySelector('fieldset#label');
    fieldset.removeChild(fieldset.lastChild);
}

function addEgress() {
    const fieldset = document.querySelector('fieldset#egress');
    const label = document.createElement('label');
    label.innerHTML = 'Egress endpoints <input type="text" name="egress[]" placeholder="IP/CIDR:Port" />';
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
    const label = document.createElement('label');
    label.innerHTML = 'Checks <input type="text" name="checks[]" placeholder="https://example.com" />';
    fieldset.appendChild(label);
}

function deleteCheck() {
    const fieldset = document.querySelector('fieldset#checks');
    fieldset.removeChild(fieldset.lastChild);
}
