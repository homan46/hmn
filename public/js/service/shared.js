

export class BaseService {
    csrfToken = ""
    

    getDefaultHeaders(){
        return {
            'Content-Type': 'application/json;charset=utf-8',
            'X-XSRF-TOKEN':this.csrfToken
        }
    }

    getApiUrl(version=1){
        return `/api/v${version}`
    }

    setCsrf(csrfToken){
        this.csrfToken = csrfToken
    }
}