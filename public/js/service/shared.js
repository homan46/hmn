

export class BaseService {
    getDefaultHeaders(){
        return {
            'Content-Type': 'application/json;charset=utf-8'
        }
    }

    getApiUrl(version=1){
        return `/api/v${version}`
    }
}