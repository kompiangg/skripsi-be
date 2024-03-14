import http from "k6/http";
import { URL } from "https://jslib.k6.io/url/1.0.0/index.js";

export const options = {
  vus: 5,
  iterations: 20,
};

const session = JSON.parse(open("./userSessions.json"));

export default function () {
  const url = new URL("http://localhost:8085/v1/orders-details");
  url.searchParams.append("start_date", "2023-10-03T00:00:00Z");
  url.searchParams.append("end_date", "2023-12-31T23:59:59Z");

  const res = http.get(url.toString(), {
    headers: {
      "Content-Type": "application/json",
      Authorization: "Bearer " + session.admin,
    },
    timeout: "5m",
  });
}
