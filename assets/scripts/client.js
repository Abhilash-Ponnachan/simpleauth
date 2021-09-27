
const HOST = location.protocol + "//" + window.location.hostname
+ (window.location.port ? ":" + window.location.port : "");

const API_LIST_QUOTES = HOST + "/api/listquotes";
const API_ADD_QUOTE = HOST + "/api/addquote";
const API_DELETE_QUOTE = HOST + "/api/deletequote";
const API_UPDATE_QUOTE = HOST + "/api/updatequote";
const API_RANDOM_QUOTE = HOST + "/api/randomquote";
const STATE = new State();

function isEmptyObj(obj){
    return (obj == null) || (Object.keys(obj).length === 0 && obj.constructor === Object)
}

async function getJSONData(url = '', qryParams = {}) {
    if (!url) {
        throw new Error("URl cannot be empty!!")
    }
    const urlObj = new URL(url)
    if (!isEmptyObj(qryParams)){
        urlObj.search = new URLSearchParams(qryParams).toString();
    }
    const response = await fetch(urlObj)
        .catch(err => console.log(err));
    return await response.json()
}

async function postJSONData(url = '', data = {}) {
    if (!url) {
        throw new Error("URl cannot be empty!!")
    }

    const response = await fetch(url, {
        method: "POST",
        body: JSON.stringify(data),
        headers: {"Content-type": "application/json; charset=UTF-8"}
      })
      .catch(err => console.log(err));
    return response;
}

function loadQuotes(){
    const quotes = getJSONData(API_LIST_QUOTES, null);
    quotes.then(respData => {
        if (respData.error){
            console.log(respData.msg)
        } else {
            // populate table with response data
            const table = getQuoteTable();
            emptyTable(table);
            populateTable(table, respData);
        }
    }).catch(err => {
            console.log(err);
    });
}

function emptyTable(table){
    const rowCount = table.rows.length;
    for (let i=rowCount-1; i > 0; i--){
        table.deleteRow(i);
    }
}

function populateTable(table, data){
    let rc = table.rows.length;
    for(let rec of data){
        addRow(table, rc, rec);
        rc+=1;
    }
}

function addRow(table, rowCount, dataRec){
    const row = table.insertRow(rowCount);
    const cell0 = row.insertCell(0);
    const cell1 = row.insertCell(1);
    const cell2 = row.insertCell(2);
    const cell3 = row.insertCell(3);
    cell0.className = "col_0";
    cell1.className = "col_1";
    cell2.className = "col_2";
    cell3.className = "col_3";
    cell0.innerHTML = dataRec.Id;
    cell1.innerHTML = `<q>${dataRec.Quote}</q>`;
    cell2.innerHTML = dataRec.Author;
    // add edit & delete buttons in the last cell
    const btnEdt = document.createElement('input');
    btnEdt.type = "button";
    btnEdt.className = "btn-edit";
    btnEdt.value = "Edit (-)";
    btnEdt.onclick = editRec;
    // custom data attribute for rec id
    btnEdt.dataRecId = dataRec.Id;
    
    const btnDel = document.createElement('input');
    btnDel.type = "button";
    btnDel.className = "btn-delete";
    btnDel.value = "Delete (x)";
    btnDel.onclick = deleteRec;
    // custom data attribute for rec id
    btnDel.dataRecId = dataRec.Id;
    
    cell3.appendChild(btnEdt);
    cell3.appendChild(btnDel);
}

function setRow(table, rowIndex, data){
    if (table?.rows && table.rows.length > rowIndex){
        const row = table.rows[rowIndex];
        const colQuote = row.childNodes[1];
        const colAuthor = row.childNodes[2];
        colQuote.innerHTML = `<q>${data.Quote}</q>`;
        colAuthor.innerText = data.Author;
    }
}

function editRec(event){
    const trgt = event?.target;
    const typ = trgt?.type;
    if (typ && typ === "button"){
        // get rowId from custom attribute
        const rowId = trgt.dataRecId;
        if (rowId != null){
            const [author, quoute] = getRowValues(trgt);
            const rowIndex = getRowIndex(trgt)
            loadRecord(rowId, rowIndex, author, quoute);
            const divModal = document.getElementById('edit-modal');
            STATE.setState(ST_VAL_UPD);
            openModal(divModal);
        }
    }
}

