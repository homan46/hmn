import { BaseService } from "./shared.js"

class AuthService extends BaseService{
    constructor(){
        super()
    }

    

    logout(){
        let url = this.getApiUrl() + `/login/`
        return fetch(url,{
            method: 'DELETE',
            headers: this.getDefaultHeaders()
        })
    }
}


export const authService = new AuthService()