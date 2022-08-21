let data = ['Ram', 'Shyam', 'Sita', 'Gita' ];

let list = document.getElementById("local-machine-folders");

data.forEach((item)=>{
    let li = document.createElement("li");
    li.innerText = item;
    list.appendChild(li);
})
