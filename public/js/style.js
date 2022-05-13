 var mainContainer = {
    position: "absolute",
    top: 0,
    bottom: 0,
    left: 0,
    right: 0,
    display: "grid",
    gridTemplateColumns:  "200px auto",
}

var navigationBar = {
    border:"solid 1px black",
    overflowX: "hidden",
    position:"relative"
}

var navigationFunctionBar = {
    width:"100%",
    position:"absolute",
    bottom:0,
    display: "grid",
    gridTemplateColumns:  "50% 50%",
}

var navigationBarDeleteArea = {
    textAlign:"center",
    backgroundColor:"gray"
}

var navigationBarLogout = {
    textAlign:"center"
}

var navigationBarList= {
    margin: 0,
    padding: 0,
    paddingLeft:"1em"
}

var navigationBarListItem = {
    listStyleType: "none"
}

var noteEditor = {
    border:"solid 1px black"
}



 export default {
    mainContainer,
    navigationBar,
    navigationFunctionBar,
    navigationBarDeleteArea,
    navigationBarLogout,
    navigationBarList,
    navigationBarListItem,
    noteEditor
}