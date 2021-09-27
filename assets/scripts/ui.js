
const clsModalActive = 'modal-active';
const clsModalInActive = 'modal-inactive';

// function to show and hide the modal element
function openModal(modal){
    setSaveBtnState();
    modal.classList.remove(clsModalInActive);
    modal.classList.add(clsModalActive);
}

function setSaveBtnState(){
    const saveBtn = document.getElementById('btn-edit-save');
    if (saveBtn != null){
        saveBtn.style.display = 'initial';
        switch(STATE.getState()){
            case ST_VAL_ADD: 
                saveBtn.innerText = "Add \u2713";
                break;
            case ST_VAL_UPD:
                saveBtn.innerText = "Update \u2713";
                break;
            case ST_VAL_READ:
                saveBtn.style.display = 'none';
                break;
        }
    }
}

function closeModal(modal){
    modal.classList.remove(clsModalActive);
    modal.classList.add(clsModalInActive);
}

function resetControls(){
    const txtRecId = document.getElementById('edit-rec-id');
    const txtAuthor = document.getElementById('edit-value-author');
    const divQuote = document.getElementById('edit-value-quote');
    txtRecId.value = '';
    txtAuthor.value = '';
    while(divQuote.firstChild){
        divQuote.removeChild(divQuote.firstChild);
    }
}

function loadRecord(recId, rowIndex, author, quote){
    const txtRecId = document.getElementById('edit-rec-id');
    const txtRowIndex = document.getElementById('edit-row-index');
    const txtAuthor = document.getElementById('edit-value-author');
    const divQuote = document.getElementById('edit-value-quote');
    if (recId != null){
        txtRecId.value = recId;
    }
    if (rowIndex != null){
        txtRowIndex.value = rowIndex;
    }
    if (author != null){
        txtAuthor.value = author;
    }
    if (quote != null){
        divQuote.innerText = quote;
    }
    
}

function initUI(saveAction, randomQuote){
    // get the dom elements
    const btnAddNew = document.getElementById('btn-add-quote');
    const btnEditCancel = document.getElementById('btn-edit-cancel');
    const btnEditSave = document.getElementById('btn-edit-save');
    const btnInspireMe = document.getElementById('btn-inspire');
    const divModal = document.getElementById('edit-modal');

    // wire up the event handlers
    btnAddNew.addEventListener('click', ()=> {
        STATE.setState(ST_VAL_ADD);
        openModal(divModal);
        resetControls();
    });
    btnEditCancel.addEventListener('click', ()=> {
        closeModal(divModal);
        resetControls();
        STATE.setState(ST_VAL_READ);
    });
    if (saveAction !== null){
        btnEditSave.addEventListener('click', () => {
           saveAction();
       });
    }
    if (randomQuote !== null){
        btnInspireMe.addEventListener('click', () => {
            STATE.setState(ST_VAL_READ);
            randomQuote();
            setSaveBtnState();
            openModal(divModal);
       });
    }
};