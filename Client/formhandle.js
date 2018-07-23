window.addEventListener("load", function (){

    let form = this.document.getElementById("searchForm");
    form.addEventListener("submit", function(e){
        e.preventDefault();
        getRates();
    });

    function getRates(){
        let httpReq = new XMLHttpRequest();

        httpReq.addEventListener("load", function(e){
           document.getElementById("rate");
           rate.innerHTML = e.target.responseText;
        });

        let startTime = document.getElementById("startTime").value;
        let endTime = document.getElementById("endTime").value;
        let year = document.getElementById("year").value;
        let month = document.getElementById("month").value;
        let day = document.getElementById("day").value;

        if(month.length < 2){
            month = "0"+month;
        }
        if(day.length < 2){
            day = "0"+day;
        }

        httpReq.addEventListener("error", function(e) {
            alert("Something has gone wrong...");
        });

        httpReq.open("GET", `http://localhost:8080/rates/${year}-${month}-${day}T${startTime}:00Z/${year}-${month}-${day}T${endTime}:00Z`);

        httpReq.send();
    }
})