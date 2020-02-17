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
}
