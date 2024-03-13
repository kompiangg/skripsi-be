import http from "k6/http";
import axios from "axios";
import dayjs from "dayjs";

export const options = {
  vus: 100,
  iterations: 1000,
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
  const payload = [
    {
      payment_id: "P001",
      customer_id: null,
      currency: "BDT",
      payment_date: dayjs(`2023-${getRandom(6, 12)}-${getRandom(1, 30)} ${getRandom(0, 23)}:${getRandom(0, 59)}:${getRandom(0, 59)}`).format("YYYY-MM-DDTHH:mm:ssZ"),
      order_details: [
        {
          item_id: "I00045",
          quantity: 2,
        },
      ],
    },
  ];

  http.post("http://localhost:8082/v1/order", JSON.stringify(payload.orderPayload), {
    headers: {
      "Content-Type": "application/json",
      Authorization: "Bearer " + data.accessToken,
    },
  });
}

function getRandom(min, max) {
  min = Math.ceil(min);
  max = Math.floor(max);
  return Math.floor(Math.random() * (max - min + 1) + min);
}
