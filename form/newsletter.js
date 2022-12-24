const form = document.getElementById("newsletter-form");

form.addEventListener("submit", (event) => {
  event.preventDefault();
const projectID = 4;
  const email = document.getElementById("email").value;
  const requestBody = {
    email,
  };
  const requestOptions = {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(requestBody),
  };
  fetch(
    `http://localhost:8080/projects/${projectID}/subscriptions`,
    requestOptions
  )
    .then((response) => response.json())
    .then((data) => {
      console.log(data);
    })
    .catch((error) => {
      console.error(error);
    });
});
