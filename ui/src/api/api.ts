import axios from "axios";

export default class Api {
    static async get(url: string) {
        return await axios.create({
            validateStatus: function (status) {
                return status === 200;
            }
        }).get(url, {
            headers: {
                "Accept": "application/json"
            }
        });
    }

    static async post(url: string, body: string) {
        return await axios.create({
            validateStatus: function (status) {
                return status === 200;
            }
        }).post(url, body,{
            headers: {
                "Accept": "application/json"
            }
        });
    }
}
