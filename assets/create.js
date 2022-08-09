$(document).ready(function() {
    $('#create-task').on('submit', async function(e) {
        e.preventDefault()
        
        var elements = document.getElementById("create-task").elements;

        var data ={};
        for(var i = 0 ; i < elements.length ; i++){
            var item = elements.item(i);
            data[item.name] = item.value;
        }
        
        $.ajax({
            url: `${window.location.origin}/task`,
            type: 'POST',
            dataType: 'json',
            contentType: 'application/x-www-form-urlencoded',
            headers: {
                'X-CSRF-TOKEN': $('meta[name="csrf-token"]').attr('content'),
                "Access-Control-Allow-Origin":"*"
            },
            data,
            success: function(response) {
                if (response.message == "Success") {
                    alert('Tambah Task Sukses!!')
                    location.href="/"
                } else {
                    alert("Error")
                }
            }
        })

        return false; // return false to cancel form action
    });
})
