
export class EditorState {
    constructor(noteId,content){
        this.noteId = noteId
        this.content = content
    }


}

export class NavigationBarState {
    constructor(noteId,content){
        this.titleTree = null
    }
}


export class AppViewModel {
    
    constructor(){
        this.navigationBarViewModel
        this.openedNoteId = 0
    }


}

class NavigationBarViewModel {
    constructor(){

    }
}

class NavigationIndexViewModel {
    constructor(noteId){
        
    }
}




export class EventBus {
    constructor(){
        this.eventDictionary = {}
    }

    emit(eventName,parameter){
        for(var func of this.eventDictionary[eventName]){
            func.call(null,parameter)
        }
    }
    subscribe(eventName,callback){
        if (this.eventDictionary[eventName] == undefined) {
            this.eventDictionary[eventName] = []
        }
        this.eventDictionary[eventName].append(callback)
    }
}