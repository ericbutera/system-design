import http from "k6/http";
import { check } from "k6";

const API_HOST = __ENV.API_HOST || "http://api:8080";

export const options = {
	thresholds: {
		http_req_failed: ["rate<0.01"], // http errors should be less than 1%
	},
	scenarios: {
		contacts: {
			executor: "constant-vus",
			vus: 75,
			duration: "20s",
		},
	},
	noConnectionReuse: true,
};

export default function () {
	const res = http.get(API_HOST);
	const success = check(res, {
		"status is 200": (r) => r.status === 200,
	});
}
