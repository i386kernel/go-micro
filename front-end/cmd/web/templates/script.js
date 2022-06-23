let brokerBtn = document.getElementById("brokerBtn")
let authBrokerBtn = document.getElementById("authBrokerBtn")
let logBtn = document.getElementById("logBtn")
let output = document.getElementById("output")
let sent = document.getElementById("payload")
let received = document.getElementById("received")

authBrokerBtn.addEventListener("click", function () {
    const payload = {
        action: "auth",
        auth: {
            email: "admin@example.com",
            password: "verysecret",
        }
    }
    const headers = new Headers();
    headers.append("Content-Type", "applciation/json");
    const body = {
        method: 'POST',
        body: JSON.stringify(payload),
        headers: headers,
    }
    fetch("http:\/\/localhost:8080/handle", body)
        .then((response) => response.json())
        .then((data) => {
            sent.innerHTML = JSON.stringify(payload, undefined, 4);
            received.innerHTML = JSON.stringify(data, undefined, 4);
            if (data.error) {
                output.innerHTML += `<br><strong>Error:</strong> ${data.message}`;
            } else {
                output.innerHTML += `<br><strong>Response from broker service</strong>: ${data.message}`;
            }
        })
        .catch((error) => {
            output.innerHTML += "<br><br>Error: " + error;
        })
})

brokerBtn.addEventListener("click", function () {
    const body = {
        method: 'POST',
    }
    fetch("http:\/\/localhost:8080", body)
        .then((response) => response.json())
        .then((data) => {
            sent.innerHTML = "empty post request";
            received.innerHTML = JSON.stringify(data, undefined, 4);
            if (data.error) {
                console.log(data.message);
            } else {
                output.innerHTML += `<br><strong>Response from broker service</strong>: ${data.message}`;
            }
        })
        .catch((error) => {
            output.innerHTML += "<br><br>Error: " + error;
        })
})

logBtn.addEventListener("click", function (){
    const payload = {
        action: "log",
        log:{
            name: "event",
            data: "Some kind of data"
        }
    }

    const headers = new Headers();
    headers.append("Content-Type", "application/json")

    const body = {
        method: "POST",
        body: JSON.stringify(payload),
        headers: headers,
    }

    fetch("http:\/\/localhost:8080/handle", body)
        .then((response) => response.json())
        .then((data) => {
            sent.innerHTML = "empty post request";
            received.innerHTML = JSON.stringify(data, undefined, 4);
            if (data.error) {
                console.log(data.message);
            } else {
                output.innerHTML += `<br><strong>Response from broker service</strong>: ${data.message}`;
            }
        })
        .catch((error) => {
            output.innerHTML += "<br><br>Error: " + error;
        })

})