function removeRow(table, rowIndex){
    console.log(`Removing row @ index = ${rowIndex}`);
    if (table?.rows && table.rows.length > rowIndex){
        const row = table.rows[rowIndex];
        row.classList.add('delete_row_effect');
        setTimeout(() => {
            row.classList.remove('delete_row_effect');
            table.deleteRow(rowIndex);
        },
        300);
    }
}

function deleteRec(event){
    const trgt = event?.target;
    const typ = trgt?.type;
    if (typ && typ === "button"){
        // get rowId from custom attribute
        const rowId = trgt.dataRecId;
        if (rowId != null){
            //console.log(rowId);
            const reqData = {
                RowId: rowId
            };
            postJSONData(API_DELETE_QUOTE, reqData)
            .then(resp => {
                if (resp.ok){
                    const table = getQuoteTable();
                    // use actual table row index
                    // to remove row from table 
                    removeRow(table, getRowIndex(trgt));
                }
            }).catch(err => {
                console.log(err);
            });
        }
    }
}

function addRecAPI(){
    // ignore id & index when destructuring
    const [, , txtAuthor, divQuote] = getEditValues();
    const data = {
        Id: 0,  // dummy Id for new record
        Author: txtAuthor.value,
        Quote: divQuote.innerText
    };
    postJSONData(API_ADD_QUOTE, data)
    .then(resp => {
        if (resp.ok){
            const contType = resp.headers.get("content-type");
            if (contType 
                && contType.indexOf("application/json") !== -1){
                    resp.json()
                        .then(respData => {
                            //console.log(respData);
                            const table = getQuoteTable();
                            const rc = table.rows.length;
                            data.Id = respData.LastRowId;
                            addRow(table, rc, data);
                        })
                }
        }
    });
}

function updateRecAPI(){
    const [txtRecId, txtRowIndex, txtAuthor, divQuote] = getEditValues();
    const data = {
        Id: parseInt(txtRecId.value),
        Author: txtAuthor.value,
        Quote: divQuote.innerText
    };
    //console.log(data);
    postJSONData(API_UPDATE_QUOTE, data)
    .then(resp => {
        if (resp.ok){
            const contType = resp.headers.get("content-type");
            if (contType 
                && contType.indexOf("application/json") !== -1){
                    resp.json()
                        .then(respData => {
                            //console.log(respData);
                            const table = getQuoteTable();
                            // TO DO : update row with saved data
                            // if rows affected in
                            const rowIndex = parseInt(txtRowIndex.value);
                            setRow(table, rowIndex, data);
                        })
                }
        }
    });
}

function saveData(){
    switch (STATE.getState()){
        case ST_VAL_ADD: 
            addRecAPI();
            break;
        case ST_VAL_UPD:
            updateRecAPI();
            break;
    }

}

function randomQuote(){
    const quote = getJSONData(API_RANDOM_QUOTE, null);
    quote.then(respData => {
        if (respData.error){
            console.log(respData.msg)
        } else {
            // show modal window with random quote
           //console.log(respData);
           const id = respData.Id;
           let author = respData.Author;
           let quote = respData.Quote;
           if (!quote){
               quote = "<Works better, if you add some quotes!>";
               author = "<The App>";
           }
           
           loadRecord(id, null, author, quote);
           //console.log(id, author, quote);
        }
    }).catch(err => {
            console.log(err);
    });
}

function getEditValues(){
    const txtId = document.getElementById('edit-rec-id');
    const txtRowIndex = document.getElementById('edit-row-index');
    const txtAuthor = document.getElementById('edit-value-author');
    const divQuote = document.getElementById('edit-value-quote');
    return [txtId, txtRowIndex, txtAuthor, divQuote];
}

function getRowValues(btn){
    let author, quote;
    if ((btn.parentNode?.nodeName == 'TD')
       && (btn.parentNode.parentNode?.nodeName == 'TR')){
           row = btn.parentNode.parentNode;
           if (row.childNodes?.length > 2){
               const col01 = row.childNodes[1];
               const col02 = row.childNodes[2]
               author = col02.innerText;
               quote = col01.innerText;
           }
    }
    return [author, quote];
}

function getQuoteTable(){
    const table = document.getElementById("quotes");
    return table;
}

function getRowIndex(btn){
    let index;
    if ((btn?.type == 'button')
       && (btn.parentNode?.nodeName == 'TD')
       && (btn.parentNode.parentNode?.nodeName == 'TR')){
           index = btn.parentNode.parentNode.rowIndex;
    }
    return index;
}

// things to do on load
window.onload = () => {
    loadQuotes();
    initUI(saveData, randomQuote);
}