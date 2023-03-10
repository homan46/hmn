
import { h, Component, render, createRef } from 'https://unpkg.com/preact?module';
import htm from 'https://unpkg.com/htm?module';
import  'https://cdn.jsdelivr.net/simplemde/latest/simplemde.min.js'


import { noteService } from './service/service.js';
import { authService } from './service/auth_service.js';
import style from './style.js';

const html = htm.bind(h);

export class App extends Component {
    constructor (){
        super()
        this.state = {

        }
    }


    render (props,state){    
        return html`<${MainPage}  viewModel=${state.viewModel}/>`
    }
}

class MainPage extends Component {
    constructor(){
        super()

        this.state = {
            noteTree: null,
            currentNoteId:1,
            currentNote:{id:null}
        }
    }

    render(props,state) {

        return html`
        <div class='main-container' style=${style.mainContainer}>
            <${NavigationBar} openNoteHandler=${this.openNoteHandler}
                updateNoteTitleHandler=${this.updateNoteTitleHandler} />
            <${NoteEditor} note=${state.currentNote} />
        </div>`
    }

    //
    // Callback for Note editor
    //
    openNoteHandler = (note)=>{
        //this.setState({currentNoteId:noteId})
        this.setState({currentNote:note});
    }

    updateNoteTitleHandler = (note)=>{
        //this.setState({currentNoteId:noteId})
        if(note.id == this.state.currentNote.id) {
            var temp = this.state.currentNote
            temp.title = note.title;
            this.setState({currentNote:temp});
        }
        
    }
}


class NavigationBar extends Component {
    constructor(){
        super()

        this.state = {
            treeData:null
        }

        this.downloadIndex()
    }

    downloadIndex(){
        noteService.getNoteTree()
        .then((json)=>{
            this.setState({treeData:json})
        })
    }

    refreshHandler = ()=>{
        this.downloadIndex()
    }

    render(props,state){
        if (state.treeData == null){
            return html`<p>loading</p>`
        } else {
            return html`
            <div class="navigation-bar" style=${style.navigationBar}>
                <${NavigationNote} 
                    openNoteHandler=${props.openNoteHandler} 
                    refreshHandler=${this.refreshHandler}
                    updateNoteTitleHandler=${props.updateNoteTitleHandler}
                    treeData=${state.treeData}/>
                    <${NavigationFunctionBar} refreshHandler=${this.refreshHandler}/>
            </div>`
        }
    }

    
}

class NavigationFunctionBar extends Component {
    constructor(){
        super()
    }

    render(props,state){
        return html`
        <div style=${style.navigationFunctionBar}>
            <div style=${style.navigationBarDeleteArea}
                onDragOver=${(e)=>{e.preventDefault()}} 
                onDrop=${this.dropOnDeleteHandler}>
                <span>trash area</span>
            </div>
            
            <input type="button" value="Logout" onClick=${this.logoutHandler}/>
            
        </div>`   
    }

    dropOnDeleteHandler = (e) => {
        //noteId from the dragged note
        var noteId = e.dataTransfer.getData("noteId");
        console.log(noteId)
     
        noteService.deleteNote(noteId).then(()=>{
            this.props.refreshHandler()
        })
    }

    logoutHandler = () => {
        authService.logout()
    }
}


class NavigationNote extends Component {

    titleEditor = createRef()

    constructor(props){
        super()

        this.state = {
            noteId: props.treeData.id,
            title : props.treeData.title,
            content : props.treeData.content,
            parentId : props.treeData.parentId,
            index : props.treeData.index,

            showNewNoteInput:false,
            enableTitleEditing:false
        }

        this.singleClickTimer
    }

    render(props,state) {
        var child = html``;
        if (props.treeData.children != null) {
            child = props.treeData.children.map((val)=>{
                return html`
                <hr 
                    onDragOver=${(e)=>{e.preventDefault()}} 
                    onDrop=${(e)=>{this.dropHandler(e,val.parentId,val.index)}}
                />
                
                <${NavigationNote}
                    refreshHandler=${this.props.refreshHandler} 
                    openNoteHandler=${props.openNoteHandler} 
                    updateNoteTitleHandler=${props.updateNoteTitleHandler}
                    treeData=${val}/>
                `
            })
        }
        
        var titleSection = null;
        

        if (state.enableTitleEditing){
            console.log(`this.value is ${props.treeData.title}`);
            titleSection = html`
                <input type="text" ref=${ dom => {
                    if(dom != null){
                        dom.focus()
                    } 
                    this.titleEditor = dom
                }}  
                onBlur=${this.titleBlurHandler} value=${props.treeData.title}/>`
        }else{
            titleSection = html`
                <span draggable="true"  
                    onDragStart=${this.dragHandler} 
                    onDragOver=${(e)=>{e.preventDefault()}} 
                    onDrop=${(e)=>{this.dropHandler(e,props.treeData.id,props.treeData.children.length)}}
                    onClick=${this.titleClickHandler}>${props.treeData.id}:${props.treeData.title}
                </span>
                <span onClick=${this.addUnderHandler}>+</span>`
        }
        
        return html`
        <li  style=${style.navigationBarListItem}>
            ${titleSection}
            <ul class="children" style=${style.navigationBarList}>
                ${child}
            </ul>
        </li>`
    }

