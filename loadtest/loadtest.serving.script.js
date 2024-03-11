import http from "k6/http";
import { URL } from "https://jslib.k6.io/url/1.0.0/index.js";

export const options = {
  vus: 20,
  iterations: 5000,
};

export async function setup() {
  const authRes = await axios.post(
    "http://localhost:8081/v1/session",
    {
      username: "loadtest",
      password: "loadtest",
    },
    {
      headers: {
        "Content-Type": "application/json",
      },
    }
  );

  return {
    accessToken: authRes.data.data.token,
  };
}

export default function (data) {
  const url = new URL("http://localhost:8085/v1/orders");
  url.searchParams.append("start_date", "2023-09-02T00:00:00Z");
  url.searchParams.append("end_date", "2023-12-31T23:59:59Z");

  http.post(url.toString(), JSON.stringify(payload.orderPayload), {
    headers: {
      "Content-Type": "application/json",
      Authorization: "Bearer " + data.accessToken,
    },
  });
}
