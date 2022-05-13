import { BaseService } from "./shared.js"

class AuthService extends BaseService{
    constructor(){
        super()
    }

    

    logout(){
        let url = this.getApiUrl() + `/session`
        fetch(url,{
            method: 'DELETE',
            headers: this.getDefaultHeaders()
        }).then(()=>{
            window.location.reload(true)
        })
        
    }
}


export const authService = new AuthService()