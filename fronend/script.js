getData();

setInterval(getData, 10000);


// Валидация поля ввода

function validateExpression(input) {
  // Используем регулярное выражение для проверки
  var regex = /^[0-9+\-*/.() ]*$/;
  return regex.test(input) && input.length > 2;
}

let exprInput = document.querySelector(".expr-input");
exprInput.addEventListener('input', ()=>{
    if (!validateExpression(exprInput.value)){
      document.querySelector(".input-container").style.borderBottom = "2px solid red";
      sendButton.style.opacity = "0.5";
    }else{
      document.querySelector(".input-container").style.borderBottom = "2px solid black";
      sendButton.style.opacity = "1";
    }
});


// отправка выражения на сервер

let sendButton = document.querySelector(".send-btn")

sendButton.addEventListener('click', sendData);


function sendData() {

  if (!validateExpression(exprInput.value)){
    alert("Invalid expression");
    return
  }
  let data = {
    expression: document.querySelector(".expr-input").value
  };

  // Отправка данных на сервер
  fetch('http://localhost:8080/postData', {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      'Access-Control-Allow-Methods': 'PUT',
      'Access-Control-Allow-Origin': '*',
    },
    body: JSON.stringify(data),
  })
  .then(response => {
    // Обработка ответа от сервера
    if (response.status == 200) {
      getData();
    } else {
      alert("Failed to handle data from server");
    }
  })
  .catch(error => {
    alert("Failed to send data to server");
  });
}

function getData() {
      // Получение данных с сервера
    fetch('http://localhost:8080/getData')
    .then(response => response.json())
    .then(data => {
      // Создание списка на HTML-странице
      let exprContainer = document.querySelector(".data-container");
      exprContainer.innerHTML = "";
      console.log(data)
      data.forEach(item => {
          CreateExression(item.Status, item.expression, item.Result);
      });
    })
    .catch(error => {
      let exprContainer = document.querySelector(".data-container");
      exprContainer.innerHTML = "failed to fetch data"+error;
      exprContainer.style.color = "red";
    });
}

function CreateExression(status, expr, result){
    let exprContainer = document.createElement("div");
    let exprImg = document.createElement("img");
    exprImg.classList.add("expr-img-status");
    if (status=="done"){
        exprImg.src = "imgs/tick.png";
    }else if (status == "waiting"){
        exprImg.src = "imgs/time.png";
    }else if (status == "failed"){
        exprImg.src = "imgs/x.png";
    }
    let exprText = document.createElement("p");
    exprText.classList.add("expr-text")
    exprText.innerHTML = expr+"="+result;
    exprContainer.appendChild(exprImg);
    exprContainer.appendChild(exprText);
    exprContainer.classList.add("expression");
    document.querySelector(".data-container").appendChild(exprContainer);
}