<template>
    <div class="text-box">
      <textarea
        v-model="text"
        placeholder="Describe your task or issue here..."
        rows="5"
      ></textarea>
    </div>
  </template>
  
  <script>
  export default {
    name: 'TextBox',
    data() {
      return {
        text: ''
      };
    }
  };
  </script>
  
  <style scoped>
  .text-box {
    margin: 20px 0;
    background-color: #1e1e1e;
    border-radius: 10px;
    padding: 15px;
    box-shadow: 0 2px 5px rgba(0, 0, 0, 0.3);
  }
  
  .text-box textarea {
    width: 100%;
    border: none;
    border-radius: 5px;
    padding: 10px;
    background-color: #333;
    color: #fff;
    resize: none;
    font-family: inherit; /* Use inherited font */
  }
  
  .text-box textarea::placeholder {
    color: #bbb;
  }
  </style>
  