<template>
  <!-- Footer -->
  <footer class="footer">
    <span>Copyright: Elias Braun 2024</span>
    <a href="https://github.com/EliasB-NU/showmaster" target="_blank" class="footer-link">
      GitHub
    </a>
  </footer>
</template>

<style>
.footer {
  background-color: #fff;
  padding: 20px;
  border-radius: 15px;
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.1);
  max-width: 450px;
  margin-bottom: 10px;
  display: flex;
  justify-content: space-between;
  text-align: center;
  position: fixed; bottom: 20px; left: 50%; transform: translateX(-50%);
  min-width: fit-content;
  width: 20%;
}

.footer span {
  font-size: 14px;
  color: #555;
}

.footer-link {
  color: #007bff;
  text-decoration: none;
}

.footer-link:hover {
  text-decoration: underline;
}
</style>