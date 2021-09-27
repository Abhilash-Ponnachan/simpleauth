const ST_VAL_READ = 'READ';
const ST_VAL_ADD = 'ADD';
const ST_VAL_UPD = 'UPDATE';
const ST_VAL_DEL = 'DELETE';
class State{
    constructor(){
        this.value = ST_VAL_READ;
        // default value
    }

    setState(value){
        if ((value === ST_VAL_READ) 
            || (value === ST_VAL_ADD) 
            || (value === ST_VAL_UPD) 
            || (value === ST_VAL_DEL)) {
                this.value = value;

        } else{
            throw new Error(`Inavlid state '${value}'.`);
        }
    }

    getState(){
        return this.value;
    }
}