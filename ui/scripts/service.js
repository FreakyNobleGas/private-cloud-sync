let data = ['Ram', 'Shyam', 'Sita', 'Gita' ];

let list = document.getElementById("local-machine-folders");

httpOptions = {
    mode: 'cors',
    headers: {
        'Content-Type': 'application/json',
    }
}

fetch('http://localhost:5000/getAllFiles', httpOptions)
    .then( (response) => response.json() )
    .then((data) => {

        data['Files'].forEach((item)=>{
            let li = document.createElement("li");
            li.innerText = item;
            list.appendChild(li);
        });
    });



