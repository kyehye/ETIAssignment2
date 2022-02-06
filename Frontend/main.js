async function displayAPIData() {
    //get API data
    const response = await fetch("http://10.31.11.11:8011/tutors");
    data = await response.json();
  
    //generate HTML code
    const tableData = data
      .map(function (value) {
        return `<tr>
              <td>${value.TutorID}</td>
              <td>${value.Name}</td>
              <td>${value.Description}</td>
          </tr>`;
      })
      .join("");
  
    //set tableBody to new HTML code
    const tableBody = document.querySelector("#tableBody");
    tableBody.innerHTML = tableData;
}
  
  displayAPIData();
  