    //route single/double click to openNoteHandler or titleDoubleClickHandler
    titleClickHandler = (event)=>{
        if(event.detail === 1) {
            this.singleClickTimer = setTimeout(() => {
                this.props.openNoteHandler(this.props.treeData)
            }, 300);
        }else if (event.detail === 2){
            clearTimeout(this.singleClickTimer);
            //TODO: maybe need to clean up timer
            this.titleDoubleClickHandler()
        }
    }

    
    titleDoubleClickHandler =() => {
        this.setState({enableTitleEditing:true})
    }

    titleBlurHandler =() => {
        var changed = false
        
        if (this.titleEditor.value != this.state.title){
            changed = true
        }

        if(changed) {
            var newTitle = this.titleEditor.value;
            noteService.updateTitle(this.state.noteId,this.titleEditor.value).then(()=>{
                this.props.refreshHandler();
                
                this.props.updateNoteTitleHandler({
                    id:this.state.noteId,
                    title:newTitle
                });
            })
        }
        this.setState({enableTitleEditing:false})
    }

    addUnderHandler = () => {
        noteService.createNote(this.state.noteId,"new title","").then(()=>{
            this.props.refreshHandler()
        })
    }

    dragHandler = (e) => {
        //set noteId of the dragged note 
        console.log(this.state.noteId)
        e.dataTransfer.setData("noteId",this.state.noteId)
    }

    dropHandler = (e,parentId,index) => {
        //noteId from the dragged note
        var noteId = e.dataTransfer.getData("noteId");
        
        noteService.moveNote(noteId,parentId,index) .then(()=>{
            this.props.refreshHandler()
        })
    }

}

class NoteEditor extends Component {
    textAreaRef = createRef()
    titleRef = createRef()
    simplemde
    dirty=false
    
    constructor(props){
        super(props)
    }

    initializateEditor(){
        this.simplemde = new SimpleMDE({ element: this.textAreaRef.current });

        //assume data is not ready when first initialize
        //this.downloadContent()
        this.simplemde.codemirror.on("change",(instance,change)=>{
            this.dirty=true
            this.uploadContent()
        })


        this.simplemde.codemirror.options.readOnly = true
        this.titleRef.current.innerHTML = "please select a note to edit"
    }

    downloadContent(){
        noteService.getNote(this.props.note.id).then(json => {
            this.simplemde.value(json.content)
            this.titleRef.current.innerHTML = json.title
        })
    }

    uploadContent(){
        if(this.dirty) {
            noteService.updateContent(this.props.note.id,this.simplemde.value())
            this.dirty = false
        }
    }

    shouldComponentUpdate(nextProps, nextState){
        if(nextProps.note.id > 0) { //TODO: improve valid noteId checking
            this.simplemde.codemirror.options.readOnly = false
        }
        
        noteService.getNote(nextProps.note.id).then(json => {
            this.simplemde.value(json.content)
            this.titleRef.current.innerHTML = json.title
        })
        return false
    }

    componentDidMount(){
        this.initializateEditor()
    }

    componentWillUnmount(){
        this.simplemde.toTextArea();
    }

    render(props,state){
        return html`

        <div class="note-editor" style=${style.noteEditor}>
            <p class="note-editor-title" ref=${this.titleRef}></p>
            <textarea class="note-editor-textarea" ref=${this.textAreaRef} >
            </textarea>
        </div>`
    }
}


class LoginWidget extends Component {
    constructor(){
        super()
        this.state = {
            username:"",
            password:""
        }
    }
    render(){
        return html`<p>This is login widget</p>`
    }
}

//import { html/*, render, Component*/} from 'https://unpkg.com/htm/preact/index.mjs?module'
//import * as bb from 'https://unpkg.com/htm/preact/index.mjs?module'
//import { h, Component, render, createRef } from 'https://unpkg.com/preact?module';


