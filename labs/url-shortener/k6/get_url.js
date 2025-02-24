import http from "k6/http";
import { group, check, sleep } from "k6";

const API_HOST = __ENV.API_HOST || "http://api:8080";
const SLEEP_DURATION = __ENV.SLEEP_DURATION || 1;

// https://grafana.com/docs/k6/latest/testing-guides/test-types/stress-testing/
export let options = {
	stages: [
		{ duration: "10s", target: 2_000 }, // ramp up 1 -> 200 over 10s
		{ duration: "30s", target: 2_000 }, // stay at 200 for 30s
		{ duration: "10s", target: 0 }, // ramp down 200 -> 0 over 10s
	],
};

export default function () {
	group("/v1/urls", () => {
		{
			let url = API_HOST + `/v1/test`;
			let request = http.get(url, { redirects: 0 });

			check(request, {
				OK: (r) => r.status === 302 && r.headers["Location"],
			});
		}
	});
}
