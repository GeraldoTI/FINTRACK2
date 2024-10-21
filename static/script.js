document.getElementById("calc-form").addEventListener("submit", function(event) {
  event.preventDefault();

  const ec2Cost = parseFloat(document.getElementById("ec2_cost").value);
  const rdsCost = parseFloat(document.getElementById("rds_cost").value);
  const vpcCost = parseFloat(document.getElementById("vpc_cost").value);
  const totalCost = ec2Cost + rdsCost + vpcCost;

  // Update result display
  document.getElementById("result").innerText = `Total Cost: $${totalCost.toFixed(2)}`;

  // Prepare data for chart
  const labels = ["EC2", "RDS", "VPC"];
  const data = [ec2Cost, rdsCost, vpcCost];

  // Render chart
  renderChart(labels, data);
});

document.getElementById("generate-pdf").addEventListener("click", function() {
  const ec2Cost = document.getElementById("ec2_cost").value;
  const rdsCost = document.getElementById("rds_cost").value;
  const vpcCost = document.getElementById("vpc_cost").value;

  const params = new URLSearchParams({ ec2Cost, rdsCost, vpcCost });
  window.location.href = "/generate-pdf?" + params; // Triggers PDF download
});

function renderChart(labels, data) {
  const ctx = document.getElementById("costChart").getContext("2d");
  const costChart = new Chart(ctx, {
      type: 'bar',
      data: {
          labels: labels,
          datasets: [{
              label: 'Cost ($)',
              data: data,
              backgroundColor: [
                  'rgba(255, 99, 132, 0.2)',
                  'rgba(54, 162, 235, 0.2)',
                  'rgba(255, 206, 86, 0.2)',
              ],
              borderColor: [
                  'rgba(255, 99, 132, 1)',
                  'rgba(54, 162, 235, 1)',
                  'rgba(255, 206, 86, 1)',
              ],
              borderWidth: 1
          }]
      },
      options: {
          scales: {
              y: {
                  beginAtZero: true
              }
          }
      }
  });
}
