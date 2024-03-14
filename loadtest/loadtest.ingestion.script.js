import http, { get } from "k6/http";

export const options = {
  vus: 10000,
  iterations: 10000,
};

const session = JSON.parse(open("./userSessions.json"));

function getRandom(min, max) {
  min = Math.ceil(min);
  max = Math.floor(max);
  return Math.floor(Math.random() * (max - min + 1) + min);
}

export default function () {
  const payload = [
    {
      payment_id: "P001",
      customer_id: null,
      payment_date: new Date(2023, getRandom(10, 12) - 1, getRandom(4, 30), getRandom(0, 23) - 16, getRandom(0, 59), getRandom(0, 59), getRandom(0, 999)).toISOString(),
      order_details: [
        {
          item_id: "I00001",
          quantity: 10,
        },
      ],
    },
  ];

  http.post("http://localhost:8082/v1/order", JSON.stringify(payload), {
    headers: {
      "Content-Type": "application/json",
      Authorization: "Bearer " + session.loadTestCashier,
    },
  });
}
