document.getElementById('sendMessage').addEventListener('click', sendMessage);

function sendMessage() {
  const userMessage = document.getElementById('userMessage').value;
  const model = document.getElementById('model').value;
  const userId = getUserId(); // Implement this function to get/generate a unique user ID

  fetch(`http://localhost:8080/chat/${model}`, {  // Replace with your API server URL
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'User-ID': userId,
    },
    body: JSON.stringify({ message: userMessage }),
  })
  .then(response => response.json())
  .then(data => {
    document.getElementById('response').innerText = data.response; // Display the response
  })
  .catch(error => {
    console.error('Error:', error);
    document.getElementById('response').innerText = "Error communicating with the API.";
  });
}

function getUserId() {
  // Implement logic to get or generate a user ID.
  // You can use chrome.storage.sync or chrome.storage.local
  // to store the user ID persistently.

    // Example using chrome.storage.local:
    return new Promise((resolve, reject) => {
       chrome.storage.local.get(['userId'], function(result) {
         if (result.userId) {
            resolve(result.userId);
         } else {
           const newUserId = generateUUID(); // Replace with your ID generation logic
           chrome.storage.local.set({userId: newUserId}, function() {
             resolve(newUserId);
           });
         }
       });
     });
}

function generateUUID() { // Example UUID generation (replace with a more robust solution if needed)
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
    var r = Math.random() * 16 | 0, v = c == 'x' ? r : (r & 0x3 | 0x8);
    return v.toString(16);
  });
}