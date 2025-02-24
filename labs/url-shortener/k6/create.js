import http from "k6/http";
import { group, check, sleep } from "k6";

const BASE_URL = __ENV.API_HOST || "http://api:8080";
const SLEEP_DURATION = __ENV.SLEEP_DURATION || 1;

// https://grafana.com/docs/k6/latest/testing-guides/test-types/stress-testing/
export let options = {
	stages: [
		{ duration: "10s", target: 200 }, // ramp up 1 -> 200 over 10s
		{ duration: "30s", target: 200 }, // stay at 200 for 30s
		{ duration: "10s", target: 0 }, // ramp down 200 -> 0 over 10s
	],
};

function randomString(length, charset = "") {
	if (!charset) charset = "abcdefghijklmnopqrstuvwxyz";
	let res = "";
	while (length--) res += charset[(Math.random() * charset.length) | 0];
	return res;
}

export default function () {
	group("/v1/urls", () => {
		{
			let url = BASE_URL + `/v1/urls`;
			let body = { long: "https://localhost/k6/" + randomString() };
			let params = {
				headers: {
					"Content-Type": "application/json",
					Accept: "application/json",
				},
			};
			let request = http.post(url, JSON.stringify(body), params);

			check(request, {
				OK: (r) => r.status === 200 && r.body.includes("short"),
			});

			sleep(SLEEP_DURATION);
		}
	});
}
