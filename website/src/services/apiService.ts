import { api } from "../proto/generated/common"
import { google }from "../proto/generated/google/protobuf/timestamp"

class APIService {
    private client = new api.ApiClient('http://localhost:1337')

    private targets = {
        auth: "user-service",
        publisher: "expense-publisher",
        reader: "expense-reader"
    }
    
    public register(username: string, password: string, requestId: string): Promise<api.RegisterResponse> {
        let msg = {
            username: username,
            password: password,
            requestId: requestId
        }
        let meta = {
            "target": this.targets.auth
        }
        let request = new api.RegisterRequest(msg)
        return new Promise((resolve, reject) => {
            this.client.Register(request, meta, (err, response) => {
                if (err) {
                    reject(err)
                } else {
                    resolve(response)
                }
            })
        })
    }

    public login(username: string, password: string, requestId: string): Promise<api.LoginResponse> {
        let msg = {
            username: username,
            password: password,
            requestId: requestId
        }
        let meta = {
            "target": this.targets.auth
        }
        let request = new api.LoginRequest(msg)
        return new Promise((resolve, reject) => {
            this.client.Login(request, meta, (err, response) => {
                if (err) {
                    reject(err)
                } else {
                    resolve(response)
                }
            })
        })
    }

    public createBill(name: string, amount: number, category: string, user_id: string, ts: number, jwt: string): Promise<boolean> {
        let timestamp = new google.protobuf.Timestamp({
            seconds: ts,
            nanos: 0
        })
        let msg = {
            name: name,
            amount: amount,
            category: category,
            user_id: user_id,
            timestamp: timestamp,
            jwt: jwt
        }
        let meta = {
            "target": this.targets.publisher
        }
        let request = new api.CreateBillMessage(msg)
        return new Promise((resolve, reject) => {
            this.client.CreateBill(request, meta, (err, _) => {
                if (err) {
                    reject(err)
                } else {
                    resolve(true)
                }
            })
        })
    }


    public getReport(jwt: string): Promise<api.GetReportResponse> {
        let msg = {
            jwt: jwt
        }
        let meta = {
            "target": this.targets.reader
        }
        let request = new api.GetReportRequest(msg)
        return new Promise((resolve, reject) => {
            this.client.GetReport(request, meta, (err, request) => {
                if (err) {
                    reject(err)
                } else {
                    resolve(request)
                }
            })
        })
    }

    public getBills(jwt: string, category: string): Promise<api.GetBillsResponse> {
        let msg = {
            jwt: jwt,
            category: category
        }
        let meta = {
            "target": this.targets.reader
        }
        let request = new api.GetBillsRequest(msg)
        return new Promise((resolve, reject) => {
            this.client.GetBills(request, meta, (err, request) => {
                if (err) {
                    reject(err)
                } else {
                    resolve(request)
                }
            })
        })
    }
}

export const apiService: APIService = new APIService() 