import { BaseService } from "./shared.js"

class NoteService extends BaseService{
    constructor(){
        super()
    }

    getNoteTree(rootNoteId){
        let url = this.getApiUrl() + `/note?tree=1`
        if(rootNoteId !== null && rootNoteId !== undefined){
            url += `&rootId=${rootNoteId}`
        }

        return fetch(url,{
            method: 'GET',
            headers: this.getDefaultHeaders()
        }).then(res=> res.json())
    }

    getNote(noteId){
        return fetch(`/api/v1/note/${noteId}`,{
            method: 'GET',
            headers: this.getDefaultHeaders(),
        }).then(res=>res.json())
    }

    moveNote(noteId) {

    }

    createNote(parentId,title,content){
        return fetch("/api/v1/note",{
            method: 'POST',
            headers: this.getDefaultHeaders(),
            body: JSON.stringify({
                parentId,
                title,
                content
            })
        })
    }

    updateTitle(noteId,newTitle){
        return fetch(`/api/v1/note/${noteId}`,{
            method: 'PATCH',
            headers: this.getDefaultHeaders(),
            body: JSON.stringify({
                title:newTitle,
            })
        })
    }

    updateContent(noteId,newContent){
        return fetch(`/api/v1/note/${noteId}`,{
            method: 'PATCH',
            headers: this.getDefaultHeaders(),
            body: JSON.stringify({
                content:newContent,
            })
        })
    }
}


export const noteService = new NoteService()

