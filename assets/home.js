$(document).ready(function () {
    $('#tasks').DataTable();
});

function changeStatusDone (id) {
    $.ajax({
        url: `${window.location.origin}/task/${id}`,
        type: 'PATCH',
        dataType: 'json',
        contentType: 'application/x-www-form-urlencoded',
        headers: {
            'X-CSRF-TOKEN': $('meta[name="csrf-token"]').attr('content'),
            "Access-Control-Allow-Origin":"*"
        },
        success: function(response) {
            if (response.message == "Success") {
                alert("Ubah status task berhasil!");
                location.reload();
            } else {
                alert("Error")
            }
        }
    })
}