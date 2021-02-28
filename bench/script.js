import http from 'k6/http';
import { check, sleep } from 'k6';


export default function () {
    var url = `http://host.docker.internal:8080/${__ENV.FUNC_NAME}`;
    var payload = JSON.stringify({
        name: "item-" + __VU,
    });
    var params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };
    let res = http.post(url, payload, params);
    check(res, { 'status was 200': (r) => r.status == 200 });
}
