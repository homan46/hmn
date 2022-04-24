import { BaseService } from "./shared.js"

export class NoteService extends BaseService{
    constructor(){
        super()
    }

    getNoteTree(rootNoteId){
        let url = this.getApiUrl() + `/note?tree=1`
        if(rootNoteId !== null && rootNoteId !== undefined){
            url += `&rootId=${rootNoteId}`
        }

        return fetch(url,{
            method: 'POST',
            headers: this.getDefaultHeaders()
        }).then(res=> res.json())
    }

    getNote(){}

    updateNote(){

    }

    createNote(parentId){

    }

    updateTitle(noteId,newTitle){
        fetch(`/api/v1/note/${noteId}`,{
                method: 'PATCH',
                headers: this.getDefaultHeaders(),
                body: JSON.stringify({
                    title:newTitle,
                    extra:1882
                })
            }).then((res)=>{
                return res.json()
            })
    }
}