
import { h, Component, render, createRef } from 'https://unpkg.com/preact?module';
import htm from 'https://unpkg.com/htm?module';
import  'https://cdn.jsdelivr.net/simplemde/latest/simplemde.min.js'

import {AppViewModel, EditorState,EventBus} from './extra.js'

const html = htm.bind(h);

const bus = new EventBus()


export class App extends Component {
    constructor (){
        super()
        this.state = {

        }
    }

    componentDidMount(){
        
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
            currentNoteId:2
        }
    }

    render(props,state) {

        return html`
        <div class='main-container'>
            <${NavigationBar}/>
            <${NoteEditor} 
                noteId=${state.currentNoteId}
                 />
        </div>`
    }

    //
    // Callback for Note editor
    //

    onNoteContentChange = (noteID,newContent) => {        
        this.setState({editing:new EditorState(noteID,newContent)})
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
        fetch(`/api/v1/note?tree=1`,{
            method:'get',
        })
        .then((respon)=>{
            console.log("bar load")
            return respon.json()
        }).catch((err)=>{
            console.log(err)
        })
        .then((json)=>{
            console.log("bar load2")
            console.log(json)
            this.setState({treeData:json})
        })
    }


    render(props,state){
        if (state.treeData == null){
            return html`<p>loading</p>`
        } else {
            return html`
            <div class="navigation-bar">
                <${NavigationNote} treeData=${state.treeData}/>
            </div>`
        }
        
    }
}

class NavigationNote extends Component {

    //noteId

    constructor(props){
        super()
        this.noteId = props.treeData.id
        this.title = props.treeData.title
        this.content = props.treeData.content
        this.parentId = props.treeData.parentId
        this.index = props.treeData.index

    }

    render(props,state) {   
        var child
        if (props.treeData.children == null) {
            child = html``
        }else {
            child = props.treeData.children.map((val)=>{
                return html`
                <${NavigationNote} treeData=${val}/>`
            })
        }
        
        return html`
        <li>
            <span onClick=${props.openNoteHandler}>${props.treeData.title}</span>
            <span onClick=${()=>this.addUnderHandler()}>+</span>
            <ul class="children">
                ${child}
            </ul>
        </li>`
    }


    addUnderHandler = () => {
        fetch("/api/v1/note",{
            method: 'POST',
            headers: {
                'Content-Type': 'application/json;charset=utf-8'
            },
            body: JSON.stringify({
                title:this.title,
                content:this.content,
                parentId: this.id,
            })
        }).then((res)=>{
            console.log(res)
        })        
    }

    changeTitleHnadler = () =>{

    }   
}

class NoteEditor extends Component {
    ref = createRef()
    simplemde
    dirty=false
    
    constructor(props){
        super(props)
    }

    initializateEditor(){
        this.simplemde = new SimpleMDE({ element: this.ref.current });
        this.refresh()
        this.simplemde.codemirror.on("change",(instance,change)=>{
            this.dirty=true
            this.uploadContent()
        })
    }

    refresh(){
        console.log("download")
        fetch(`/api/v1/note/${this.props.noteId}`,{
            method:'get',
        }).then((respon)=>{
            return respon.json()
        }).then((json)=>{
            this.simplemde.value(json.Content)
        })
    }

    uploadContent(){
        if(this.dirty) {
            fetch(`/api/v1/note/${this.props.noteId}`,{
                method:'put',
                headers: {
                    'Content-Type': 'application/json'
                  },
                body:JSON.stringify({ID:0,Content:this.simplemde.value()})
            })
            this.dirty = false
        }
    }

    componentDidMount(){
        this.initializateEditor()
    }

    componentWillUnmount(){
        this.simplemde.toTextArea();
    }

    render(props,state){
        return html`
        <div class="note-editor">
            <textarea class="note-editor-textarea" ref=${this.ref} >
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


